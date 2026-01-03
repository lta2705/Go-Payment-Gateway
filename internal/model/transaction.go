package model

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
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

func (t *Transaction) ToBeautifiedString() string {
	return fmt.Sprintf(
		"--- Transaction Report ---\n"+
			"ID            : %s\n"+
			"Type          : %s\n"+
			"Status        : [%s]\n"+
			"Amount        : %s %.2f (Tip: %.2f)\n"+
			"POS ID        : %s\n"+
			"Txn ID        : %s\n"+
			"Approval No   : %s\n"+
			"Error         : %s (%s)\n"+
			"Created At    : %s\n"+
			"--------------------------",
		t.ID,
		t.TransactionType,
		t.Status,
		t.CurrCd, t.TotTrAmt, t.TipAmt,
		t.PcPosId,
		t.TransactionId,
		t.AprvNo,
		t.ErrorCode, t.ErrorDetail,
		t.CreatedAt.Format("2006-01-02 15:04:05"),
	)
}

func (t *Transaction) ToBeautifiedJson() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return "Error formatting transaction"
	}
	return string(b)
}
