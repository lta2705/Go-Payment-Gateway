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

type VoidService interface {
	CreateVoidTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
}

type VoidServiceImpl struct {
	txRepo         repository.TransactionRepository
	pollingService PollingService
	sender         worker.KafkaProducerWorker
}

func (v *VoidServiceImpl) CreateVoidTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {

	existingVoid, _ := v.txRepo.FindByPcPosIdAndTransactionId(dto.PcPosId, dto.TransactionId)
	if existingVoid != nil {
		logger.Info("Void transaction already exists", "PcPosId", dto.PcPosId, "TransactionId", dto.TransactionId)
		dto.Status = "FAILED"
		dto.ErrorCode = "01"
		dto.ErrorDetail = "Void transaction already exists"
		return dto, nil
	}

	orgTransaction, err := v.txRepo.FindByPcPosIdAndTransactionId(dto.PcPosId, dto.OrgPcPosTxnId)
	if err != nil {
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeTcpServerError
		dto.ErrorDetail = constant.ErrDetailCode3
		return dto, err
	}

	if orgTransaction == nil {
		logger.Warn("Original transaction not found for void", "PcPosId", dto.PcPosId, "OrgPcPosTxnId", dto.OrgPcPosTxnId)
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeNotFoundOriginTx
		dto.ErrorDetail = constant.ErrDetailCode7
		return dto, nil
	} else if orgTransaction.Status != constant.TxStatusSuccess {
		logger.Warn("Original transaction not successful for void", "PcPosId", dto.PcPosId, "OrgPcPosTxnId", dto.OrgPcPosTxnId)
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeTxNotSuccess
		dto.ErrorDetail = constant.ErrDetailCode13
		return dto, nil
	} else if orgTransaction.Status == constant.TxStatusVoided {
		logger.Warn("Original transaction already voided", "PcPosId", dto.PcPosId, "OrgPcPosTxnId", dto.OrgPcPosTxnId)
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeTxVoided
		dto.ErrorDetail = constant.ErrDetailCode14
		return dto, nil
	}

	// 3. Khởi tạo giao dịch Void mới từ DTO
	newVoidTx := &model.Transaction{}
	err = copier.Copy(newVoidTx, dto)
	if err != nil {
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeCannotMapping
		dto.ErrorDetail = constant.ErrDetailCode16
		return nil, err
	}

	newVoidTx.ID = uuid.New()
	newVoidTx.Status = constant.TxStatusStarted
	newVoidTx.ErrorCode = constant.ErrCodeNoErr
	newVoidTx.ErrorDetail = constant.ErrDetailCode0
	newVoidTx.UpdatedBy = "SERVER"

	err = v.txRepo.CreateTransaction(newVoidTx)
	if err != nil {
		logger.Error("Error creating void transaction in DB", err)
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeTcpServerError
		dto.ErrorDetail = constant.ErrDetailCode3
		return dto, err
	}

	jsonData, err := json.Marshal(newVoidTx)
	if err != nil {
		logger.Error("Failed to marshal void transaction", err)
	} else {
		senderErr := v.sender.SendMessage(string(jsonData))
		if senderErr != nil {
			logger.Error("Failed to produce void message to Kafka", senderErr)
			return nil, senderErr
		}
		logger.Info("Successfully sent void transaction to Kafka", "ID", newVoidTx.ID.String())
	}

	logger.Info("Starting polling for void transaction status update...")
	updatedTransaction := v.pollingService.Poll(newVoidTx, "VOID")

	// 7. Map kết quả cuối cùng trả về DTO
	err = copier.Copy(dto, updatedTransaction)
	if err != nil {
		logger.Error("Error copying final void model to DTO", err)
		return nil, err
	}

	return dto, nil
}

func NewVoidService(txRepo repository.TransactionRepository, pollingService PollingService, sender worker.KafkaProducerWorker) VoidService {
	return &VoidServiceImpl{
		txRepo:         txRepo,
		pollingService: pollingService,
		sender:         sender,
	}
}
