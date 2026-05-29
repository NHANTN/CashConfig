package model

import "time"

type User struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"column:username;size:50;not null;uniqueIndex" json:"username"`
	PasswordHash string    `gorm:"column:password_hash;size:255;not null" json:"-"`
	Name         string    `gorm:"column:name;size:100;not null" json:"name"`
	Email        string    `gorm:"column:email;size:200" json:"email"`
	RoleID       int64     `gorm:"column:role_id;not null;index" json:"role_id"`
	AuthSource   string    `gorm:"column:auth_source;size:20;default:'local'" json:"auth_source"`
	LDAPDN       string    `gorm:"column:ldap_dn;size:500" json:"-"`
	Status       int       `gorm:"column:status;default:1" json:"status"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	Role         *Role     `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

func (User) TableName() string {
	return "users"
}
