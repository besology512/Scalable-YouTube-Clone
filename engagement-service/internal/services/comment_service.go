package services

import (
	"engagement-service/internal/models"
	"engagement-service/internal/repository"
	"errors"
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

func (s *CommentService) UpdateComment(userID, commentID, newContent string) error {
	comment, err := s.repo.GetByID(commentID)
	if err != nil {
		return err
	}
	if comment.UserID != userID {
		return errors.New("unauthorized")
	}
	return s.repo.UpdateContent(commentID, newContent)
}

func (s *CommentService) DeleteComment(userID, commentID string) error {
	comment, err := s.repo.GetByID(commentID)
	if err != nil {
		return err
	}
	if comment.UserID != userID {
		return errors.New("unauthorized")
	}
	return s.repo.Delete(commentID)
}
