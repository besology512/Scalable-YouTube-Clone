package handlers

import (
	"engagement-service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	service *services.CommentService
}

func NewCommentHandler(service *services.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

type CommentRequest struct {
	Content string `json:"content"`
}

type UpdateRequest struct {
	Content string `json:"content"`
}

// PostComment godoc
// @Summary Add comment
// @Description Add a comment to a video if it exists
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path string true "Video ID"
// @Param X-User-ID header string true "User ID"
// @Param comment body CommentRequest true "Comment content"
// @Success 200 {object} models.Comment
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /videos/{id}/comments [post]
// @Security ApiKeyAuth
// @Security BearerAuth
func (h *CommentHandler) PostComment(c *gin.Context) {
	videoID := c.Param("id")
	userID := c.GetHeader("X-User-ID")

	var req CommentRequest
	if userID == "" || c.BindJSON(&req) != nil || req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment"})
		return
	}

	comment, err := h.service.PostComment(videoID, userID, req.Content)
	if err != nil {
		if err.Error() == "video does not exist" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Video does not exist"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, comment)
}

// GetComments godoc
// @Summary Get comments
// @Description Get all comments for a video
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path string true "Video ID"
// @Success 200 {array} models.Comment
// @Failure 500 {object} map[string]string
// @Router /videos/{id}/comments [get]
// @Security ApiKeyAuth
// @Security BearerAuth
func (h *CommentHandler) GetComments(c *gin.Context) {
	videoID := c.Param("id")
	comments, err := h.service.GetComments(videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comments)
}

// UpdateComment godoc
// @Summary Update comment
// @Description Update a user's own comment
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path string true "Video ID"
// @Param commentId path string true "Comment ID"
// @Param X-User-ID header string true "User ID"
// @Param content body UpdateRequest true "Updated comment content"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /videos/{id}/comments/{commentId} [put]
// @Security ApiKeyAuth
// @Security BearerAuth
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	commentID := c.Param("commentId")
	userID := c.GetHeader("X-User-ID")

	var req UpdateRequest
	if userID == "" || c.BindJSON(&req) != nil || req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content"})
		return
	}

	err := h.service.UpdateComment(userID, commentID, req.Content)
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment updated"})
}

// DeleteComment godoc
// @Summary Delete comment
// @Description Delete a user's own comment
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path string true "Video ID"
// @Param commentId path string true "Comment ID"
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /videos/{id}/comments/{commentId} [delete]
// @Security ApiKeyAuth
// @Security BearerAuth
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentID := c.Param("commentId")
	userID := c.GetHeader("X-User-ID")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID"})
		return
	}

	err := h.service.DeleteComment(userID, commentID)
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}
