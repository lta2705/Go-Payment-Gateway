package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lta2705/Go-Payment-Gateway/internal/handler"
	"go.uber.org/zap"
)

var logger *zap.Logger

func NewRouter(txHandler *handler.TransactionHandlerImpl) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/transactions", txHandler.CreateTransaction)
	}

	return r
}
