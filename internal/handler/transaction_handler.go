package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lta2705/Go-Payment-Gateway/internal/constant"
	"github.com/lta2705/Go-Payment-Gateway/internal/dto"
	"github.com/lta2705/Go-Payment-Gateway/internal/middleware"
	"github.com/lta2705/Go-Payment-Gateway/internal/service"
	"go.uber.org/zap"
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
	case constant.TxTypeSale:
		logger.Info("Processing SALE transaction")
		go s.transactionService.CreateSaleTransaction(&transactionDto)
	case constant.TxTypeVoid:
		logger.Info("Processing VOID transaction")
		go s.transactionService.CreateVoidTransaction(&transactionDto)
	case constant.TxTypeQRRefund:
		logger.Info("Processing REFUND transaction")
		go s.transactionService.CreateRefundTransaction(&transactionDto)
	case constant.TxTypeQR:
		logger.Info("Processing QR transaction")
		go s.transactionService.CreateQRTransaction(&transactionDto)
	case constant.TxTypeCheckStatus:
		g	
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
