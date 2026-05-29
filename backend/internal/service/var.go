package service

import (
	"bytes"
	"encoding/csv"

	"gorm.io/gorm"

	"github.com/cashier-config/server/internal/model"
)

type VarService struct {
	db *gorm.DB
}

func NewVarService(db *gorm.DB) *VarService {
	return &VarService{db: db}
}

func (s *VarService) List(env, varName string) ([]model.Var, error) {
	var list []model.Var
	q := s.db.Order("id ASC")
	if env != "" {
		q = q.Where("env = ?", env)
	}
	if varName != "" {
		q = q.Where("var_name LIKE ?", "%"+varName+"%")
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *VarService) GetByID(id int64) (*model.Var, error) {
	var m model.Var
	if err := s.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *VarService) Create(m *model.Var) error {
	return s.db.Create(m).Error
}

func (s *VarService) Update(m *model.Var) error {
	return s.db.Save(m).Error
}

func (s *VarService) Delete(id int64) error {
	return s.db.Delete(&model.Var{}, id).Error
}

func (s *VarService) ImportCSV(records [][]string) (int64, error) {
	var count int64
	for i, row := range records {
		if i == 0 {
			continue
		}
		if len(row) < 4 {
			continue
		}
		m := model.Var{
			VarName: row[0],
			Value:   row[1],
			Env:     row[2],
			Matcher: row[3],
		}
		if err := s.db.Create(&m).Error; err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func (s *VarService) ExportCSV(env, varName string) ([]byte, string, error) {
	list, err := s.List(env, varName)
	if err != nil {
		return nil, "", err
	}
	var buf bytes.Buffer
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(&buf)
	w.Write([]string{"Var_Name", "Value", "Env", "Matcher"})
	for _, m := range list {
		w.Write([]string{m.VarName, m.Value, m.Env, m.Matcher})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), "Var.csv", nil
}
