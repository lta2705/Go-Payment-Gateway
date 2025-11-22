package handler

import (
	"github.com/lta2705/Go-Payment-Gateway/internal/constant"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lta2705/Go-Payment-Gateway/internal/dto"
	"github.com/lta2705/Go-Payment-Gateway/internal/service"
	"go.uber.org/zap"
)

type TransactionHandler interface {
	CreateTransaction()
}
type TransactionHandlerImpl struct {
	cardService        service.CardService
	qrService          service.QRService
	voidService        service.VoidService
	refundService      service.RefundService
	checkStatusService service.CheckStatusService
	logger             *zap.Logger
}

func (s *TransactionHandlerImpl) CreateTransaction(c *gin.Context) {

	transactionDto := &dto.TransactionDTO{}

	// Parse JSON
	if err := c.ShouldBindJSON(transactionDto); err != nil {
		s.logger.Error("Error parsing request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var (
		transaction interface{}
		err         error
	)

	// Dispatch based on transactionType
	switch transactionDto.TransactionType {
	case constant.TxTypeSale:
		s.logger.Info("Processing SALE transaction")
		transaction, err = s.cardService.CreateCardTransaction(transactionDto)

	case constant.TxTypeVoid:
		s.logger.Info("Processing VOID transaction")
		s.logger.Info("original Transaction ID", zap.String("OrgPcPosTxnId", transactionDto.OrgPcPosTxnId))
		if transactionDto.OrgPcPosTxnId == "" {
			s.logger.Error("Original Transaction ID is required for VOID transactions")
			transactionDto.Status = constant.TxStatusFailed
			transactionDto.ErrorCode = "17"
			transactionDto.ErrorDetail = "Original Transaction ID not found"
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Original Transaction ID is required for VOID transactions",
				"data":  transactionDto,
			})
			return
		}
		transaction, err = s.voidService.CreateVoidTransaction(transactionDto)

	case constant.TxTypeQRRefund:
		s.logger.Info("Processing REFUND transaction")
		transaction, err = s.refundService.CreateRefundTransaction(transactionDto)

	case constant.TxTypeQR:
		s.logger.Info("Processing QR transaction")
		transaction, err = s.qrService.CreateQRTransaction(transactionDto)

	case constant.TxTypeCheckStatus:
		s.logger.Info("Check Transaction Status")
		transaction, err = s.checkStatusService.CheckTransactionStatus(transactionDto)

	default:
		s.logger.Error("Unsupported transaction type", zap.String("TransactionType", transactionDto.TransactionType))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported transaction type"})
		return
	}

	if err != nil {
		s.logger.Error("Transaction failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Return value from service
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    transaction,
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
