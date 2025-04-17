package kafka

import (
	"fmt"
	"log"
	"video-upload-service/internal/config"

	"github.com/segmentio/kafka-go"
)

func NewProducer() *kafka.Writer {
	conf := config.Load()

	writer := &kafka.Writer{
		Addr:     kafka.TCP(conf.KafkaBrokers),
		Topic:    conf.KafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}

	return writer
}

func ProduceMessage(message string) error {
	writer := NewProducer()

	err := writer.WriteMessages(
		nil, // context
		kafka.Message{
			Key:   []byte("video-upload"),
			Value: []byte(message),
		},
	)

	if err != nil {
		log.Printf("Failed to write messages: %v", err)
		return err
	}

	fmt.Println("Message produced successfully to Kafka topic:", message)
	return nil
}
