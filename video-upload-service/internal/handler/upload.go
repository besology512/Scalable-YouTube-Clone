package handler

import (
    "context"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "your-project/internal/storage"
    "your-project/internal/kafka"
    "your-project/internal/db"
    "your-project/shared-libs/models"
    "your-project/shared-libs/events"
)

type UploadHandler struct {
    Storage *storage.MinioClient
    Kafka   *kafka.Producer
    DB      *db.MongoRepo
}

func (h *UploadHandler) Upload(c *gin.Context) {
    fileHeader, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
        return
    }

    title := c.PostForm("title")
    description := c.PostForm("description")
    userID := c.PostForm("user_id")

    file, err := fileHeader.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not open file"})
        return
    }
    defer file.Close()

    videoID := uuid.New().String()
    fileName := videoID + "-" + fileHeader.Filename

    path, err := h.Storage.UploadFile(context.Background(), file, fileName)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
        return
    }

    // Save metadata to Mongo
    meta := models.Video{
        ID:          videoID,
        Title:       title,
        Description: description,
        Filename:    fileName,
        UserID:      userID,
        UploadedAt:  time.Now(),
        Status:      "processing",
    }
    _ = h.DB.SaveVideo(meta)

    // Emit Kafka event
    event := events.VideoUploadedEvent{
        VideoID:    videoID,
        Filename:   fileName,
        UserID:     userID,
        OriginalURL: path,
        UploadedAt: time.Now().Unix(),
    }
    _ = h.Kafka.Produce("video.uploaded", event)

    c.JSON(http.StatusOK, gin.H{"video_id": videoID})
}
