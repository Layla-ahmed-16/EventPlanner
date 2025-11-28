package invitation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"event-planner/internal/auth"
)

// Handler handles HTTP requests for invitations
type Handler struct {
	service *Service
}

// NewHandler creates a new invitation handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// SendInvitation handles POST /invitations
func (h *Handler) SendInvitation(w http.ResponseWriter, r *http.Request) {
	// Get inviter ID from context (set by auth middleware)
	inviterID, ok := auth.GetUserID(r.Context())
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req SendInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	invitation, err := h.service.SendInvitation(r.Context(), &req, inviterID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "invitation sent successfully",
		"data":    invitation,
	})
}

// GetMyInvitations handles GET /invitations/my
func (h *Handler) GetMyInvitations(w http.ResponseWriter, r *http.Request) {
	// Get user email from context - we need to modify auth to include email
	// For now, we'll get it from query parameter as a workaround
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, `{"error": "email parameter is required"}`, http.StatusBadRequest)
		return
	}

	invitations, err := h.service.GetMyInvitations(r.Context(), email)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": invitations,
	})
}

// GetEventInvitations handles GET /events/{id}/invitations
func (h *Handler) GetEventInvitations(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "invalid event ID"}`, http.StatusBadRequest)
		return
	}

	invitations, err := h.service.GetEventInvitations(r.Context(), eventID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": invitations,
	})
}

// RespondToInvitation handles PUT /invitations/{id}/respond
func (h *Handler) RespondToInvitation(w http.ResponseWriter, r *http.Request) {
	// For now, we'll get email from query parameter
	// In a complete implementation, you'd get this from the auth context
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, `{"error": "email parameter is required"}`, http.StatusBadRequest)
		return
	}

	idStr := r.PathValue("id")
	invitationID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "invalid invitation ID"}`, http.StatusBadRequest)
		return
	}

	var req RespondToInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	err = h.service.RespondToInvitation(r.Context(), invitationID, req.Status, email)
	if err != nil {
		if err.Error() == "you are not authorized to respond to this invitation" {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusForbidden)
		} else {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "invitation response recorded successfully",
	})
}
