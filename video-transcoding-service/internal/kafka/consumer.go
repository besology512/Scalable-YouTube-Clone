package kafka

import (
	"github.com/segmentio/kafka-go"
	"log"
	"video-transcoding-service/internal/config"
	"video-transcoding-service/internal/minio"
)

type Consumer struct {
	Broker  string
	GroupID string
}

func NewConsumer(cfg config.KafkaConfig) *Consumer {
	return &Consumer{
		Broker:  cfg.Brokers,
		GroupID: cfg.GroupID,
	}
}

func (c *Consumer) Start(minioClient *minio.Client, transcodeFunc func(string, string) error) {
	// Create Kafka reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{c.Broker},
		GroupID: c.GroupID,
		Topic:   "video-uploaded",
	})

	for {
		msg, err := reader.ReadMessage(nil)
		if err != nil {
			log.Println("Error reading message: ", err)
			continue
		}

		// Process video
		videoFile := string(msg.Value)
		transcodedFile := "/tmp/transcoded_" + videoFile

		err = transcodeFunc(videoFile, transcodedFile)
		if err != nil {
			log.Println("Error transcoding video: ", err)
			continue
		}

		// Upload transcoded file to MinIO
		err = minioClient.UploadFile(transcodedFile)
		if err != nil {
			log.Println("Error uploading to MinIO: ", err)
			continue
		}

		// Produce completion message to Kafka
		// ...
	}
}
