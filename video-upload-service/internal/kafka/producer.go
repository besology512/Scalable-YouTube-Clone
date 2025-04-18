package kafka

import (
	"fmt"
	"log"
	"video-upload-service/internal/config"

	"github.com/IBM/sarama"
)

func NewProducer() (sarama.SyncProducer, error) {
	conf := config.Load()

	// Create a new Sarama configuration
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 5
	saramaConfig.Producer.Return.Successes = true

	// Create a new Sarama producer
	producer, err := sarama.NewSyncProducer([]string{conf.KafkaBrokers}, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return producer, nil
}

func ProduceMessage(message string) error {
	producer, err := NewProducer()
	if err != nil {
		log.Printf("Failed to create producer: %v", err)
		return err
	}
	defer producer.Close()

	conf := config.Load()

	// Create a new Kafka message
	msg := &sarama.ProducerMessage{
		Topic: conf.KafkaTopic,
		Key:   sarama.StringEncoder("video-upload"),
		Value: sarama.StringEncoder(message),
	}

	// Send the message
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}

	fmt.Printf("Message produced successfully to Kafka topic %s (partition: %d, offset: %d): %s\n", conf.KafkaTopic, partition, offset, message)
	return nil
}
