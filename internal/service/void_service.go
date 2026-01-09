package service

import (
	"errors"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"

	"github.com/lta2705/Go-Payment-Gateway/internal/constant"
	"github.com/lta2705/Go-Payment-Gateway/internal/dto"
	"github.com/lta2705/Go-Payment-Gateway/internal/model"
	"github.com/lta2705/Go-Payment-Gateway/internal/repository"
	"github.com/lta2705/Go-Payment-Gateway/internal/worker"
	"github.com/lta2705/Go-Payment-Gateway/utils"
)

type VoidService interface {
	CreateVoidTransaction(req *dto.TransactionDTO) (*dto.TransactionDTO, error)
}

type VoidServiceImpl struct {
	txRepo               repository.TransactionRepository
	merchantTerminalRepo repository.MerchantTerminalRepository
	pollingService       PollingService
	sender               worker.KafkaProducerWorker
}

func (s *VoidServiceImpl) CreateVoidTransaction(req *dto.TransactionDTO) (*dto.TransactionDTO, error) {

	// 1. Idempotency
	existingVoid, err := s.txRepo.FindByPcPosIdAndTransactionId(req.PcPosId, req.TransactionId)
	if err != nil {
		logger.Error("Failed to check existing void transaction", err)
		return nil, err
	}

	if existingVoid != nil {
		logger.Info(
			"Void transaction already exists",
			"PcPosId", req.PcPosId,
			"TransactionId", req.TransactionId,
		)

		req.Status = constant.TxStatusFailed
		req.ErrorCode = constant.ErrCodeDuplicateTx
		req.ErrorDetail = "Void transaction already exists"
		return req, nil
	}

	// 2. Validate original transaction
	orgTx, err := s.txRepo.FindByPcPosIdAndTransactionId(req.PcPosId, req.OrgPcPosTxnId)
	if err != nil {
		logger.Error("Failed to query original transaction for void", err)
		req.Status = constant.TxStatusFailed
		req.ErrorCode = constant.ErrCodeTcpServerError
		req.ErrorDetail = constant.ErrDetailCode3
		return req, err
	}

	if orgTx == nil {
		req.Status = constant.TxStatusFailed
		req.ErrorCode = constant.ErrCodeNotFoundOriginTx
		req.ErrorDetail = constant.ErrDetailCode7
		return req, nil
	}

	if orgTx.Status != constant.TxStatusSuccess {
		req.Status = constant.TxStatusFailed
		req.ErrorCode = constant.ErrCodeTxNotSuccess
		req.ErrorDetail = constant.ErrDetailCode13
		return req, nil
	}

	if orgTx.Status == constant.TxStatusVoided {
		req.Status = constant.TxStatusFailed
		req.ErrorCode = constant.ErrCodeTxVoided
		req.ErrorDetail = constant.ErrDetailCode14
		return req, nil
	}

	// 3. Map DTO → Model
	voidTx := &model.Transaction{}
	if err := copier.Copy(voidTx, req); err != nil {
		logger.Error("Failed to copy void DTO to model", err)
		return nil, err
	}

	voidTx.ID = uuid.New()
	voidTx.Status = constant.TxStatusStarted
	voidTx.ErrorCode = constant.ErrCodeNoErr
	voidTx.ErrorDetail = constant.ErrDetailCode0
	voidTx.UpdatedBy = "SERVER"

	// 4. Persist
	if err := s.txRepo.CreateTransaction(voidTx); err != nil {
		logger.Error("Failed to create void transaction in DB", err)
		req.Status = constant.TxStatusFailed
		req.ErrorCode = constant.ErrCodeTcpServerError
		req.ErrorDetail = constant.ErrDetailCode3
		return req, err
	}

	// 5. Resolve terminalId
	terminalId, err := s.merchantTerminalRepo.FindTerminalIdByPcPosId(voidTx.PcPosId)
	if err != nil {
		logger.Error("Failed to find terminalId for void", err)
		return nil, err
	}

	// 6. Enrich & send Kafka
	payload, err := utils.EnrichTransactionToJSON(voidTx, terminalId)
	if err != nil {
		logger.Error("Failed to enrich void payload", err)
		return nil, err
	}

	if err := s.sender.SendMessage(payload); err != nil {
		logger.Error("Failed to send void message to Kafka", err)
		return nil, err
	}

	logger.Info(
		"Void transaction sent to Kafka",
		"TransactionId", voidTx.ID.String(),
		"TerminalId", terminalId,
	)

	// 7. Polling
	updatedTx := s.pollingService.Poll(voidTx, "VOID")
	if updatedTx == nil {
		return nil, errors.New("polling returned nil transaction")
	}

	// 8. Map result
	if err := copier.Copy(req, updatedTx); err != nil {
		logger.Error("Failed to map final void tx to DTO", err)
		return nil, err
	}

	return req, nil
}

func NewVoidService(
	txRepo repository.TransactionRepository,
	merchantTerminalRepo repository.MerchantTerminalRepository,
	pollingService PollingService,
	sender worker.KafkaProducerWorker,
) VoidService {
	return &VoidServiceImpl{
		txRepo:               txRepo,
		merchantTerminalRepo: merchantTerminalRepo,
		pollingService:       pollingService,
		sender:               sender,
	}
}
