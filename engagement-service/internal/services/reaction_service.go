package services

import (
	"errors"
	"fmt"
	"net/http"

	"engagement-service/internal/repository"
)

type ReactionService struct {
	repo                repository.ReactionRepository
	streamingServiceURL string
}

func NewReactionService(repo repository.ReactionRepository, streamingServiceURL string) *ReactionService {
	return &ReactionService{
		repo:                repo,
		streamingServiceURL: streamingServiceURL,
	}
}

func (s *ReactionService) ToggleReaction(videoID, userID, reactionType string) error {
	err := s.checkVideoExists(videoID)
	if err != nil {
		return err
	}
	return s.repo.ToggleReaction(videoID, userID, reactionType)
}

func (s *ReactionService) CountReactions(videoID string) (int64, int64, error) {
	return s.repo.CountReactions(videoID)
}

func (s *ReactionService) checkVideoExists(videoID string) error {
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
