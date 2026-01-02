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

type RefundService interface {
	CreateRefundTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
}

type RefundServiceImpl struct {
	txRepo         repository.TransactionRepository
	pollingService PollingService
	sender         *worker.KafkaProducerWorker
}

func (r *RefundServiceImpl) CreateRefundTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {
	pcPosId := dto.PcPosId
	transactionId := dto.TransactionId
	orgPcPosTxnId := dto.OrgPcPosTxnId

	existingRefund, _ := r.txRepo.FindByPcPosIdAndTransactionId(pcPosId, transactionId)
	if existingRefund != nil {
		logger.Info("Refund transaction already exists", "PcPosId", pcPosId, "TransactionId", transactionId)
		dto.Status = "FAILED"
		dto.ErrorCode = "01"
		dto.ErrorDetail = "Refund transaction already exists"
		return dto, nil
	}

	orgTransaction, err := r.txRepo.FindByPcPosIdAndTransactionId(pcPosId, orgPcPosTxnId)
	if err != nil || orgTransaction == nil {
		logger.Error("Original transaction not found for refund", err, "OrgId", orgPcPosTxnId)
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeNotFoundOriginTx
		dto.ErrorDetail = constant.ErrDetailCode7
		return dto, err
	}

	logger.Info("Creating new refund transaction", "PcPosId", pcPosId, "TransactionId", transactionId)

	// 3. Mapping DTO sang Model
	newRefund := &model.Transaction{}
	err = copier.Copy(newRefund, dto)
	if err != nil {
		logger.Error("Error copying refund DTO to model", err)
		return nil, err
	}

	newRefund.UpdatedBy = "SERVER"
	newRefund.ID = uuid.New()

	// 4. Lưu vào Database
	createErr := r.txRepo.CreateTransaction(newRefund)
	if createErr != nil {
		logger.Error("Error creating refund transaction in DB", createErr, "TransactionId", transactionId)
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeTcpServerError
		dto.ErrorDetail = constant.ErrDetailCode3
		return dto, createErr
	}

	jsonData, err := json.Marshal(newRefund)
	if err != nil {
		logger.Error("Failed to marshal refund transaction", err)
	} else {
		senderErr := r.sender.SendMessage(string(jsonData))
		if senderErr != nil {
			logger.Error("Failed to produce refund message to Kafka", senderErr)
			return nil, senderErr
		}
		logger.Info("Successfully sent refund to Kafka", "ID", newRefund.ID.String())
	}

	logger.Info("Waiting for refund processing...")
	updatedTransaction := r.pollingService.Poll(newRefund, "REFUND")

	err = copier.Copy(dto, updatedTransaction)
	if err != nil {
		logger.Error("Error copying final refund model to DTO", err)
		return nil, err
	}

	return dto, nil
}

func NewRefundService(txRepo repository.TransactionRepository, pollingService PollingService, sender *worker.KafkaProducerWorker) RefundService {
	return &RefundServiceImpl{
		txRepo:         txRepo,
		pollingService: pollingService,
		sender:         sender,
	}
}
