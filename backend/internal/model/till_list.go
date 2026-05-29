package model

type TillList struct {
	ID          int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	HostName    string `gorm:"column:host_name;size:100;not null;index;unique" json:"host_name"`
	MacAddress  string `gorm:"column:mac_address;size:20;not null" json:"mac_address"`
	Location    string `gorm:"column:location;size:10;index" json:"location"`
	StoreNumber string `gorm:"column:store_number;size:20" json:"store_number"`
	Env         string `gorm:"column:env;size:10;index" json:"env"`
	Name        string `gorm:"column:name;size:100" json:"name"`
	Ip          string `gorm:"column:ip;size:50" json:"ip"`
	HardwareModel string `gorm:"column:hardware_model;size:50" json:"hardware_model"`
	GroupID     int64  `gorm:"column:group_id" json:"group_id"`
	RequestBody string `gorm:"column:request_body;type:text" json:"request_body"`
	LastSeen    string `gorm:"column:last_seen;size:30" json:"last_seen"`
}

func (TillList) TableName() string {
	return "till_lists"
}