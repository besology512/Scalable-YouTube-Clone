package main

import (
	"fmt"
	"log"
	"net/http"

	"streaming-service/internal/clients"
	"streaming-service/internal/config"
	"streaming-service/internal/handlers"

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
	// Health check endpoint:
	// Stream endpoint:
	// GET /stream/{name}
	// e.g. /stream/video123.mp4
	router.HandleFunc("/stream/{name}", handlers.StreamHandler(minioClient, cfg.MinioBucket)).Methods("GET")
	router.HandleFunc("/videos/{id}/exists", handlers.VideoExistsHandler(minioClient, cfg.MinioBucket)).Methods("GET")

	// 4) Start HTTP server
	fmt.Printf("Streaming service running on port %s\n", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, router))
}
