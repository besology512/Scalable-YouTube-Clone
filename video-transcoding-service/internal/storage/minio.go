package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"video-transcoding-service/internal/config"
)

type MinioClient struct {
	client *minio.Client
	config *config.Config
}

// NewMinioClient initializes and returns a new MinioClient
func NewMinioClient(cfg *config.Config) (*MinioClient, error) {
	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	return &MinioClient{
		client: client,
		config: cfg,
	}, nil
}

// UploadFile uploads a file to the specified bucket
func (mc *MinioClient) UploadFile(ctx context.Context, bucketName, objectName string, fileData []byte, contentType string) error {
	// Ensure the bucket exists
	exists, err := mc.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check if bucket exists: %w", err)
	}
	if !exists {
		err = mc.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	// Upload the file
	_, err = mc.client.PutObject(ctx, bucketName, objectName, bytes.NewReader(fileData), int64(len(fileData)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	log.Printf("File %s uploaded successfully to bucket %s", objectName, bucketName)
	return nil
}

// DownloadFile downloads a file from the specified bucket and returns its content
func (mc *MinioClient) DownloadFile(ctx context.Context, bucketName, objectName string) ([]byte, error) {
	object, err := mc.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer object.Close()

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, object)
	if err != nil {
		return nil, fmt.Errorf("failed to read object data: %w", err)
	}

	log.Printf("File %s downloaded successfully from bucket %s", objectName, bucketName)
	return buffer.Bytes(), nil
}
