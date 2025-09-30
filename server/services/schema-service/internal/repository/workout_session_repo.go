package repository

import (
	"context"
	"time"

	"github.com/tdmdh/fit-up-server/services/schema-service/internal/types"
)

// =============================================================================
// WORKOUT SESSION REPOSITORY IMPLEMENTATION
// =============================================================================

func (s *Store) StartWorkoutSession(ctx context.Context, userID int, workoutID int) (*types.WorkoutSession, error) {
	// First check if there's already an active session
	checkQuery := `
		SELECT session_id 
		FROM workout_sessions 
		WHERE user_id = $1 AND status = 'active'
	`

	var existingSessionID int
	err := s.db.QueryRow(ctx, checkQuery, userID).Scan(&existingSessionID)
	if err == nil {
		// Active session exists, return error or end previous session
		endQuery := `
			UPDATE workout_sessions 
			SET status = 'abandoned', end_time = NOW() 
			WHERE session_id = $1
		`
		_, _ = s.db.Exec(ctx, endQuery, existingSessionID)
	}

	// Get total exercises count for this workout
	exerciseCountQuery := `
		SELECT COUNT(*) 
		FROM workout_exercises 
		WHERE workout_id = $1
	`

	var totalExercises int
	err = s.db.QueryRow(ctx, exerciseCountQuery, workoutID).Scan(&totalExercises)
	if err != nil {
		totalExercises = 0
	}

	// Create new session
	q := `
		INSERT INTO workout_sessions (user_id, workout_id, start_time, status, total_exercises, completed_exercises, total_volume, notes)
		VALUES ($1, $2, NOW(), 'active', $3, 0, 0.0, '')
		RETURNING session_id, user_id, workout_id, start_time, end_time, status, total_exercises, completed_exercises, total_volume, notes
	`

	var session types.WorkoutSession
	err = s.db.QueryRow(ctx, q,
		userID,
		workoutID,
		totalExercises,
	).Scan(
		&session.SessionID,
		&session.UserID,
		&session.WorkoutID,
		&session.StartTime,
		&session.EndTime,
		&session.Status,
		&session.TotalExercises,
		&session.CompletedExercises,
		&session.TotalVolume,
		&session.Notes,
	)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *Store) CompleteWorkoutSession(ctx context.Context, sessionID int, summary *types.SessionSummary) (*types.WorkoutSession, error) {
	q := `
		UPDATE workout_sessions 
		SET 
			end_time = NOW(),
			status = 'completed',
			completed_exercises = $1,
			total_volume = $2,
			notes = $3
		WHERE session_id = $4
		RETURNING session_id, user_id, workout_id, start_time, end_time, status, total_exercises, completed_exercises, total_volume, notes
	`

	var session types.WorkoutSession
	err := s.db.QueryRow(ctx, q,
		summary.ExercisesCompleted,
		summary.TotalVolume,
		summary.Notes,
		sessionID,
	).Scan(
		&session.SessionID,
		&session.UserID,
		&session.WorkoutID,
		&session.StartTime,
		&session.EndTime,
		&session.Status,
		&session.TotalExercises,
		&session.CompletedExercises,
		&session.TotalVolume,
		&session.Notes,
	)

	if err != nil {
		return nil, err
	}

	// Log individual exercise performances
	if len(summary.Exercises) > 0 {
		exerciseQuery := `
			INSERT INTO session_exercise_performances (session_id, exercise_id, sets_completed, best_reps, best_weight, total_volume, rpe, notes)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`

		for _, exercise := range summary.Exercises {
			_, err = s.db.Exec(ctx, exerciseQuery,
				sessionID,
				exercise.ExerciseID,
				exercise.SetsCompleted,
				exercise.BestSet.Reps,
				exercise.BestSet.Weight,
				exercise.TotalVolume,
				exercise.RPE,
				exercise.Notes,
			)
			// Continue even if individual exercise logging fails
		}
	}

	return &session, nil
}

func (s *Store) SkipWorkout(ctx context.Context, userID int, workoutID int, reason string) (*types.SkippedWorkout, error) {
	q := `
		INSERT INTO skipped_workouts (user_id, workout_id, skip_date, reason)
		VALUES ($1, $2, NOW(), $3)
		RETURNING skip_id, user_id, workout_id, skip_date, reason
	`

	var skipped types.SkippedWorkout
	err := s.db.QueryRow(ctx, q,
		userID,
		workoutID,
		reason,
	).Scan(
		&skipped.SkipID,
		&skipped.UserID,
		&skipped.WorkoutID,
		&skipped.SkipDate,
		&skipped.Reason,
	)

	if err != nil {
		return nil, err
	}

	return &skipped, nil
}

func (s *Store) LogExercisePerformance(ctx context.Context, sessionID int, exerciseID int, performance *types.ExercisePerformance) error {
	// Check if performance already exists for this session and exercise
	checkQuery := `
		SELECT performance_id 
		FROM session_exercise_performances 
		WHERE session_id = $1 AND exercise_id = $2
	`

	var existingID int
	err := s.db.QueryRow(ctx, checkQuery, sessionID, exerciseID).Scan(&existingID)

	if err == nil {
		// Update existing performance
		updateQuery := `
			UPDATE session_exercise_performances 
			SET sets_completed = $1, best_reps = $2, best_weight = $3, total_volume = $4, rpe = $5, notes = $6
			WHERE performance_id = $7
		`
		_, err = s.db.Exec(ctx, updateQuery,
			performance.SetsCompleted,
			performance.BestSet.Reps,
			performance.BestSet.Weight,
			performance.TotalVolume,
			performance.RPE,
			performance.Notes,
			existingID,
		)
		return err
	}

	// Insert new performance
	insertQuery := `
		INSERT INTO session_exercise_performances (session_id, exercise_id, sets_completed, best_reps, best_weight, total_volume, rpe, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err = s.db.Exec(ctx, insertQuery,
		sessionID,
		exerciseID,
		performance.SetsCompleted,
		performance.BestSet.Reps,
		performance.BestSet.Weight,
		performance.TotalVolume,
		performance.RPE,
		performance.Notes,
	)

	return err
}

func (s *Store) GetActiveSession(ctx context.Context, userID int) (*types.WorkoutSession, error) {
	q := `
		SELECT session_id, user_id, workout_id, start_time, end_time, status, total_exercises, completed_exercises, total_volume, notes
		FROM workout_sessions
		WHERE user_id = $1 AND status = 'active'
		ORDER BY start_time DESC
		LIMIT 1
	`

	var session types.WorkoutSession
	err := s.db.QueryRow(ctx, q, userID).Scan(
		&session.SessionID,
		&session.UserID,
		&session.WorkoutID,
		&session.StartTime,
		&session.EndTime,
		&session.Status,
		&session.TotalExercises,
		&session.CompletedExercises,
		&session.TotalVolume,
		&session.Notes,
	)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *Store) GetSessionHistory(ctx context.Context, userID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutSession], error) {
	// Count total sessions
	countQuery := `
		SELECT COUNT(*) 
		FROM workout_sessions 
		WHERE user_id = $1 AND status IN ('completed', 'abandoned')
	`

	var totalCount int
	err := s.db.QueryRow(ctx, countQuery, userID).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	// Get sessions with pagination
	q := `
		SELECT session_id, user_id, workout_id, start_time, end_time, status, total_exercises, completed_exercises, total_volume, notes
		FROM workout_sessions
		WHERE user_id = $1 AND status IN ('completed', 'abandoned')
		ORDER BY start_time DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.Query(ctx, q, userID, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []types.WorkoutSession
	for rows.Next() {
		var session types.WorkoutSession
		err := rows.Scan(
			&session.SessionID,
			&session.UserID,
			&session.WorkoutID,
			&session.StartTime,
			&session.EndTime,
			&session.Status,
			&session.TotalExercises,
			&session.CompletedExercises,
			&session.TotalVolume,
			&session.Notes,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	// Calculate pagination info
	totalPages := (totalCount + pagination.PageSize - 1) / pagination.PageSize

	return &types.PaginatedResponse[types.WorkoutSession]{
		Data:       sessions,
		TotalCount: totalCount,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *Store) GetSessionMetrics(ctx context.Context, sessionID int) (*types.SessionMetrics, error) {
	// Get basic session info
	sessionQuery := `
		SELECT 
			EXTRACT(EPOCH FROM (COALESCE(end_time, NOW()) - start_time))::int as duration,
			total_volume,
			completed_exercises,
			total_exercises
		FROM workout_sessions
		WHERE session_id = $1
	`

	var metrics types.SessionMetrics
	var completedExercises, totalExercises int

	err := s.db.QueryRow(ctx, sessionQuery, sessionID).Scan(
		&metrics.Duration,
		&metrics.TotalVolume,
		&completedExercises,
		&totalExercises,
	)

	if err != nil {
		return nil, err
	}

	metrics.SessionID = sessionID

	// Calculate completion rate
	if totalExercises > 0 {
		metrics.CompletionRate = float64(completedExercises) / float64(totalExercises) * 100
	}

	// Get average RPE from exercise performances
	rpeQuery := `
		SELECT AVG(rpe) 
		FROM session_exercise_performances 
		WHERE session_id = $1 AND rpe > 0
	`

	err = s.db.QueryRow(ctx, rpeQuery, sessionID).Scan(&metrics.RPE)
	if err != nil {
		metrics.RPE = 0 // Default if no RPE data
	}

	// Estimate calories (rough calculation: 5 calories per kg of volume)
	metrics.CaloriesBurned = int(metrics.TotalVolume * 5)

	// Calculate average intensity (percentage of estimated 1RM)
	// This is a simplified calculation
	metrics.AverageIntensity = 0.75 // Default moderate intensity

	return &metrics, nil
}

func (s *Store) GetWeeklySessionStats(ctx context.Context, userID int, weekStart time.Time) (*types.WeeklySessionStats, error) {
	weekEnd := weekStart.AddDate(0, 0, 7)

	// Get planned sessions for the week
	plannedQuery := `
		SELECT COUNT(DISTINCT w.workout_id)
		FROM weekly_schemas ws
		JOIN workouts w ON ws.schema_id = w.schema_id
		WHERE ws.user_id = $1 AND ws.week_start = $2 AND ws.active = true
	`

	var sessionsPlanned int
	err := s.db.QueryRow(ctx, plannedQuery, userID, weekStart).Scan(&sessionsPlanned)
	if err != nil {
		sessionsPlanned = 0
	}

	// Get completed sessions for the week
	statsQuery := `
		SELECT 
			COUNT(*) as sessions_completed,
			COALESCE(SUM(total_volume), 0) as total_volume,
			COALESCE(AVG(rpe.avg_rpe), 0) as average_rpe
		FROM workout_sessions ws
		LEFT JOIN (
			SELECT session_id, AVG(rpe) as avg_rpe
			FROM session_exercise_performances
			WHERE rpe > 0
			GROUP BY session_id
		) rpe ON ws.session_id = rpe.session_id
		WHERE ws.user_id = $1 
		AND ws.start_time >= $2 
		AND ws.start_time < $3 
		AND ws.status = 'completed'
	`

	var stats types.WeeklySessionStats
	stats.WeekStart = weekStart
	stats.SessionsPlanned = sessionsPlanned

	err = s.db.QueryRow(ctx, statsQuery, userID, weekStart, weekEnd).Scan(
		&stats.SessionsCompleted,
		&stats.TotalVolume,
		&stats.AverageRPE,
	)

	if err != nil {
		return nil, err
	}

	// Calculate completion rate
	if sessionsPlanned > 0 {
		stats.CompletionRate = float64(stats.SessionsCompleted) / float64(sessionsPlanned) * 100
	}

	return &stats, nil
}
