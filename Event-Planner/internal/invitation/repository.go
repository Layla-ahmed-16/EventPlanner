package invitation

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

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

func (r *Repository) GetInvitationsByEmail(ctx context.Context, email string) ([]InvitationWithDetails, error) {
	query := `
        SELECT 
            i.id,
            i.event_id,
            i.inviter_id,
            i.invitee_email,
            i.invitee_id,
            i.role,
            i.status,
            i.message,
            i.created_at,
            i.responded_at,
            e.title,
            to_char(e.date, 'YYYY-MM-DD') AS event_date,
            to_char(e.time, 'HH24:MI:SS') AS event_time,
            e.location,
            u.email AS inviter_email
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
		inv := InvitationWithDetails{}
		err := rows.Scan(
			&inv.ID,
			&inv.EventID,
			&inv.InviterID,
			&inv.InviteeEmail,
			&inv.InviteeID,
			&inv.Role,
			&inv.Status,
			&inv.Message,
			&inv.CreatedAt,
			&inv.RespondedAt,
			&inv.EventTitle,
			&inv.EventDate,
			&inv.EventTime,
			&inv.EventLocation,
			&inv.InviterEmail,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan invitation: %w", err)
		}
		invitations = append(invitations, inv)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating invitations: %w", err)
	}

	if invitations == nil {
		invitations = []InvitationWithDetails{}
	}

	return invitations, nil
}

func (r *Repository) GetInvitationsByEventID(ctx context.Context, eventID int) ([]InvitationWithDetails, error) {
	query := `
        SELECT 
            i.id,
            i.event_id,
            i.inviter_id,
            i.invitee_email,
            i.invitee_id,
            i.role,
            i.status,
            i.message,
            i.created_at,
            i.responded_at,
            e.title,
            to_char(e.date, 'YYYY-MM-DD') AS event_date,
            to_char(e.time, 'HH24:MI:SS') AS event_time,
            e.location,
            u.email AS inviter_email
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
		inv := InvitationWithDetails{}
		err := rows.Scan(
			&inv.ID,
			&inv.EventID,
			&inv.InviterID,
			&inv.InviteeEmail,
			&inv.InviteeID,
			&inv.Role,
			&inv.Status,
			&inv.Message,
			&inv.CreatedAt,
			&inv.RespondedAt,
			&inv.EventTitle,
			&inv.EventDate,
			&inv.EventTime,
			&inv.EventLocation,
			&inv.InviterEmail,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan invitation: %w", err)
		}
		invitations = append(invitations, inv)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating invitations: %w", err)
	}

	if invitations == nil {
		invitations = []InvitationWithDetails{}
	}

	return invitations, nil
}

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

