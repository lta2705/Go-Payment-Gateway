package handler

import (
	"github.com/lta2705/Go-Payment-Gateway/internal/service"
	"go.uber.org/zap"
)

var logger *zap.Logger

type MerchantCredentialsHandler struct {
	merchantCredentialsService service.MerchantCredentialsService
}

func NewMerchantCredentialsHandler(merchantCredentialsService service.MerchantCredentialsService) *MerchantCredentialsHandler {
	return &MerchantCredentialsHandler{
		merchantCredentialsService: merchantCredentialsService,
	}
}

func (h *MerchantCredentialsHandler) Authenticate(apiKey string) *string {
	return h.merchantCredentialsService.Authenticate(apiKey)
}
