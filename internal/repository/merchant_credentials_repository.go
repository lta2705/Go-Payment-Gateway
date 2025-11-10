package repository

import (
	"github.com/lta2705/Go-Payment-Gateway/internal/model"
	"gorm.io/gorm"
)

type MerchantCredentialsRepository interface {
	FindApiKeyByMerchantID(merchantID string) (string, error)
}
type MerchantCredentialsRepositoryImpl struct {
	db *gorm.DB // Uncomment and use this if you need database access
}

func (mcri *MerchantCredentialsRepositoryImpl) FindApiKeyByMerchantID(merchantID string) (string, error) {
	var credentials model.TransactionCredentials
	err := mcri.db.Where("merchant_id = ?", merchantID).First(&credentials).Error
	if err != nil {
		return "", err
	}

	return credentials.ApiKey, nil
}

func NewMerchantCredentialsRepository(db *gorm.DB) MerchantCredentialsRepository {
	return &MerchantCredentialsRepositoryImpl{db: db}
}
