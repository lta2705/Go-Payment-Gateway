package service

import (
	"github.com/lta2705/Go-Payment-Gateway/internal/dto"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(dto *dto.TransactionDTO) error // Create a new transaction record

}

type TransactionRepositoryImpl struct {
	db *gorm.DB
}

func (tri *TransactionRepositoryImpl) CreateTransaction(dto *dto.TransactionDTO) error {
	return tri.db.Create(dto).Error
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &TransactionRepositoryImpl{db: db}
}
