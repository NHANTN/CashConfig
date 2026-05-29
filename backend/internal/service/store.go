package service

import (
	"bytes"
	"encoding/csv"

	"gorm.io/gorm"

	"github.com/cashier-config/server/internal/model"
)

type StoreService struct {
	db *gorm.DB
}

func NewStoreService(db *gorm.DB) *StoreService {
	return &StoreService{db: db}
}

func (s *StoreService) List(location, eft string) ([]model.Store, error) {
	var list []model.Store
	q := s.db.Order("id ASC")
	if location != "" {
		q = q.Where("location = ?", location)
	}
	if eft != "" {
		q = q.Where("eft = ?", eft)
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *StoreService) GetByID(id int64) (*model.Store, error) {
	var m model.Store
	if err := s.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *StoreService) Create(m *model.Store) error {
	return s.db.Create(m).Error
}

func (s *StoreService) Update(m *model.Store) error {
	return s.db.Save(m).Error
}

func (s *StoreService) Delete(id int64) error {
	return s.db.Delete(&model.Store{}, id).Error
}

func (s *StoreService) ImportCSV(records [][]string) (int64, error) {
	var count int64
	for i, row := range records {
		if i == 0 {
			continue
		}
		if len(row) < 7 {
			continue
		}
		m := model.Store{
			StoreNumber:    row[0],
			NetworkSegment: row[1],
			WebposEnv:      row[2],
			EFT:            row[3],
			Location:       row[4],
			RfServer:       row[5],
			CashtillSegGW:  row[6],
		}
		if err := s.db.Create(&m).Error; err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func (s *StoreService) ExportCSV(location, eft string) ([]byte, string, error) {
	list, err := s.List(location, eft)
	if err != nil {
		return nil, "", err
	}
	var buf bytes.Buffer
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(&buf)
	w.Write([]string{"Store_Number", "Network_Segment", "Webpos_Env", "EFT", "Location", "RF_Server", "Cashtill_Seg_GW"})
	for _, m := range list {
		w.Write([]string{m.StoreNumber, m.NetworkSegment, m.WebposEnv, m.EFT, m.Location, m.RfServer, m.CashtillSegGW})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), "Store.csv", nil
}
