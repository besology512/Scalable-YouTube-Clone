package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reaction struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	VideoID   string `gorm:"index;not null"`
	UserID    string `gorm:"index;not null"`
	Type      string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" swaggerignore:"true"`
}

func (r *Reaction) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New().String()
	return
}
