package worker

import (
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/segmentio/kafka-go"
)

type KafkaProducerWorker struct {
	Writer *kafka.Writer
}

func (pw *KafkaProducerWorker) SendMessage(value string) error {
	logger.Info("Received message from terminal", value)
	err := pw.Writer.WriteMessages(context.Background(),
		kafka.Message{
			Value:  []byte(value),
			Offset: kafka.SeekCurrent,
		},
	)

	logger.Info("Message successfully to Kafka", value, "with offset", kafka.Message{Offset: kafka.SeekCurrent})
	return err
}

func NewProducerWorker(writer *kafka.Writer) *KafkaProducerWorker {
	return &KafkaProducerWorker{
		Writer: writer,
	}
}
