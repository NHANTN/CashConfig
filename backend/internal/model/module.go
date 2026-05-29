package model

type Module struct {
	ID       int64  `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name     string `gorm:"column:name;size:255;not null;index" json:"name"`
	Location string `gorm:"column:location;size:10;not null" json:"location"`
	Modules  string `gorm:"column:modules;type:text;not null" json:"modules"`
	Env      string `gorm:"column:env;size:10;not null;index" json:"env"`
}

func (Module) TableName() string {
	return "modules"
}
