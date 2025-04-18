package main

import (
	"fmt"
	"log"

	"github.com/besology512/api-gateway/internal/config"
	"github.com/besology512/api-gateway/internal/gateway"
	"github.com/besology512/api-gateway/internal/logger"
)

func main() {
	config := config.Load()
	router := gateway.APIGateWayServer(config)
	addr := fmt.Sprintf(":%s", config.Port)
	logger.Log.Infof("starting API Gateway on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to start gateway: %v", err)

	}
}
