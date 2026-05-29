package model

import "time"

type SyncReport struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TillListID      int64     `gorm:"column:till_list_id;not null;index" json:"till_list_id"`
	RequestBody     string    `gorm:"column:request_body;type:text" json:"request_body"`
	ModuleExecution string    `gorm:"column:module_execution;type:text" json:"module_execution"`
	Status          int       `gorm:"column:status;default:0" json:"status"`
	Duration        int       `gorm:"column:duration;default:0" json:"duration"`
	SyncTime        string    `gorm:"column:sync_time;size:30" json:"sync_time"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (SyncReport) TableName() string { return "sync_reports" }
