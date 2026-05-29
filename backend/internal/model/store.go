package model

type Store struct {
	ID             int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	StoreNumber    string `gorm:"column:store_number;size:50;not null;index;unique" json:"store_number"`
	NetworkSegment string `gorm:"column:network_segment;size:50" json:"network_segment"`
	WebposEnv      string `gorm:"column:webpos_env;size:10" json:"webpos_env"`
	EFT            string `gorm:"column:eft;size:100" json:"eft"`
	Location       string `gorm:"column:location;size:10" json:"location"`
	RfServer       string `gorm:"column:rf_server;size:50" json:"rf_server"`
	CashtillSegGW  string `gorm:"column:cashtill_seg_gw;size:50" json:"cashtill_seg_gw"`
}

func (Store) TableName() string {
	return "stores"
}
