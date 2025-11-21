package service

import (
	"github.com/jinzhu/copier"
	"github.com/lta2705/Go-Payment-Gateway/internal/constant"
	"github.com/lta2705/Go-Payment-Gateway/internal/dto"
	"github.com/lta2705/Go-Payment-Gateway/internal/model"
	"github.com/lta2705/Go-Payment-Gateway/internal/repository"
	"go.uber.org/zap"
)

type CheckStatusService interface {
	CheckTransactionStatus(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
}
type CheckStatusServiceImpl struct {
	TxRepo repository.TransactionRepository
	logger *zap.Logger
}
	
func (t CheckStatusServiceImpl) CheckTransactionStatus(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {
	var transaction = &model.Transaction{}

	err := copier.Copy(transaction, dto)
	if err != nil {
		dto.ErrorCode = constant.ErrCodeTcpServerError
		dto.ErrorDetail = constant.ErrDetailCode3
		dto.Status = constant.TxStatusFailed
		return dto, err
	}

	transaction, err = t.TxRepo.FindByTransactionId(transaction.TransactionId)
	if err != nil {
		dto.ErrorCode = constant.ErrCodeTcpServerError
		dto.ErrorDetail = constant.ErrDetailCode3
		dto.Status = constant.TxStatusFailed
		return dto, err
	}

	transaction.Status = constant.TxStatusSuccess
	transaction.ErrorCode = constant.ErrCodeNoErr
	transaction.ErrorDetail = constant.ErrDetailCode0

	_ = copier.Copy(dto, transaction)

	return dto, nil
}

func NewCheckStatusService(TxRepo repository.TransactionRepository, logger *zap.Logger) CheckStatusService {
	return &CheckStatusServiceImpl{
		TxRepo: TxRepo,
		logger: logger,
	}
}