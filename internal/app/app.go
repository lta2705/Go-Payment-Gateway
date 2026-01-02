package app

import (
	"github.com/gin-gonic/gin"
	"github.com/lta2705/Go-Payment-Gateway/internal/functionality"
)

type App struct {
	Router   *gin.Engine
	Producer functionality.ProducerService
	Consumer functionality.ConsumerService
}

func NewApp(r *gin.Engine, producer functionality.ProducerService, consumer functionality.ConsumerService) *App {
	return &App{
		Router:   r,
		Producer: producer,
		Consumer: consumer,
	}
}
