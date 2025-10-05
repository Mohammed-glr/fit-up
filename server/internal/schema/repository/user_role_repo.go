package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

// GetUserRole haalt de rol van een gebruiker op uit de cache
func (s *Store) GetUserRole(ctx context.Context, authUserID string) (types.UserRole, error) {
	query := `
		SELECT role 
		FROM user_roles_cache 
		WHERE auth_user_id = $1
	`

	var role types.UserRole
	err := s.db.QueryRow(ctx, query, authUserID).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			// Als niet in cache, default naar 'user'
			return types.RoleUser, nil
		}
		return "", err
	}

	return role, nil
}

// UpsertUserRole update of insert een user role in de cache
func (s *Store) UpsertUserRole(ctx context.Context, authUserID string, role types.UserRole) error {
	query := `
		INSERT INTO user_roles_cache (auth_user_id, role, last_synced_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (auth_user_id) 
		DO UPDATE SET 
			role = EXCLUDED.role,
			last_synced_at = EXCLUDED.last_synced_at
	`

	_, err := s.db.Exec(ctx, query, authUserID, role, time.Now())
	return err
}

// BatchUpsertUserRoles voegt meerdere user roles toe in één keer
func (s *Store) BatchUpsertUserRoles(ctx context.Context, roles map[string]types.UserRole) error {
	if len(roles) == 0 {
		return nil
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO user_roles_cache (auth_user_id, role, last_synced_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (auth_user_id) 
		DO UPDATE SET 
			role = EXCLUDED.role,
			last_synced_at = EXCLUDED.last_synced_at
	`

	now := time.Now()
	for authUserID, role := range roles {
		_, err := tx.Exec(ctx, query, authUserID, role, now)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// DeleteUserRole verwijdert een user role uit de cache
func (s *Store) DeleteUserRole(ctx context.Context, authUserID string) error {
	query := `DELETE FROM user_roles_cache WHERE auth_user_id = $1`
	_, err := s.db.Exec(ctx, query, authUserID)
	return err
}

// GetStaleRoles haalt rollen op die langer dan de opgegeven duur niet zijn gesynchroniseerd
func (s *Store) GetStaleRoles(ctx context.Context, staleDuration time.Duration) ([]types.UserRoleCache, error) {
	query := `
		SELECT auth_user_id, role, last_synced_at
		FROM user_roles_cache
		WHERE last_synced_at < $1
	`

	staleTime := time.Now().Add(-staleDuration)
	rows, err := s.db.Query(ctx, query, staleTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var staleRoles []types.UserRoleCache
	for rows.Next() {
		var roleCache types.UserRoleCache
		err := rows.Scan(
			&roleCache.AuthUserID,
			&roleCache.Role,
			&roleCache.LastSyncedAt,
		)
		if err != nil {
			return nil, err
		}
		staleRoles = append(staleRoles, roleCache)
	}

	return staleRoles, rows.Err()
}
