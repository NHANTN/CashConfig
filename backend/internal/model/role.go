package model

import "time"

type Role struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"column:name;size:50;not null" json:"name"`
	Code        string    `gorm:"column:code;size:50;not null;uniqueIndex" json:"code"`
	Permissions string    `gorm:"column:permissions;type:text" json:"permissions"`
	Status      int       `gorm:"column:status;default:1" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (Role) TableName() string {
	return "roles"
}
