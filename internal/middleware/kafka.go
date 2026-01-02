package middleware

import (
	"github.com/lta2705/Go-Payment-Gateway/pkg/config"
	"github.com/lta2705/Go-Payment-Gateway/utils"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func CreateKafkaProducer(cfg *config.KafkaProducerConfig) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(cfg.BootstrapServers...),
		Topic:    cfg.ProducerTopic,
		Balancer: &kafka.Hash{},

		// Reliability
		RequiredAcks: kafka.RequiredAcks(utils.ParseAcks(cfg.Acks)),
		MaxAttempts:  10,

		// Timeout
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,

		// Batching
		BatchSize:    100,
		BatchTimeout: 10 * time.Millisecond,

		AllowAutoTopicCreation: true,
	}
}

func CreateKafkaConsumer(cfg *config.KafkaConsumerConfig) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.BootstrapServers,
		GroupID: cfg.ConsumerGroupID,
		Topic:   cfg.ConsumerTopic,

		// Fetch behavior
		MinBytes: 1e3,
		MaxBytes: 10e6,
		MaxWait:  1 * time.Second,

		// Consumer group stability
		SessionTimeout:   30 * time.Second,
		RebalanceTimeout: 60 * time.Second,

		// Offset
		StartOffset: kafka.FirstOffset,

		// Commit
		CommitInterval: 0,

		// Isolation
		IsolationLevel: kafka.ReadCommitted,

		Logger:      nil,
		ErrorLogger: kafka.LoggerFunc(log.Printf),
	})
}
