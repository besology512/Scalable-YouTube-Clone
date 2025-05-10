package models

import "time"

type Channel struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Subscribers int       `json:"subscribers"`
	// Add other fields
}

type Subscription struct {
	UserID    string    `json:"user_id"`
	ChannelID string    `json:"channel_id"`
	CreatedAt time.Time `json:"created_at"`
}
