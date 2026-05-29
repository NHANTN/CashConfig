package model

import "time"

type OperationLog struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int64     `gorm:"column:user_id;index" json:"user_id"`
	Username   string    `gorm:"column:username;size:100" json:"username"`
	Action     string    `gorm:"column:action;size:50" json:"action"`
	TargetType string    `gorm:"column:target_type;size:50" json:"target_type"`
	TargetID   string    `gorm:"column:target_id;size:50" json:"target_id"`
	Detail     string    `gorm:"column:detail;type:text" json:"detail"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (OperationLog) TableName() string {
	return "operation_logs"
}
