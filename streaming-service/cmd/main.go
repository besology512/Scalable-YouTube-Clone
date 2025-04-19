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
	cfg := config.Load()

	minioClient := clients.NewMinioClient(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
	)

	router := mux.NewRouter()
	streamHandler := handlers.NewStreamHandler(minioClient, cfg.MinioBucket)

	router.Handle("/stream", streamHandler).Methods("GET")

	fmt.Printf("Streaming service running on port %s\n", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, router))
}
