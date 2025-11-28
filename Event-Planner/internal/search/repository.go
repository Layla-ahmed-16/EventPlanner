package search

import (
	"context"
	"fmt"

	"event-planner/internal/event"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SearchEvents(ctx context.Context, f *EventsFilter) ([]event.EventWithAttendeeInfo, error) {
	query := `
		SELECT 
			e.id,
			e.title,
			e.description,
			e.date,
			e.time,
			e.location,
			e.organizer_id,
			e.created_at,
			ea.role,
			ea.status
		FROM events e
		JOIN event_attendees ea ON e.id = ea.event_id
		WHERE ea.user_id = $1
	`

	args := []interface{}{f.UserID}
	argIdx := 2

	if f.Query != "" {
		query += fmt.Sprintf(" AND (e.title ILIKE $%d OR e.description ILIKE $%d)", argIdx, argIdx+1)
		likeVal := "%" + f.Query + "%"
		args = append(args, likeVal, likeVal)
		argIdx += 2
	}

	if f.DateFrom != "" {
		query += fmt.Sprintf(" AND e.date >= $%d", argIdx)
		args = append(args, f.DateFrom)
		argIdx++
	}

	if f.DateTo != "" {
		query += fmt.Sprintf(" AND e.date <= $%d", argIdx)
		args = append(args, f.DateTo)
		argIdx++
	}

	// Role
	if f.Role != "" {
		query += fmt.Sprintf(" AND ea.role = $%d", argIdx)
		args = append(args, f.Role)
		argIdx++
	}

	// Status
	if f.Status != "" {
		query += fmt.Sprintf(" AND ea.status = $%d", argIdx)
		args = append(args, f.Status)
		argIdx++
	}

	query += " ORDER BY e.date DESC, e.time DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search events: %w", err)
	}
	defer rows.Close()

	var eventsWithInfo []event.EventWithAttendeeInfo
	for rows.Next() {
		var e event.EventWithAttendeeInfo
		if err := rows.Scan(
			&e.ID,
			&e.Title,
			&e.Description,
			&e.Date,
			&e.Time,
			&e.Location,
			&e.OrganizerID,
			&e.CreatedAt,
			&e.Role,
			&e.Status,
		); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		eventsWithInfo = append(eventsWithInfo, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating search results: %w", err)
	}

	if eventsWithInfo == nil {
		eventsWithInfo = []event.EventWithAttendeeInfo{}
	}

	return eventsWithInfo, nil
}

