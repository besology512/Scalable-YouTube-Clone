package storage

import "sync"

var (
	mu        sync.Mutex
	reactions = make(map[string]map[string]string) // videoID -> userID -> "like"/"dislike"
)

func React(videoID, userID, reaction string) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := reactions[videoID]; !exists {
		reactions[videoID] = make(map[string]string)
	}

	currentReaction, reacted := reactions[videoID][userID]

	if reacted && currentReaction == reaction {
		delete(reactions[videoID], userID)
		return
	}

	reactions[videoID][userID] = reaction
}

func CountReactions(videoID string) (likes int, dislikes int) {
	mu.Lock()
	defer mu.Unlock()

	for _, r := range reactions[videoID] {
		if r == "like" {
			likes++
		} else if r == "dislike" {
			dislikes++
		}
	}
	return
}
