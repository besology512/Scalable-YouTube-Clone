package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"video-upload-service/internal/kafka"
	"video-upload-service/internal/model"
	"video-upload-service/internal/storage"
)

func UploadVideoHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Parse video file from the request
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Parse metadata from form
	title := r.FormValue("title")
	creator := r.FormValue("creator")

	// Upload the video to MinIO
	videoURL, err := storage.UploadToMinIO(file, header.Filename)
	if err != nil {
		http.Error(w, "Failed to upload video", http.StatusInternalServerError)
		return
	}

	// Log video upload success
	log.Printf("Video uploaded successfully to MinIO: %s", videoURL)

	// Create metadata object
	meta := model.VideoMetadata{
		VideoURL:   videoURL,
		Title:      title,
		Creator:    creator,
		UploadedAt: time.Now(),
	}

	// Convert metadata to JSON
	jsonData, err := json.Marshal(meta)
	if err != nil {
		http.Error(w, "Failed to serialize metadata", http.StatusInternalServerError)
		return
	}
	// Send message to Kafka to notify that the video is uploaded and ready for editing the metadata
	err = kafka.ProduceMessage(string(jsonData))
	if err != nil {
		http.Error(w, "Failed to notify Kafka", http.StatusInternalServerError)
		return
	}
	log.Printf("Meta Data Sent via kafka : %s", string(jsonData))

	// Send success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Video uploaded successfully with Metadata"))
}
