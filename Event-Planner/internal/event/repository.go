package event

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles all database operations for events
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new event repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// CreateEvent inserts a new event into the database
func (r *Repository) CreateEvent(ctx context.Context, event *Event) error {
	query := `
		INSERT INTO events (title, description, date, time, location, organizer_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query,
		event.Title,
		event.Description,
		event.Date,
		event.Time,
		event.Location,
		event.OrganizerID,
	).Scan(&event.ID, &event.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

// GetEventByID retrieves a single event by ID
func (r *Repository) GetEventByID(ctx context.Context, eventID int) (*Event, error) {
	query := `
		SELECT id, title, description, date, time, location, organizer_id, created_at
		FROM events
		WHERE id = $1
	`

	event := &Event{}
	err := r.db.QueryRow(ctx, query, eventID).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.Date,
		&event.Time,
		&event.Location,
		&event.OrganizerID,
		&event.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return event, nil
}

// GetAllEvents retrieves all events from the database
func (r *Repository) GetAllEvents(ctx context.Context) ([]Event, error) {
	query := `
		SELECT id, title, description, date, time, location, organizer_id, created_at
		FROM events
		ORDER BY date DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		event := Event{}
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.Date,
			&event.Time,
			&event.Location,
			&event.OrganizerID,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating events: %w", err)
	}

	return events, nil
}

// GetEventsByOrganizerID retrieves all events created by a specific user
func (r *Repository) GetEventsByOrganizerID(ctx context.Context, organizerID int) ([]Event, error) {
	query := `
		SELECT id, title, description, date, time, location, organizer_id, created_at
		FROM events
		WHERE organizer_id = $1
		ORDER BY date DESC
	`

	rows, err := r.db.Query(ctx, query, organizerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organizer events: %w", err)
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		event := Event{}
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.Date,
			&event.Time,
			&event.Location,
			&event.OrganizerID,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating organizer events: %w", err)
	}

	return events, nil
}

// UpdateEvent updates an existing event
func (r *Repository) UpdateEvent(ctx context.Context, eventID int, updates *UpdateEventRequest) (*Event, error) {
	// First, get the current event
	currentEvent, err := r.GetEventByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// Apply updates (only non-empty fields)
	if updates.Title != "" {
		currentEvent.Title = updates.Title
	}
	if updates.Description != "" {
		currentEvent.Description = updates.Description
	}
	if updates.Date != "" {
		eventDate, err := time.Parse("2006-01-02", updates.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %w", err)
		}
		currentEvent.Date = eventDate
	}
	if updates.Time != "" {
		eventTime, err := time.Parse("15:04:05", updates.Time)
		if err != nil {
			return nil, fmt.Errorf("invalid time format: %w", err)
		}
		currentEvent.Time = eventTime
	}
	if updates.Location != "" {
		currentEvent.Location = updates.Location
	}

	query := `
		UPDATE events
		SET title = $1, description = $2, date = $3, time = $4, location = $5
		WHERE id = $6
		RETURNING id, title, description, date, time, location, organizer_id, created_at
	`

	err = r.db.QueryRow(ctx, query,
		currentEvent.Title,
		currentEvent.Description,
		currentEvent.Date,
		currentEvent.Time,
		currentEvent.Location,
		eventID,
	).Scan(
		&currentEvent.ID,
		&currentEvent.Title,
		&currentEvent.Description,
		&currentEvent.Date,
		&currentEvent.Time,
		&currentEvent.Location,
		&currentEvent.OrganizerID,
		&currentEvent.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	return currentEvent, nil
}

// DeleteEvent removes an event from the database
func (r *Repository) DeleteEvent(ctx context.Context, eventID int) error {
	query := `DELETE FROM events WHERE id = $1`

	result, err := r.db.Exec(ctx, query, eventID)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("event not found")
	}

	return nil
}

// JoinEvent adds a user as an attendee to an event
func (r *Repository) JoinEvent(ctx context.Context, userID, eventID int) error {
	query := `
		INSERT INTO event_attendees (user_id, event_id, role, status)
		VALUES ($1, $2, 'attendee', 'going')
	`

	_, err := r.db.Exec(ctx, query, userID, eventID)
	if err != nil {
		return fmt.Errorf("failed to join event: %w", err)
	}

	return nil
}

// GetEventsByAttendeeID retrieves all events where the user is an attendee (including as organizer)
func (r *Repository) GetEventsByAttendeeID(ctx context.Context, userID int) ([]EventWithAttendeeInfo, error) {
	query := `
		SELECT e.id, e.title, e.description, e.date, e.time, e.location, e.organizer_id, e.created_at, ea.role, ea.status
		FROM events e
		JOIN event_attendees ea ON e.id = ea.event_id
		WHERE ea.user_id = $1
		ORDER BY e.date DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendee events: %w", err)
	}
	defer rows.Close()

	var events []EventWithAttendeeInfo
	for rows.Next() {
		event := EventWithAttendeeInfo{}
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.Date,
			&event.Time,
			&event.Location,
			&event.OrganizerID,
			&event.CreatedAt,
			&event.Role,
			&event.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan attendee event: %w", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating attendee events: %w", err)
	}

	return events, nil
}

// GetMyOrganizedEvents retrieves all events organized by a specific user
func (r *Repository) GetMyOrganizedEvents(ctx context.Context, organizerID int) ([]Event, error) {
	query := `
		SELECT id, title, description, date, time, location, organizer_id, created_at
		FROM events
		WHERE organizer_id = $1
		ORDER BY date DESC
	`

	rows, err := r.db.Query(ctx, query, organizerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organized events: %w", err)
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		event := Event{}
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.Date,
			&event.Time,
			&event.Location,
			&event.OrganizerID,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan organized event: %w", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating organized events: %w", err)
	}

	return events, nil
}

// AddOrganizerAsAttendee automatically adds the organizer as an attendee when creating an event
func (r *Repository) AddOrganizerAsAttendee(ctx context.Context, userID, eventID int) error {
	query := `
		INSERT INTO event_attendees (user_id, event_id, role, status)
		VALUES ($1, $2, 'organizer', 'going')
	`

	_, err := r.db.Exec(ctx, query, userID, eventID)
	if err != nil {
		return fmt.Errorf("failed to add organizer as attendee: %w", err)
	}

	return nil
}

// AddAttendee adds a user to an event
func (r *Repository) AddAttendee(ctx context.Context, eventID, userID int, role string) error {
	query := `
		INSERT INTO event_attendees (user_id, event_id, role, status)
		VALUES ($1, $2, $3, 'going')
	`

	_, err := r.db.Exec(ctx, query, userID, eventID, role)
	if err != nil {
		return fmt.Errorf("failed to add user to event: %w", err)
	}

	return nil
}

// UpdateAttendanceStatus updates a user's attendance status for an event
func (r *Repository) UpdateAttendanceStatus(ctx context.Context, userID, eventID int, status string) error {
	query := `
		UPDATE event_attendees
		SET status = $1
		WHERE user_id = $2 AND event_id = $3
	`

	result, err := r.db.Exec(ctx, query, status, userID, eventID)
	if err != nil {
		return fmt.Errorf("failed to update attendance status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("attendance record not found")
	}

	return nil
}

// GetEventAttendees retrieves all attendees for an event
func (r *Repository) GetEventAttendees(ctx context.Context, eventID int) ([]EventAttendee, error) {
	query := `
		SELECT id, user_id, event_id, role, status, created_at
		FROM event_attendees
		WHERE event_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event attendees: %w", err)
	}
	defer rows.Close()

	var attendees []EventAttendee
	for rows.Next() {
		attendee := EventAttendee{}
		err := rows.Scan(
			&attendee.ID,
			&attendee.UserID,
			&attendee.EventID,
			&attendee.Role,
			&attendee.Status,
			&attendee.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan attendee: %w", err)
		}
		attendees = append(attendees, attendee)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating attendees: %w", err)
	}

	return attendees, nil
}
