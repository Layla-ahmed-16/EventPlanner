package invitation

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)
type EventAttendeeService interface {
	AddAttendee(ctx context.Context, eventID, userID int, role string) error
}

type Service struct {
	repo            *Repository
	attendeeService EventAttendeeService
}

func NewService(repo *Repository, attendeeService EventAttendeeService) *Service {
	return &Service{
		repo:            repo,
		attendeeService: attendeeService,
	}
}

//validates and sends an invitation
func (s *Service) SendInvitation(ctx context.Context, req *SendInvitationRequest, inviterID int) (*Invitation, error) {
	// Validate input
	if err := s.validateSendInvitationRequest(req); err != nil {
		return nil, err
	}

	// Check if invitee user exists
	inviteeID, err := s.repo.GetUserIDByEmail(ctx, req.InviteeEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to check invitee: %w", err)
	}

	invitation := &Invitation{
		EventID:      req.EventID,
		InviterID:    inviterID,
		InviteeEmail: req.InviteeEmail,
		InviteeID:    inviteeID,
		Role:         req.Role,
		Message:      req.Message,
		Status:       "pending",
	}

	if err := s.repo.SendInvitation(ctx, invitation); err != nil {
		return nil, err
	}

	return invitation, nil
}

//retrieves all invitations for a user by email
func (s *Service) GetMyInvitations(ctx context.Context, email string) ([]InvitationWithDetails, error) {
	invitations, err := s.repo.GetInvitationsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Return empty slice instead of nil for JSON response
	if invitations == nil {
		invitations = []InvitationWithDetails{}
	}

	return invitations, nil
}

//retrieves all invitations for a specific event
func (s *Service) GetEventInvitations(ctx context.Context, eventID int) ([]InvitationWithDetails, error) {
	invitations, err := s.repo.GetInvitationsByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if invitations == nil {
		invitations = []InvitationWithDetails{}
	}

	return invitations, nil
}
func (s *Service) RespondToInvitation(ctx context.Context, invitationID int, status string, userEmail string) error {
	// Validate status
	if status != "accepted" && status != "declined" {
		return fmt.Errorf("invalid status: must be 'accepted' or 'declined'")
	}
	
	invitation, err := s.repo.GetInvitationByID(ctx, invitationID)
	if err != nil {
		return fmt.Errorf("invitation not found: %w", err)
	}

	if invitation.InviteeEmail != userEmail {
		return fmt.Errorf("you are not authorized to respond to this invitation")
	}

	if invitation.Status != "pending" {
		return fmt.Errorf("invitation has already been responded to")
	}

	if err := s.repo.UpdateInvitationStatus(ctx, invitationID, status); err != nil {
		return err
	}

	if status == "accepted" && invitation.InviteeID != nil && s.attendeeService != nil {
		if err := s.attendeeService.AddAttendee(ctx, invitation.EventID, *invitation.InviteeID, invitation.Role); err != nil {
			return fmt.Errorf("failed to add invitee as attendee: %w", err)
		}
	}

	return nil
}

func (s *Service) validateSendInvitationRequest(req *SendInvitationRequest) error {
	if req.EventID <= 0 {
		return fmt.Errorf("invalid event ID")
	}

	if req.InviteeEmail == "" {
		return fmt.Errorf("invitee email is required")
	}

	if !s.isValidEmail(req.InviteeEmail) {
		return fmt.Errorf("invalid email format")
	}

	if req.Role != "attendee" && req.Role != "collaborator" && req.Role != "organizer" {
		return fmt.Errorf("invalid role: must be 'attendee', 'collaborator', or 'organizer'")
	}

	if len(req.Message) > 500 {
		return fmt.Errorf("message must not exceed 500 characters")
	}

	return nil
}

func (s *Service) isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if len(email) > 254 {
		return false
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

