package service

import (
	"github.com/lta2705/Go-Payment-Gateway/internal/repository"
	"go.uber.org/zap"
)

type MerchantCredentialsService interface {
}

type MerchantCredentialsServiceImpl struct {
	logger    *zap.Logger
	TxCreRepo repository.MerchantCredentialsRepository
}

func (t MerchantCredentialsServiceImpl) Authenticate(apiKey string) *string {
	merchantID, err := t.TxCreRepo.FindMerchantIDByApiKey(apiKey)
	if err != nil {
		t.logger.Error("Error authenticating merchant", zap.Error(err))
	}
	if merchantID != "" {
		t.logger.Info("Successfully authenticated merchant", zap.String("MerchantID", merchantID))
		return &merchantID
	}

	t.logger.Warn("Cannot find the exist merchantID by:", zap.String("apiKey", apiKey))
	return nil
}
