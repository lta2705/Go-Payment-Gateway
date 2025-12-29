package service

import (
	"context"
	"encoding/json"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/lta2705/Go-Payment-Gateway/internal/constant"
	"github.com/lta2705/Go-Payment-Gateway/internal/dto"
	"github.com/lta2705/Go-Payment-Gateway/internal/model"
	"github.com/lta2705/Go-Payment-Gateway/internal/repository"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type CardService interface {
	CreateCardTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
}

type CardServiceImpl struct {
	txRepo         repository.TransactionRepository
	pollingService PollingService
	kafkaProducer  *kafka.Writer
}

func (c *CardServiceImpl) CreateCardTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {

	pcPosId := dto.PcPosId
	transactionId := dto.TransactionId

	transaction, _ := c.txRepo.FindByPcPosIdAndTransactionId(pcPosId, transactionId)

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

	logger.Info("New transaction before insert:", zap.Any("Transaction", newTransaction))

	error := c.txRepo.CreateTransaction(newTransaction)
	if error != nil {
		logger.Error("Error creating new sale transaction in DB", zap.Error(error), zap.String("TransactionId", transactionId))
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeTcpServerError
		dto.ErrorDetail = constant.ErrDetailCode3

		return dto, error
	}

	jsonData, err := json.Marshal(newTransaction)
	if err != nil {
		logger.Error("Failed to marshal transaction to JSON", zap.Error(err))
	} else {
		// 2. Gửi qua Kafka
		err = c.kafkaProducer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(newTransaction.ID.String()),
			Value: jsonData,
		})

		if err != nil {
			logger.Error("Failed to send message to Kafka", zap.Error(err))
			// Tùy nghiệp vụ mà bạn có trả về lỗi hay không,
			// thông thường nếu DB đã lưu thì vẫn tiếp tục polling.
		} else {
			logger.Info("Successfully sent transaction to Kafka", zap.String("ID", newTransaction.ID.String()))
		}
	}
	// --- KẾT THÚC LOGIC KAFKA ---

	logger.Info("Starting polling for transaction status update...")

	// Đợi service khác xử lý và cập nhật DB, polling sẽ bắt được thay đổi này
	updatedTransaction := c.pollingService.Poll(newTransaction, "CHANGE")

	err = copier.Copy(dto, updatedTransaction)
	if err != nil {
		logger.Error("Error copying final model to DTO", zap.Error(err))
		return nil, err
	}

	return dto, nil
}

func NewCardService(txRepo repository.TransactionRepository, pollingService PollingService, kafkaProducer *kafka.Writer) CardService {
	return &CardServiceImpl{
		txRepo:         txRepo,
		pollingService: NewPollingService(txRepo),
		kafkaProducer:  kafkaProducer,
	}
}
