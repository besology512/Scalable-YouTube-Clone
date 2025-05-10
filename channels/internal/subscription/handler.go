package subscription

import (
	"channel-service/internal/models"
	"channel-service/pkg/database"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type Handler struct {
	DB *database.PostgresDB
}

func NewHandler(db *database.PostgresDB) *Handler {
	return &Handler{DB: db}
}

// Extract channel ID from URL path (e.g., "/subscriptions/abc123/subscribe" -> "abc123")
func getChannelIDFromURL(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return ""
	}
	return parts[3] // Adjust index based on your route structure
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Subscribe(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	channelID := getChannelIDFromURL(r.URL.Path)

	// Validate channel exists
	var channel models.Channel
	if err := h.DB.DB.First(&channel, "id = ?", channelID).Error; err != nil {
		http.Error(w, "Channel not found", http.StatusNotFound)
		return
	}

	// Create subscription
	sub := models.Subscription{
		UserID:    userID,
		ChannelID: channelID,
	}

	if err := h.DB.DB.Create(&sub).Error; err != nil {
		http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
		return
	}

	// Update subscriber count
	h.DB.DB.Model(&models.Channel{}).
		Where("id = ?", channelID).
		Update("subscribers", gorm.Expr("subscribers + 1"))

	w.WriteHeader(http.StatusCreated)
}
