package minio

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"video-transcoding-service/internal/config"
)

type Client struct {
	Client *minio.Client
	Bucket string
}

func NewClient(cfg config.MinioConfig) (*Client, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		Client: client,
		Bucket: cfg.Bucket,
	}, nil
}

func (c *Client) UploadFile(filePath string) error {
	_, err := c.Client.FPutObject(nil, c.Bucket, "video.mp4", filePath, minio.PutObjectOptions{})
	if err != nil {
		log.Println("Error uploading file to MinIO: ", err)
		return err
	}
	return nil
}
