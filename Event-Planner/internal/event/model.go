package event

import (
	"encoding/json"
	"time"
)

// Event represents an event in the system
type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"-"` // Will be handled by custom marshaling
	Time        time.Time `json:"-"` // Will be handled by custom marshaling
	Location    string    `json:"location"`
	OrganizerID int       `json:"organizer_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// MarshalJSON custom marshaling for Event to format date and time properly
func (e Event) MarshalJSON() ([]byte, error) {
	type Alias Event
	return json.Marshal(&struct {
		Date string `json:"date"`
		Time string `json:"time"`
		*Alias
	}{
		Date:  e.Date.Format("2006-01-02"),
		Time:  e.Time.Format("15:04:05"),
		Alias: (*Alias)(&e),
	})
}

// CreateEventRequest is the request payload for creating an event
type CreateEventRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Date        string `json:"date" binding:"required"` // YYYY-MM-DD
	Time        string `json:"time" binding:"required"` // HH:MM:SS
	Location    string `json:"location" binding:"required"`
}

// UpdateEventRequest is the request payload for updating an event
type UpdateEventRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	Location    string `json:"location"`
}

// EventAttendee represents an attendee of an event
type EventAttendee struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	EventID   int       `json:"event_id"`
	Role      string    `json:"role"`   // 'organizer', 'attendee', 'collaborator'
	Status    string    `json:"status"` // 'going', 'maybe', 'not_going'
	CreatedAt time.Time `json:"created_at"`
}

// EventWithAttendeeInfo represents an event with the user's attendance info
type EventWithAttendeeInfo struct {
	Event
	Role   string `json:"role"`
	Status string `json:"status"`
}

// AddAttendeeRequest is the request payload for inviting someone to an event
type AddAttendeeRequest struct {
	UserID int    `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required"` // 'attendee', 'collaborator', or 'organizer'
}

// UpdateAttendanceRequest is the request payload for updating attendance status
type UpdateAttendanceRequest struct {
	Status string `json:"status" binding:"required"` // 'going', 'maybe', 'not_going'
}
