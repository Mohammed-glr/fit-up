package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tdmdh/fit-up-server/internal/auth/types"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	query := `
		SELECT id, username, name, bio, email, image, password, role, is_two_factor_enabled, created_at, updated_at
		FROM users 
		WHERE email = $1
	`

	var user types.User
	err := s.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Name,
		&user.Bio,
		&user.Email,
		&user.Image,
		&user.PasswordHash,
		&user.Role,
		&user.IsTwoFactorEnabled,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, types.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *Store) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	query := `
		SELECT id, username, name, bio, email, image, password, role, is_two_factor_enabled, created_at, updated_at
		FROM users 
		WHERE id = $1
	`

	var user types.User
	err := s.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Name,
		&user.Bio,
		&user.Email,
		&user.Image,
		&user.PasswordHash,
		&user.Role,
		&user.IsTwoFactorEnabled,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, types.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *Store) GetUserByUsername(ctx context.Context, username string) (*types.User, error) {
	query := `
		SELECT id, username, name, bio, email, image, password, role, is_two_factor_enabled, created_at, updated_at
		FROM users 
		WHERE username = $1
	`

	var user types.User
	err := s.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Name,
		&user.Bio,
		&user.Email,
		&user.Image,
		&user.PasswordHash,
		&user.Role,
		&user.IsTwoFactorEnabled,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, types.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *Store) CreateUser(ctx context.Context, user *types.User) error {
	query := `
		INSERT INTO users (id, username, name, bio, email, image, password, role, is_two_factor_enabled)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := s.db.Exec(ctx, query,
		user.ID,
		user.Username,
		user.Name,
		user.Bio,
		user.Email,
		user.Image,
		user.PasswordHash,
		user.Role,
		user.IsTwoFactorEnabled,
	)

	return err
}

func (s *Store) UpdateUser(ctx context.Context, id string, updates *types.UpdateUserRequest) error {
	query := `
		UPDATE users 
		SET name = $2, bio = $3, updated_at = NOW()
		WHERE id = $1
	`

	_, err := s.db.Exec(ctx, query, id, updates.Name, updates.Bio)
	return err
}

func (s *Store) UpdateUserPassword(ctx context.Context, userID string, hashedPassword string) error {
	query := `
		UPDATE users 
		SET password = $2, updated_at = NOW()
		WHERE id = $1
	`

	_, err := s.db.Exec(ctx, query, userID, hashedPassword)
	return err
}

func (s *Store) CreatePasswordResetToken(ctx context.Context, email string, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO password_reset_tokens (email, token, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) 
		DO UPDATE SET token = $2, expires_at = $3, created_at = NOW()
	`

	_, err := s.db.Exec(ctx, query, email, token, expiresAt)
	return err
}

func (s *Store) GetPasswordResetToken(ctx context.Context, token string) (*types.PasswordResetToken, error) {
	query := `
		SELECT email, token, expires_at, used
		FROM password_reset_tokens 
		WHERE token = $1
	`

	var resetToken types.PasswordResetToken
	err := s.db.QueryRow(ctx, query, token).Scan(
		&resetToken.Email,
		&resetToken.Token,
		&resetToken.Expires,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, types.ErrPasswordResetTokenNotFound
		}
		return nil, err
	}

	return &resetToken, nil
}

func (s *Store) GetUserByPasswordResetToken(ctx context.Context, token string) (*types.User, error) {
	query := `
		SELECT u.id, u.username, u.name, u.bio, u.email, u.image, u.password, u.role, u.is_two_factor_enabled, u.created_at, u.updated_at
		FROM users u
		INNER JOIN password_reset_tokens prt ON u.email = prt.email
		WHERE prt.token = $1 AND prt.expires_at > NOW() AND prt.used = false
	`

	var user types.User
	err := s.db.QueryRow(ctx, query, token).Scan(
		&user.ID,
		&user.Username,
		&user.Name,
		&user.Bio,
		&user.Email,
		&user.Image,
		&user.PasswordHash,
		&user.Role,
		&user.IsTwoFactorEnabled,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, types.ErrPasswordResetTokenNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *Store) DeletePasswordResetToken(ctx context.Context, token string) error {
	query := `DELETE FROM password_reset_tokens WHERE token = $1`
	_, err := s.db.Exec(ctx, query, token)
	return err
}

func (s *Store) MarkPasswordResetTokenAsUsed(ctx context.Context, token string) error {
	query := `
		UPDATE password_reset_tokens 
		SET used = true 
		WHERE token = $1
	`

	_, err := s.db.Exec(ctx, query, token)
	return err
}

// Refresh Token methods
func (s *Store) CreateRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time, accessTokenJTI string) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, access_token_jti, expires_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := s.db.Exec(ctx, query, userID, token, accessTokenJTI, expiresAt)
	return err
}

func (s *Store) GetRefreshToken(ctx context.Context, token string) (*types.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, access_token_jti, expires_at, created_at, last_used_at, is_revoked, revoked_at
		FROM refresh_tokens 
		WHERE token_hash = $1
	`

	var refreshToken types.RefreshToken
	err := s.db.QueryRow(ctx, query, token).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.TokenHash,
		&refreshToken.AccessTokenJTI,
		&refreshToken.ExpiresAt,
		&refreshToken.CreatedAt,
		&refreshToken.LastUsedAt,
		&refreshToken.IsRevoked,
		&refreshToken.RevokedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, types.ErrRefreshTokenNotFound
		}
		return nil, err
	}

	return &refreshToken, nil
}

func (s *Store) DeleteRefreshToken(ctx context.Context, token string) error {
	query := `DELETE FROM refresh_tokens WHERE token_hash = $1`
	_, err := s.db.Exec(ctx, query, token)
	return err
}

func (s *Store) CleanupExpiredRefreshTokens(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`
	_, err := s.db.Exec(ctx, query)
	return err
}

func (s *Store) RevokeRefreshToken(ctx context.Context, token string) error {
	query := `
		UPDATE refresh_tokens 
		SET is_revoked = true, revoked_at = NOW() 
		WHERE token_hash = $1
	`

	_, err := s.db.Exec(ctx, query, token)
	return err
}

func (s *Store) RevokeAllUserRefreshTokens(ctx context.Context, userID string) error {
	query := `
		UPDATE refresh_tokens 
		SET is_revoked = true, revoked_at = NOW() 
		WHERE user_id = $1 AND is_revoked = false
	`

	_, err := s.db.Exec(ctx, query, userID)
	return err
}

func (s *Store) UpdateRefreshTokenLastUsed(ctx context.Context, token string) error {
	query := `
		UPDATE refresh_tokens 
		SET last_used_at = NOW() 
		WHERE token_hash = $1
	`

	_, err := s.db.Exec(ctx, query, token)
	return err
}
