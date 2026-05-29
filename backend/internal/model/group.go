package model

type Group struct {
	ID    int64  `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name  string `gorm:"column:name;size:255;not null;index;unique" json:"name"`
	Steps string `gorm:"column:steps;type:text;not null" json:"steps"` // JSON array of {name, path}
}

func (Group) TableName() string {
	return "groups"
}
