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

func (h *CommentHandler) GetComments(c *gin.Context) {
	videoID := c.Param("id")
	comments, err := h.service.GetComments(videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comments)
}

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
