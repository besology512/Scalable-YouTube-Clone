package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"video-upload-service/internal/config"
	"video-upload-service/internal/db"
	"video-upload-service/internal/model"

	"github.com/IBM/sarama"
)

func NewConsumer() (sarama.ConsumerGroup, error) {
	conf := config.Load()

	// Create a new Sarama configuration
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	// Create a new consumer group
	consumerGroup, err := sarama.NewConsumerGroup([]string{conf.KafkaBrokers}, "video-upload-group", saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return consumerGroup, nil
}

type ConsumerHandler struct{}

func (handler *ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (handler *ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (handler *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		// Parse the message
		var videoMetadata model.VideoMetadata
		err := json.Unmarshal(message.Value, &videoMetadata)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// Insert video metadata into MongoDB
		err = db.InsertVideoMetadata(videoMetadata)
		if err != nil {
			log.Printf("Failed to insert metadata to MongoDB: %v", err)
			continue
		}

		// Mark the message as processed
		session.MarkMessage(message, "")

		// Log the success
		fmt.Printf("Successfully processed and saved video metadata: %s\n", videoMetadata.VideoURL)
		err = NotifyTranscodingService(videoMetadata)
		if err != nil {
			log.Printf("Failed to notify transcoding service: %v", err)
			continue
		}
	}
	return nil
}

func ConsumeMessages() {
	// Initialize MongoDB
	_, err := db.InitializeMongoDB()
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
		return
	}

	consumerGroup, err := NewConsumer()
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer group: %v", err)
		return
	}
	defer consumerGroup.Close()

	handler := &ConsumerHandler{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
		<-sigchan
		cancel()
	}()

	for {
		if err := consumerGroup.Consume(ctx, []string{config.Load().KafkaTopic}, handler); err != nil {
			log.Printf("Error consuming messages: %v", err)
		}

		// Check if context was canceled
		if ctx.Err() != nil {
			return
		}
	}
}
func NotifyTranscodingService(videoMetadata model.VideoMetadata) error {
	conf := config.Load()

	// Create the message payload
	messagePayload := map[string]string{
		"video_url": videoMetadata.VideoURL,
	}
	messageBytes, err := json.Marshal(messagePayload)
	if err != nil {
		return fmt.Errorf("failed to marshal message payload: %w", err)
	}

	// Produce the message to the Kafka transcode topic
	err = ProduceMessageToTopic(conf.KafkaTranscodeTopic, string(messageBytes))
	if err != nil {
		return fmt.Errorf("failed to produce message to Kafka: %w", err)
	}

	fmt.Printf("Successfully notified transcoding service for video: %s\n", videoMetadata.VideoURL)
	return nil
}

func ProduceMessageToTopic(topic, message string) error {
	producer, err := NewProducer()
	if err != nil {
		log.Printf("Failed to create producer: %v", err)
		return err
	}
	defer producer.Close()

	// Create a new Kafka message
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder("video-transcode"),
		Value: sarama.StringEncoder(message),
	}

	// Send the message
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}

	fmt.Printf("Message produced successfully to Kafka topic %s (partition: %d, offset: %d): %s\n", topic, partition, offset, message)
	return nil
}
