package handlers

import (
	"engagement-service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReactionHandler struct {
	service *services.ReactionService
}

func NewReactionHandler(service *services.ReactionService) *ReactionHandler {
	return &ReactionHandler{service: service}
}

func (h *ReactionHandler) HandleLike(c *gin.Context) {
	videoID := c.Param("id")
	userID := c.GetHeader("X-User-ID")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID"})
		return
	}

	err := h.service.ToggleReaction(videoID, userID, "like")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	likes, dislikes, err := h.service.CountReactions(videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Liked!",
		"video_id": videoID,
		"likes":    likes,
		"dislikes": dislikes,
	})
}

func (h *ReactionHandler) HandleDislike(c *gin.Context) {
	videoID := c.Param("id")
	userID := c.GetHeader("X-User-ID")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID"})
		return
	}

	err := h.service.ToggleReaction(videoID, userID, "dislike")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	likes, dislikes, err := h.service.CountReactions(videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Disliked!",
		"video_id": videoID,
		"likes":    likes,
		"dislikes": dislikes,
	})
}
