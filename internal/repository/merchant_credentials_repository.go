package repository

import (
	"github.com/lta2705/Go-Payment-Gateway/internal/model"
	"gorm.io/gorm"
)

type MerchantCredentialsRepository interface {
	FindMerchantIDByApiKey(apiKey string) (string, error)
}
type MerchantCredentialsRepositoryImpl struct {
	db *gorm.DB
}

func (m *MerchantCredentialsRepositoryImpl) FindMerchantIDByApiKey(apiKey string) (string, error) {
	var credentials model.MerchantCredential

	err := m.db.Select("merchant_id").Where("api_key = ?", apiKey).First(&credentials).Error

	if err != nil {
		return "", err
	}

	return credentials.MerchantId, nil
}

func NewMerchantCredentialsRepository(db *gorm.DB) MerchantCredentialsRepository {
	return &MerchantCredentialsRepositoryImpl{db: db}
}
