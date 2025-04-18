package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"path/filepath"
	"strings"
	"video-upload-service/internal/config"
)

func UploadToMinIO(file io.Reader, originalFilename string) (string, error) {
	conf := config.Load()

	// Initialize MinIO client
	minioClient, err := minio.New(conf.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.MinioAccessKey, conf.MinioSecretKey, ""),
		Secure: false, // If you're using http (not https), set this to false
	})
	if err != nil {
		log.Printf("Failed to initialize MinIO client: %v", err)
		return "", err
	}

	// Upload the file to MinIO
	ext := filepath.Ext(originalFilename) // get extension like ".mp4"
	uniqueName := uuid.New().String() + ext
	objectName := uniqueName
	bucketName := conf.MinioBucket
	ext_lowered := strings.ToLower(filepath.Ext(originalFilename))
	var contentType string
	switch ext_lowered {
	case ".mp4":
		contentType = "video/mp4"
	case ".mov":
		contentType = "video/quicktime"
	case ".avi":
		contentType = "video/x-msvideo"
	case ".mkv":
		contentType = "video/x-matroska"
	default:
		contentType = "application/octet-stream" // fallback
	}

	_, err = minioClient.PutObject(context.Background(), bucketName, objectName, file, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Printf("Failed to upload video to MinIO: %v", err)
		return "", err
	}

	// Return the URL or path of the uploaded video
	videoURL := "http://" + conf.MinioEndpoint + "/" + bucketName + "/" + objectName
	return videoURL, nil
}
