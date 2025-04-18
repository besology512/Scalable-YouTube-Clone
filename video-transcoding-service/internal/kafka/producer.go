package kafka

import (
	"log"

	"video-transcoding-service/internal/config"

	"github.com/IBM/sarama"
)

type Producer struct {
	syncProducer sarama.SyncProducer
}

func NewProducer(cfg *config.Config) (*Producer, error) {
	// Configure Sarama
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Retry.Max = 5
	producerConfig.Producer.Return.Successes = true

	// Create a new SyncProducer
	producer, err := sarama.NewSyncProducer(cfg.KafkaBrokers, producerConfig)
	if err != nil {
		return nil, err
	}

	return &Producer{syncProducer: producer}, nil
}

func (p *Producer) ProduceMessage(topic string, objectName string) error {
	// Create a new Kafka message
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(objectName),
	}

	// Send the message
	partition, offset, err := p.syncProducer.SendMessage(message)
	if err != nil {
		return err
	}

	log.Printf("Message sent to topic %s [partition: %d, offset: %d]\n ObjectName: %s \n", topic, partition, offset, objectName)
	return nil
}

func (p *Producer) Close() error {
	return p.syncProducer.Close()
}
