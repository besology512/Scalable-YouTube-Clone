package channel

import (
	"channel-service/internal/models"
	"channel-service/pkg/database"
	"encoding/json"
	"net/http"
	"strings"
)

type Handler struct {
	DB *database.PostgresDB
}

func NewHandler(db *database.PostgresDB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateChannel(w, r)
	case http.MethodGet:
		h.GetChannel(w, r)
		// Add other methods
	}
}

func (h *Handler) GetChannel(w http.ResponseWriter, r *http.Request) {
	channelID := strings.Split(r.URL.Path, "/")[2] // Extract from /channels/{id}
	var ch models.Channel

	if err := h.DB.DB.First(&ch, "id = ?", channelID).Error; err != nil {
		http.Error(w, "Channel not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(ch)
}

func (h *Handler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	var ch models.Channel
	if err := json.NewDecoder(r.Body).Decode(&ch); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get user ID from JWT
	userID := r.Context().Value("userID").(string)
	ch.UserID = userID

	if err := h.DB.CreateChannel(&ch); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ch)
}
