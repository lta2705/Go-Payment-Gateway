package functionality

import (
	"context"
	"github.com/lta2705/Go-Payment-Gateway/internal/handler"
	"github.com/lta2705/Go-Payment-Gateway/internal/worker"
	"github.com/segmentio/kafka-go"
)

type ConsumerService interface {
	ReadTransaction(ctx context.Context)
}

type ConsumerServiceImpl struct {
	Consumer *worker.KafkaConsumerWorker
	Handler  handler.TransactionRespHandler
}

func (cs *ConsumerServiceImpl) ReadTransaction(ctx context.Context) {
	go func() {
		cs.Consumer.ConsumeMessage(func(msg kafka.Message) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			return cs.Handler.HandleTransaction(msg.Value)
		})
	}()
}

func NewConsumerService(consumer *worker.KafkaConsumerWorker, handler handler.TransactionRespHandler) ConsumerService {
	return &ConsumerServiceImpl{
		Consumer: consumer,
		Handler:  handler,
	}
}
