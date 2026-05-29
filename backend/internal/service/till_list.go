package service

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/cashier-config/server/internal/model"
)

type TillListService struct {
	db *gorm.DB
}

func NewTillListService(db *gorm.DB) *TillListService {
	return &TillListService{db: db}
}

func (s *TillListService) List(hostName, location, env string) ([]model.TillList, error) {
	var list []model.TillList
	q := s.db.Order("last_seen DESC, id DESC")
	if hostName != "" {
		q = q.Where("host_name LIKE ?", "%"+hostName+"%")
	}
	if location != "" {
		q = q.Where("location = ?", location)
	}
	if env != "" {
		q = q.Where("env = ?", env)
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *TillListService) GetByID(id int64) (*model.TillList, error) {
	var m model.TillList
	if err := s.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *TillListService) Create(m *model.TillList) error {
	return s.db.Create(m).Error
}

func (s *TillListService) Update(m *model.TillList) error {
	return s.db.Save(m).Error
}

func (s *TillListService) Delete(id int64) error {
	return s.db.Delete(&model.TillList{}, id).Error
}

func (s *TillListService) ImportCSV(records [][]string) (int64, error) {
	var count int64
	for i, row := range records {
		if i == 0 {
			continue
		}
		if len(row) < 2 {
			continue
		}
		m := model.TillList{
			HostName:   row[0],
			MacAddress: row[1],
		}
		if err := s.db.Create(&m).Error; err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func (s *TillListService) ExportCSV(hostName, location, env string) ([]byte, string, error) {
	list, err := s.List(hostName, location, env)
	if err != nil {
		return nil, "", err
	}
	var buf bytes.Buffer
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(&buf)
	w.Write([]string{"ID", "HostName", "MAC", "Location", "Store", "Env", "IP", "Hardware", "GroupID", "LastSeen"})
	for _, m := range list {
		w.Write([]string{
			fmt.Sprintf("%d", m.ID), m.HostName, m.MacAddress, m.Location,
			m.StoreNumber, m.Env, m.Ip, m.HardwareModel,
			fmt.Sprintf("%d", m.GroupID), m.LastSeen,
		})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), "TillList.csv", nil
}

func (s *TillListService) ListReports(tillListID int64) ([]model.SyncReport, error) {
	var reports []model.SyncReport
	if err := s.db.Where("till_list_id = ?", tillListID).Order("id DESC").Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

func (s *TillListService) GetReport(reportID int64) (*model.SyncReport, error) {
	var r model.SyncReport
	if err := s.db.First(&r, reportID).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

type DeviceWithReports struct {
	TillList model.TillList   `json:"till_list"`
	Reports  []model.SyncReport `json:"reports"`
}

func (s *TillListService) QueryReportsByDevice(hostName, macAddress string) ([]DeviceWithReports, error) {
	q := s.db.Model(&model.TillList{})
	if hostName != "" {
		q = q.Where("host_name LIKE ?", "%"+hostName+"%")
	}
	if macAddress != "" {
		q = q.Where("mac_address = ?", macAddress)
	}
	if hostName == "" && macAddress == "" {
		return nil, fmt.Errorf("host_name or mac_address is required")
	}

	var devices []model.TillList
	if err := q.Find(&devices).Error; err != nil {
		return nil, err
	}

	result := make([]DeviceWithReports, 0, len(devices))
	for _, d := range devices {
		var reports []model.SyncReport
		if err := s.db.Where("till_list_id = ?", d.ID).Order("id DESC").Find(&reports).Error; err != nil {
			return nil, err
		}
		result = append(result, DeviceWithReports{TillList: d, Reports: reports})
	}
	return result, nil
}

func (s *TillListService) CheckIn(raw string) error {
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return err
	}

	hostName, _ := parsed["Name"].(string)
	macAddress, _ := parsed["MACAddress"].(string)
	// Fallback: extract MACAddress from Fact array if not present at top level
	if macAddress == "" {
		if facts, ok := parsed["Fact"].([]interface{}); ok {
			for _, f := range facts {
				if fact, ok := f.(map[string]interface{}); ok {
					if key, _ := fact["Key"].(string); key == "MACAddress" {
						if val, _ := fact["Value"].(string); val != "" {
							macAddress = val
							break
						}
					}
				}
			}
		}
	}
	location, _ := parsed["Location"].(string)
	storeNumber, _ := parsed["StoreNumber"].(string)
	env, _ := parsed["Env"].(string)
	ip, _ := parsed["Ip"].(string)
	hardwareModel, _ := parsed["HardwareModel"].(string)
	name, _ := parsed["Name"].(string)
	var groupID int64
	if g, ok := parsed["Group"].(float64); ok {
		groupID = int64(g)
	}

	execStatus := 0
	duration := 0
	syncTime := time.Now().Format("2006-01-02 15:04:05")
	if exec, ok := parsed["Execution"].(map[string]interface{}); ok {
		if s, ok := exec["Status"].(float64); ok {
			execStatus = int(s)
		}
		if d, ok := exec["Duration"].(float64); ok {
			duration = int(d)
		}
		if start, ok := exec["StartTime"].(string); ok {
			syncTime = start
		}
	}

	var moduleExecBytes []byte
	if exec, ok := parsed["Execution"].(map[string]interface{}); ok {
		if me, ok := exec["ModuleExecution"]; ok {
			moduleExecBytes, _ = json.Marshal(me)
		}
	}

	var existing model.TillList
	err := s.db.Where("host_name = ?", hostName).First(&existing).Error

	tillListID := int64(0)
	if err == gorm.ErrRecordNotFound {
		m := model.TillList{
			HostName:      hostName,
			MacAddress:    macAddress,
			Location:      location,
			StoreNumber:   storeNumber,
			Env:           env,
			Name:          name,
			Ip:            ip,
			HardwareModel: hardwareModel,
			GroupID:       groupID,
			RequestBody:   raw,
			LastSeen:      syncTime,
		}
		if err := s.db.Create(&m).Error; err != nil {
			return err
		}
		tillListID = m.ID
	} else if err != nil {
		return err
	} else {
		tillListID = existing.ID
		existing.MacAddress = macAddress
		existing.Location = location
		existing.StoreNumber = storeNumber
		existing.Env = env
		existing.Name = name
		existing.Ip = ip
		existing.HardwareModel = hardwareModel
		existing.GroupID = groupID
		existing.RequestBody = raw
		existing.LastSeen = syncTime
		if err := s.db.Save(&existing).Error; err != nil {
			return err
		}
	}

	report := model.SyncReport{
		TillListID:      tillListID,
		RequestBody:     raw,
		ModuleExecution: string(moduleExecBytes),
		Status:          execStatus,
		Duration:        duration,
		SyncTime:        syncTime,
	}
	return s.db.Create(&report).Error
}