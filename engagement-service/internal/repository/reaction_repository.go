package repository

type ReactionRepository interface {
	ToggleReaction(videoID, userID, reactionType string) error
	CountReactions(videoID string) (likes int64, dislikes int64, err error)
}
