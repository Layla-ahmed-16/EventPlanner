package search

import (
	"encoding/json"
	"fmt"
	"net/http"

	"event-planner/internal/auth"
)

// Handler handles HTTP requests for search
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// SearchEvents handles GET /events/search
func (h *Handler) SearchEvents(w http.ResponseWriter, r *http.Request) {
	// Get current user ID from context
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	q := r.URL.Query()

	filter := &EventsFilter{
		Query:    q.Get("q"),
		DateFrom: q.Get("date_from"),
		DateTo:   q.Get("date_to"),
		Role:     q.Get("role"),
		Status:   q.Get("status"),
		UserID:   userID,
	}

	events, err := h.service.SearchEvents(r.Context(), filter)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": events,
	})
}
