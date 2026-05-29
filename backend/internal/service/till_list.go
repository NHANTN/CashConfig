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
	db     *gorm.DB
	writer *SyncReportWriter
}

func NewTillListService(db *gorm.DB, writer *SyncReportWriter) *TillListService {
	return &TillListService{db: db, writer: writer}
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

func (s *TillListService) extractFields(parsed map[string]interface{}) (fields struct {
	hostName      string
	macAddress    string
	location      string
	storeNumber   string
	env           string
	ip            string
	hardwareModel string
	name          string
	groupID       int64
	lastSeen      string
	execStatus    int
	duration      int
	syncTime      string
	moduleExec    string
}) {
	fields.hostName, _ = parsed["Name"].(string)
	fields.macAddress, _ = parsed["MACAddress"].(string)
	if fields.macAddress == "" {
		if facts, ok := parsed["Fact"].([]interface{}); ok {
			for _, f := range facts {
				if fact, ok := f.(map[string]interface{}); ok {
					if key, _ := fact["Key"].(string); key == "MACAddress" {
						if val, _ := fact["Value"].(string); val != "" {
							fields.macAddress = val
							break
						}
					}
				}
			}
		}
	}
	fields.location, _ = parsed["Location"].(string)
	fields.storeNumber, _ = parsed["StoreNumber"].(string)
	fields.env, _ = parsed["Env"].(string)
	fields.ip, _ = parsed["Ip"].(string)
	fields.hardwareModel, _ = parsed["HardwareModel"].(string)
	fields.name, _ = parsed["Name"].(string)
	if g, ok := parsed["Group"].(float64); ok {
		fields.groupID = int64(g)
	}
	fields.syncTime = time.Now().Format("2006-01-02 15:04:05")
	if exec, ok := parsed["Execution"].(map[string]interface{}); ok {
		if s, ok := exec["Status"].(float64); ok {
			fields.execStatus = int(s)
		}
		if d, ok := exec["Duration"].(float64); ok {
			fields.duration = int(d)
		}
		if start, ok := exec["StartTime"].(string); ok {
			fields.syncTime = start
		}
		if me, ok := exec["ModuleExecution"]; ok {
			b, _ := json.Marshal(me)
			fields.moduleExec = string(b)
		}
	}
	return
}

func (s *TillListService) CheckIn(raw string) error {
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return err
	}
	f := s.extractFields(parsed)

	var tillID int64
	if err := s.db.Raw(
		`INSERT INTO till_lists (host_name, mac_address, location, store_number, env, name, ip, hardware_model, group_id, request_body, last_seen)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT (host_name) DO UPDATE SET
		     mac_address = EXCLUDED.mac_address,
		     location = EXCLUDED.location,
		     store_number = EXCLUDED.store_number,
		     env = EXCLUDED.env,
		     name = EXCLUDED.name,
		     ip = EXCLUDED.ip,
		     hardware_model = EXCLUDED.hardware_model,
		     group_id = EXCLUDED.group_id,
		     request_body = EXCLUDED.request_body,
		     last_seen = EXCLUDED.last_seen
		 RETURNING id`,
		f.hostName, f.macAddress, f.location, f.storeNumber, f.env,
		f.name, f.ip, f.hardwareModel, f.groupID, raw, f.syncTime,
	).Scan(&tillID).Error; err != nil {
		return err
	}

	s.writer.Write(model.SyncReport{
		TillListID:      tillID,
		RequestBody:     raw,
		ModuleExecution: f.moduleExec,
		Status:          f.execStatus,
		Duration:        f.duration,
		SyncTime:        f.syncTime,
	})
	return nil
}