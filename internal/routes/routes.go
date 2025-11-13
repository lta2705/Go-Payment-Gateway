package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lta2705/Go-Payment-Gateway/internal/handler"
	"github.com/lta2705/Go-Payment-Gateway/internal/service"
	"go.uber.org/zap"
)

var logger *zap.Logger

func SetupRoutes(txServ service.TransactionService, merchantCredServ service.MerchantCredentialsService) {
	r := gin.Default()

	txHandler := handler.NewTransactionHandler(txServ)
	merchantCredHandler := handler.NewMerchantCredentialsHandler(merchantCredServ)

	api := r.Group("")
	{
		api.POST("/transactions", txHandler.CreateTransaction)
	}
}
