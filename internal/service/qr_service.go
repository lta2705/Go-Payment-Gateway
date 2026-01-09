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

type QRService interface {
	CreateQRTransaction(req *dto.TransactionDTO) (*dto.TransactionDTO, error)
}

type QRServiceImpl struct {
	txRepo               repository.TransactionRepository
	merchantTerminalRepo repository.MerchantTerminalRepository
	pollingService       PollingService
	sender               worker.KafkaProducerWorker
}

func (s *QRServiceImpl) CreateQRTransaction(req *dto.TransactionDTO) (*dto.TransactionDTO, error) {

	// 1. Idempotency
	existingTx, err := s.txRepo.FindByPcPosIdAndTransactionId(req.PcPosId, req.TransactionId)
	if err != nil {
		logger.Error("Failed to check existing QR transaction", err)
		return nil, err
	}

	if existingTx != nil {
		req.Status = constant.TxStatusFailed
		req.ErrorCode = constant.ErrCodeDuplicateTx
		req.ErrorDetail = "QR transaction already exists"
		return req, nil
	}

	// 2. Map DTO → Model
	qrTx := &model.Transaction{}
	if err := copier.Copy(qrTx, req); err != nil {
		logger.Error("Failed to copy QR DTO to model", err)
		return nil, err
	}

	qrTx.ID = uuid.New()
	qrTx.Status = constant.TxStatusStarted
	qrTx.ErrorCode = constant.ErrCodeNoErr
	qrTx.ErrorDetail = constant.ErrDetailCode0
	qrTx.UpdatedBy = "SERVER"

	// 3. Persist
	if err := s.txRepo.CreateTransaction(qrTx); err != nil {
		logger.Error("Failed to create QR transaction in DB", err)
		req.Status = constant.TxStatusFailed
		req.ErrorCode = constant.ErrCodeTcpServerError
		req.ErrorDetail = constant.ErrDetailCode3
		return req, err
	}

	// 4. Resolve terminalId
	terminalId, err := s.merchantTerminalRepo.FindTerminalIdByPcPosId(qrTx.PcPosId)
	if err != nil {
		logger.Error("Failed to find terminalId for QR", err)
		return nil, err
	}

	// 5. Enrich & send Kafka
	payload, err := utils.EnrichTransactionToJSON(qrTx, terminalId)
	if err != nil {
		logger.Error("Failed to enrich QR payload", err)
		return nil, err
	}

	if err := s.sender.SendMessage(payload); err != nil {
		logger.Error("Failed to send QR message to Kafka", err)
		return nil, err
	}

	logger.Info(
		"QR transaction sent to Kafka",
		"TransactionId", qrTx.ID.String(),
		"TerminalId", terminalId,
	)

	// 6. Polling
	updatedTx := s.pollingService.Poll(qrTx, "CHANGE")
	if updatedTx == nil {
		return nil, errors.New("polling returned nil transaction")
	}

	// 7. Map result
	if err := copier.Copy(req, updatedTx); err != nil {
		logger.Error("Failed to map final QR tx to DTO", err)
		return nil, err
	}

	return req, nil
}

func NewQRService(
	txRepo repository.TransactionRepository,
	merchantTerminalRepo repository.MerchantTerminalRepository,
	pollingService PollingService,
	sender worker.KafkaProducerWorker,
) QRService {
	return &QRServiceImpl{
		txRepo:               txRepo,
		merchantTerminalRepo: merchantTerminalRepo,
		pollingService:       pollingService,
		sender:               sender,
	}
}
