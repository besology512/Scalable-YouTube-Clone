package storage

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"video-upload-service/internal/config"
)

func UploadToMinIO(file io.Reader) (string, error) {
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
	objectName := "uploaded_video.mp4" // Can be dynamic based on input
	bucketName := conf.MinioBucket
	contentType := "video/mp4" // Set content type based on video format

	_, err = minioClient.PutObject(context.Background(), bucketName, objectName, file, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Printf("Failed to upload video to MinIO: %v", err)
		return "", err
	}

	// Return the URL or path of the uploaded video
	videoURL := "http://" + conf.MinioEndpoint + "/" + bucketName + "/" + objectName
	return videoURL, nil
}
