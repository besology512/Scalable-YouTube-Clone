package gateway

import (
	"github.com/besology512/api-gateway/internal/gateway/clients"
	"github.com/gin-gonic/gin"
)

func SetRoutes(program *gin.Engine, authclient *clients.AuthClient, fnClient *clients.FunctionClient) {

	public := program.Group("/auth")
	{
		public.POST("/register", authclient.LogOut)
		public.POST("/login", authclient.Login)
		public.POST("/refresh", authclient.Refresh)
	}
}
