package app

import (
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

type App struct {
	Router   *gin.Engine
	Producer *kafka.Writer
	Consumer *kafka.Reader
}

func NewApp(r *gin.Engine, p *kafka.Writer, c *kafka.Reader) *App {
	return &App{
		Router:   r,
		Producer: p,
		Consumer: c,
	}
}
