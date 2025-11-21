package service

import (
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/lta2705/Go-Payment-Gateway/internal/constant"
	"github.com/lta2705/Go-Payment-Gateway/internal/dto"
	"github.com/lta2705/Go-Payment-Gateway/internal/model"
	"github.com/lta2705/Go-Payment-Gateway/internal/repository"
	"go.uber.org/zap"
)

type CardService interface {
	CreateCardTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
}

type CardServiceImpl struct {
	txRepo repository.TransactionRepository
	logger *zap.Logger
	pollingService PollingService
}


func (c *CardServiceImpl) CreateCardTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {
	defer c.logger.Sync()

	pcPosId := dto.PcPosId
	transactionId := dto.TransactionId

	transaction, _ := c.txRepo.FindByPcPosIdAndTransactionId(pcPosId, transactionId)

	if transaction != nil {
		c.logger.Info("Sale transaction already exists", zap.String("PcPosId", pcPosId), zap.String("TransactionId", transactionId))

		dto.Status = "FAILED"
		dto.ErrorCode = "01"
		dto.ErrorDetail = "Transaction already exists"
		return dto, nil
	}

	c.logger.Info("Creating new sale transaction", zap.String("PcPosId", pcPosId), zap.String("TransactionId", transactionId))

	newTransaction := &model.Transaction{}

	err := copier.Copy(newTransaction, dto)
	if err != nil {
		c.logger.Error("Error copying transaction DTO to model", zap.Error(err))
		return nil, err
	}

	newTransaction.UpdatedBy = "SERVER"
	newTransaction.ID = uuid.New()

	c.logger.Info("New transaction before insert:", zap.Any("Transaction", newTransaction))

	error := c.txRepo.CreateTransaction(newTransaction)
	if error != nil {
		c.logger.Error("Error creating new sale transaction in DB", zap.Error(error), zap.String("TransactionId", transactionId))
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeTcpServerError
		dto.ErrorDetail = constant.ErrDetailCode3

		return dto, error
	}

	c.logger.Info("Successfully created new sale transaction in DB", zap.Any("Payload", &transaction))

	updatedTransaction := c.pollingService.Poll(newTransaction, "CHANGE")

	updatedTransaction.ErrorCode = constant.ErrCodeNoErr
	updatedTransaction.ErrorDetail = constant.ErrDetailCode0
	updatedTransaction.Status = constant.TxStatusSuccess

	err = copier.Copy(dto, updatedTransaction)
	if err != nil {
		c.logger.Error("Error copying final model to DTO", zap.Error(err))
		return nil, err
	}

	return dto, nil
}

func NewCardService(txRepo repository.TransactionRepository, logger *zap.Logger, pollingService PollingService) CardService {
	return &CardServiceImpl{
		txRepo:         txRepo,
		logger:         logger,
		pollingService: NewPollingService(txRepo, logger),
	}
}