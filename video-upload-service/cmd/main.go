package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"video-upload-service/internal/handler"
)

func main() {
	// Initialize router
	router := mux.NewRouter()

	// Set up the routes
	router.HandleFunc("/upload", handler.UploadVideoHandler).Methods("POST")

	// Start server
	port := "8080"
	fmt.Printf("Server started on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
