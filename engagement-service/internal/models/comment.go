package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	VideoID   string
	UserID    string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()
	return
}
