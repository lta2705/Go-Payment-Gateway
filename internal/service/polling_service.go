package service

import (
	"github.com/bytedance/gopkg/util/logger"
	"github.com/joho/godotenv"
	"github.com/lta2705/Go-Payment-Gateway/internal/constant"
	_ "github.com/lta2705/Go-Payment-Gateway/internal/dto"
	"github.com/lta2705/Go-Payment-Gateway/internal/model"
	"github.com/lta2705/Go-Payment-Gateway/internal/repository"
	"os"
	"strconv"
	"strings"
	"time"
)

type PollingService interface {
	Poll(model *model.Transaction, t string) *model.Transaction
}

type PollingServiceImpl struct {
	txRepo repository.TransactionRepository
}

func (p PollingServiceImpl) getTimeout() int {

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file", err)
	}
	timeoutStr := os.Getenv("POLLING_MAX_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		logger.Error("Error converting POLLING_MAX_TIMEOUT to int", err)
		return 60 // default timeout
	}
	return timeout
}

func (p PollingServiceImpl) Poll(model *model.Transaction, mode string) *model.Transaction {
	startTime := time.Now()
	timeout := time.Duration(p.getTimeout()) * time.Millisecond
	transactionId := model.TransactionId
	for time.Since(startTime) < timeout {
		logger.Info("Polling for transaction status...")
		// Here you would add the logic to check transaction statuses
		pendingTransaction, err := p.txRepo.FindByTransactionId(transactionId)
		if err != nil {
			logger.Error("Error fetching transaction during polling", err)
		}
		if p.isUpdated(pendingTransaction, mode) {
			logger.Info("Transaction status updated", "TransactionId", transactionId)
			//pendingTransaction.Status = constant.TxStatusSuccess
			//pendingTransaction.ErrorCode = constant.ErrCodeNoErr
			//pendingTransaction.ErrorDetail = constant.ErrDetailCode0

			return pendingTransaction
		}
		time.Sleep(2 * time.Second) // Poll every 2 seconds
	}

	logger.Warn("Polling timeout reached without status update", "TransactionId", transactionId)

	model.ErrorCode = constant.ErrCodeTrmNotResponse
	model.ErrorDetail = constant.ErrDetailCode11
	model.Status = constant.TxStatusFailed
	return model
}

func (p PollingServiceImpl) isUpdated(model *model.Transaction, mode string) bool {
	updatedBy := strings.ToUpper(model.UpdatedBy)
	status := strings.ToUpper(model.Status)
	errorCode := model.ErrorCode

	switch mode {
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

func NewPollingService(txRepo repository.TransactionRepository) PollingService {
	return &PollingServiceImpl{
		txRepo: txRepo,
	}
}
