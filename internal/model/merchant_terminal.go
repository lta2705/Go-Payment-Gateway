package model

type MerchantTerminal struct {
	PcPosId    string `gorm:"type:varchar(50)"`
	TerminalId string `gorm:"type:varchar(50);primaryKey"`
}
