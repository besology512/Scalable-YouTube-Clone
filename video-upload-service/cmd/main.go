package main

import (
    "github.com/gin-gonic/gin"
    "your-project/internal/handler"
    "your-project/internal/storage"
    "your-project/internal/db"
    "your-project/internal/kafka"
)

func main() {
    minioClient, _ := storage.NewMinioClient("http://localhost:9000", "minio", "minio123", "videos")
    mongoRepo := db.NewMongoRepo("mongodb://localhost:27017", "video_db")
    kafkaProducer := kafka.NewProducer("localhost:9092")

    h := handler.UploadHandler{
        Storage: minioClient,
        DB:      mongoRepo,
        Kafka:   kafkaProducer,
    }

    r := gin.Default()
    r.POST("/upload", h.Upload)
    r.Run(":8080")
}
