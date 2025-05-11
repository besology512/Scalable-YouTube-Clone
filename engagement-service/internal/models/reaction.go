package models

import "gorm.io/gorm"

type Reaction struct {
	gorm.Model
	VideoID string `gorm:"index;not null"`
	UserID  string `gorm:"index;not null"`
	Type    string `gorm:"not null"` // "like" or "dislike"
}
