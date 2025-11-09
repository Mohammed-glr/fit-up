package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type CoachInvitation struct {
	ID               string     `json:"id"`
	CoachID          string     `json:"coach_id"`
	Email            string     `json:"email"`
	FirstName        *string    `json:"first_name"`
	LastName         *string    `json:"last_name"`
	InvitationToken  string     `json:"invitation_token"`
	Status           string     `json:"status"`
	CustomMessage    *string    `json:"custom_message"`
	ExpiresAt        time.Time  `json:"expires_at"`
	CreatedAt        time.Time  `json:"created_at"`
	AcceptedAt       *time.Time `json:"accepted_at"`
	AcceptedByUserID *string    `json:"accepted_by_user_id"`
}


func (s *Store) CreateInvitation(ctx context.Context, inv *CoachInvitation) error {
	query := `
		INSERT INTO coach_invitations (
			id, coach_id, email, first_name, last_name, 
			invitation_token, status, custom_message, expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at
	`

	err := s.db.QueryRow(ctx, query,
		inv.ID,
		inv.CoachID,
		inv.Email,
		inv.FirstName,
		inv.LastName,
		inv.InvitationToken,
		inv.Status,
		inv.CustomMessage,
		inv.ExpiresAt,
	).Scan(&inv.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create invitation: %w", err)
	}

	return nil
}

func (s *Store) GetInvitationByToken(ctx context.Context, token string) (*CoachInvitation, error) {
	query := `
		SELECT 
			id, coach_id, email, first_name, last_name,
			invitation_token, status, custom_message, expires_at,
			created_at, accepted_at, accepted_by_user_id
		FROM coach_invitations
		WHERE invitation_token = $1
	`

	var inv CoachInvitation
	err := s.db.QueryRow(ctx, query, token).Scan(
		&inv.ID,
		&inv.CoachID,
		&inv.Email,
		&inv.FirstName,
		&inv.LastName,
		&inv.InvitationToken,
		&inv.Status,
		&inv.CustomMessage,
		&inv.ExpiresAt,
		&inv.CreatedAt,
		&inv.AcceptedAt,
		&inv.AcceptedByUserID,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("invitation not found")
		}
		return nil, fmt.Errorf("failed to get invitation: %w", err)
	}

	return &inv, nil
}

func (s *Store) GetInvitationsByCoachID(ctx context.Context, coachID string) ([]*CoachInvitation, error) {
	query := `
		SELECT 
			id, coach_id, email, first_name, last_name,
			invitation_token, status, custom_message, expires_at,
			created_at, accepted_at, accepted_by_user_id
		FROM coach_invitations
		WHERE coach_id = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(ctx, query, coachID)
	if err != nil {
		return nil, fmt.Errorf("failed to query invitations: %w", err)
	}
	defer rows.Close()

	var invitations []*CoachInvitation
	for rows.Next() {
		var inv CoachInvitation
		err := rows.Scan(
			&inv.ID,
			&inv.CoachID,
			&inv.Email,
			&inv.FirstName,
			&inv.LastName,
			&inv.InvitationToken,
			&inv.Status,
			&inv.CustomMessage,
			&inv.ExpiresAt,
			&inv.CreatedAt,
			&inv.AcceptedAt,
			&inv.AcceptedByUserID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan invitation: %w", err)
		}
		invitations = append(invitations, &inv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating invitations: %w", err)
	}

	return invitations, nil
}

func (s *Store) UpdateInvitationStatus(ctx context.Context, id, status string) error {
	query := `
		UPDATE coach_invitations
		SET status = $1
		WHERE id = $2
	`

	result, err := s.db.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update invitation status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("invitation not found")
	}

	return nil
}

func (s *Store) AcceptInvitation(ctx context.Context, id, userID string) error {
	query := `
		UPDATE coach_invitations
		SET 
			status = 'accepted',
			accepted_at = CURRENT_TIMESTAMP,
			accepted_by_user_id = $1
		WHERE id = $2 AND status = 'pending'
	`

	result, err := s.db.Exec(ctx, query, userID, id)
	if err != nil {
		return fmt.Errorf("failed to accept invitation: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("invitation not found or already processed")
	}

	return nil
}

func (s *Store) DeleteInvitation(ctx context.Context, id string) error {
	query := `DELETE FROM coach_invitations WHERE id = $1`

	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete invitation: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("invitation not found")
	}

	return nil
}

func (s *Store) ExpireOldInvitations(ctx context.Context) (int64, error) {
	query := `
		UPDATE coach_invitations
		SET status = 'expired'
		WHERE status = 'pending' AND expires_at < CURRENT_TIMESTAMP
	`

	result, err := s.db.Exec(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to expire invitations: %w", err)
	}

	return result.RowsAffected(), nil
}

func (s *Store) GetInvitationByCoachAndEmail(ctx context.Context, coachID, email string) (*CoachInvitation, error) {
	query := `
		SELECT 
			id, coach_id, email, first_name, last_name,
			invitation_token, status, custom_message, expires_at,
			created_at, accepted_at, accepted_by_user_id
		FROM coach_invitations
		WHERE coach_id = $1 AND LOWER(email) = LOWER($2) AND status = 'pending'
		ORDER BY created_at DESC
		LIMIT 1
	`

	var inv CoachInvitation
	err := s.db.QueryRow(ctx, query, coachID, email).Scan(
		&inv.ID,
		&inv.CoachID,
		&inv.Email,
		&inv.FirstName,
		&inv.LastName,
		&inv.InvitationToken,
		&inv.Status,
		&inv.CustomMessage,
		&inv.ExpiresAt,
		&inv.CreatedAt,
		&inv.AcceptedAt,
		&inv.AcceptedByUserID,
	)

	if err != nil {
		if err == pgx.ErrNoRows || err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get invitation: %w", err)
	}

	return &inv, nil
}

