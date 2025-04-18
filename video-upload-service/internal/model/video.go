package model

import "time"

type VideoMetadata struct {
	VideoURL   string    `bson:"video_url"`
	Title      string    `bson:"title"`
	Creator    string    `bson:"creator"`
	UploadedAt time.Time `bson:"uploaded_at"`
}
