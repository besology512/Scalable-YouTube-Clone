package handler

import (
	"fmt"
	"log"
	"net/http"
	"video-upload-service/internal/kafka"
	"video-upload-service/internal/storage"
)

func UploadVideoHandler(w http.ResponseWriter, r *http.Request) {
	// Parse video file from the request
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Upload the video to MinIO
	videoURL, err := storage.UploadToMinIO(file)
	if err != nil {
		http.Error(w, "Failed to upload video", http.StatusInternalServerError)
		return
	}

	// Log video upload success
	log.Printf("Video uploaded successfully to MinIO: %s", videoURL)

	// Send message to Kafka to notify that the video is uploaded and ready for transcoding
	message := fmt.Sprintf("Video uploaded: %s", videoURL)
	err = kafka.ProduceMessage(message)
	if err != nil {
		http.Error(w, "Failed to notify Kafka", http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Video uploaded successfully"))
}
