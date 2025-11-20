package dto

type MerchantCredentialsDTO struct {
	Id         string `json:"id" binding:"required"`
	MerchantId string `json:"merchant_id" binding:"required"`
	ApiKey     string `headers:"X-API-KEY" binding:"required"`
}
