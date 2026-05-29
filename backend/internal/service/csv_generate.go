package service

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/cashier-config/server/internal/model"
)

type CsvGenerateService struct {
	db        *gorm.DB
	moduleSvc *ModuleService
	ruleSvc   *RuleService
	storeSvc  *StoreService
	tillSvc   *TillListService
	varSvc    *VarService
	outputDir string
}

func NewCsvGenerateService(
	db *gorm.DB,
	moduleSvc *ModuleService,
	ruleSvc *RuleService,
	storeSvc *StoreService,
	tillSvc *TillListService,
	varSvc *VarService,
	outputDir string,
) *CsvGenerateService {
	return &CsvGenerateService{
		db:        db,
		moduleSvc: moduleSvc,
		ruleSvc:   ruleSvc,
		storeSvc:  storeSvc,
		tillSvc:   tillSvc,
		varSvc:    varSvc,
		outputDir: outputDir,
	}
}

type CsvFile struct {
	Type     string
	Filename string
	Data     []byte
}

func (s *CsvGenerateService) generateAll() ([]CsvFile, error) {
	var files []CsvFile

	if data, _, err := s.moduleSvc.ExportCSV("", ""); err == nil {
		files = append(files, CsvFile{Type: "module", Filename: "Module.csv", Data: data})
	}
	if data, _, err := s.ruleSvc.ExportCSV("", "", ""); err == nil {
		files = append(files, CsvFile{Type: "rule", Filename: "Rule.csv", Data: data})
	}
	if data, _, err := s.storeSvc.ExportCSV("", ""); err == nil {
		files = append(files, CsvFile{Type: "store", Filename: "Store.csv", Data: data})
	}
	if data, _, err := s.tillSvc.ExportCSV("", "", ""); err == nil {
		files = append(files, CsvFile{Type: "till", Filename: "TillList.csv", Data: data})
	}
	if data, _, err := s.varSvc.ExportCSV("", ""); err == nil {
		files = append(files, CsvFile{Type: "var", Filename: "Var.csv", Data: data})
	}
	return files, nil
}

func (s *CsvGenerateService) Generate(t string) ([]CsvFile, error) {
	if t == "" {
		return s.generateAll()
	}
	files, err := s.generateAll()
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.Type == t {
			return []CsvFile{f}, nil
		}
	}
	return nil, fmt.Errorf("unknown type: %s", t)
}

func (s *CsvGenerateService) SaveGeneration(files []CsvFile, operator string) (string, error) {
	timestamp := time.Now().Format("20060102150405")
	dir := filepath.Join(s.outputDir, timestamp)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(dir, f.Filename), f.Data, 0644); err != nil {
			return "", err
		}
	}

	types := make([]string, len(files))
	for i, f := range files {
		types[i] = f.Type
	}

	log := model.CsvGenerationLog{
		FileType:  strings.Join(types, ","),
		FileCount: len(files),
		Operator:  operator,
		Status:    "success",
		Detail:    fmt.Sprintf("generated %d files in %s", len(files), timestamp),
	}
	if err := s.db.Create(&log).Error; err != nil {
		return "", err
	}
	return timestamp, nil
}

func (s *CsvGenerateService) GetLatestDir() (string, error) {
	entries, err := os.ReadDir(s.outputDir)
	if err != nil {
		return "", err
	}
	if len(entries) == 0 {
		return "", fmt.Errorf("no generations found")
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() > entries[j].Name()
	})
	return filepath.Join(s.outputDir, entries[0].Name()), nil
}

func (s *CsvGenerateService) GetFileData(fileType string) ([]byte, string, error) {
	dir, err := s.GetLatestDir()
	if err != nil {
		return nil, "", err
	}
	filename := csvFilename(fileType)
	data, err := os.ReadFile(filepath.Join(dir, filename))
	if err != nil {
		return nil, "", err
	}
	return data, filename, nil
}

func (s *CsvGenerateService) GetAllAsZip() ([]byte, string, error) {
	dir, err := s.GetLatestDir()
	if err != nil {
		return nil, "", err
	}
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, "", err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		f, err := w.Create(e.Name())
		if err != nil {
			continue
		}
		f.Write(data)
	}
	w.Close()
	return buf.Bytes(), "csv-files.zip", nil
}

func (s *CsvGenerateService) GetHistory(fileType string) ([]model.CsvGenerationLog, error) {
	var logs []model.CsvGenerationLog
	q := s.db.Order("generated_at DESC").Limit(50)
	if fileType != "" {
		q = q.Where("file_type LIKE ?", "%"+fileType+"%")
	}
	if err := q.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (s *CsvGenerateService) Diff(fileType, from, to string) (map[string]interface{}, error) {
	getFile := func(timestamp string) ([]byte, error) {
		dir := filepath.Join(s.outputDir, timestamp)
		return os.ReadFile(filepath.Join(dir, csvFilename(fileType)))
	}
	fromData, err := getFile(from)
	if err != nil {
		return nil, fmt.Errorf("from version not found: %s", from)
	}
	toData, err := getFile(to)
	if err != nil {
		return nil, fmt.Errorf("to version not found: %s", to)
	}
	return map[string]interface{}{
		"from":      string(fromData),
		"to":        string(toData),
		"from_time": from,
		"to_time":   to,
	}, nil
}

func csvFilename(fileType string) string {
	switch fileType {
	case "module":
		return "Module.csv"
	case "rule":
		return "Rule.csv"
	case "store":
		return "Store.csv"
	case "till":
		return "TillList.csv"
	case "var":
		return "Var.csv"
	default:
		return fileType + ".csv"
	}
}
