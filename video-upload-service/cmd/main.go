package main

import (
	"fmt"
	"log"
	"net/http"
	"video-upload-service/internal/handler"
	"video-upload-service/internal/kafka"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize router
	router := mux.NewRouter()

	// start kafka consumer
	go kafka.ConsumeMessages()
	// Set up the routes
	router.HandleFunc("/upload", handler.UploadVideoHandler).Methods("POST")

	// Start server
	port := "8080"
	fmt.Printf("Server started on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
