package handler

import (
	"github.com/lta2705/Go-Payment-Gateway/internal/service"
	"go.uber.org/zap"
)

type MerchantCredentialsHandler struct {
	merchantCredentialsService service.MerchantCredentialsService
}

func NewMerchantCredentialsHandler(merchantCredentialsService service.MerchantCredentialsService) *MerchantCredentialsHandler {
	return &MerchantCredentialsHandler{
		merchantCredentialsService: merchantCredentialsService,
	}
}

func (h *MerchantCredentialsHandler) Authenticate(apiKey string) *string {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	return h.merchantCredentialsService.Authenticate(apiKey)
}
