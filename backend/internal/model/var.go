package model

type Var struct {
	ID      int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	VarName string `gorm:"column:var_name;size:255;not null;index" json:"var_name"`
	Value   string `gorm:"column:value;type:text;not null" json:"value"`
	Env     string `gorm:"column:env;size:10;not null" json:"env"`
	Matcher string `gorm:"column:matcher;type:text;default:'[]'" json:"matcher"`
}

func (Var) TableName() string {
	return "vars"
}
