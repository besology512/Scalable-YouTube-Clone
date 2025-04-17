package kafka

import (
	"fmt"
	"log"
	"video-upload-service/internal/config"

	"github.com/segmentio/kafka-go"
)

func NewConsumer() *kafka.Reader {
	conf := config.Load()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{conf.KafkaBrokers},
		Topic:   conf.KafkaTopic,
		GroupID: "video-upload-group",
	})

	return reader
}

func ConsumeMessages() {
	reader := NewConsumer()

	for {
		message, err := reader.ReadMessage(nil)
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			continue
		}

		// Process the video upload message (transcoding logic)
		fmt.Printf("Received message: %s\n", string(message.Value))
	}
}
