package service

import (
	"github.com/lta2705/Go-Payment-Gateway/internal/middleware"
	"github.com/lta2705/Go-Payment-Gateway/internal/repository"
	"go.uber.org/zap"
)

type MerchantCredentialsService interface {
	Authenticate(apiKey string) *string
}

type MerchantCredentialsServiceImpl struct {
	TxCreRepo repository.MerchantCredentialsRepository
}

func (t MerchantCredentialsServiceImpl) Authenticate(apiKey string) *string {
	logger := middleware.CreateLogger()
	defer logger.Sync()

	merchantID, err := t.TxCreRepo.FindMerchantIDByApiKey(apiKey)
	if err != nil {
		logger.Error("Error authenticating merchant", zap.Error(err))
	}
	if merchantID != "" {
		logger.Info("Successfully authenticated merchant", zap.String("MerchantID", merchantID))
		return &merchantID
	}

	logger.Warn("Cannot find the exist merchantID by:", zap.String("apiKey", apiKey))
	return nil
}

func NewMerchantCredentialsService(txCreRepo repository.MerchantCredentialsRepository) MerchantCredentialsService {
	return &MerchantCredentialsServiceImpl{
		TxCreRepo: txCreRepo,
	}
}
