package handler

import (
	"github.com/lta2705/Go-Payment-Gateway/internal/service"
	"go.uber.org/zap"
)

type MerchantCredentialsHandler struct {
	logger                     *zap.SugaredLogger
	merchantCredentialsService service.MerchantCredentialsService
}
