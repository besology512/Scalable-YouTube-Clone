package router

import (
	"search-service/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/search", handler.Search)
	return r
}
