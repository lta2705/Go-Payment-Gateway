package functionality

import "github.com/lta2705/Go-Payment-Gateway/internal/worker"

type ProducerService interface {
	ProduceMessage(msg string) error
}

type ProducerServiceImpl struct {
	producer worker.KafkaProducerWorker
}

func (ps *ProducerServiceImpl) ProduceMessage(msg string) error {
	return ps.producer.SendMessage(msg)
}

func NewProduceService(producer worker.KafkaProducerWorker) ProducerService {
	return &ProducerServiceImpl{
		producer: producer,
	}
}
