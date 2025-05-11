package services

import (
	"engagement-service/internal/models"
	"engagement-service/internal/repository"
)

type CommentService struct {
	repo repository.CommentRepository
}

func NewCommentService(repo repository.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) PostComment(videoID, userID, content string) (*models.Comment, error) {
	comment := &models.Comment{
		VideoID: videoID,
		UserID:  userID,
		Content: content,
	}
	err := s.repo.Create(comment)
	return comment, err
}

func (s *CommentService) GetComments(videoID string) ([]models.Comment, error) {
	return s.repo.GetByVideoID(videoID)
}

func (s *CommentService) UpdateComment(commentID, newContent string) error {
	return s.repo.UpdateContent(commentID, newContent)
}

func (s *CommentService) DeleteComment(commentID string) error {
	return s.repo.Delete(commentID)
}
