package kafka

import (
	"github.com/segmentio/kafka-go"
	"log"
)

func ProduceTranscodingComplete(broker string, message string) error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   "transcoding-complete",
	})

	err := writer.WriteMessages(nil, kafka.Message{
		Value: []byte(message),
	})
	if err != nil {
		log.Println("Failed to write message to Kafka: ", err)
		return err
	}

	return nil
}
