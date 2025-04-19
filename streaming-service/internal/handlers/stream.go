package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
)

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
