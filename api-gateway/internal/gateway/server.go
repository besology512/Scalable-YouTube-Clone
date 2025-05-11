// internal/gateway/server.go
package gateway

import (
	"github.com/besology512/api-gateway/internal/config"
	"github.com/besology512/api-gateway/internal/gateway/clients"
	"github.com/gin-gonic/gin"
)

func APIGateWayServer(config *config.Config) *gin.Engine {

	//data, err := os.ReadFile(config.AuthPublicKeyPath)

	/*if err != nil {
		log.Fatalf("unable to read auth public key: %v", err)
	}*/

	//publicKey, err := jwt.ParseRSAPublicKeyFromPEM(data)

	// if err != nil {
	// 	log.Fatalf("Invalid Public Key: %v", err)
	// }

	// redisClient := redis.NewClient(&redis.Options{
	// 	Addr:     config.RedisServerAdress,
	// 	Password: config.RedisPassword,
	// 	DB:       config.RedisDBindex,
	// })

	// // rateLimiter, err := middleware.NewRateLimitMiddleware(
	// // 	redisClient,
	// // 	"token-bucket",
	// // 	config.MaxRequests,
	// // 	config.RateLimitTTL.String(),
	// // )

	// if err != nil {
	// 	return nil
	// }

	// concurrencyLimiter := middleware.NewConcurrencyMiddleware(
	// 	redisClient,
	// 	config.ConcurrencyLimit,
	// 	config.RateLimitTTL,
	// )

	authClient := clients.NewAuthClient(config.AuthServiceUrl)
	funcClient := clients.NewFunctionClient(config.MicroServiceUrl)

	program := gin.New()
	program.Use(gin.Logger(), gin.Recovery())
	// program.Use(
	// 	gin.Logger(),
	// 	gin.Recovery(),
	// 	rateLimiter,
	// 	concurrencyLimiter,
	// 	//middleware.AuthMiddleware(publicKey),
	// )

	SetRoutes(program, authClient, funcClient)
	return program

}
