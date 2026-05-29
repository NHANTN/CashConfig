package service

import (
	"bytes"
	"encoding/csv"

	"gorm.io/gorm"

	"github.com/cashier-config/server/internal/model"
)

type RuleService struct {
	db *gorm.DB
}

func NewRuleService(db *gorm.DB) *RuleService {
	return &RuleService{db: db}
}

func (s *RuleService) List(typ, location, envName string) ([]model.Rule, error) {
	var list []model.Rule
	q := s.db.Order("sort ASC, id ASC")
	if typ != "" {
		q = q.Where("type = ?", typ)
	}
	if location != "" {
		q = q.Where("location = ?", location)
	}
	if envName != "" {
		q = q.Where("env_name = ?", envName)
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *RuleService) GetByID(id int64) (*model.Rule, error) {
	var m model.Rule
	if err := s.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *RuleService) Create(m *model.Rule) error {
	return s.db.Create(m).Error
}

func (s *RuleService) Update(m *model.Rule) error {
	return s.db.Save(m).Error
}

func (s *RuleService) Delete(id int64) error {
	return s.db.Delete(&model.Rule{}, id).Error
}

func (s *RuleService) ImportCSV(records [][]string) (int64, error) {
	var count int64
	for i, row := range records {
		if i == 0 {
			continue
		}
		if len(row) < 6 {
			continue
		}
		m := model.Rule{
			Name:   row[0],
			Type:   row[1],
			Location: row[2],
			EnvName: row[3],
			Rule:   row[4],
			Result: row[5],
		}
		if err := s.db.Create(&m).Error; err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func (s *RuleService) TestRule(hostName string) ([]model.Rule, error) {
	var all []model.Rule
	if err := s.db.Order("sort ASC, id ASC").Find(&all).Error; err != nil {
		return nil, err
	}
	var matched []model.Rule
	for _, r := range all {
		if matchRule(hostName, r.Rule) {
			matched = append(matched, r)
		}
	}
	return matched, nil
}

func (s *RuleService) UpdateSort(id int64, sort int) error {
	return s.db.Model(&model.Rule{}).Where("id = ?", id).Update("sort", sort).Error
}

func matchRule(hostName, pattern string) bool {
	return pattern == "" || pattern == "*" || pattern == hostName || matchGlob(hostName, pattern)
}

func matchGlob(val, pattern string) bool {
	var px, vx int
	n := len(pattern)
	m := len(val)
	var starPx, starVx int = -1, -1
	for vx < m {
		if px < n && (pattern[px] == val[vx] || pattern[px] == '?') {
			px++
			vx++
			continue
		}
		if px < n && pattern[px] == '*' {
			starPx = px
			starVx = vx
			px++
			continue
		}
		if starPx >= 0 {
			px = starPx + 1
			vx = starVx + 1
			starVx++
			continue
		}
		return false
	}
	for px < n && pattern[px] == '*' {
		px++
	}
	return px == n
}

func (s *RuleService) ExportCSV(typ, location, envName string) ([]byte, string, error) {
	list, err := s.List(typ, location, envName)
	if err != nil {
		return nil, "", err
	}
	var buf bytes.Buffer
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(&buf)
	w.Write([]string{"Name", "Type", "Location", "ENV_NAME", "Rule", "Result"})
	for _, m := range list {
		w.Write([]string{m.Name, m.Type, m.Location, m.EnvName, m.Rule, m.Result})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), "Rule.csv", nil
}
