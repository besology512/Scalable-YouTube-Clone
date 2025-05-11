// @title Scalable YouTube Clone API
// @version 1.0
// @description Graduation project - Video Streaming Platform with Microservices Architecture

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token for authentication
package main

import (
	"fmt"
	"log"
	"net/http"

	"streaming-service/internal/clients"
	"streaming-service/internal/config"
	"streaming-service/internal/handlers"

	_ "streaming-service/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
)

func main() {
	// 1) Load config (from .env or env vars)
	cfg := config.Load()

	// 2) Initialize MinIO client
	minioClient := clients.NewMinioClient(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
	)

	// 3) Set up router
	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	router.HandleFunc("/stream/{name}", handlers.StreamHandler(minioClient, cfg.MinioBucket)).Methods("GET")
	router.HandleFunc("/videos/{id}/exists", handlers.VideoExistsHandler(minioClient, "prosessed-videos"))

	// 4) Start HTTP server
	fmt.Printf("Streaming service running on port %s\n", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, router))
}
