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

type RefundService interface {
	CreateRefundTransaction(req *dto.TransactionDTO) (*dto.TransactionDTO, error)
}

type RefundServiceImpl struct {
	txRepo               repository.TransactionRepository
	merchantTerminalRepo repository.MerchantTerminalRepository
	pollingService       PollingService
	sender               worker.KafkaProducerWorker
}

func (s *RefundServiceImpl) CreateRefundTransaction(req *dto.TransactionDTO) (*dto.TransactionDTO, error) {

	// 1. Idempotency check
	existedRefund, err := s.txRepo.FindByPcPosIdAndTransactionId(req.PcPosId, req.TransactionId)
	if err != nil {
		logger.Error("Failed to check existing refund transaction", err)
		return nil, err
	}

	if existedRefund != nil {
		logger.Info(
			"Refund transaction already exists",
			"PcPosId", req.PcPosId,
			"TransactionId", req.TransactionId,
		)

		req.Status = constant.TxStatusFailed
		req.ErrorCode = constant.ErrCodeDuplicateTx
		req.ErrorDetail = "Refund transaction already exists"

		return req, nil
	}

	// 2. Validate original transaction
	orgTx, err := s.txRepo.FindByPcPosIdAndTransactionId(req.PcPosId, req.OrgPcPosTxnId)
	if err != nil {
		logger.Error("Failed to query original transaction for refund", err)
		req.Status = constant.TxStatusFailed
		req.ErrorCode = constant.ErrCodeTcpServerError
		req.ErrorDetail = constant.ErrDetailCode3
		return req, err
	}

	if orgTx == nil {
		logger.Warn(
			"Original transaction not found for refund",
			"PcPosId", req.PcPosId,
			"OrgPcPosTxnId", req.OrgPcPosTxnId,
		)

		req.Status = constant.TxStatusFailed
		req.ErrorCode = constant.ErrCodeNotFoundOriginTx
		req.ErrorDetail = constant.ErrDetailCode7

		return req, nil
	}

	// (Giữ nguyên logic gốc: không kiểm tra SUCCESS)

	// 3. Map DTO → Refund Transaction Model
	refundTx := &model.Transaction{}
	if err := copier.Copy(refundTx, req); err != nil {
		logger.Error("Failed to copy refund DTO to model", err)
		return nil, err
	}

	refundTx.ID = uuid.New()
	refundTx.UpdatedBy = "SERVER"

	logger.Info("Creating new refund transaction", "Transaction", refundTx)

	// 4. Persist refund transaction
	if err := s.txRepo.CreateTransaction(refundTx); err != nil {
		logger.Error(
			"Failed to create refund transaction in DB",
			err,
			"TransactionId", req.TransactionId,
		)

		req.Status = constant.TxStatusFailed
		req.ErrorCode = constant.ErrCodeTcpServerError
		req.ErrorDetail = constant.ErrDetailCode3

		return req, err
	}

	// 5. Resolve TerminalId
	terminalId, err := s.merchantTerminalRepo.FindTerminalIdByPcPosId(refundTx.PcPosId)
	if err != nil {
		logger.Error(
			"Failed to find terminalId by PcPosId",
			err,
			"PcPosId", refundTx.PcPosId,
		)
		return nil, err
	}

	// 6. Enrich payload & publish to Kafka (GIỐNG CardService)
	payload, err := utils.EnrichTransactionToJSON(refundTx, terminalId)
	if err != nil {
		logger.Error("Failed to enrich refund transaction payload", err)
		return nil, err
	}

	if err := s.sender.SendMessage(payload); err != nil {
		logger.Error(
			"Failed to send refund transaction to Kafka",
			err,
			"TransactionId", refundTx.ID.String(),
		)
		return nil, err
	}

	logger.Info(
		"Refund transaction sent to Kafka successfully",
		"TransactionId", refundTx.ID.String(),
		"TerminalId", terminalId,
	)

	// 7. Polling for refund transaction status
	logger.Info(
		"Polling for refund transaction status update",
		"TransactionId", refundTx.ID.String(),
	)

	updatedTx := s.pollingService.Poll(refundTx, "REFUND")
	if updatedTx == nil {
		return nil, errors.New("polling returned nil transaction")
	}

	// 8. Map Model → DTO
	if err := copier.Copy(req, updatedTx); err != nil {
		logger.Error("Failed to copy final refund model to DTO", err)
		return nil, err
	}

	return req, nil
}

func NewRefundService(
	txRepo repository.TransactionRepository,
	merchantTerminalRepo repository.MerchantTerminalRepository,
	pollingService PollingService,
	sender worker.KafkaProducerWorker,
) RefundService {
	return &RefundServiceImpl{
		txRepo:               txRepo,
		merchantTerminalRepo: merchantTerminalRepo,
		pollingService:       pollingService,
		sender:               sender,
	}
}
