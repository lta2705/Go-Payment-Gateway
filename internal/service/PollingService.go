package service

import (
	"github.com/joho/godotenv"
	_ "github.com/lta2705/Go-Payment-Gateway/internal/dto"
	"github.com/lta2705/Go-Payment-Gateway/internal/model"
	"github.com/lta2705/Go-Payment-Gateway/internal/repository"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
	"time"
)

type PollingService interface {
	Poll(model *model.Transaction, t string) *model.Transaction
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

func (p PollingServiceImpl) Poll(model *model.Transaction, t string) *model.Transaction {
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
		if p.isUpdate(pendingTransaction, t) {
			p.logger.Info("Transaction status updated", zap.String("TransactionId", transactionId))
			return pendingTransaction
		}
		time.Sleep(2 * time.Second) // Poll every 2 seconds
	}

	return model
}

func (p PollingServiceImpl) isUpdate(model *model.Transaction, t string) bool {
	updatedBy := strings.ToUpper(model.UpdatedBy)
	status := strings.ToUpper(model.Status)
	errorCode := model.ErrorCode

	switch t {
	case "CHANGE":
		return updatedBy != "SERVER"

	case "REFUND":
		isUpdatedByTerminal := updatedBy == "TERMINAL"
		isRefunded := status == "REFUNDED"
		hasError := errorCode != "0"
		isUpdatedByTcpServer := updatedBy == "TCP_SERVER"

		return (isUpdatedByTerminal && (isRefunded || hasError)) ||
			isUpdatedByTcpServer

	case "VOID":
		isUpdatedByTerminal := updatedBy == "TERMINAL"
		isStarted := status == "STARTED"
		hasError := errorCode != "0"
		isUpdatedByTcpServer := updatedBy == "TCP_SERVER"
		isUpdatedByNotify := updatedBy == "NOTIFY"

		return (isUpdatedByTerminal && (isStarted || hasError)) ||
			isUpdatedByTcpServer ||
			isUpdatedByNotify

	default:
		return false
	}
}

func NewPollingService(logger *zap.Logger, txRepo repository.TransactionRepository) PollingService {
	return &PollingServiceImpl{
		logger: logger,
		txRepo: txRepo,
	}
}
