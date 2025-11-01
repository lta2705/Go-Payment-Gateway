package main

import (
	"fmt"
	"gin-demo/dto"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	fmt.Println("Hello, World!")

	r := gin.Default()

	tx := &dto.TransactionDTO{}
	tx.SetTxType("QR")
	tx.SetPcPosId("id12345")
	tx.SetPcPosTxnId("txn67890")

	ctx := gin.Context{}
	ctx.BindHeader(``)

	r.POST("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/transaction", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"txType":     tx.GetTxType(),
			"pcPosId":    tx.GetPcPosId(),
			"pcPosTxnId": tx.GetPcPosTxnId(),
		})
	})

	log.Println("Starting server on port 8084")
	r.Run(":8084")
}
