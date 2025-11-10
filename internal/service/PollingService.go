package service

import (
	"github.com/joho/godotenv"
	_ "github.com/lta2705/Go-Payment-Gateway/internal/dto"
	"github.com/lta2705/Go-Payment-Gateway/internal/model"
	"github.com/lta2705/Go-Payment-Gateway/internal/repository"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
)

type PollingService interface {
	Poll(model *model.Transaction) model.Transaction
}

type PollingServiceImpl struct {
	logger *zap.Logger
	txRepo repository.TransactionRepository
}

func (p PollingServiceImpl) getTimeout() int {
	err := godotenv.Load()
	if err != nil {
		p.logger.Error("Error loading .env file", zap.Error(err))
	}
	timeoutStr := os.Getenv("POLLING_MAX_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		p.logger.Error("Error converting POLLING_MAX_TIMEOUT to int", zap.Error(err))
		return 60 // default timeout
	}
	return timeout
}

func (p PollingServiceImpl) Poll(model *model.Transaction) model.Transaction {
	startTime := time.Now()
	timeout := time.Duration(p.getTimeout()) * time.Millisecond
	transactionId := model.TransactionId
	for time.Since(startTime) < timeout {
		p.logger.Info("Polling for transaction status...")
		// Here you would add the logic to check transaction statuses
		pendingTransaction, err := p.txRepo.FindByTransactionId(&transactionId)
		if err != nil {
			p.logger.Error("Error fetching transaction during polling", zap.Error(err))
		}
		if pendingTransaction.Status != "pending" {
			p.logger.Info("Transaction status updated", zap.String("TransactionId", transactionId), zap.String("Status", pendingTransaction.Status))
			return *pendingTransaction
		}
		time.Sleep(2 * time.Second) // Poll every 2 seconds
	}

	return *model
}

func NewPollingService(logger *zap.Logger, txRepo repository.TransactionRepository) PollingService {
	return &PollingServiceImpl{
		logger: logger,
		txRepo: txRepo,
	}
}
