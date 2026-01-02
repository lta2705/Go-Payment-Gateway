package handler

import (
	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
	"github.com/lta2705/Go-Payment-Gateway/internal/constant"
	"github.com/lta2705/Go-Payment-Gateway/internal/dto"
	"github.com/lta2705/Go-Payment-Gateway/internal/service"
	"net/http"
)

type TransactionReqHandler interface {
	CreateTransaction(gin *gin.Context)
}
type TransactionReqHandlerImpl struct {
	cardService        service.CardService
	qrService          service.QRService
	voidService        service.VoidService
	refundService      service.RefundService
	checkStatusService service.CheckStatusService
	credentialHandler  service.MerchantCredentialsService
}

func (s *TransactionReqHandlerImpl) CreateTransaction(c *gin.Context) {

	apiKey := c.GetHeader("X-API-KEY")
	_, authenErr := s.credentialHandler.Authenticate(apiKey)
	if authenErr != nil {
		if authenErr.Error() == "invalid api key" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
		} else {
			c.JSON(500, gin.H{"error": "Internal Server Error"})
		}
		return
	}

	transactionDto := &dto.TransactionDTO{}

	// Parse JSON
	if err := c.ShouldBindJSON(transactionDto); err != nil {
		logger.Error("Error parsing request", err)
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
		logger.Info("Processing SALE transaction")
		transaction, err = s.cardService.CreateCardTransaction(transactionDto)

	case constant.TxTypeVoid:
		logger.Info("Processing VOID transaction")
		logger.Info("original Transaction ID with OrgPcPosTxnId", transactionDto.OrgPcPosTxnId)
		if transactionDto.OrgPcPosTxnId == "" {
			logger.Error("Original Transaction ID is required for VOID transactions")
			transactionDto.Status = constant.TxStatusFailed
			transactionDto.ErrorCode = constant.ErrCodeNotFoundOriginTx
			transactionDto.ErrorDetail = constant.ErrDetailCode7
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Original Transaction ID is required for VOID transactions",
				"data":  transactionDto,
			})
			return
		}
		transaction, err = s.voidService.CreateVoidTransaction(transactionDto)

	case constant.TxTypeQRRefund:
		logger.Info("Processing REFUND transaction")
		transaction, err = s.refundService.CreateRefundTransaction(transactionDto)

	case constant.TxTypeQR:
		logger.Info("Processing QR transaction")
		transaction, err = s.qrService.CreateQRTransaction(transactionDto)

	case constant.TxTypeCheckStatus:
		logger.Info("Check Transaction Status")
		transaction, err = s.checkStatusService.CheckTransactionStatus(transactionDto)

	default:
		logger.Error("Unsupported transaction type wiht TransactionType", transactionDto.TransactionType)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported transaction type"})
		return
	}

	if err != nil {
		logger.Error("Transaction failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Return value from service
	c.JSON(http.StatusOK, gin.H{
		"data": transaction,
	})
}

func NewTransactionHandler(cardSvc service.CardService, qrSvc service.QRService,
	voidSvc service.VoidService, refundSvc service.RefundService,
	checkStatusSvc service.CheckStatusService, credentialHandler service.MerchantCredentialsService) TransactionReqHandler {
	return &TransactionReqHandlerImpl{
		cardService:        cardSvc,
		qrService:          qrSvc,
		voidService:        voidSvc,
		refundService:      refundSvc,
		checkStatusService: checkStatusSvc,
		credentialHandler:  credentialHandler,
	}
}
