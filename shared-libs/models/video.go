package models

import "time"

type Video struct {
	ID          string         `bson:"_id"` // use UUID
	Title       string         `bson:"title"`
	Description string         `bson:"description"`
	UserID      string         `bson:"user_id"`
	Filename    string         `bson:"filename"`
	Status      string         `bson:"status"` // "processing", "ready"
	UploadedAt  time.Time      `bson:"uploaded_at"`
	Variants    []VideoVariant `bson:"variants"`
}

type VideoVariant struct {
	Resolution string `bson:"resolution"`
	URL        string `bson:"url"`
}
