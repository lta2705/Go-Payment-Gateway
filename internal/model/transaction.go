package model

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID              uuid.UUID `gorm:"type:uuid;uuid_generate_v4();primaryKey"`
	TransactionType string    `gorm:"type:varchar(50);not null"`
	CurrCd          string    `gorm:"type:varchar(10)"`
	TotTrAmt        float64   `gorm:"type:numeric(15,2)"`
	TipAmt          float64   `gorm:"type:numeric(15,2)"`
	PcPosId         string    `gorm:"type:varchar(50);not null"`
	TransactionId   string    `gorm:"type:varchar(50);not null"`
	OrgPcPosTxnId   string    `gorm:"type:varchar(50)"`
	AprvNo          string    `gorm:"type:varchar(50)"`
	MsgType         string    `gorm:"type:varchar(50)"`
	Status          string    `gorm:"type:varchar(50)"`
	ErrorCode       string    `gorm:"type:varchar(50)"`
	ErrorDetail     string    `gorm:"type:varchar(50)"`
	CreatedAt       time.Time `gorm:"type:timestamp"`
	UpdatedBy       string    `gorm:"type:varchar(50)"`
}
