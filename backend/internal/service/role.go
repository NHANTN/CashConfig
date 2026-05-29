package service

import (
	"gorm.io/gorm"

	"github.com/cashier-config/server/internal/model"
)

type RoleService struct {
	db *gorm.DB
}

func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{db: db}
}

func (s *RoleService) List() ([]model.Role, error) {
	var list []model.Role
	if err := s.db.Order("id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *RoleService) GetByID(id int64) (*model.Role, error) {
	var role model.Role
	if err := s.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *RoleService) Create(m *model.Role) error {
	return s.db.Create(m).Error
}

func (s *RoleService) Update(m *model.Role) error {
	return s.db.Save(m).Error
}

func (s *RoleService) Delete(id int64) error {
	return s.db.Delete(&model.Role{}, id).Error
}
