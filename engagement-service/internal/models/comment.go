package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	VideoID string
	UserID  string
	Content string
}
