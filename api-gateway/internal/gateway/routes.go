package gateway

import (
	"github.com/besology512/api-gateway/internal/gateway/clients"
	"github.com/besology512/api-gateway/internal/gateway/handlers"
	"github.com/gin-gonic/gin"
)

func SetRoutes(program *gin.Engine, authclient *clients.AuthClient, fnClient *clients.FunctionClient) {

	public := program.Group("/auth")
	{
		public.GET("/login", authclient.Login)
		public.POST("/refresh", authclient.Refresh)
		public.POST("/logout", authclient.LogOut)
		public.GET("/health", authclient.Health)
	}

	protected := program.Group("/")
	{

		functionServiceHandler := handlers.NewFunctionHandler(fnClient)

		protected.Any("/functions/*path", functionServiceHandler.Proxy)
		protected.Any("/jobs/*path", functionServiceHandler.Proxy)
	}

}
