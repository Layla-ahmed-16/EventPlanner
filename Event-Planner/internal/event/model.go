package event

import (
	"encoding/json"
	"time"
)
type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"-"` 
	Time        time.Time `json:"-"` 
	Location    string    `json:"location"`
	OrganizerID int       `json:"organizer_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// format date and time
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
type CreateEventRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Date        string `json:"date" binding:"required"` // YYYY-MM-DD
	Time        string `json:"time" binding:"required"` // HH:MM:SS
	Location    string `json:"location" binding:"required"`
}

type UpdateEventRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	Location    string `json:"location"`
}

type EventAttendee struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	EventID   int       `json:"event_id"`
	Role      string    `json:"role"`   // 'organizer', 'attendee', 'collaborator'
	Status    string    `json:"status"` // 'going', 'maybe', 'not_going'
	CreatedAt time.Time `json:"created_at"`
}

type EventWithAttendeeInfo struct {
	Event
	Role   string `json:"role"`
	Status string `json:"status"`
}

type AddAttendeeRequest struct {
	UserID int    `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required"` // 'attendee', 'collaborator', or 'organizer'
}

type UpdateAttendanceRequest struct {
	Status string `json:"status" binding:"required"` // 'going', 'maybe', 'not_going'
}

