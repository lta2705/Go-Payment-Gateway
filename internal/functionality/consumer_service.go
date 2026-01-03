package functionality

import (
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/lta2705/Go-Payment-Gateway/internal/handler"
	"github.com/lta2705/Go-Payment-Gateway/internal/worker"
	"github.com/segmentio/kafka-go"
)

type ConsumerService interface {
	ReadTransaction(ctx context.Context)
}

type ConsumerServiceImpl struct {
	Consumer worker.KafkaConsumerWorker
	Handler  handler.TransactionRespHandler
}

func (cs *ConsumerServiceImpl) ReadTransaction(ctx context.Context) {
	logger.Info("Consumer service started and waiting for messages...")

	// Bỏ go func() ở đây. Để hàm ConsumeMessage giữ Goroutine này lại.
	err := cs.Consumer.ConsumeMessage(func(msg kafka.Message) error {
		logger.Info("Message received!")
		return cs.Handler.HandleTransaction(msg.Value)
	})

	if err != nil {
		logger.Error("Consumer stopped!", err)
	}
}

func NewConsumerService(consumer worker.KafkaConsumerWorker, handler handler.TransactionRespHandler) ConsumerService {
	return &ConsumerServiceImpl{
		Consumer: consumer,
		Handler:  handler,
	}
}
