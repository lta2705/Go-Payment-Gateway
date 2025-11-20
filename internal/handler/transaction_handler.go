package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lta2705/Go-Payment-Gateway/internal/constant"
	"github.com/lta2705/Go-Payment-Gateway/internal/dto"
	"github.com/lta2705/Go-Payment-Gateway/internal/service"
	"go.uber.org/zap"
)

type TransactionHandler interface {
	CreateTransaction()
}
type TransactionHandlerImpl struct {
	transactionService  service.TransactionService
	merchantCredService service.MerchantCredentialsService
	logger              *zap.Logger
}

func (s *TransactionHandlerImpl) CreateTransaction(c *gin.Context) {

	var transactionDto dto.TransactionDTO

	err := c.ShouldBindJSON(&transactionDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})

		s.logger.Error("Error parsing request", zap.Error(err))
		return
	}

	switch transactionDto.TransactionType {
	case constant.TxTypeSale:
		s.logger.Info("Processing SALE transaction")
		go s.transactionService.CreateSaleTransaction(&transactionDto)
	case constant.TxTypeVoid:
		s.logger.Info("Processing VOID transaction")
		go s.transactionService.CreateVoidTransaction(&transactionDto)
	case constant.TxTypeQRRefund:
		s.logger.Info("Processing REFUND transaction")
		go s.transactionService.CreateRefundTransaction(&transactionDto)
	case constant.TxTypeQR:
		s.logger.Info("Processing QR transaction")
		go s.transactionService.CreateQRTransaction(&transactionDto)
	case constant.TxTypeCheckStatus:
		s.logger.Info("Check Transaction Status")
		go s.transactionService.CheckTransactionStatus(&transactionDto)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unsupported transaction type"})
		s.logger.Error("Unsupported transaction type", zap.String("TransactionType", transactionDto.TransactionType))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction created successfully",
		"data":    transactionDto,
	})
}

func NewTransactionHandler(transactionServ service.TransactionService, logger *zap.Logger) *TransactionHandlerImpl {
	return &TransactionHandlerImpl{
		transactionService: transactionServ,
		logger:             logger,
	}
}
