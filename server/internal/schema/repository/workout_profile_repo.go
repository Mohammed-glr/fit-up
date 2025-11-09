package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

func (s *Store) CreateWorkoutProfile(ctx context.Context, authUserID string, profile *types.WorkoutProfileRequest) (*types.WorkoutProfile, error) {
	q := ` 
INSERT INTO workout_profiles (auth_user_id, level, goal, frequency, equipment)
VALUES ($1, $2, $3, $4, $5)
RETURNING workout_profile_id, auth_user_id, level, goal, frequency, equipment, created_at
`

	equipmentJSON, err := json.Marshal(profile.Equipment)
	if err != nil {
		return nil, err
	}

	var createdProfile types.WorkoutProfile
	err = s.db.QueryRow(ctx, q,
		authUserID,
		profile.Level,
		profile.Goal,
		profile.Frequency,
		equipmentJSON,
	).Scan(
		&createdProfile.WorkoutProfileID,
		&createdProfile.AuthUserID,
		&createdProfile.Level,
		&createdProfile.Goal,
		&createdProfile.Frequency,
		&createdProfile.Equipment,
		&createdProfile.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &createdProfile, nil
}

func (s *Store) GetWorkoutProfileByAuthID(ctx context.Context, authUserID string) (*types.WorkoutProfile, error) {
	q := `
SELECT workout_profile_id, auth_user_id, level, goal, frequency, equipment, created_at
FROM workout_profiles
WHERE auth_user_id = $1
`

	var profile types.WorkoutProfile
	err := s.db.QueryRow(ctx, q, authUserID).Scan(
		&profile.WorkoutProfileID,
		&profile.AuthUserID,
		&profile.Level,
		&profile.Goal,
		&profile.Frequency,
		&profile.Equipment,
		&profile.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (s *Store) GetWorkoutProfileByUsername(ctx context.Context, username string) (*types.WorkoutProfile, error) {
	q := `
SELECT wp.workout_profile_id, wp.auth_user_id, wp.level, wp.goal, wp.frequency, wp.equipment, wp.created_at
FROM workout_profiles wp
JOIN users u ON u.id = wp.auth_user_id
WHERE LOWER(u.username) = LOWER($1)
`

	var profile types.WorkoutProfile
	err := s.db.QueryRow(ctx, q, username).Scan(
		&profile.WorkoutProfileID,
		&profile.AuthUserID,
		&profile.Level,
		&profile.Goal,
		&profile.Frequency,
		&profile.Equipment,
		&profile.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (s *Store) GetWorkoutProfileByID(ctx context.Context, workoutProfileID int) (*types.WorkoutProfile, error) {
	q := `
SELECT workout_profile_id, auth_user_id, level, goal, frequency, equipment, created_at
FROM workout_profiles
WHERE workout_profile_id = $1
`

	var profile types.WorkoutProfile
	err := s.db.QueryRow(ctx, q, workoutProfileID).Scan(
		&profile.WorkoutProfileID,
		&profile.AuthUserID,
		&profile.Level,
		&profile.Goal,
		&profile.Frequency,
		&profile.Equipment,
		&profile.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (s *Store) UpdateWorkoutProfile(ctx context.Context, authUserID string, profile *types.WorkoutProfileRequest) (*types.WorkoutProfile, error) {
	equipmentJSON, err := json.Marshal(profile.Equipment)
	if err != nil {
		return nil, err
	}

	q := `
UPDATE workout_profiles 
SET level = $2, goal = $3, frequency = $4, equipment = $5
WHERE auth_user_id = $1
RETURNING workout_profile_id, auth_user_id, level, goal, frequency, equipment, created_at
`

	var updatedProfile types.WorkoutProfile
	err = s.db.QueryRow(ctx, q,
		authUserID,
		profile.Level,
		profile.Goal,
		profile.Frequency,
		equipmentJSON,
	).Scan(
		&updatedProfile.WorkoutProfileID,
		&updatedProfile.AuthUserID,
		&updatedProfile.Level,
		&updatedProfile.Goal,
		&updatedProfile.Frequency,
		&updatedProfile.Equipment,
		&updatedProfile.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &updatedProfile, nil
}

func (s *Store) DeleteWorkoutProfile(ctx context.Context, authUserID string) error {
	q := `DELETE FROM workout_profiles WHERE auth_user_id = $1`

	_, err := s.db.Exec(ctx, q, authUserID)
	return err
}

func (s *Store) ListWorkoutProfiles(ctx context.Context, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutProfile], error) {
	// Count total records
	countQuery := `SELECT COUNT(*) FROM workout_profiles`
	var total int
	err := s.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Calculate offset
	offset := (pagination.Page - 1) * pagination.Limit

	// Get profiles with pagination
	q := `
SELECT workout_profile_id, auth_user_id, level, goal, frequency, equipment, created_at
FROM workout_profiles
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

	rows, err := s.db.Query(ctx, q, pagination.Limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []types.WorkoutProfile
	for rows.Next() {
		var profile types.WorkoutProfile
		err := rows.Scan(
			&profile.WorkoutProfileID,
			&profile.AuthUserID,
			&profile.Level,
			&profile.Goal,
			&profile.Frequency,
			&profile.Equipment,
			&profile.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	return &types.PaginatedResponse[types.WorkoutProfile]{
		Data:       profiles,
		TotalCount: total,
		Page:       pagination.Page,
		PageSize:   pagination.Limit,
		TotalPages: (total + pagination.Limit - 1) / pagination.Limit,
	}, nil
}

func (s *Store) SearchWorkoutProfiles(ctx context.Context, query string, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutProfile], error) {
	// Count total matching records
	countQuery := `
SELECT COUNT(*) FROM workout_profiles 
WHERE auth_user_id ILIKE $1 OR level ILIKE $1 OR goal ILIKE $1
`
	searchPattern := "%" + query + "%"
	var total int
	err := s.db.QueryRow(ctx, countQuery, searchPattern).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Calculate offset
	offset := (pagination.Page - 1) * pagination.Limit

	// Get matching profiles with pagination
	q := `
SELECT workout_profile_id, auth_user_id, level, goal, frequency, equipment, created_at
FROM workout_profiles
WHERE auth_user_id ILIKE $1 OR level ILIKE $1 OR goal ILIKE $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`

	rows, err := s.db.Query(ctx, q, searchPattern, pagination.Limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []types.WorkoutProfile
	for rows.Next() {
		var profile types.WorkoutProfile
		err := rows.Scan(
			&profile.WorkoutProfileID,
			&profile.AuthUserID,
			&profile.Level,
			&profile.Goal,
			&profile.Frequency,
			&profile.Equipment,
			&profile.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	return &types.PaginatedResponse[types.WorkoutProfile]{
		Data:       profiles,
		TotalCount: total,
		Page:       pagination.Page,
		PageSize:   pagination.Limit,
		TotalPages: (total + pagination.Limit - 1) / pagination.Limit,
	}, nil
}

func (s *Store) GetProfilesByLevel(ctx context.Context, level types.FitnessLevel) ([]types.WorkoutProfile, error) {
	q := `
SELECT workout_profile_id, auth_user_id, level, goal, frequency, equipment, created_at
FROM workout_profiles
WHERE level = $1
ORDER BY created_at DESC
`

	rows, err := s.db.Query(ctx, q, level)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []types.WorkoutProfile
	for rows.Next() {
		var profile types.WorkoutProfile
		err := rows.Scan(
			&profile.WorkoutProfileID,
			&profile.AuthUserID,
			&profile.Level,
			&profile.Goal,
			&profile.Frequency,
			&profile.Equipment,
			&profile.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

func (s *Store) GetProfilesByGoal(ctx context.Context, goal types.FitnessGoal) ([]types.WorkoutProfile, error) {
	q := `
SELECT workout_profile_id, auth_user_id, level, goal, frequency, equipment, created_at
FROM workout_profiles
WHERE goal = $1
ORDER BY created_at DESC
`

	rows, err := s.db.Query(ctx, q, goal)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []types.WorkoutProfile
	for rows.Next() {
		var profile types.WorkoutProfile
		err := rows.Scan(
			&profile.WorkoutProfileID,
			&profile.AuthUserID,
			&profile.Level,
			&profile.Goal,
			&profile.Frequency,
			&profile.Equipment,
			&profile.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

func (s *Store) CountActiveProfiles(ctx context.Context) (int, error) {
	q := `SELECT COUNT(*) FROM workout_profiles`
	var count int
	err := s.db.QueryRow(ctx, q).Scan(&count)
	return count, err
}

func (s *Store) SearchUsers(ctx context.Context, query string, coachID string, limit int) ([]types.UserSearchResult, error) {
	debugQ := `
			SELECT COUNT(*) as total_users,
				COUNT(CASE WHEN wp.workout_profile_id IS NOT NULL THEN 1 END) as users_with_profiles,
				COUNT(CASE WHEN ca.assignment_id IS NOT NULL THEN 1 END) as users_with_coach
			FROM users u
			LEFT JOIN workout_profiles wp ON u.id = wp.auth_user_id
			LEFT JOIN coach_assignments ca ON ca.user_id = wp.workout_profile_id AND ca.is_active = true
			WHERE LOWER(u.name) LIKE LOWER($1) OR LOWER(u.username) LIKE LOWER($1) OR LOWER(u.email) LIKE LOWER($1)
			`
	searchPattern := fmt.Sprintf("%%%s%%", query)
	var totalUsers, usersWithProfiles, usersWithCoach int
	_ = s.db.QueryRow(ctx, debugQ, searchPattern).Scan(&totalUsers, &usersWithProfiles, &usersWithCoach)
	log.Printf("[SearchUsers DEBUG] Query: '%s' | Total matching users: %d, With profiles: %d, With coach: %d",
		query, totalUsers, usersWithProfiles, usersWithCoach)

	q := `
		SELECT 
			COALESCE(wp.workout_profile_id, 0) as workout_profile_id,
			COALESCE(wp.auth_user_id, u.id) as auth_user_id,
			u.username,
			COALESCE(SPLIT_PART(u.name, ' ', 1), '') as first_name,
			COALESCE(SPLIT_PART(u.name, ' ', 2), '') as last_name,
			u.email,
			COALESCE(wp.level, '') as level,
			COALESCE(wp.goal, '') as goal,
			CASE WHEN ca.assignment_id IS NOT NULL AND ca.is_active THEN true ELSE false END as has_coach,
			ca.coach_id as current_coach_id,
			COALESCE(wp.created_at, u.created_at) as created_at
		FROM users u
		LEFT JOIN workout_profiles wp ON u.id = wp.auth_user_id
		LEFT JOIN coach_assignments ca ON ca.user_id = wp.workout_profile_id AND ca.is_active = true
		WHERE (
			LOWER(u.username) LIKE LOWER($1) OR
			LOWER(u.name) LIKE LOWER($1) OR
			LOWER(u.email) LIKE LOWER($1)
		)
		AND (ca.coach_id IS NULL OR ca.coach_id != $2)
		ORDER BY 
			CASE WHEN ca.coach_id IS NULL THEN 0 ELSE 1 END,
			u.username ASC
		LIMIT $3
		`

	rows, err := s.db.Query(ctx, q, searchPattern, coachID, limit)
	if err != nil {
		log.Printf("[SearchUsers] Query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var results []types.UserSearchResult
	for rows.Next() {
		var result types.UserSearchResult
		err := rows.Scan(
			&result.WorkoutProfileID,
			&result.AuthUserID,
			&result.Username,
			&result.FirstName,
			&result.LastName,
			&result.Email,
			&result.FitnessLevel,
			&result.FitnessGoal,
			&result.HasCoach,
			&result.CurrentCoachID,
			&result.CreatedAt,
		)
		if err != nil {
			log.Printf("[SearchUsers] Scan error: %v", err)
			continue
		}
		results = append(results, result)
	}

	log.Printf("[SearchUsers] Found %d users matching query '%s' for coach %s", len(results), query, coachID)
	return results, nil
}
