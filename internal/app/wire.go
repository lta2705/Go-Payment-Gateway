//go:build wireinject
// +build wireinject

package app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/lta2705/Go-Payment-Gateway/internal/handler"
	"github.com/lta2705/Go-Payment-Gateway/internal/middleware"
	"github.com/lta2705/Go-Payment-Gateway/internal/repository"
	"github.com/lta2705/Go-Payment-Gateway/internal/routes"
	"github.com/lta2705/Go-Payment-Gateway/internal/service"
	"github.com/lta2705/Go-Payment-Gateway/pkg/config"
)

var repositorySet = wire.NewSet(
	repository.NewTransactionRepository,
)

var serviceSet = wire.NewSet(
	service.NewCardService,
	service.NewQRService,
	service.NewVoidService,
	service.NewRefundService,
	service.NewCheckStatusService,
	service.NewPollingService,
)

var handlerSet = wire.NewSet(
	handler.NewTransactionHandler,
)

var loggerSet = wire.NewSet(middleware.NewLogger)

func InitializeApp() (*gin.Engine, error) {
	wire.Build(
		config.LoadDBConfig,
		middleware.SetupDatabase,
		loggerSet,
		repositorySet,
		serviceSet,
		handlerSet,
		routes.NewRouter,
	)
	return nil, nil
}
