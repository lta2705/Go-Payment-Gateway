package service

import (
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/lta2705/Go-Payment-Gateway/internal/constant"
	"github.com/lta2705/Go-Payment-Gateway/internal/dto"
	"github.com/lta2705/Go-Payment-Gateway/internal/model"
	"github.com/lta2705/Go-Payment-Gateway/internal/repository"
	"go.uber.org/zap"
)

type VoidService interface {
	CreateVoidTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
}

type VoidServiceImpl struct {
	TxRepo         repository.TransactionRepository
	pollingService PollingService
	logger         *zap.Logger
}

func (t VoidServiceImpl) CreateVoidTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {
	var voidTransaction = &model.Transaction{}

	err := copier.Copy(voidTransaction, dto)
	if err != nil {
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeCannotMapping
		dto.ErrorDetail = constant.ErrDetailCode16
		return nil, err
	}
	voidTransaction, err = t.TxRepo.FindByPcPosIdAndTransactionId(voidTransaction.PcPosId, voidTransaction.OrgPcPosTxnId)
	if err != nil {
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeTcpServerError
		dto.ErrorDetail = constant.ErrDetailCode3
		return dto, err
	}
	if voidTransaction == nil {
		t.logger.Warn("Original transaction not found for void", zap.String("PcPosId", dto.PcPosId), zap.String("OrgPcPosTxnId", dto.OrgPcPosTxnId))
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeNotFoundOriginTx
		dto.ErrorDetail = constant.ErrDetailCode7
		return dto, nil
	} else if voidTransaction.Status != constant.TxStatusSuccess {
		t.logger.Warn("Original transaction not successful for void", zap.String("PcPosId", dto.PcPosId), zap.String("OrgPcPosTxnId", dto.OrgPcPosTxnId))
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeTxNotSuccess
		dto.ErrorDetail = constant.ErrDetailCode13
		return dto, nil
	} else if voidTransaction.Status == constant.TxStatusVoided {
		t.logger.Warn("Original transaction already voided", zap.String("PcPosId", dto.PcPosId), zap.String("OrgPcPosTxnId", dto.OrgPcPosTxnId))
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeTxVoided
		dto.ErrorDetail = constant.ErrDetailCode14
		return dto, nil
	}

	voidTransaction.ID = uuid.New()
	voidTransaction.Status = constant.TxStatusStarted
	voidTransaction.ErrorCode = constant.ErrCodeNoErr
	voidTransaction.ErrorDetail = constant.ErrDetailCode0
	voidTransaction.UpdatedBy = "SERVER"

	err = t.TxRepo.CreateTransaction(voidTransaction)
	if err != nil {
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeTcpServerError
		dto.ErrorDetail = constant.ErrDetailCode3
		return dto, err
	}

	updatedTransaction := t.pollingService.Poll(voidTransaction, "VOID")

	err = copier.Copy(dto, updatedTransaction)
	if err != nil {
		return nil, err
	}

	return dto, nil
}

func NewVoidService(txRepo repository.TransactionRepository, logger *zap.Logger, pollingService PollingService) VoidService {
	return &VoidServiceImpl{
		TxRepo:         txRepo,
		logger:         logger,
		pollingService: NewPollingService(txRepo, logger),
	}
}
