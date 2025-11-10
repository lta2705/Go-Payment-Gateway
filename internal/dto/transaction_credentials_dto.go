package dto

type TransactionCredentialsDTO struct {
	Id         string `json:"id" binding:"required"`
	MerchantId string `json:"merchant_id" binding:"required"`
	ApiKey     string
}
