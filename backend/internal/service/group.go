package service

import (
	"gorm.io/gorm"

	"github.com/cashier-config/server/internal/model"
)

type GroupService struct {
	db *gorm.DB
}

func NewGroupService(db *gorm.DB) *GroupService {
	return &GroupService{db: db}
}

func (s *GroupService) List(name string) ([]model.Group, error) {
	var list []model.Group
	q := s.db.Order("id ASC")
	if name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *GroupService) ListAll() ([]model.Group, error) {
	return s.List("")
}

func (s *GroupService) GetByID(id int64) (*model.Group, error) {
	var m model.Group
	if err := s.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *GroupService) Create(m *model.Group) error {
	return s.db.Create(m).Error
}

func (s *GroupService) Update(m *model.Group) error {
	return s.db.Save(m).Error
}

func (s *GroupService) Delete(id int64) error {
	return s.db.Delete(&model.Group{}, id).Error
}
