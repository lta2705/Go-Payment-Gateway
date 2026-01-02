package service

import (
	"encoding/json"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/lta2705/Go-Payment-Gateway/internal/constant"
	"github.com/lta2705/Go-Payment-Gateway/internal/dto"
	"github.com/lta2705/Go-Payment-Gateway/internal/model"
	"github.com/lta2705/Go-Payment-Gateway/internal/repository"
	"github.com/lta2705/Go-Payment-Gateway/internal/worker"
)

type QRService interface {
	CreateQRTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
}

type QRServiceImpl struct {
	txRepo         repository.TransactionRepository
	pollingService PollingService
	sender         *worker.KafkaProducerWorker
}

func (q *QRServiceImpl) CreateQRTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {

	pcPosId := dto.PcPosId
	transactionId := dto.TransactionId
	transaction, _ := q.txRepo.FindByPcPosIdAndTransactionId(pcPosId, transactionId)

	if transaction != nil {
		logger.Info("QR transaction already exists", "PcPosId", pcPosId, "TransactionId", transactionId)

		dto.Status = "FAILED"
		dto.ErrorCode = "01"
		dto.ErrorDetail = "Transaction already exists"
		return dto, nil
	}

	logger.Info("Creating new QR transaction", "PcPosId", pcPosId, "TransactionId", transactionId)

	// 2. Mapping DTO sang Model
	newTransaction := &model.Transaction{}
	err := copier.Copy(newTransaction, dto)
	if err != nil {
		logger.Error("Error copying transaction DTO to model", err)
		return nil, err
	}

	newTransaction.UpdatedBy = "SERVER"
	newTransaction.ID = uuid.New()

	logger.Info("New QR transaction before insert:", "Transaction", newTransaction)

	createErr := q.txRepo.CreateTransaction(newTransaction)
	if createErr != nil {
		logger.Error("Error creating new QR transaction in DB", createErr, "TransactionId", transactionId)
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeTcpServerError
		dto.ErrorDetail = constant.ErrDetailCode3
		return dto, createErr
	}

	jsonData, err := json.Marshal(newTransaction)
	if err != nil {
		logger.Error("Failed to marshal QR transaction to JSON", err)
	} else {
		senderErr := q.sender.SendMessage(string(jsonData))
		if senderErr != nil {
			logger.Error("Failed to produce QR message to sender server", senderErr, "TransactionId", transactionId)
			return nil, senderErr
		}
		logger.Info("Successfully sent QR transaction to Kafka", "ID", newTransaction.ID.String())
	}

	logger.Info("Starting polling for QR transaction status update...")
	updatedTransaction := q.pollingService.Poll(newTransaction, "CHANGE")

	err = copier.Copy(dto, updatedTransaction)
	if err != nil {
		logger.Error("Error copying final QR model to DTO", err)
		return nil, err
	}

	return dto, nil
}

func NewQRService(txRepo repository.TransactionRepository, pollingService PollingService, sender *worker.KafkaProducerWorker) QRService {
	return &QRServiceImpl{
		txRepo:         txRepo,
		pollingService: pollingService,
		sender:         sender,
	}
}
