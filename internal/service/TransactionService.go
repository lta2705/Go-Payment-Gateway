package service

import (
	"github.com/jinzhu/copier"
	"github.com/lta2705/Go-Payment-Gateway/internal/dto"
	"github.com/lta2705/Go-Payment-Gateway/internal/model"
	"github.com/lta2705/Go-Payment-Gateway/internal/repository"
	"go.uber.org/zap"
)

type TransactionService interface {
	CreateSaleTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
	CreateVoidTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
	CreateRefundTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
	CreateQRTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
}
type TransactionServiceImpl struct {
	TxRepo         repository.TransactionRepository
	logger         *zap.Logger
	pollingService PollingService
}

func (t *TransactionServiceImpl) CreateSaleTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {
	pcPosId := dto.PcPosId
	transactionId := dto.TransactionId

	transaction, err := t.TxRepo.FindByPcPosIdAndTransactionId(pcPosId, transactionId)

	if err != nil && transaction != nil {
		t.logger.Error("Error checking transaction existence", zap.Error(err), zap.String("PcPosId", pcPosId))
		return nil, err // Trả về lỗi hệ thống
	}

	if transaction != nil {
		t.logger.Info("Sale transaction already exists", zap.String("PcPosId", pcPosId), zap.String("TransactionId", transactionId))

		dto.Status = "FAILED"
		dto.ErrorCode = "01"
		dto.ErrorDetail = "Transaction already exists"
		return dto, nil
	}

	t.logger.Info("Creating new sale transaction", zap.String("PcPosId", pcPosId), zap.String("TransactionId", transactionId))

	newTransaction := &model.Transaction{}
	err = copier.Copy(newTransaction, dto)
	if err != nil {
		t.logger.Error("Error copying transaction DTO to model", zap.Error(err))
		return nil, err
	}

	err = t.TxRepo.CreateTransaction(newTransaction)
	if err != nil {
		t.logger.Error("Error creating new sale transaction in DB", zap.Error(err), zap.String("TransactionId", transactionId))
		return nil, err
	}

	t.logger.Info("Successfully created new sale transaction in DB", zap.String("TransactionId", transactionId))

	updatedTransaction := t.pollingService.Poll(newTransaction, "CHANGE")

	updatedTransaction.ErrorCode = "00"
	updatedTransaction.ErrorDetail = "Approval"
	updatedTransaction.Status = "SUCCESS"

	err = copier.Copy(dto, updatedTransaction)
	if err != nil {
		t.logger.Error("Error copying final model to DTO", zap.Error(err))
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

func NewTransactionService(TxRepo *repository.TransactionRepository, logger *zap.Logger) TransactionService {
	return &TransactionServiceImpl{
		TxRepo: *TxRepo,
		logger: logger,
	}
}
