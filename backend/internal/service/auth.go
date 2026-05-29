package service

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/cashier-config/server/internal/config"
	"github.com/cashier-config/server/internal/model"
)

type AuthService struct {
	db  *gorm.DB
	cfg config.JWTConfig
}

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	RoleCode string `json:"role_code"`
	jwt.RegisteredClaims
}

func NewAuthService(db *gorm.DB, cfg config.JWTConfig) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

func (s *AuthService) generateToken(user *model.User) (string, error) {
	roleCode := ""
	if user.Role != nil {
		roleCode = user.Role.Code
	}
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RoleCode: roleCode,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.cfg.TTL) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.Secret))
}

func (s *AuthService) Login(username, password string) (string, *model.User, error) {
	var user model.User
	if err := s.db.Preload("Role").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("invalid credentials")
		}
		return "", nil, err
	}
	if user.Status != 1 {
		return "", nil, errors.New("account disabled")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}
	token, err := s.generateToken(&user)
	if err != nil {
		return "", nil, err
	}
	return token, &user, nil
}

func (s *AuthService) GenerateTokenFromClaims(claims *Claims) (string, error) {
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(s.cfg.TTL) * time.Hour))
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.Secret))
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (s *AuthService) GetPermissions(userID int64) ([]string, error) {
	var user model.User
	if err := s.db.Preload("Role").First(&user, userID).Error; err != nil {
		return nil, err
	}
	if user.Role == nil {
		return nil, nil
	}
	var perms []string
	if err := json.Unmarshal([]byte(user.Role.Permissions), &perms); err != nil {
		return nil, err
	}
	return perms, nil
}
