package dto

type TransactionDTO struct {
	TransactionType string  `json:"transaction_type" binding:"required"`
	CurrCd          string  `json:"curr_cd" binding:"required"`
	TotTrAmt        float64 `json:"tot_tr_amt" binding:"required"`
	TipAmt          float64 `json:"tip_amt"`
	PcPosId         string  `json:"pc_pos_id" binding:"required"`
	TransactionId   string  `json:"transaction_id" binding:"required"`
	OrgPcPosTxnId   string  `json:"org_pc_pos_txn_id"`
	AprvNo          string  `json:"aprv_no"`
	MsgType         string  `json:"msg_type" binding:"required"`
	Status          string  `json:"status"`
	ErrorCode       string  `json:"error_code"`
	ErrorDetail     string  `json:"error_detail"`
}
