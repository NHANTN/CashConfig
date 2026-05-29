package service

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/cashier-config/server/internal/model"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) List() ([]model.User, error) {
	var list []model.User
	if err := s.db.Preload("Role").Order("id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *UserService) GetByID(id int64) (*model.User, error) {
	var user model.User
	if err := s.db.Preload("Role").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) Create(m *model.User, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	m.PasswordHash = string(hash)
	return s.db.Create(m).Error
}

func (s *UserService) Update(m *model.User) error {
	data := map[string]interface{}{
		"username": m.Username,
		"name":     m.Name,
		"role_id":  m.RoleID,
		"status":   m.Status,
	}
	return s.db.Model(&model.User{}).Where("id = ?", m.ID).Updates(data).Error
}

func (s *UserService) UpdatePassword(id int64, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.db.Model(&model.User{}).Where("id = ?", id).Update("password_hash", string(hash)).Error
}

func (s *UserService) Delete(id int64) error {
	return s.db.Delete(&model.User{}, id).Error
}
