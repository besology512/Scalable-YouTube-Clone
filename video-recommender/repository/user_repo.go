package repository

import (
	"context"
	"video-recommender/config"
	"video-recommender/models"

	"go.mongodb.org/mongo-driver/bson"
)

func GetUserByID(userID string) (*models.User, error) {
	var user models.User
	err := config.DB.Collection("users").FindOne(context.Background(), bson.M{"id": userID}).Decode(&user)
	return &user, err
}
