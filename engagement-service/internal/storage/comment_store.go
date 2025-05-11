package storage

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        string
	VideoID   string
	UserID    string
	Content   string
	Timestamp time.Time
}

var (
	commentMu   sync.Mutex
	commentsMap = make(map[string][]Comment) // videoID -> []comments
)

func AddComment(videoID, userID, content string) Comment {
	commentMu.Lock()
	defer commentMu.Unlock()

	comment := Comment{
		ID:        uuid.New().String(),
		VideoID:   videoID,
		UserID:    userID,
		Content:   content,
		Timestamp: time.Now(),
	}

	commentsMap[videoID] = append(commentsMap[videoID], comment)
	return comment
}

func GetComments(videoID string) []Comment {
	commentMu.Lock()
	defer commentMu.Unlock()

	return commentsMap[videoID]
}
func DeleteComment(videoID, commentID string) bool {
	commentMu.Lock()
	defer commentMu.Unlock()

	comments, exists := commentsMap[videoID]
	if !exists {
		return false
	}

	for i, comment := range comments {
		if comment.ID == commentID {
			commentsMap[videoID] = append(comments[:i], comments[i+1:]...)
			return true
		}
	}
	return false
}
func UpdateComment(videoID, commentID, newContent string) bool {
	commentMu.Lock()
	defer commentMu.Unlock()

	comments, exists := commentsMap[videoID]
	if !exists {
		return false
	}

	for i, comment := range comments {
		if comment.ID == commentID {
			comments[i].Content = newContent
			return true
		}
	}
	return false
}
