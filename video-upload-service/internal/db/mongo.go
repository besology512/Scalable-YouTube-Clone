package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"video-upload-service/internal/config"
	"video-upload-service/internal/model"
)

var mongoClient *mongo.Client

// InitializeMongoDB initializes the MongoDB client
func InitializeMongoDB() (*mongo.Client, error) {
	conf := config.Load()

	clientOptions := options.Client().ApplyURI(conf.MongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}

	// Ping the database to ensure connection is established
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
		return nil, err
	}

	mongoClient = client
	return mongoClient, nil
}

// InsertVideoMetadata inserts the video metadata into MongoDB
func InsertVideoMetadata(metadata model.VideoMetadata) error {
	collection := mongoClient.Database("video_upload_db").Collection("metadata")

	_, err := collection.InsertOne(context.Background(), metadata)
	if err != nil {
		log.Printf("Failed to insert video metadata: %v", err)
		return err
	}

	log.Printf("Successfully inserted metadata for video: %s", metadata.VideoURL)
	return nil
}
