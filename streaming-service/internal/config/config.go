package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	ServerPort     string
}

func Load() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found — using system environment variables.")
	}

	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}

	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	if accessKey == "" {
		accessKey = "minioadmin"
	}

	secretKey := os.Getenv("MINIO_SECRET_KEY")
	if secretKey == "" {
		secretKey = "minioadmin"
	}

	bucket := os.Getenv("MINIO_BUCKET")
	if bucket == "" {
		bucket = "prossessed-videos"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	return &Config{
		MinioEndpoint:  endpoint,
		MinioAccessKey: accessKey,
		MinioSecretKey: secretKey,
		MinioBucket:    bucket,
		ServerPort:     port,
	}
}
