package handlers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
)

// VideoExistsHandler godoc
// @Summary Check video existence
// @Description Check if a video exists in MinIO by ID
// @Tags Streaming
// @Produce plain
// @Param id path string true "Video ID"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Missing video ID"
// @Failure 404 {string} string "Video not found"
// @Security BearerAuth
// @Router /videos/{id}/exists [get]
func VideoExistsHandler(minioClient *minio.Client, bucket string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		if id == "" {
			http.Error(w, "Missing video ID", http.StatusBadRequest)
			return
		}

		_, err := minioClient.StatObject(
			context.Background(),
			bucket,
			id,
			minio.StatObjectOptions{},
		)
		if err != nil {
			http.Error(w, "Video not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
