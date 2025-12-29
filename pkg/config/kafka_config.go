package config

import (
	"github.com/joho/godotenv"
	"github.com/lta2705/Go-Payment-Gateway/utils"
	"log"
)

type KafkaProducerConfig struct {
	BootstrapServers  []string
	ProducerTopic     string
	Acks              string
	Retries           int
	MaxInFlight       int
	EnableIdempotence bool
}

type KafkaConsumerConfig struct {
	BootstrapServers  []string
	ConsumerGroupID   string
	ConsumerTopic     string
	AutoOffsetReset   string
	EnableAutoCommit  bool
	IsolationLevel    string
	HeartbeatInterval int
}

func LoadKafkaProducerConfig() *KafkaProducerConfig {
	// Load .env file
	loadConfig()

	return &KafkaProducerConfig{
		BootstrapServers:  []string{utils.String("KAFKA_BROKER", "localhost:9092")},
		ProducerTopic:     utils.String("KAFKA_PRODUCER_TOPIC", "producer_topic"),
		Acks:              utils.String("KAFKA_PRODUCER_ACKS", "all"),
		Retries:           utils.Int("KAFKA_PRODUCER_RETRIES", 5),
		MaxInFlight:       utils.Int("KAFKA_PRODUCER_MAX_IN_FLIGHT", 5),
		EnableIdempotence: utils.Bool("KAFKA_PRODUCER_ENABLE_IDEMPOTENCE", true),
	}
}

func LoadKafkaConsumerConfig() *KafkaConsumerConfig {
	// Load .env file
	loadConfig()
	return &KafkaConsumerConfig{
		BootstrapServers:  []string{utils.String("KAFKA_BROKER", "localhost:9092")},
		ConsumerGroupID:   utils.String("KAFKA_CONSUMER_GROUP_ID", "consumer_group"),
		ConsumerTopic:     utils.String("KAFKA_CONSUMER_TOPIC", "consumer_topic"),
		AutoOffsetReset:   utils.String("KAFKA_CONSUMER_AUTO_OFFSET_RESET", "earliest"),
		EnableAutoCommit:  utils.Bool("KAFKA_CONSUMER_ENABLE_AUTO_COMMIT", true),
		IsolationLevel:    utils.String("KAFKA_CONSUMER_ISOLATION_LEVEL", "read_committed"),
		HeartbeatInterval: utils.Int("KAFKA_CONSUMER_HEARTBEAT_INTERVAL_MS", 20000),
	}
}

func loadConfig() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Println(".env file not found, using environment variables")
	}

}
