package services

import (
	"engagement-service/internal/models"
	"engagement-service/internal/repository"
	"errors"
	"fmt"
	"net/http"
)

type CommentService struct {
	repo                repository.CommentRepository
	streamingServiceURL string
}

func NewCommentService(repo repository.CommentRepository, streamingServiceURL string) *CommentService {
	return &CommentService{
		repo:                repo,
		streamingServiceURL: streamingServiceURL,
	}
}

func (s *CommentService) PostComment(videoID, userID, content string) (*models.Comment, error) {
	if err := s.checkVideoExists(videoID); err != nil {
		return nil, err
	}

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

func (s *CommentService) checkVideoExists(videoID string) error {
	url := fmt.Sprintf("%s/videos/%s/exists", s.streamingServiceURL, videoID)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to verify video existence: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return errors.New("video does not exist")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response from video service: %d", resp.StatusCode)
	}

	return nil
}
