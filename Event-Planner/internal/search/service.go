package search

import (
	"context"
	"fmt"
	"time"

	"event-planner/internal/event"
)

// Service handles business logic for search
type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// SearchEvents searches events for the current user with filters
func (s *Service) SearchEvents(ctx context.Context, f *EventsFilter) ([]event.EventWithAttendeeInfo, error) {
	if f.UserID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Validate role if provided
	if f.Role != "" {
		validRoles := map[string]bool{
			"organizer":    true,
			"attendee":     true,
			"collaborator": true,
		}
		if !validRoles[f.Role] {
			return nil, fmt.Errorf("invalid role: must be 'organizer', 'attendee', or 'collaborator'")
		}
	}

	// Validate status if provided
	if f.Status != "" {
		validStatuses := map[string]bool{
			"going":     true,
			"maybe":     true,
			"not_going": true,
		}
		if !validStatuses[f.Status] {
			return nil, fmt.Errorf("invalid status: must be 'going', 'maybe', or 'not_going'")
		}
	}

	// Validate dates if provided
	if f.DateFrom != "" {
		if _, err := time.Parse("2006-01-02", f.DateFrom); err != nil {
			return nil, fmt.Errorf("invalid date_from format, use YYYY-MM-DD")
		}
	}
	if f.DateTo != "" {
		if _, err := time.Parse("2006-01-02", f.DateTo); err != nil {
			return nil, fmt.Errorf("invalid date_to format, use YYYY-MM-DD")
		}
	}

	events, err := s.repo.SearchEvents(ctx, f)
	if err != nil {
		return nil, err
	}

	if events == nil {
		events = []event.EventWithAttendeeInfo{}
	}

	return events, nil
}
