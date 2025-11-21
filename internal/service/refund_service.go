package service

import (
	"github.com/jinzhu/copier"
	"github.com/lta2705/Go-Payment-Gateway/internal/constant"
	"github.com/lta2705/Go-Payment-Gateway/internal/dto"
	"github.com/lta2705/Go-Payment-Gateway/internal/model"
	"github.com/lta2705/Go-Payment-Gateway/internal/repository"
	"go.uber.org/zap"
)

type RefundService interface {
	CreateRefundTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
}

type RefundServiceImpl struct {
	TxRepo         repository.TransactionRepository
	pollingService PollingService
	logger         *zap.Logger
}

func (t *RefundServiceImpl) CreateRefundTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {
	var refundTransaction = &model.Transaction{}

	err := copier.Copy(refundTransaction, dto)
	if err != nil {
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeTcpServerError
		dto.ErrorDetail = constant.ErrDetailCode3
		return nil, err
	}

	_, err = t.TxRepo.FindByPcPosIdAndTransactionId(refundTransaction.PcPosId, refundTransaction.OrgPcPosTxnId)

	return nil, nil
}

func NewRefundService(txRepo repository.TransactionRepository, logger *zap.Logger, pollingService PollingService) RefundService {
	return &RefundServiceImpl{
		TxRepo:         txRepo,
		logger:         logger,
		pollingService: NewPollingService(txRepo, logger),
	}
}

