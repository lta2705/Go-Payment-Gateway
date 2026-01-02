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

type CardService interface {
	CreateCardTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
}

type CardServiceImpl struct {
	txRepo         repository.TransactionRepository
	pollingService PollingService
	sender         *worker.KafkaProducerWorker
}

func (c *CardServiceImpl) CreateCardTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {

	pcPosId := dto.PcPosId
	transactionId := dto.TransactionId

	transaction, _ := c.txRepo.FindByPcPosIdAndTransactionId(pcPosId, transactionId)

	if transaction != nil {
		logger.Info("Sale transaction already exists", "PcPosId", pcPosId, "TransactionId", transactionId)

		dto.Status = "FAILED"
		dto.ErrorCode = "01"
		dto.ErrorDetail = "Transaction already exists"
		return dto, nil
	}

	logger.Info("Creating new sale transaction", "PcPosId", pcPosId, "TransactionId", transactionId)

	newTransaction := &model.Transaction{}

	err := copier.Copy(newTransaction, dto)
	if err != nil {
		logger.Error("Error copying transaction DTO to model", err)
		return nil, err
	}

	newTransaction.UpdatedBy = "SERVER"
	newTransaction.ID = uuid.New()

	logger.Info("New transaction before insert:", "Transaction", newTransaction)

	error := c.txRepo.CreateTransaction(newTransaction)
	if error != nil {
		logger.Error("Error creating new sale transaction in DB", error, "TransactionId", transactionId)
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeTcpServerError
		dto.ErrorDetail = constant.ErrDetailCode3

		return dto, error
	}

	jsonData, err := json.Marshal(newTransaction)
	if err != nil {
		logger.Error("Failed to marshal transaction to JSON", err)
	} else {
		// 2. Gửi qua Kafka
		senderErr := c.sender.SendMessage(string(jsonData))
		if senderErr != nil {
			logger.Error("Failed to produce message to sender server", senderErr, "TransactionId", transactionId)
			return nil, senderErr
		}
		logger.Info("Successfully sent transaction to Kafka", "ID", newTransaction.ID.String())
	}

	logger.Info("Starting polling for transaction status update...")

	updatedTransaction := c.pollingService.Poll(newTransaction, "CHANGE")

	err = copier.Copy(dto, updatedTransaction)
	if err != nil {
		logger.Error("Error copying final model to DTO", err)
		return nil, err
	}

	return dto, nil
}

func NewCardService(txRepo repository.TransactionRepository, pollingService PollingService, sender *worker.KafkaProducerWorker) CardService {
	return &CardServiceImpl{
		txRepo:         txRepo,
		pollingService: NewPollingService(txRepo),
		sender:         sender,
	}
}
