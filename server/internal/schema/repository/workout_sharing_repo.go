package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

type WorkoutSharingRepo interface {
	GetWorkoutShareSummary(ctx context.Context, sessionID int, userID string) (*types.WorkoutShareSummary, error)
}

// GetWorkoutShareSummary retrieves a complete workout summary for sharing
func (s *Store) GetWorkoutShareSummary(ctx context.Context, sessionID int, userID string) (*types.WorkoutShareSummary, error) {
	// First, verify the session belongs to the user and is completed
	var summary types.WorkoutShareSummary

	sessionQuery := `
		SELECT 
			ws.session_id,
			COALESCE(w.workout_name, w.focus) as workout_title,
			ws.end_time as completed_at,
			COALESCE(sm.duration_seconds, 0) as duration_seconds,
			ws.total_exercises,
			ws.completed_exercises,
			COALESCE(ws.total_volume, 0) as total_volume
		FROM workout_sessions ws
		JOIN workouts w ON ws.workout_id = w.workout_id
		LEFT JOIN session_metrics sm ON ws.session_id = sm.session_id
		WHERE ws.session_id = $1 
			AND ws.user_id = $2 
			AND ws.status = 'completed'
	`

	var durationSeconds int

	err := s.db.QueryRow(ctx, sessionQuery, sessionID, userID).Scan(
		&summary.SessionID,
		&summary.WorkoutTitle,
		&summary.CompletedAt,
		&durationSeconds,
		&summary.TotalExercises,
		&summary.TotalExercises, // Using completed_exercises
		&summary.TotalVolumeLbs,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("workout session %d not found, doesn't belong to user, or is not completed", sessionID)
		}
		return nil, fmt.Errorf("failed to get workout session: %w", err)
	}

	summary.DurationMinutes = durationSeconds / 60

	// Get exercise performances
	exercisesQuery := `
		SELECT 
			e.name as exercise_name,
			ep.sets_completed,
			COALESCE(ep.total_volume, 0) as total_volume,
			COUNT(sp.set_id) as total_reps
		FROM exercise_performances ep
		JOIN exercises e ON ep.exercise_id = e.exercise_id
		LEFT JOIN set_performances sp ON ep.performance_id = sp.performance_id
		WHERE ep.session_id = $1
		GROUP BY e.name, ep.sets_completed, ep.total_volume, ep.performance_id
		ORDER BY ep.performance_id
	`

	rows, err := s.db.Query(ctx, exercisesQuery, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exercises: %w", err)
	}
	defer rows.Close()

	summary.Exercises = make([]types.WorkoutShareExercise, 0)
	summary.TotalSets = 0
	summary.TotalReps = 0

	for rows.Next() {
		var exercise types.WorkoutShareExercise
		err := rows.Scan(
			&exercise.ExerciseName,
			&exercise.SetsCompleted,
			&exercise.TotalVolumeLbs,
			&exercise.TotalReps,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan exercise: %w", err)
		}

		summary.TotalSets += exercise.SetsCompleted
		summary.TotalReps += exercise.TotalReps

		// Get best set for this exercise
		bestSetQuery := `
			SELECT weight, reps
			FROM set_performances sp
			JOIN exercise_performances ep ON sp.performance_id = ep.performance_id
			WHERE ep.session_id = $1 
				AND ep.exercise_id = (SELECT exercise_id FROM exercises WHERE name = $2 LIMIT 1)
			ORDER BY (weight * reps) DESC
			LIMIT 1
		`

		var bestSet types.BestSet
		err = s.db.QueryRow(ctx, bestSetQuery, sessionID, exercise.ExerciseName).Scan(
			&bestSet.Weight,
			&bestSet.Reps,
		)
		if err == nil {
			exercise.BestSet = &bestSet
		}

		// TODO: Check if PR was achieved (requires historical data comparison)
		exercise.PRAchieved = false

		summary.Exercises = append(summary.Exercises, exercise)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating exercises: %w", err)
	}

	// Get user info
	userQuery := `
		SELECT 
			COALESCE(first_name || ' ' || last_name, email) as name,
			profile_photo_url
		FROM users
		WHERE id = $1
	`

	var photoURL *string
	err = s.db.QueryRow(ctx, userQuery, userID).Scan(
		&summary.UserName,
		&photoURL,
	)
	if err == nil && photoURL != nil {
		summary.UserPhotoURL = *photoURL
	}

	// Count PRs achieved (simplified - just count exercises with best performance)
	summary.PRsAchieved = 0

	return &summary, nil
}
