package repository

import (
	"engagement-service/internal/db"
	"engagement-service/internal/models"

	"gorm.io/gorm"
)

type GormCommentRepository struct {
	db *gorm.DB
}

func NewGormCommentRepository() *GormCommentRepository {
	return &GormCommentRepository{
		db: db.DB,
	}
}

func (r *GormCommentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *GormCommentRepository) GetByID(commentID string) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.First(&comment, "id = ?", commentID).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *GormCommentRepository) GetByVideoID(videoID string) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.Where("video_id = ?", videoID).Find(&comments).Error
	return comments, err
}

func (r *GormCommentRepository) UpdateContent(commentID string, newContent string) error {
	return r.db.Model(&models.Comment{}).Where("id = ?", commentID).Update("content", newContent).Error
}

func (r *GormCommentRepository) Delete(commentID string) error {
	return r.db.Delete(&models.Comment{}, commentID).Error
}
