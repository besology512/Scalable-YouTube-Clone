package tests

import (
	"testing"
	"video-recommender/config"
	"video-recommender/service"
)

func TestRecommendVideos(t *testing.T) {
	config.ConnectDB()

	videos, err := service.RecommendVideos("user1")
	if err != nil {
		t.Fatal("Failed to get recommendations:", err)
	}

	if len(videos) == 0 {
		t.Error("Expected some recommendations, got 0")
	}

	for _, v := range videos {
		t.Logf("Recommended: %s", v.Title)
	}
}
