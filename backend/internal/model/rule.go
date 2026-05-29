package model

type Rule struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"column:name;size:255;not null;index" json:"name"`
	Type     string `gorm:"column:type;size:10;not null" json:"type"`
	Location string `gorm:"column:location;size:10;not null" json:"location"`
	EnvName  string `gorm:"column:env_name;size:100;not null" json:"env_name"`
	Rule     string `gorm:"column:rule;type:text;not null" json:"rule"`
	Result   string `gorm:"column:result;size:255;not null" json:"result"`
	Sort     int    `gorm:"column:sort;default:0" json:"sort"`
}

func (Rule) TableName() string {
	return "rules"
}
