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

// @Summary Like a video
// @Description Toggle like reaction for a video. Creates, deletes or updates like state for a user.
// @Tags Reactions
// @Accept json
// @Produce json
// @Param id path string true "Video ID"
// @Param X-User-ID header string true "User ID from Auth"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /videos/{id}/like [post]
// @Security ApiKeyAuth
// @Security BearerAuth
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

// HandleDislike godoc
// @Summary Toggle dislike reaction
// @Description Dislike or remove dislike on a video
// @Tags Reactions
// @Accept json
// @Produce json
// @Param id path string true "Video ID"
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /videos/{id}/dislike [post]
// @Security ApiKeyAuth
// @Security BearerAuth
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
