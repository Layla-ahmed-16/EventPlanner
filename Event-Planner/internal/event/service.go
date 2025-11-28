package event

import (
	"context"
	"fmt"
	"time"
)

// Service handles business logic for events
type Service struct {
	repo *Repository
}

// NewService creates a new event service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateEvent validates and creates a new event
func (s *Service) CreateEvent(ctx context.Context, req *CreateEventRequest, organizerID int) (*Event, error) {
	// Validate required fields
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Validate date and time formats
	if err := s.validateDateTime(req.Date, req.Time); err != nil {
		return nil, err
	}

	// Parse dates to ensure the event is in the future
	if err := s.validateFutureEvent(req.Date, req.Time); err != nil {
		return nil, err
	}

	// Parse date and time strings
	eventDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	eventTime, err := time.Parse("15:04:05", req.Time)
	if err != nil {
		return nil, fmt.Errorf("invalid time format: %w", err)
	}

	event := &Event{
		Title:       req.Title,
		Description: req.Description,
		Date:        eventDate,
		Time:        eventTime,
		Location:    req.Location,
		OrganizerID: organizerID,
	}

	if err := s.repo.CreateEvent(ctx, event); err != nil {
		return nil, err
	}

	// Automatically add the organizer as an attendee with 'organizer' role
	if err := s.repo.AddOrganizerAsAttendee(ctx, organizerID, event.ID); err != nil {
		return nil, fmt.Errorf("failed to add organizer as attendee: %w", err)
	}

	return event, nil
}

// GetEventByID retrieves an event by ID
func (s *Service) GetEventByID(ctx context.Context, eventID int) (*Event, error) {
	if eventID <= 0 {
		return nil, fmt.Errorf("invalid event ID")
	}

	event, err := s.repo.GetEventByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// GetAllEvents retrieves all events
func (s *Service) GetAllEvents(ctx context.Context) ([]Event, error) {
	events, err := s.repo.GetAllEvents(ctx)
	if err != nil {
		return nil, err
	}

	// Return empty slice instead of nil for JSON response
	if events == nil {
		events = []Event{}
	}

	return events, nil
}

// GetEventsByOrganizerID retrieves all events created by a specific user
func (s *Service) GetEventsByOrganizerID(ctx context.Context, organizerID int) ([]Event, error) {
	if organizerID <= 0 {
		return nil, fmt.Errorf("invalid organizer ID")
	}

	events, err := s.repo.GetEventsByOrganizerID(ctx, organizerID)
	if err != nil {
		return nil, err
	}

	// Return empty slice instead of nil for JSON response
	if events == nil {
		events = []Event{}
	}

	return events, nil
}

// UpdateEvent validates and updates an event
func (s *Service) UpdateEvent(ctx context.Context, eventID int, req *UpdateEventRequest, organizerID int) (*Event, error) {
	if eventID <= 0 {
		return nil, fmt.Errorf("invalid event ID")
	}

	// Get the event to check ownership
	event, err := s.repo.GetEventByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// Check if the user is the organizer
	if event.OrganizerID != organizerID {
		return nil, fmt.Errorf("you are not authorized to update this event")
	}

	// Validate date and time if provided
	if req.Date != "" && req.Time != "" {
		if err := s.validateDateTime(req.Date, req.Time); err != nil {
			return nil, err
		}
		if err := s.validateFutureEvent(req.Date, req.Time); err != nil {
			return nil, err
		}
	} else if req.Date != "" {
		if err := s.validateDateFormat(req.Date); err != nil {
			return nil, err
		}
	} else if req.Time != "" {
		if err := s.validateTimeFormat(req.Time); err != nil {
			return nil, err
		}
	}

	updatedEvent, err := s.repo.UpdateEvent(ctx, eventID, req)
	if err != nil {
		return nil, err
	}

	return updatedEvent, nil
}

// DeleteEvent validates and deletes an event
func (s *Service) DeleteEvent(ctx context.Context, eventID int, organizerID int) error {
	if eventID <= 0 {
		return fmt.Errorf("invalid event ID")
	}

	// Get the event to check ownership
	event, err := s.repo.GetEventByID(ctx, eventID)
	if err != nil {
		return err
	}

	// Check if the user is the organizer
	if event.OrganizerID != organizerID {
		return fmt.Errorf("you are not authorized to delete this event")
	}

	if err := s.repo.DeleteEvent(ctx, eventID); err != nil {
		return err
	}

	return nil
}

// Validation helper functions

func (s *Service) validateCreateRequest(req *CreateEventRequest) error {
	if req.Title == "" {
		return fmt.Errorf("event title is required")
	}

	if req.Date == "" {
		return fmt.Errorf("event date is required")
	}

	if req.Time == "" {
		return fmt.Errorf("event time is required")
	}

	if req.Location == "" {
		return fmt.Errorf("event location is required")
	}

	if len(req.Title) > 255 {
		return fmt.Errorf("event title must not exceed 255 characters")
	}

	if len(req.Description) > 1000 {
		return fmt.Errorf("event description must not exceed 1000 characters")
	}

	return nil
}

func (s *Service) validateDateTime(date, timeStr string) error {
	if err := s.validateDateFormat(date); err != nil {
		return err
	}

	if err := s.validateTimeFormat(timeStr); err != nil {
		return err
	}

	return nil
}

func (s *Service) validateDateFormat(date string) error {
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return fmt.Errorf("invalid date format, use YYYY-MM-DD")
	}
	return nil
}

func (s *Service) validateTimeFormat(timeStr string) error {
	_, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return fmt.Errorf("invalid time format, use HH:MM:SS")
	}
	return nil
}

func (s *Service) validateFutureEvent(date, timeStr string) error {
	eventDateTime, err := time.Parse("2006-01-02 15:04:05", date+" "+timeStr)
	if err != nil {
		return fmt.Errorf("failed to parse event date and time")
	}

	if eventDateTime.Before(time.Now()) {
		return fmt.Errorf("event date and time must be in the future")
	}

	return nil
}

// JoinEvent allows a user to join an event as an attendee
func (s *Service) JoinEvent(ctx context.Context, userID, eventID int) error {
	if userID <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	if eventID <= 0 {
		return fmt.Errorf("invalid event ID")
	}

	// Check if event exists
	_, err := s.repo.GetEventByID(ctx, eventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}

	if err := s.repo.JoinEvent(ctx, userID, eventID); err != nil {
		return err
	}

	return nil
}

// GetMyAttendingEvents retrieves all events where the user is an attendee (including as organizer)
func (s *Service) GetMyAttendingEvents(ctx context.Context, userID int) ([]EventWithAttendeeInfo, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	events, err := s.repo.GetEventsByAttendeeID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Return empty slice instead of nil for JSON response
	if events == nil {
		events = []EventWithAttendeeInfo{}
	}

	return events, nil
}

// InviteUserToEvent invites a user to an event (only event creator can invite)
func (s *Service) InviteUserToEvent(ctx context.Context, eventID, inviterID int, req *AddAttendeeRequest) error {
	if eventID <= 0 {
		return fmt.Errorf("invalid event ID")
	}

	if inviterID <= 0 {
		return fmt.Errorf("invalid inviter ID")
	}

	if req.UserID <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	// Validate role - now allowing 'organizer' role for co-organizers
	if req.Role != "attendee" && req.Role != "collaborator" && req.Role != "organizer" {
		return fmt.Errorf("invalid role: must be 'attendee', 'collaborator', or 'organizer'")
	}

	// Check if event exists and get event details
	event, err := s.repo.GetEventByID(ctx, eventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}

	// Check if the inviter is the event creator
	if event.OrganizerID != inviterID {
		return fmt.Errorf("only the event creator can invite users to this event")
	}

	// Check if user is trying to invite themselves
	if req.UserID == inviterID {
		return fmt.Errorf("you cannot invite yourself to the event")
	}

	if err := s.repo.AddAttendee(ctx, eventID, req.UserID, req.Role); err != nil {
		return err
	}

	return nil
} // UpdateAttendanceStatus updates a user's attendance status for an event
func (s *Service) UpdateAttendanceStatus(ctx context.Context, userID, eventID int, status string) error {
	if userID <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	if eventID <= 0 {
		return fmt.Errorf("invalid event ID")
	}

	// Validate status
	validStatuses := []string{"going", "maybe", "not_going"}
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("invalid status: must be 'going', 'maybe', or 'not_going'")
	}

	if err := s.repo.UpdateAttendanceStatus(ctx, userID, eventID, status); err != nil {
		return err
	}

	return nil
}

// GetEventAttendees retrieves all attendees for an event
func (s *Service) GetEventAttendees(ctx context.Context, eventID int) ([]EventAttendee, error) {
	if eventID <= 0 {
		return nil, fmt.Errorf("invalid event ID")
	}

	attendees, err := s.repo.GetEventAttendees(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// Return empty slice instead of nil for JSON response
	if attendees == nil {
		attendees = []EventAttendee{}
	}

	return attendees, nil
}

// GetMyOrganizedEvents retrieves all events organized by a specific user
func (s *Service) GetMyOrganizedEvents(ctx context.Context, organizerID int) ([]Event, error) {
	if organizerID <= 0 {
		return nil, fmt.Errorf("invalid organizer ID")
	}

	events, err := s.repo.GetMyOrganizedEvents(ctx, organizerID)
	if err != nil {
		return nil, err
	}

	// Return empty slice instead of nil for JSON response
	if events == nil {
		events = []Event{}
	}

	return events, nil
}
