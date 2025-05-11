package model

type Video struct {
	VideoID     string   `json:"video_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}
