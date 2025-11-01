package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	fmt.Println("Hello, World!")

	r := gin.Default()

	ctx := gin.Context{}
	ctx.BindHeader(``)

	r.POST("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/transaction", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "transaction received",
		})
	})

	log.Println("Starting server on port 8084")
	r.Run(":8084")
}
