package service

import (
	"errors"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/lta2705/Go-Payment-Gateway/internal/repository"
)

type MerchantCredentialsService interface {
	Authenticate(apiKey string) (string, error)
}

type MerchantCredentialsServiceImpl struct {
	TxCreRepo repository.MerchantCredentialsRepository
}

func (t MerchantCredentialsServiceImpl) Authenticate(apiKey string) (string, error) {
	merchantID, err := t.TxCreRepo.FindMerchantIDByApiKey(apiKey)

	if err != nil {
		logger.Error("Database error during merchant authentication")
		return "", err
	}

	if merchantID == "" {
		logger.Warn("Authentication failed: API Key not found", "apiKey", apiKey)
		return "", errors.New("invalid api key")
	}

	logger.Info("Successfully authenticated merchant", "MerchantID", merchantID)
	return merchantID, nil
}

func NewMerchantCredentialsService(txCreRepo repository.MerchantCredentialsRepository) MerchantCredentialsService {
	return &MerchantCredentialsServiceImpl{
		TxCreRepo: txCreRepo,
	}
}
