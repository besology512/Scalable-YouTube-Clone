package sync

import (
	"context"
	"log"
	"search-service/indexer"
	"search-service/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SyncFromMongo(mongoURI string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("video_db").Collection("videos")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("MongoDB find error: %v", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var v model.Video
		if err := cursor.Decode(&v); err != nil {
			log.Printf("Failed to decode: %v", err)
			continue
		}

		err = indexer.IndexVideo(v)
		if err != nil {
			log.Printf("Failed to index video %s: %v", v.VideoID, err)
		} else {
			log.Printf("Indexed video: %s", v.VideoID)
		}
	}
}
