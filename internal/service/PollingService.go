package service

import (
	"github.com/lta2705/Go-Payment-Gateway/internal/repository"
	"go.uber.org/zap"
)

type PollingService interface {
	Poll()
}

type PollingServiceImpl struct {
	logger *zap.Logger
	txRepo repository.TransactionRepository
}

func (p PollingServiceImpl) Poll() {

}

func NewPollingService(logger *zap.Logger, txRepo repository.TransactionRepository) PollingService {
	return &PollingServiceImpl{
		logger: logger,
		txRepo: txRepo,
	}
}
