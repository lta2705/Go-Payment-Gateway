package worker

import (
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumerWorker interface {
	ConsumeMessage(handler func(msg kafka.Message) error) error
}

type KafkaConsumerWorkerImpl struct {
	Reader *kafka.Reader
}

func (cw *KafkaConsumerWorkerImpl) ConsumeMessage(handler func(msg kafka.Message) error) error {
	for {
		msg, err := cw.Reader.FetchMessage(context.Background())
		if err != nil {
			logger.Error("fetch error", err)
			continue
		}
		logger.Info("Consumer message from Kafka", string(msg.Value))
		// Logic Processing
		if err := handler(msg); err != nil {
			logger.Error("handler error", err)
			continue
		}

		//commit after processing
		if err := cw.Reader.CommitMessages(context.Background(), msg); err != nil {
			logger.Error("commit error", err)
		}
	}
}

func NewConsumerWorker(reader *kafka.Reader) KafkaConsumerWorker {
	return &KafkaConsumerWorkerImpl{
		Reader: reader,
	}
}
