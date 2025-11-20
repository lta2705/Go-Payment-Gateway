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

type TransactionService interface {
	CreateSaleTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
	CreateVoidTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
	CreateRefundTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
	CreateQRTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
	CheckTransactionStatus(dto *dto.TransactionDTO) (*dto.TransactionDTO, error)
}
type TransactionServiceImpl struct {
	TxRepo         repository.TransactionRepository
	pollingService PollingService
	logger         *zap.Logger
}

func (t *TransactionServiceImpl) CreateSaleTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {
	defer t.logger.Sync()

	pcPosId := dto.PcPosId
	transactionId := dto.TransactionId

	transaction, _ := t.TxRepo.FindByPcPosIdAndTransactionId(pcPosId, transactionId)

	if transaction != nil {
		t.logger.Info("Sale transaction already exists", zap.String("PcPosId", pcPosId), zap.String("TransactionId", transactionId))

		dto.Status = "FAILED"
		dto.ErrorCode = "01"
		dto.ErrorDetail = "Transaction already exists"
		return dto, nil
	}

	t.logger.Info("Creating new sale transaction", zap.String("PcPosId", pcPosId), zap.String("TransactionId", transactionId))

	newTransaction := &model.Transaction{}

	err := copier.Copy(newTransaction, dto)
	if err != nil {
		t.logger.Error("Error copying transaction DTO to model", zap.Error(err))
		return nil, err
	}

	newTransaction.UpdatedBy = "SERVER"
	newTransaction.ID = uuid.New()

	t.logger.Info("New transaction before insert:", zap.Any("Transaction", newTransaction))

	error := t.TxRepo.CreateTransaction(newTransaction)
	if error != nil {
		t.logger.Error("Error creating new sale transaction in DB", zap.Error(error), zap.String("TransactionId", transactionId))
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeTcpServerError
		dto.ErrorDetail = constant.ErrDetailCode3

		return dto, error
	}

	t.logger.Info("Successfully created new sale transaction in DB", zap.Any("Payload", &transaction))

	updatedTransaction := t.pollingService.Poll(newTransaction, "CHANGE")

	updatedTransaction.ErrorCode = constant.ErrCodeNoErr
	updatedTransaction.ErrorDetail = constant.ErrDetailCode0
	updatedTransaction.Status = constant.TxStatusSuccess

	err = copier.Copy(dto, updatedTransaction)
	if err != nil {
		t.logger.Error("Error copying final model to DTO", zap.Error(err))
		return nil, err
	}

	return dto, nil
}

func (t TransactionServiceImpl) CreateVoidTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {
	var transaction = &model.Transaction{}

	err := copier.Copy(transaction, dto)
	if err != nil {
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeTcpServerError
		dto.ErrorDetail = constant.ErrDetailCode3
		return nil, err
	}
	transaction, err = t.TxRepo.FindByPcPosIdAndTransactionId(transaction.PcPosId, transaction.OrgPcPosTxnId)
	if err != nil {
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeTcpServerError
		dto.ErrorDetail = constant.ErrDetailCode3
		return dto, err
	}
	if transaction != nil {
		t.logger.Warn("Original transaction not found for void", zap.String("PcPosId", dto.PcPosId), zap.String("OrgPcPosTxnId", dto.OrgPcPosTxnId))
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeNotFoundOriginTx
		dto.ErrorDetail = constant.ErrDetailCode7
		return dto, nil
	} else if transaction.Status != constant.TxStatusSuccess {
		t.logger.Warn("Original transaction not successful for void", zap.String("PcPosId", dto.PcPosId), zap.String("OrgPcPosTxnId", dto.OrgPcPosTxnId))
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeTxNotSuccess
		dto.ErrorDetail = constant.ErrDetailCode13
		return dto, nil
	} else if transaction.Status == constant.TxStatusVoided {
		t.logger.Warn("Original transaction already voided", zap.String("PcPosId", dto.PcPosId), zap.String("OrgPcPosTxnId", dto.OrgPcPosTxnId))
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeTxVoided
		dto.ErrorDetail = constant.ErrDetailCode14
		return dto, nil
	}

	err = t.TxRepo.CreateTransaction(transaction)
	if err != nil {
		dto.Status = constant.TxStatusFailed
		dto.ErrorCode = constant.ErrCodeTcpServerError
		dto.ErrorDetail = constant.ErrDetailCode3
		return dto, err
	}

	updatedTransaction := t.pollingService.Poll(transaction, "VOID")

	updatedTransaction.ErrorCode = constant.ErrCodeNoErr
	updatedTransaction.ErrorDetail = constant.ErrDetailCode0
	updatedTransaction.Status = constant.TxStatusSuccess

	err = copier.Copy(dto, updatedTransaction)
	if err != nil {
		return nil, err
	}

	return dto, nil
}

func (t TransactionServiceImpl) CreateRefundTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {

	return nil, nil
}

func (t TransactionServiceImpl) CreateQRTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {
	return nil, nil
}

func (t TransactionServiceImpl) createTransaction(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {
	return nil, nil
}

func (t TransactionServiceImpl) CheckTransactionStatus(dto *dto.TransactionDTO) (*dto.TransactionDTO, error) {
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

func NewTransactionService(TxRepo repository.TransactionRepository, logger *zap.Logger) TransactionService {
	return &TransactionServiceImpl{
		TxRepo:         TxRepo,
		logger:         logger,
		pollingService: NewPollingService(TxRepo, logger),
	}
}
