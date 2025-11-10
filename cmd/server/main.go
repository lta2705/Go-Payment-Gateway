package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lta2705/Go-Payment-Gateway/internal/middleware"
)

func main() {
	logger := middleware.SetupLogger()

	r := gin.Default()

	r.POST("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/transaction", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "transaction received",
		})
		logger.Info("Transaction endpoint hit")
	})

	r.Run(":8084")
}
