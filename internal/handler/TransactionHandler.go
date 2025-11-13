package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lta2705/Go-Payment-Gateway/internal/dto"
	"github.com/lta2705/Go-Payment-Gateway/internal/service"
)

type TransactionHandler interface {
	CreateTransaction()
}
type TransactionHandlerImpl struct {
	transactionService service.TransactionService
}

func (h *TransactionHandlerImpl) CreateTransaction(c *gin.Context) {
	var transactionDTO *dto.TransactionDTO

	err := c.ShouldBindJSON(&transactionDTO)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
	}

}

func NewTransactionHandler(transactionServ service.TransactionService) *TransactionHandlerImpl {
	return &TransactionHandlerImpl{
		transactionService: transactionServ,
	}
}
