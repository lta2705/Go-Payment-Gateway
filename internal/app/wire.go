//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/lta2705/Go-Payment-Gateway/internal/functionality"
	"github.com/lta2705/Go-Payment-Gateway/internal/handler"
	"github.com/lta2705/Go-Payment-Gateway/internal/middleware"
	"github.com/lta2705/Go-Payment-Gateway/internal/repository"
	"github.com/lta2705/Go-Payment-Gateway/internal/routes"
	"github.com/lta2705/Go-Payment-Gateway/internal/service"
	"github.com/lta2705/Go-Payment-Gateway/internal/worker"
	"github.com/lta2705/Go-Payment-Gateway/pkg/config"
)

var repositorySet = wire.NewSet(
	repository.NewTransactionRepository,
	repository.NewMerchantCredentialsRepository,
)

var serviceSet = wire.NewSet(
	service.NewCardService,
	service.NewQRService,
	service.NewVoidService,
	service.NewRefundService,
	service.NewCheckStatusService,
	service.NewPollingService,
	service.NewMerchantCredentialsService,
)

var ProducerSet = wire.NewSet(
	config.LoadKafkaProducerConfig,
	middleware.CreateKafkaProducer,
	worker.NewProducerWorker,
	functionality.NewProduceService,
)

var ConsumerSet = wire.NewSet(
	config.LoadKafkaConsumerConfig,
	middleware.CreateKafkaConsumer,
	worker.NewConsumerWorker,
	functionality.NewConsumerService,
)

var handlerSet = wire.NewSet(
	handler.NewTransactionHandler,
)

var databaseSet = wire.NewSet(
	config.LoadDBConfig,
	middleware.SetupDatabase)

func InitializeApp() (*App, error) {
	wire.Build(
		databaseSet,
		ProducerSet,
		ConsumerSet,
		repositorySet,
		serviceSet,
		handlerSet,
		routes.NewRouter,
		NewApp,
	)
	return &App{}, nil
}
