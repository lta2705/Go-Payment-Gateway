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
	cardService 	  service.CardService
	qrService	   service.QRService
	voidService    service.VoidService
	refundService service.RefundService
	checkStatusService service.CheckStatusService
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
		go s.cardService.CreateCardTransaction(&transactionDto)
	case constant.TxTypeVoid:
		s.logger.Info("Processing VOID transaction")
		go s.voidService.CreateVoidTransaction(&transactionDto)
	case constant.TxTypeQRRefund:
		s.logger.Info("Processing REFUND transaction")
		go s.refundService.CreateRefundTransaction(&transactionDto)
	case constant.TxTypeQR:
		s.logger.Info("Processing QR transaction")
		go s.qrService.CreateQRTransaction(&transactionDto)
	case constant.TxTypeCheckStatus:
		s.logger.Info("Check Transaction Status")
		go s.checkStatusService.CheckTransactionStatus(&transactionDto)
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

func NewTransactionHandler(cardSvc service.CardService, qrSvc service.QRService, voidSvc service.VoidService, refundSvc service.RefundService, checkStatusSvc service.CheckStatusService, logger *zap.Logger) *TransactionHandlerImpl {
	return &TransactionHandlerImpl{
		cardService:        cardSvc,
		qrService:          qrSvc,
		voidService:        voidSvc,
		refundService:      refundSvc,
		checkStatusService: checkStatusSvc,
		logger:             logger,
	}
}
