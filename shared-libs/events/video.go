package events

type VideoUploadedEvent struct {
    VideoID       string `json:"video_id"`
    Filename      string `json:"filename"`
    UserID        string `json:"user_id"`
    OriginalURL   string `json:"original_url"`  // from MinIO
    UploadedAt    int64  `json:"uploaded_at"`
}

type VideoTranscodedEvent struct {
    VideoID    string              `json:"video_id"`
    Variants   []TranscodedVariant `json:"variants"`
    CompletedAt int64              `json:"completed_at"`
}

type TranscodedVariant struct {
    Resolution string `json:"resolution"` // 360p, 720p, etc
    URL        string `json:"url"`        // MinIO path
}
