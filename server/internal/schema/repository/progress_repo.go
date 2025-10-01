package repository

import (
	"context"
	"time"

	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

func (s *Store) CreateProgressLog(ctx context.Context, progress *types.ProgressLogRequest) (*types.ProgressLog, error) {
	q := `
		INSERT INTO progress_logs (user_id, exercise_id, date, sets_completed, reps_completed, weight_used, duration_seconds)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING log_id, user_id, exercise_id, date, sets_completed, reps_completed, weight_used, duration_seconds
	`

	var log types.ProgressLog
	err := s.db.QueryRow(ctx, q,
		progress.UserID,
		progress.ExerciseID,
		progress.Date,
		progress.SetsCompleted,
		progress.RepsCompleted,
		progress.WeightUsed,
		progress.DurationSeconds,
	).Scan(
		&log.LogID,
		&log.UserID,
		&log.ExerciseID,
		&log.Date,
		&log.SetsCompleted,
		&log.RepsCompleted,
		&log.WeightUsed,
		&log.DurationSeconds,
	)
	if err != nil {
		return nil, err
	}

	return &log, nil
}

func (s *Store) GetProgressLogByID(ctx context.Context, logID int) (*types.ProgressLog, error) {
	q := `
		SELECT log_id, user_id, exercise_id, date, sets_completed, reps_completed, weight_used, duration_seconds
		FROM progress_logs
		WHERE log_id = $1
	`

	var log types.ProgressLog
	err := s.db.QueryRow(ctx, q, logID).Scan(
		&log.LogID,
		&log.UserID,
		&log.ExerciseID,
		&log.Date,
		&log.SetsCompleted,
		&log.RepsCompleted,
		&log.WeightUsed,
		&log.DurationSeconds,
	)
	if err != nil {
		return nil, err
	}

	return &log, nil
}

func (s *Store) UpdateProgressLog(ctx context.Context, logID int, progress *types.ProgressLogRequest) (*types.ProgressLog, error) {
	q := `
		UPDATE progress_logs
		SET user_id = $1, exercise_id = $2, date = $3, sets_completed = $4, reps_completed = $5, weight_used = $6, duration_seconds = $7
		WHERE log_id = $8
		RETURNING log_id, user_id, exercise_id, date, sets_completed, reps_completed, weight_used, duration_seconds
	`

	var log types.ProgressLog
	err := s.db.QueryRow(ctx, q,
		progress.UserID,
		progress.ExerciseID,
		progress.Date,
		progress.SetsCompleted,
		progress.RepsCompleted,
		progress.WeightUsed,
		progress.DurationSeconds,
		logID,
	).Scan(
		&log.LogID,
		&log.UserID,
		&log.ExerciseID,
		&log.Date,
		&log.SetsCompleted,
		&log.RepsCompleted,
		&log.WeightUsed,
		&log.DurationSeconds,
	)
	if err != nil {
		return nil, err
	}

	return &log, nil
}

func (s *Store) DeleteProgressLog(ctx context.Context, logID int) error {
	q := `
		DELETE FROM progress_logs
		WHERE log_id = $1
	`
	_, err := s.db.Exec(ctx, q, logID)
	return err
}

func (s *Store) GetProgressLogsByUserID(ctx context.Context, userID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.ProgressLog], error) {
	q := `
		SELECT log_id, user_id, exercise_id, date, sets_completed, reps_completed, weight_used, duration_seconds
		FROM progress_logs
		WHERE user_id = $1
		ORDER BY date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.Query(ctx, q, userID, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []types.ProgressLog
	for rows.Next() {
		var log types.ProgressLog
		err := rows.Scan(
			&log.LogID,
			&log.UserID,
			&log.ExerciseID,
			&log.Date,
			&log.SetsCompleted,
			&log.RepsCompleted,
			&log.WeightUsed,
			&log.DurationSeconds,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	countQuery := `SELECT COUNT(*) FROM progress_logs WHERE user_id = $1`
	var totalCount int
	err = s.db.QueryRow(ctx, countQuery, userID).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	totalPages := (totalCount + pagination.PageSize - 1) / pagination.PageSize

	return &types.PaginatedResponse[types.ProgressLog]{
		Data:       logs,
		TotalCount: totalCount,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *Store) GetProgressLogsByUserAndExercise(ctx context.Context, userID int, exerciseID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.ProgressLog], error) {
	q := `
		SELECT log_id, user_id, exercise_id, date, sets_completed, reps_completed, weight_used, duration_seconds
		FROM progress_logs
		WHERE user_id = $1 AND exercise_id = $2
		ORDER BY date DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := s.db.Query(ctx, q, userID, exerciseID, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []types.ProgressLog
	for rows.Next() {
		var log types.ProgressLog
		err := rows.Scan(
			&log.LogID,
			&log.UserID,
			&log.ExerciseID,
			&log.Date,
			&log.SetsCompleted,
			&log.RepsCompleted,
			&log.WeightUsed,
			&log.DurationSeconds,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	countQuery := `SELECT COUNT(*) FROM progress_logs WHERE user_id = $1 AND exercise_id = $2`
	var totalCount int
	err = s.db.QueryRow(ctx, countQuery, userID, exerciseID).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	totalPages := (totalCount + pagination.PageSize - 1) / pagination.PageSize

	return &types.PaginatedResponse[types.ProgressLog]{
		Data:       logs,
		TotalCount: totalCount,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *Store) GetProgressLogsByUserAndDateRange(ctx context.Context, userID int, startDate, endDate time.Time) ([]types.ProgressLog, error) {
	q := `
		SELECT log_id, user_id, exercise_id, date, sets_completed, reps_completed, weight_used, duration_seconds
		FROM progress_logs
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date DESC
	`

	rows, err := s.db.Query(ctx, q, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []types.ProgressLog
	for rows.Next() {
		var log types.ProgressLog
		err := rows.Scan(
			&log.LogID,
			&log.UserID,
			&log.ExerciseID,
			&log.Date,
			&log.SetsCompleted,
			&log.RepsCompleted,
			&log.WeightUsed,
			&log.DurationSeconds,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// Placeholder implementations for complex methods - would need more sophisticated queries
func (s *Store) FilterProgressLogs(ctx context.Context, filter types.ProgressFilter, pagination types.PaginationParams) (*types.PaginatedResponse[types.ProgressLog], error) {
	// TODO: Implement filtering logic based on ProgressFilter
	return s.GetProgressLogsByUserID(ctx, filter.UserID, pagination)
}

func (s *Store) GetUserProgressSummary(ctx context.Context, userID int) (*types.UserProgressSummary, error) {
	// TODO: Implement complex aggregation queries
	return &types.UserProgressSummary{
		UserID:        userID,
		TotalWorkouts: 0,
		CurrentStreak: 0,
		LastWorkout:   nil,
		PersonalBests: []types.PersonalBest{},
	}, nil
}

func (s *Store) GetPersonalBests(ctx context.Context, userID int) ([]types.PersonalBest, error) {
	// TODO: Implement personal best calculation
	return []types.PersonalBest{}, nil
}

func (s *Store) GetProgressTrend(ctx context.Context, userID int, exerciseID int, days int) ([]types.ProgressLog, error) {
	startDate := time.Now().AddDate(0, 0, -days)
	endDate := time.Now()

	q := `
		SELECT log_id, user_id, exercise_id, date, sets_completed, reps_completed, weight_used, duration_seconds
		FROM progress_logs
		WHERE user_id = $1 AND exercise_id = $2 AND date >= $3 AND date <= $4
		ORDER BY date ASC
	`

	rows, err := s.db.Query(ctx, q, userID, exerciseID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []types.ProgressLog
	for rows.Next() {
		var log types.ProgressLog
		err := rows.Scan(
			&log.LogID,
			&log.UserID,
			&log.ExerciseID,
			&log.Date,
			&log.SetsCompleted,
			&log.RepsCompleted,
			&log.WeightUsed,
			&log.DurationSeconds,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (s *Store) GetWorkoutStreak(ctx context.Context, userID int) (int, error) {
	// TODO: Implement streak calculation
	return 0, nil
}

func (s *Store) BulkCreateProgressLogs(ctx context.Context, logs []types.ProgressLogRequest) ([]types.ProgressLog, error) {
	var results []types.ProgressLog

	for _, logReq := range logs {
		log, err := s.CreateProgressLog(ctx, &logReq)
		if err != nil {
			return nil, err
		}
		results = append(results, *log)
	}

	return results, nil
}

func (s *Store) GetLatestProgressLogsForUser(ctx context.Context, userID int) ([]types.ProgressLog, error) {
	q := `
		SELECT log_id, user_id, exercise_id, date, sets_completed, reps_completed, weight_used, duration_seconds
		FROM progress_logs
		WHERE user_id = $1
		ORDER BY date DESC
		LIMIT 10
	`

	rows, err := s.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []types.ProgressLog
	for rows.Next() {
		var log types.ProgressLog
		err := rows.Scan(
			&log.LogID,
			&log.UserID,
			&log.ExerciseID,
			&log.Date,
			&log.SetsCompleted,
			&log.RepsCompleted,
			&log.WeightUsed,
			&log.DurationSeconds,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}
