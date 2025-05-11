package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
)

// StreamHandler godoc
// @Summary Stream video
// @Description Streams a video file from MinIO by name
// @Tags Streaming
// @Produce video/mp4
// @Param name path string true "Video file name"
// @Success 200 {file} file
// @Failure 400 {string} string "Missing video name"
// @Failure 404 {string} string "Video not found"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /stream/{name} [get]
func StreamHandler(minioClient *minio.Client, bucket string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		name := mux.Vars(r)["name"]
		if name == "" {
			http.Error(w, "Missing video name", http.StatusBadRequest)
			return
		}

		info, err := minioClient.StatObject(
			context.Background(),
			bucket,
			name,
			minio.StatObjectOptions{},
		)
		if err != nil {
			log.Printf("StatObject(%s) error: %v", name, err)
			http.Error(w, "Video not found", http.StatusNotFound)
			return
		}

		obj, err := minioClient.GetObject(
			context.Background(),
			bucket,
			name,
			minio.GetObjectOptions{},
		)
		if err != nil {
			log.Printf("GetObject(%s) error: %v", name, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer obj.Close()

		w.Header().Set("Content-Type", "video/mp4")
		http.ServeContent(w, r, name, info.LastModified, obj)
	}
}
