package repository

import (
	"context"

	"github.com/tdmdh/fit-up-server/services/schema-service/internal/types"
)


 

func (s *Store) CreateUser(ctx context.Context, user *types.WorkoutUserRequest) (*types.WorkoutUser, error) {
	q := ` 
		INSERT INTO users (email, name, equipment, frequency, goal, level)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, email, equipment, frequency, goal, level, name, id AS user_id
	`

	var createdUser types.WorkoutUser
	err := s.db.QueryRow(ctx, q,
		user.Email,
		user.Name,
		user.Equipment,
		user.Frequency,
		user.Goal,
		user.Level,
	).Scan(
		&createdUser.CreatedAt,
		&createdUser.Email,
		&createdUser.Equipment,
		&createdUser.Frequency,
		&createdUser.Goal,
		&createdUser.Level,
		&createdUser.Name,
		&createdUser.UserID,
	)

	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (s *Store)	GetUserByID(ctx context.Context, userID int) (*types.WorkoutUser, error) {
	q := `
		SELECT created_at, email, equipment, frequency, goal, level, name, id AS user_id
		FROM users
		WHERE id = $1
	`

	var user types.WorkoutUser
	err := s.db.QueryRow(ctx, q, userID).Scan(
		&user.CreatedAt,
		&user.Email,
		&user.Equipment,
		&user.Frequency,
		&user.Goal,
		&user.Level,
		&user.Name,
		&user.UserID,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, types.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*types.WorkoutUser, error) {
	q := `
		SELECT created_at, email, equipment, frequency, goal, level, name, id AS user_id
		FROM users
		WHERE email = $1
	`

	var user types.WorkoutUser
	err := s.db.QueryRow(ctx, q, email).Scan(
		&user.CreatedAt,
		&user.Email,
		&user.Equipment,
		&user.Frequency,
		&user.Goal,
		&user.Level,
		&user.Name,
		&user.UserID,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, types.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}


func (s *Store) UpdateUser(ctx context.Context, userID int, user *types.WorkoutUserRequest) (*types.WorkoutUser, error) {
	q := `
		UPDATE users
		SET email = $1, name = $2, equipment = $3, frequency = $4, goal = $5, level = $6
		WHERE id = $7
		RETURNING created_at, email, equipment, frequency, goal, level, name, id AS user_id
	`

	var updatedUser types.WorkoutUser
	err := s.db.QueryRow(ctx, q,
		user.Email,
		user.Name,
		user.Equipment,
		user.Frequency,
		user.Goal,
		user.Level,
		userID,
	).Scan(
		&updatedUser.CreatedAt,
		&updatedUser.Email,
		&updatedUser.Equipment,
		&updatedUser.Frequency,
		&updatedUser.Goal,
		&updatedUser.Level,
		&updatedUser.Name,
		&updatedUser.UserID,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, types.ErrUserNotFound
		}
		return nil, err
	}

	return &updatedUser, nil
}

func (s *Store) DeleteUser(ctx context.Context, userID int) error {
	q := `
		DELETE FROM users
		WHERE id = $1
	`

	cmdTag, err := s.db.Exec(ctx, q, userID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return types.ErrUserNotFound
	}

	return nil
}

func (s *Store) ListUsers(ctx context.Context, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutUser], error) {
	q := `
		SELECT created_at, email, equipment, frequency, goal, level, name, id AS user_id
		FROM users
		ORDER BY id
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Query(ctx, q, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []types.WorkoutUser
	for rows.Next() {
		var user types.WorkoutUser
		err := rows.Scan(
			&user.CreatedAt,
			&user.Email,
			&user.Equipment,
			&user.Frequency,
			&user.Goal,
			&user.Level,
			&user.Name,
			&user.UserID,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	var total int
	countQuery := `SELECT COUNT(*) FROM users`
	err = s.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, err
	}


	totalPages := (total + pagination.Limit - 1) / pagination.Limit
	page := (pagination.Offset / pagination.Limit) + 1
	pageSize := pagination.Limit


	return &types.PaginatedResponse[types.WorkoutUser]{
		Data:       users,
		TotalCount: total,
		TotalPages: totalPages,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}


func (s *Store) SearchUsers(ctx context.Context, query string, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutUser], error) {
	q := `	
		SELECT created_at, email, equipment, frequency, goal, level, name, id AS user_id
		FROM users
		WHERE name ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%'
		ORDER BY id
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.Query(ctx, q, query, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []types.WorkoutUser
	for rows.Next() {
		var user types.WorkoutUser
		err := rows.Scan(
			&user.CreatedAt,
			&user.Email,
			&user.Equipment,
			&user.Frequency,
			&user.Goal,
			&user.Level,
			&user.Name,
			&user.UserID,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM users
		WHERE name ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%'
	`
	err = s.db.QueryRow(ctx, countQuery, query).Scan(&total)
	if err != nil {
		return nil, err
	}

	totalPages := (total + pagination.Limit - 1) / pagination.Limit
	page := (pagination.Offset / pagination.Limit) + 1
	pageSize := pagination.Limit

	return &types.PaginatedResponse[types.WorkoutUser]{
		Data:       users,
		TotalCount: total,
		TotalPages: totalPages,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *Store) GetUsersByLevel(ctx context.Context, level types.FitnessLevel) ([]types.WorkoutUser, error) {
	q := `
		SELECT created_at, email, equipment, frequency, goal, level, name, id AS user_id
		FROM users
		WHERE level = $1
	`



	rows, err := s.db.Query(ctx, q, level)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []types.WorkoutUser
	for rows.Next() {
		var user types.WorkoutUser
		err := rows.Scan(
			&user.CreatedAt,
			&user.Email,
			&user.Equipment,
			&user.Frequency,
			&user.Goal,
			&user.Level,
			&user.Name,
			&user.UserID,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return users, nil
}


func (s *Store) GetUsersByGoal(ctx context.Context, goal types.FitnessGoal) ([]types.WorkoutUser, error) {
	q := `
		SELECT created_at, email, equipment, frequency, goal, level, name, id AS user_id
		FROM users
		WHERE goal = $1
	`

	rows, err := s.db.Query(ctx, q, goal)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []types.WorkoutUser
	for rows.Next() {
		var user types.WorkoutUser
		err := rows.Scan(
			&user.CreatedAt,
			&user.Email,
			&user.Equipment,
			&user.Frequency,
			&user.Goal,
			&user.Level,
			&user.Name,
			&user.UserID,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return users, nil
}


func (s *Store) CountActiveUsers(ctx context.Context) (int, error) {
	q := `
		SELECT COUNT(*)
		FROM users
		WHERE frequency > 0
	`

	var count int
	err := s.db.QueryRow(ctx, q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

