package main

import (
	"engagement-service/internal/db"
	"engagement-service/internal/handlers"
	"engagement-service/internal/models"
	"engagement-service/internal/repository"
	"engagement-service/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()
	db.DB.AutoMigrate(&models.Comment{}, &models.Reaction{})

	commentRepo := repository.NewGormCommentRepository()
	reactionRepo := repository.NewGormReactionRepository()

	commentService := services.NewCommentService(commentRepo)
	reactionService := services.NewReactionService(reactionRepo)

	commentHandler := handlers.NewCommentHandler(commentService)
	reactionHandler := handlers.NewReactionHandler(reactionService)

	r := gin.Default()

	r.POST("/videos/:id/like", reactionHandler.HandleLike)
	r.POST("/videos/:id/dislike", reactionHandler.HandleDislike)

	r.POST("/videos/:id/comments", commentHandler.PostComment)
	r.GET("/videos/:id/comments", commentHandler.GetComments)
	r.PUT("/videos/:id/comments/:commentId", commentHandler.UpdateComment)
	r.DELETE("/videos/:id/comments/:commentId", commentHandler.DeleteComment)

	r.Run(":8087")
}
