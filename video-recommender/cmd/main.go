package main

import (
	"log"
	"net/http"

	"video-recommender/api"
	"video-recommender/config"

	"github.com/go-chi/chi/v5"
)

func main() {
	config.ConnectDB()

	r := chi.NewRouter()
	r.Get("/recommendations/{userId}", api.GetRecommendationsHandler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
