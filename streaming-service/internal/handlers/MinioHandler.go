package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/minio/minio-go/v7"
)

type StreamHandler struct {
	MinioClient *minio.Client
	BucketName  string
}

func NewStreamHandler(client *minio.Client, bucket string) *StreamHandler {
	return &StreamHandler{
		MinioClient: client,
		BucketName:  bucket,
	}
}

func (h *StreamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	videoName := r.URL.Query().Get("name")
	if videoName == "" {
		http.Error(w, "Missing video name in query", http.StatusBadRequest)
		return
	}

	log.Printf("Requested video: %s", videoName)

	object, err := h.MinioClient.GetObject(
		context.Background(),
		h.BucketName,
		videoName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		log.Printf("MinIO GetObject error: %v", err)
		log.Printf("Failed to get video from MinIO: %v", err)
		http.Error(w, "Error retrieving video", http.StatusInternalServerError)
		return
	}
	defer object.Close()

	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Content-Disposition", "inline")
	http.ServeContent(w, r, videoName, time.Now(), object)
}
