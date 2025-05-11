package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"video-recommender/service"
)

func GetRecommendationsHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	videos, err := service.RecommendVideos(userID)
	if err != nil {
		http.Error(w, "Failed to get recommendations: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}
