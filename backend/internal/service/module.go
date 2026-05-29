package service

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"

	"github.com/cashier-config/server/internal/model"
)

type ModuleService struct {
	db *gorm.DB
}

func NewModuleService(db *gorm.DB) *ModuleService {
	return &ModuleService{db: db}
}

func (s *ModuleService) List(env, location string) ([]model.Module, error) {
	var list []model.Module
	q := s.db.Order("id ASC")
	if env != "" {
		q = q.Where("env = ?", env)
	}
	if location != "" {
		q = q.Where("location = ?", location)
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *ModuleService) GetByID(id int64) (*model.Module, error) {
	var m model.Module
	if err := s.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *ModuleService) Create(m *model.Module) error {
	return s.db.Create(m).Error
}

func (s *ModuleService) Update(m *model.Module) error {
	return s.db.Save(m).Error
}

func (s *ModuleService) Delete(id int64) error {
	return s.db.Delete(&model.Module{}, id).Error
}

func (s *ModuleService) resolveModuleJSON(modulesField string) (string, error) {
	var ids []int64
	if err := json.Unmarshal([]byte(modulesField), &ids); err != nil {
		return modulesField, nil
	}
	if len(ids) == 0 {
		return "[]", nil
	}
	var groups []model.Group
	if err := s.db.Where("id IN ?", ids).Find(&groups).Error; err != nil {
		return "[]", err
	}
	idMap := make(map[int64]model.Group, len(groups))
	for _, g := range groups {
		idMap[g.ID] = g
	}
	type stepItem struct {
		Name string `json:"name"`
		Path string `json:"path"`
	}
	type groupItem struct {
		Name  string     `json:"name"`
		Steps []stepItem `json:"steps"`
	}
	var result []groupItem
	for _, id := range ids {
		if g, ok := idMap[id]; ok {
			var steps []stepItem
			json.Unmarshal([]byte(g.Steps), &steps)
			result = append(result, groupItem{Name: g.Name, Steps: steps})
		}
	}
	if len(result) == 0 {
		return "[]", nil
	}
	raw, err := json.Marshal(result)
	if err != nil {
		return "[]", err
	}
	return string(raw), nil
}

func (s *ModuleService) ImportCSV(records [][]string) (int64, error) {
	var count int64
	for i, row := range records {
		if i == 0 {
			continue
		}
		if len(row) < 5 {
			continue
		}
		m := model.Module{
			Name:     row[0],
			Location: row[1],
			Modules:  row[2],
			Env:      row[3],
		}
		if err := s.db.Create(&m).Error; err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func (s *ModuleService) ExportCSV(env, location string) ([]byte, string, error) {
	list, err := s.List(env, location)
	if err != nil {
		return nil, "", err
	}
	var buf bytes.Buffer
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(&buf)
	w.Write([]string{"Name", "Location", "Modules", "Env", "ID"})
	for _, m := range list {
		modulesJSON, err := s.resolveModuleJSON(m.Modules)
		if err != nil {
			return nil, "", err
		}
		w.Write([]string{
			m.Name, m.Location, modulesJSON, m.Env, fmt.Sprint(m.ID),
		})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), "Module.csv", nil
}
