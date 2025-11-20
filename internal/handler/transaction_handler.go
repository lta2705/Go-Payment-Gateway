package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lta2705/Go-Payment-Gateway/internal/dto"
	"github.com/lta2705/Go-Payment-Gateway/internal/middleware"
	"github.com/lta2705/Go-Payment-Gateway/internal/service"
	"go.uber.org/zap"
	"net/http"
)

type TransactionHandler interface {
	CreateTransaction()
}
type TransactionHandlerImpl struct {
	transactionService  service.TransactionService
	merchantCredService service.MerchantCredentialsService
}

func (s *TransactionHandlerImpl) CreateTransaction(c *gin.Context) {
	logger := middleware.CreateLogger()
	defer logger.Sync()

	var transactionDto dto.TransactionDTO

	err := c.ShouldBindJSON(&transactionDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})

		logger.Error("Error parsing request", zap.Error(err))
		return
	}

	switch transactionDto.TransactionType {
	case "SALE":
		logger.Info("Processing SALE transaction")
		s.transactionService.CreateSaleTransaction(&transactionDto)
	case "VOID":
		logger.Info("Processing VOID transaction")
		s.transactionService.CreateVoidTransaction(&transactionDto)
	case "REFUND":
		logger.Info("Processing REFUND transaction")
		s.transactionService.CreateRefundTransaction(&transactionDto)
	case "QR":
		logger.Info("Processing QR transaction")
		s.transactionService.CreateQRTransaction(&transactionDto)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unsupported transaction type"})
		logger.Error("Unsupported transaction type", zap.String("TransactionType", transactionDto.TransactionType))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction created successfully",
		"data":    transactionDto,
	})
}

func NewTransactionHandler(transactionServ service.TransactionService) *TransactionHandlerImpl {
	return &TransactionHandlerImpl{
		transactionService: transactionServ,
	}
}
