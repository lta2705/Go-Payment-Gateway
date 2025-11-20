package repository

import (
	"errors"
	"github.com/lta2705/Go-Payment-Gateway/internal/model"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(tx *model.Transaction) error
	UpdateTransaction(tx *model.Transaction) error
	GetDB() *gorm.DB
	FindByTransactionId(transactionId string) (*model.Transaction, error)
	FindByPcPosIdAndTransactionId(transactionId, pcPosId string) (*model.Transaction, error)
}

type TransactionRepositoryImpl struct {
	db *gorm.DB
}

func (tri *TransactionRepositoryImpl) CreateTransaction(tx *model.Transaction) error {
	return tri.db.Create(tx).Error
}

func (tri *TransactionRepositoryImpl) UpdateTransaction(tx *model.Transaction) error {
	return tri.db.Save(tx).Error
}

func (tri *TransactionRepositoryImpl) GetDB() *gorm.DB {
	return tri.db
}

func (tri *TransactionRepositoryImpl) FindByTransactionId(transactionId string) (*model.Transaction, error) {
	var tx model.Transaction
	err := tri.db.Where(
		"transaction_id = ?",
		transactionId,
	).First(&tx).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // correct
	}

	if err != nil {
		return nil, err
	}

	return &tx, nil
}

func (tri *TransactionRepositoryImpl) FindByPcPosIdAndTransactionId(pcPosId, transactionId string) (*model.Transaction, error) {
	var tx model.Transaction
	err := tri.db.Where(
		"pc_pos_id = ? AND transaction_id = ?",
		pcPosId, transactionId,
	).First(&tx).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // correct
	}

	if err != nil {
		return nil, err
	}

	return &tx, nil
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &TransactionRepositoryImpl{db: db}
}
