package services

import (
	"engagement-service/internal/repository"
)

type ReactionService struct {
	repo repository.ReactionRepository
}

func NewReactionService(repo repository.ReactionRepository) *ReactionService {
	return &ReactionService{repo: repo}
}

func (s *ReactionService) ToggleReaction(videoID, userID, reactionType string) error {
	return s.repo.ToggleReaction(videoID, userID, reactionType)
}

func (s *ReactionService) CountReactions(videoID string) (int64, int64, error) {
	return s.repo.CountReactions(videoID)
}
