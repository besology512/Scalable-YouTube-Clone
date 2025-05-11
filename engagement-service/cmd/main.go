// @title Scalable YouTube Clone API
// @version 1.0
// @description Graduation project - Video Streaming Platform with Microservices Architecture

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token for authentication
package main

import (
	"engagement-service/internal/db"
	"engagement-service/internal/handlers"
	"engagement-service/internal/models"
	"engagement-service/internal/repository"
	"engagement-service/internal/services"

	_ "engagement-service/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	db.Init()
	db.DB.AutoMigrate(&models.Comment{}, &models.Reaction{})
	//db.DB.Migrator().DropTable(&models.Comment{}, &models.Reaction{})

	commentRepo := repository.NewGormCommentRepository()
	reactionRepo := repository.NewGormReactionRepository()

	commentService := services.NewCommentService(commentRepo, "http://host.docker.internal:8083")

	reactionService := services.NewReactionService(reactionRepo, "http://host.docker.internal:8083")

	commentHandler := handlers.NewCommentHandler(commentService)
	reactionHandler := handlers.NewReactionHandler(reactionService)

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/videos/:id/like", reactionHandler.HandleLike)
	r.POST("/videos/:id/dislike", reactionHandler.HandleDislike)

	r.POST("/videos/:id/comments", commentHandler.PostComment)
	r.GET("/videos/:id/comments", commentHandler.GetComments)
	r.PUT("/videos/:id/comments/:commentId", commentHandler.UpdateComment)
	r.DELETE("/videos/:id/comments/:commentId", commentHandler.DeleteComment)

	r.Run(":8087")
}
