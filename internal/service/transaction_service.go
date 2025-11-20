package service

import (
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/lta2705/Go-Payment-Gateway/internal/dto"
	"github.com/lta2705/Go-Payment-Gateway/internal/middleware"
	"github.com/lta2705/Go-Payment-Gateway/internal/model"
	"github.com/lta2705/Go-Payment-Gateway/internal/repository"
	"go.uber.org/zap"
	"time"
)

type TransactionService interface {
	CreateSaleTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
	CreateVoidTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
	CreateRefundTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
	CreateQRTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
}
type TransactionServiceImpl struct {
	TxRepo         repository.TransactionRepository
	pollingService PollingService
}

func (t *TransactionServiceImpl) CreateSaleTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {
	logger := middleware.CreateLogger()
	defer logger.Sync()

	pcPosId := dto.PcPosId
	transactionId := dto.TransactionId

	transaction, _ := t.TxRepo.FindByPcPosIdAndTransactionId(pcPosId, transactionId)

	if transaction != nil {
		logger.Info("Sale transaction already exists", zap.String("PcPosId", pcPosId), zap.String("TransactionId", transactionId))

		dto.Status = "FAILED"
		dto.ErrorCode = "01"
		dto.ErrorDetail = "Transaction already exists"
		return dto, nil
	}

	logger.Info("Creating new sale transaction", zap.String("PcPosId", pcPosId), zap.String("TransactionId", transactionId))

	newTransaction := &model.Transaction{}

	err := copier.Copy(newTransaction, dto)
	if err != nil {
		logger.Error("Error copying transaction DTO to model", zap.Error(err))
		return nil, err
	}

	newTransaction.UpdatedBy = "SERVER"
	newTransaction.ID = uuid.New()
	newTransaction.CreatedAt = time.Now()

	logger.Info("New transaction before insert:", zap.Any("Transaction", newTransaction))

	error := t.TxRepo.CreateTransaction(newTransaction)
	if error != nil {
		logger.Error("Error creating new sale transaction in DB", zap.Error(error), zap.String("TransactionId", transactionId))
		return nil, error
	}

	logger.Info("Successfully created new sale transaction in DB", zap.String("TransactionId", transactionId))

	updatedTransaction := t.pollingService.Poll(newTransaction, "CHANGE")

	updatedTransaction.ErrorCode = "00"
	updatedTransaction.ErrorDetail = "Approval"
	updatedTransaction.Status = "SUCCESS"

	err = copier.Copy(dto, updatedTransaction)
	if err != nil {
		logger.Error("Error copying final model to DTO", zap.Error(err))
		return nil, err
	}

	return dto, nil
}

func (t TransactionServiceImpl) CreateVoidTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {
	return nil, nil
}

func (t TransactionServiceImpl) CreateRefundTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {

	return nil, nil
}

func (t TransactionServiceImpl) CreateQRTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {
	return nil, nil
}

func (t TransactionServiceImpl) createTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {
	return nil, nil
}

func NewTransactionService(TxRepo repository.TransactionRepository) TransactionService {
	return &TransactionServiceImpl{
		TxRepo:         TxRepo,
		pollingService: NewPollingService(TxRepo),
	}
}
