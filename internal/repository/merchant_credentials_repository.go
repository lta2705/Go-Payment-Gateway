package repository

import (
	"github.com/lta2705/Go-Payment-Gateway/internal/model"
	"gorm.io/gorm"
)

type MerchantCredentialsRepository interface {
	FindMerchantIDByApiKey(merchantID string) (string, error)
}
type MerchantCredentialsRepositoryImpl struct {
	db *gorm.DB
}

func (mcri *MerchantCredentialsRepositoryImpl) FindMerchantIDByApiKey(merchantID string) (string, error) {
	var credentials model.MerchantCredentials
	err := mcri.db.Where("merchant_id = ?", merchantID).First(&credentials).Error
	if err != nil {
		return "", err
	}

	return credentials.MerchantId, nil
}

func NewMerchantCredentialsRepository(db *gorm.DB) MerchantCredentialsRepository {
	return &MerchantCredentialsRepositoryImpl{db: db}
}
