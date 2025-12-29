package service

import (
	"github.com/bytedance/gopkg/util/logger"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/lta2705/Go-Payment-Gateway/internal/constant"
	"github.com/lta2705/Go-Payment-Gateway/internal/dto"
	"github.com/lta2705/Go-Payment-Gateway/internal/model"
	"github.com/lta2705/Go-Payment-Gateway/internal/repository"
	"go.uber.org/zap"
)

type QRService interface {
	CreateQRTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
}

type QRServiceImpl struct {
	TxRepo         repository.TransactionRepository
	pollingService PollingService
}

func (t *QRServiceImpl) CreateQRTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {

	pcPosId := dto.PcPosId
	transactionId := dto.TransactionId

	QRTransaction, _ := t.TxRepo.FindByPcPosIdAndTransactionId(pcPosId, transactionId)

	if QRTransaction != nil {
		logger.Info("QR transaction already exists", pcPosId, transactionId)

		dto.Status = "FAILED"
		dto.ErrorCode = "01"
		dto.ErrorDetail = "Transaction already exists"
		return dto, nil
	}

	logger.Info("Creating new QR transaction", zap.String("PcPosId", pcPosId), zap.String("TransactionId", transactionId))

	newTransaction := &model.Transaction{}

	err := copier.Copy(newTransaction, dto)
	if err != nil {
		logger.Error("Error copying transaction DTO to model", zap.Error(err))
		return nil, err
	}

	newTransaction.UpdatedBy = "SERVER"
	newTransaction.ID = uuid.New()

	logger.Info("New transaction before insert:", zap.Any("Transaction", newTransaction))

	error := t.TxRepo.CreateTransaction(newTransaction)
	if error != nil {
		logger.Error("Error creating new QR transaction in DB", zap.Error(error), zap.String("TransactionId", transactionId))
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeTcpServerError
		dto.ErrorDetail = constant.ErrDetailCode3

		return dto, error
	}

	logger.Info("Successfully created new QR transaction in DB", zap.Any("Payload", &QRTransaction))

	updatedTransaction := t.pollingService.Poll(newTransaction, "CHANGE")

	err = copier.Copy(dto, updatedTransaction)
	if err != nil {
		logger.Error("Error copying final model to DTO", zap.Error(err))
		return nil, err
	}

	return dto, nil
}

func NewQRService(txRepo repository.TransactionRepository, pollingService PollingService) QRService {
	return &QRServiceImpl{
		TxRepo:         txRepo,
		pollingService: NewPollingService(txRepo),
	}
}
