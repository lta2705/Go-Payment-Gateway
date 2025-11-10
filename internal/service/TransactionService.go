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
	TxRepo repository.TransactionRepository
	logger *zap.Logger
}

func (t TransactionServiceImpl) CreateSaleTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {
	pcPosId := dto.PcPosId
	transactionId := dto.TransactionId
	transaction, err := t.TxRepo.FindByPcPosIdAndTransactionId(pcPosId, transactionId)
	if err != nil && transaction == nil {
		t.logger.Info("Creating new sale transaction", zap.String("PcPosId", pcPosId), zap.String("TransactionId", transactionId))
		var newTransaction *model.Transaction
		// Map dto to model.Transaction
		err := copier.Copy(newTransaction, &transaction)
		if err != nil {
			t.logger.Error("Error copying transaction DTO to model", zap.Error(err))
			return nil, err
		}
		t.TxRepo.CreateTransaction(newTransaction)
		t.logger.Info("Created new sale transaction", zap.String("TransactionId", transactionId))

		// Call a polling to check the transaction status
		return nil, nil
	}
	t.logger.Info("Sale transaction already exists", zap.String("PcPosId", pcPosId), zap.String("TransactionId", transactionId))
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

func NewTransactionService(TxRepo *repository.TransactionRepository, logger *zap.Logger) TransactionService {
	return &TransactionServiceImpl{
		TxRepo: *TxRepo,
		logger: logger,
	}
}
