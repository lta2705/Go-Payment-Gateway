package worker

import (
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/segmentio/kafka-go"
)

type KafkaProducerWorker interface {
	SendMessage(value string) error
}
type KafkaProducerWorkerImpl struct {
	Writer *kafka.Writer
}

func (pw *KafkaProducerWorkerImpl) SendMessage(value string) error {
	logger.Info("Received message from terminal", value)
	msg := kafka.Message{
		Value: []byte(value),
	}

	err := pw.Writer.WriteMessages(context.Background(), msg)

	logger.Info("Message successfully to Kafka", value, "with offset", msg.Offset)
	return err
}

func NewProducerWorker(writer *kafka.Writer) KafkaProducerWorker {
	return &KafkaProducerWorkerImpl{
		Writer: writer,
	}
}
