package invitation

import "time"

// Invitation represents an invitation to an event
type Invitation struct {
	ID           int        `json:"id"`
	EventID      int        `json:"event_id"`
	InviterID    int        `json:"inviter_id"`
	InviteeEmail string     `json:"invitee_email"`
	InviteeID    *int       `json:"invitee_id,omitempty"`
	Role         string     `json:"role"`   // 'attendee' or 'collaborator'
	Status       string     `json:"status"` // 'pending', 'accepted', 'declined'
	Message      string     `json:"message,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	RespondedAt  *time.Time `json:"responded_at,omitempty"`
}

// InvitationWithDetails includes event and inviter details
type InvitationWithDetails struct {
	Invitation
	EventTitle    string `json:"event_title"`
	EventDate     string `json:"event_date"`
	EventTime     string `json:"event_time"`
	EventLocation string `json:"event_location"`
	InviterEmail  string `json:"inviter_email"`
}

// SendInvitationRequest is the request payload for sending invitations
type SendInvitationRequest struct {
	EventID      int    `json:"event_id" binding:"required"`
	InviteeEmail string `json:"invitee_email" binding:"required"`
	Role         string `json:"role" binding:"required"` // 'attendee' or 'collaborator'
	Message      string `json:"message,omitempty"`
}

// RespondToInvitationRequest is the request payload for responding to invitations
type RespondToInvitationRequest struct {
	Status string `json:"status" binding:"required"` // 'accepted' or 'declined'
}
