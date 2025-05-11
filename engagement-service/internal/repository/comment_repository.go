package repository

import "engagement-service/internal/models"

type CommentRepository interface {
	Create(comment *models.Comment) error
	GetByVideoID(videoID string) ([]models.Comment, error)
	UpdateContent(commentID string, newContent string) error
	Delete(commentID string) error
}
