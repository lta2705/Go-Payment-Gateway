package repository

import (
	"github.com/lta2705/Go-Payment-Gateway/internal/model"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(model *model.Transaction) error                                               // Create a new transaction record
	UpdateTransaction(model *model.Transaction) error                                               // Update an existing transaction record
	GetDB() *gorm.DB                                                                                // Get the underlying gorm DB instance
	FindByTransactionId(id *string) (*model.Transaction, error)                                     // Find a transaction by its ID
	FindByPcPosIdAndTransactionId(transactionId string, pcPosId string) (*model.Transaction, error) // Find a transaction by PcPosId and TransactionId
}

type TransactionRepositoryImpl struct {
	db *gorm.DB
}

func (tri *TransactionRepositoryImpl) CreateTransaction(model *model.Transaction) error {
	return tri.db.Create(model).Error
}

func (tri *TransactionRepositoryImpl) UpdateTransaction(model *model.Transaction) error {
	return tri.db.Save(model).Error
}

func (tri *TransactionRepositoryImpl) GetDB() *gorm.DB {
	return tri.db
}

func (tri *TransactionRepositoryImpl) FindByTransactionId(id *string) (*model.Transaction, error) {
	var transaction model.Transaction
	err := tri.db.Where("transaction_id = ?", id).First(&transaction).Error
	return &transaction, err
}

func (tri *TransactionRepositoryImpl) FindByPcPosIdAndTransactionId(transactionId string, pcPosId string) (*model.Transaction, error) {
	var transaction model.Transaction
	err := tri.db.Where("pc_pos_id = ? AND transaction_id = ?", pcPosId, transactionId).First(&transaction).Error
	return &transaction, err
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &TransactionRepositoryImpl{db: db}
}
