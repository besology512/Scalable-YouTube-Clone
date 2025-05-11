package repository

import (
	"context"
	"video-recommender/config"
	"video-recommender/models"

	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

// GetVideosByIDs retrieves videos by their IDs.
func GetVideosByIDs(ids []string) ([]models.Video, error) {
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := config.DB.Collection("videos").Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find videos by IDs: %w", err)
	}
	defer func() {
		if cerr := cursor.Close(context.Background()); cerr != nil {
			log.Printf("error closing cursor: %v", cerr)
		}
	}()
	var videos []models.Video
	if err := cursor.All(context.Background(), &videos); err != nil {
		return nil, fmt.Errorf("failed to decode videos: %w", err)
	}
	return videos, nil
}

// GetVideosByTags retrieves videos that match any of the specified tags.
func GetVideosByTags(tags []string) ([]models.Video, error) {
	if len(tags) == 0 {
		return []models.Video{}, nil // or return an error if appropriate
	}

	filter := bson.M{"tags": bson.M{"$in": tags}}
	cursor, err := config.DB.Collection("videos").Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find videos by tags: %w", err)
	}
	defer func() {
		if cerr := cursor.Close(context.Background()); cerr != nil {
			log.Printf("error closing cursor: %v", cerr)
		}
	}()
	var videos []models.Video
	if err := cursor.All(context.Background(), &videos); err != nil {
		return nil, fmt.Errorf("failed to decode videos: %w", err)
	}
	return videos, nil
}

// GetVideosByCategories retrieves videos that match any of the specified categories.
func GetVideosByCategories(categories []string) ([]models.Video, error) {
	filter := bson.M{"categories": bson.M{"$in": categories}}
	cursor, err := config.DB.Collection("videos").Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find videos by categories: %w", err)
	}
	defer func() {
		if cerr := cursor.Close(context.Background()); cerr != nil {
			log.Printf("error closing cursor: %v", cerr)
		}
	}()
	var videos []models.Video
	if err := cursor.All(context.Background(), &videos); err != nil {
		return nil, fmt.Errorf("failed to decode videos: %w", err)
	}
	return videos, nil
}

// GetVideosByDurationRange retrieves videos within a specific duration range.
func GetVideosByDurationRange(minDuration, maxDuration int) ([]models.Video, error) {
	filter := bson.M{"duration": bson.M{"$gte": minDuration, "$lte": maxDuration}}
	cursor, err := config.DB.Collection("videos").Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find videos by duration range: %w", err)
	}
	defer func() {
		if cerr := cursor.Close(context.Background()); cerr != nil {
			log.Printf("error closing cursor: %v", cerr)
		}
	}()
	var videos []models.Video
	if err := cursor.All(context.Background(), &videos); err != nil {
		return nil, fmt.Errorf("failed to decode videos: %w", err)
	}
	return videos, nil
}

// GetVideosByRating retrieves videos with a rating greater than or equal to the specified value.
func GetVideosByRating(minRating float64) ([]models.Video, error) {
	filter := bson.M{"rating": bson.M{"$gte": minRating}}
	cursor, err := config.DB.Collection("videos").Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find videos by rating: %w", err)
	}
	defer func() {
		if cerr := cursor.Close(context.Background()); cerr != nil {
			log.Printf("error closing cursor: %v", cerr)
		}
	}()
	var videos []models.Video
	if err := cursor.All(context.Background(), &videos); err != nil {
		return nil, fmt.Errorf("failed to decode videos: %w", err)
	}
	return videos, nil
}

// GetVideosByUploadDate retrieves videos uploaded on or after a specific date.
func GetVideosByUploadDate(startDate string) ([]models.Video, error) {
	filter := bson.M{"uploadDate": bson.M{"$gte": startDate}}
	cursor, err := config.DB.Collection("videos").Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find videos by upload date: %w", err)
	}
	defer func() {
		if cerr := cursor.Close(context.Background()); cerr != nil {
			log.Printf("error closing cursor: %v", cerr)
		}
	}()
	var videos []models.Video
	if err := cursor.All(context.Background(), &videos); err != nil {
		return nil, fmt.Errorf("failed to decode videos: %w", err)
	}
	return videos, nil
}

// GetVideosByWatchTime retrieves videos with a watch time greater than or equal to the specified value.
func GetVideosByWatchTime(minWatchTime int) ([]models.Video, error) {
	filter := bson.M{"watchTime": bson.M{"$gte": minWatchTime}}
	cursor, err := config.DB.Collection("videos").Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find videos by watch time: %w", err)
	}
	defer func() {
		if cerr := cursor.Close(context.Background()); cerr != nil {
			log.Printf("error closing cursor: %v", cerr)
		}
	}()
	var videos []models.Video
	if err := cursor.All(context.Background(), &videos); err != nil {
		return nil, fmt.Errorf("failed to decode videos: %w", err)
	}
	return videos, nil
}

// GetVideosByBoostedTags retrieves videos with boosted scores for specific tags.
func GetVideosByBoostedTags(tag string, minBoost int) ([]models.Video, error) {
	filter := bson.M{"tagsBoost." + tag: bson.M{"$gte": minBoost}}
	cursor, err := config.DB.Collection("videos").Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find videos by boosted tags: %w", err)
	}
	defer func() {
		if cerr := cursor.Close(context.Background()); cerr != nil {
			log.Printf("error closing cursor: %v", cerr)
		}
	}()
	var videos []models.Video
	if err := cursor.All(context.Background(), &videos); err != nil {
		return nil, fmt.Errorf("failed to decode videos: %w", err)
	}
	return videos, nil
}

// GetVideosByBoostedCategories retrieves videos with boosted scores for specific categories.
func GetVideosByBoostedCategories(category string, minBoost int) ([]models.Video, error) {
	filter := bson.M{"categoriesBoost." + category: bson.M{"$gte": minBoost}}
	cursor, err := config.DB.Collection("videos").Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find videos by boosted categories: %w", err)
	}
	defer func() {
		if cerr := cursor.Close(context.Background()); cerr != nil {
			log.Printf("error closing cursor: %v", cerr)
		}
	}()
	var videos []models.Video
	if err := cursor.All(context.Background(), &videos); err != nil {
		return nil, fmt.Errorf("failed to decode videos: %w", err)
	}
	return videos, nil
}

// GetVideosByChannelName retrieves videos uploaded by a specific channel name.
func GetVideosByChannelName(channelName string) ([]models.Video, error) {
	filter := bson.M{"channelName": channelName}
	cursor, err := config.DB.Collection("videos").Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find videos by channel name: %w", err)
	}
	defer func() {
		if cerr := cursor.Close(context.Background()); cerr != nil {
			log.Printf("error closing cursor: %v", cerr)
		}
	}()
	var videos []models.Video
	if err := cursor.All(context.Background(), &videos); err != nil {
		return nil, fmt.Errorf("failed to decode videos: %w", err)
	}
	return videos, nil
}

// FindVideosByTagsAndCategories retrieves videos that match any of the specified tags and categories.
func FindVideosByTagsAndCategories(tags, categories []string, watchHistory []string) ([]models.Video, error) {
	filter := bson.M{
		"$and": []bson.M{
			{"tags": bson.M{"$in": tags}},
			{"categories": bson.M{"$in": categories}},
			{"_id": bson.M{"$nin": watchHistory}}, // Exclude videos in the user's watch history
		},
	}
	cursor, err := config.DB.Collection("videos").Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find videos by tags, categories, and watch history: %w", err)
	}
	defer func() {
		if cerr := cursor.Close(context.Background()); cerr != nil {
			log.Printf("error closing cursor: %v", cerr)
		}
	}()
	var videos []models.Video
	if err := cursor.All(context.Background(), &videos); err != nil {
		return nil, fmt.Errorf("failed to decode videos: %w", err)
	}
	return videos, nil
}
