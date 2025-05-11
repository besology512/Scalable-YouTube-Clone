package handlers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
)

func VideoExistsHandler(minioClient *minio.Client, bucket string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		videoID := mux.Vars(r)["id"]
		if videoID == "" {
			http.Error(w, "Missing video ID", http.StatusBadRequest)
			return
		}

		_, err := minioClient.StatObject(context.Background(), bucket, videoID, minio.StatObjectOptions{})
		if err != nil {
			http.Error(w, "Video not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
