package model

import "time"

type SSOState struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	State     string    `gorm:"column:state;size:100;not null;uniqueIndex" json:"state"`
	Nonce     string    `gorm:"column:nonce;size:100" json:"-"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (SSOState) TableName() string {
	return "sso_states"
}
