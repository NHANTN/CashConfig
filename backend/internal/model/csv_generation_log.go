package model

import "time"

type CsvGenerationLog struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	GeneratedAt time.Time `gorm:"column:generated_at;autoCreateTime" json:"generated_at"`
	FileType    string    `gorm:"column:file_type;size:50" json:"file_type"`
	FileCount   int       `gorm:"column:file_count" json:"file_count"`
	Operator    string    `gorm:"column:operator;size:100" json:"operator"`
	Status      string    `gorm:"column:status;size:20" json:"status"`
	Detail      string    `gorm:"column:detail;type:text" json:"detail"`
}

func (CsvGenerationLog) TableName() string {
	return "csv_generation_logs"
}
