package main

import (
	"log"
	"video-transcoding-service/internal/config"
	"video-transcoding-service/internal/ffmpeg"
	"video-transcoding-service/internal/kafka"
	"video-transcoding-service/internal/minio"
)

func main() {
	// Load configurations
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize MinIO client
	minioClient, err := minio.NewClient(cfg.Minio)
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	// Start Kafka consumer
	kafkaConsumer := kafka.NewConsumer(cfg.Kafka)
	go kafkaConsumer.Start(minioClient, ffmpeg.Transcode)

	// Wait for Kafka messages (can be extended to more services)
	select {}
}
