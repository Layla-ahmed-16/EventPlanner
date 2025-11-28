package invitation

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles all database operations for invitations
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new invitation repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// SendInvitation creates a new invitation
func (r *Repository) SendInvitation(ctx context.Context, invitation *Invitation) error {
	query := `
		INSERT INTO invitations (event_id, inviter_id, invitee_email, invitee_id, role, message)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query,
		invitation.EventID,
		invitation.InviterID,
		invitation.InviteeEmail,
		invitation.InviteeID,
		invitation.Role,
		invitation.Message,
	).Scan(&invitation.ID, &invitation.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to send invitation: %w", err)
	}

	return nil
}

// GetInvitationByID retrieves a single invitation by ID
func (r *Repository) GetInvitationByID(ctx context.Context, invitationID int) (*Invitation, error) {
	query := `
		SELECT id, event_id, inviter_id, invitee_email, invitee_id, role, status, message, created_at, responded_at
		FROM invitations
		WHERE id = $1
	`

	invitation := &Invitation{}
	err := r.db.QueryRow(ctx, query, invitationID).Scan(
		&invitation.ID,
		&invitation.EventID,
		&invitation.InviterID,
		&invitation.InviteeEmail,
		&invitation.InviteeID,
		&invitation.Role,
		&invitation.Status,
		&invitation.Message,
		&invitation.CreatedAt,
		&invitation.RespondedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get invitation: %w", err)
	}

	return invitation, nil
}

// GetInvitationsByEmail retrieves all invitations for a specific email
func (r *Repository) GetInvitationsByEmail(ctx context.Context, email string) ([]InvitationWithDetails, error) {
	query := `
		SELECT 
			i.id, i.event_id, i.inviter_id, i.invitee_email, i.invitee_id, i.role, i.status, i.message, i.created_at, i.responded_at,
			e.title, e.date, e.time, e.location,
			u.email as inviter_email
		FROM invitations i
		JOIN events e ON i.event_id = e.id
		JOIN users u ON i.inviter_id = u.id
		WHERE i.invitee_email = $1
		ORDER BY i.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get invitations by email: %w", err)
	}
	defer rows.Close()

	var invitations []InvitationWithDetails
	for rows.Next() {
		invitation := InvitationWithDetails{}
		err := rows.Scan(
			&invitation.ID,
			&invitation.EventID,
			&invitation.InviterID,
			&invitation.InviteeEmail,
			&invitation.InviteeID,
			&invitation.Role,
			&invitation.Status,
			&invitation.Message,
			&invitation.CreatedAt,
			&invitation.RespondedAt,
			&invitation.EventTitle,
			&invitation.EventDate,
			&invitation.EventTime,
			&invitation.EventLocation,
			&invitation.InviterEmail,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan invitation: %w", err)
		}
		invitations = append(invitations, invitation)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating invitations: %w", err)
	}

	return invitations, nil
}

// GetInvitationsByEventID retrieves all invitations for a specific event
func (r *Repository) GetInvitationsByEventID(ctx context.Context, eventID int) ([]InvitationWithDetails, error) {
	query := `
		SELECT 
			i.id, i.event_id, i.inviter_id, i.invitee_email, i.invitee_id, i.role, i.status, i.message, i.created_at, i.responded_at,
			e.title, e.date, e.time, e.location,
			u.email as inviter_email
		FROM invitations i
		JOIN events e ON i.event_id = e.id
		JOIN users u ON i.inviter_id = u.id
		WHERE i.event_id = $1
		ORDER BY i.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invitations by event: %w", err)
	}
	defer rows.Close()

	var invitations []InvitationWithDetails
	for rows.Next() {
		invitation := InvitationWithDetails{}
		err := rows.Scan(
			&invitation.ID,
			&invitation.EventID,
			&invitation.InviterID,
			&invitation.InviteeEmail,
			&invitation.InviteeID,
			&invitation.Role,
			&invitation.Status,
			&invitation.Message,
			&invitation.CreatedAt,
			&invitation.RespondedAt,
			&invitation.EventTitle,
			&invitation.EventDate,
			&invitation.EventTime,
			&invitation.EventLocation,
			&invitation.InviterEmail,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan invitation: %w", err)
		}
		invitations = append(invitations, invitation)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating invitations: %w", err)
	}

	return invitations, nil
}

// UpdateInvitationStatus updates the status of an invitation
func (r *Repository) UpdateInvitationStatus(ctx context.Context, invitationID int, status string) error {
	query := `
		UPDATE invitations
		SET status = $1, responded_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(ctx, query, status, time.Now(), invitationID)
	if err != nil {
		return fmt.Errorf("failed to update invitation status: %w", err)
	}

	return nil
}

// GetUserIDByEmail retrieves user ID by email (helper function)
func (r *Repository) GetUserIDByEmail(ctx context.Context, email string) (*int, error) {
	query := `SELECT id FROM users WHERE email = $1`

	var userID int
	err := r.db.QueryRow(ctx, query, email).Scan(&userID)
	if err != nil {
		// User doesn't exist, return nil (not an error)
		return nil, nil
	}

	return &userID, nil
}
