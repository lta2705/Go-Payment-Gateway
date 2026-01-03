package dto

import (
	"github.com/google/uuid"
	"time"
)

type TerminalTransactionDTO struct {
	ID              uuid.UUID `json:"ID"`
	TransactionType string    `json:"TransactionType"`
	CurrCd          string    `json:"CurrCd"`
	TotTrAmt        float64   `json:"TotTrAmt"`
	TipAmt          float64   `json:"TipAmt"`
	PcPosId         string    `json:"PcPosId"`
	TransactionId   string    `json:"TransactionId"`
	OrgPcPosTxnId   string    `json:"OrgPcPosTxnId"`
	AprvNo          string    `json:"AprvNo"`
	MsgType         string    `json:"MsgType"`
	Status          string    `json:"Status"`
	ErrorCode       string    `json:"ErrorCode"`
	ErrorDetail     string    `json:"ErrorDetail"`
	CreatedAt       time.Time `json:"CreatedAt"`
	UpdatedBy       string    `json:"UpdatedBy"`
	TerminalId      string    `json:"terminalId"`
}
