package service

import (
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/go-ldap/ldap/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/cashier-config/server/internal/config"
	"github.com/cashier-config/server/internal/model"
)

type LDAPService struct {
	db  *gorm.DB
	cfg *config.LDAPConfig
}

func NewLDAPService(db *gorm.DB, cfg *config.LDAPConfig) *LDAPService {
	return &LDAPService{db: db, cfg: cfg}
}

func (s *LDAPService) Authenticate(username, password string) (*model.User, error) {
	conn, err := s.dial()
	if err != nil {
		return nil, fmt.Errorf("ldap dial failed: %w", err)
	}
	defer conn.Close()

	if s.cfg.BindDN != "" {
		if err := conn.Bind(s.cfg.BindDN, s.cfg.BindPass); err != nil {
			return nil, fmt.Errorf("ldap bind failed: %w", err)
		}
	}

	filter := fmt.Sprintf("(%s=%s)", s.cfg.UIDAttr, ldap.EscapeFilter(username))
	searchReq := ldap.NewSearchRequest(
		s.cfg.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{s.cfg.UIDAttr, s.cfg.NameAttr, s.cfg.MailAttr, "dn"},
		nil,
	)
	result, err := conn.Search(searchReq)
	if err != nil {
		return nil, fmt.Errorf("ldap search failed: %w", err)
	}
	if len(result.Entries) == 0 {
		return nil, errors.New("invalid credentials")
	}

	entry := result.Entries[0]
	if err := conn.Bind(entry.DN, password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	user, err := s.findOrCreateUser(username, entry)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *LDAPService) dial() (*ldap.Conn, error) {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	if s.cfg.StartTLS {
		conn, err := ldap.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
		if err := conn.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
			conn.Close()
			return nil, err
		}
		return conn, nil
	}
	return ldap.Dial("tcp", addr)
}

func (s *LDAPService) findOrCreateUser(username string, entry *ldap.Entry) (*model.User, error) {
	var user model.User
	err := s.db.Preload("Role").Where("username = ?", username).First(&user).Error
	if err == nil {
		if user.Status != 1 {
			return nil, errors.New("account disabled")
		}
		user.LDAPDN = entry.DN
		s.db.Model(&user).Update("ldap_dn", entry.DN)
		return &user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	name := entry.GetAttributeValue(s.cfg.NameAttr)
	mail := entry.GetAttributeValue(s.cfg.MailAttr)
	if name == "" {
		name = username
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("ldap-"+username), bcrypt.DefaultCost)
	user = model.User{
		Username:     username,
		PasswordHash: string(hash),
		Name:         name,
		Email:        mail,
		RoleID:       s.cfg.DefaultRoleID,
		AuthSource:   "ldap",
		LDAPDN:       entry.DN,
		Status:       1,
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}
	s.db.Preload("Role").First(&user, user.ID)
	return &user, nil
}
