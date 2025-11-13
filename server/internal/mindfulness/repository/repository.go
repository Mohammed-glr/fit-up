package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tdmdh/fit-up-server/internal/mindfulness/types"
)

type MindfulnessRepo interface {
	CreateMindfulnessSession(ctx context.Context, userID string, req *types.CreateMindfulnessSessionRequest) (*types.MindfulnessSession, error)
	GetMindfulnessSessions(ctx context.Context, userID string, limit int) ([]types.MindfulnessSession, error)
	GetMindfulnessStats(ctx context.Context, userID string) (*types.MindfulnessStats, error)

	CreateBreathingExercise(ctx context.Context, userID string, req *types.CreateBreathingExerciseRequest) (*types.BreathingExercise, error)
	GetBreathingExercises(ctx context.Context, userID string, limit int) ([]types.BreathingExercise, error)
	GetBreathingStats(ctx context.Context, userID string) (*types.BreathingStats, error)

	CreateGratitudeEntry(ctx context.Context, userID string, req *types.CreateGratitudeEntryRequest) (*types.GratitudeEntry, error)
	GetGratitudeEntries(ctx context.Context, userID string, limit int) ([]types.GratitudeEntry, error)
	DeleteGratitudeEntry(ctx context.Context, userID string, entryID int) error

	GetReflectionPrompts(ctx context.Context, category *string) ([]types.ReflectionPrompt, error)
	CreateReflectionResponse(ctx context.Context, userID string, req *types.CreateReflectionResponseRequest) (*types.ReflectionResponse, error)
	GetReflectionResponses(ctx context.Context, userID string, limit int) ([]types.ReflectionResponse, error)

	GetOrCreateStreak(ctx context.Context, userID string) (*types.MindfulnessStreak, error)
	UpdateStreak(ctx context.Context, userID string) error
}

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) CreateMindfulnessSession(ctx context.Context, userID string, req *types.CreateMindfulnessSessionRequest) (*types.MindfulnessSession, error) {
	query := `
		INSERT INTO mindfulness_sessions (user_id, session_type, duration_seconds, notes, mood_before, mood_after)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING session_id, user_id, session_type, duration_seconds, completed_at, notes, mood_before, mood_after
	`

	var session types.MindfulnessSession
	err := s.db.QueryRow(ctx, query, userID, req.SessionType, req.Duration, req.Notes, req.MoodBefore, req.MoodAfter).Scan(
		&session.SessionID,
		&session.UserID,
		&session.SessionType,
		&session.Duration,
		&session.CompletedAt,
		&session.Notes,
		&session.MoodBefore,
		&session.MoodAfter,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create mindfulness session: %w", err)
	}

	go s.UpdateStreak(context.Background(), userID)

	return &session, nil
}

func (s *Store) GetMindfulnessSessions(ctx context.Context, userID string, limit int) ([]types.MindfulnessSession, error) {
	query := `
		SELECT session_id, user_id, session_type, duration_seconds, completed_at, notes, mood_before, mood_after
		FROM mindfulness_sessions
		WHERE user_id = $1
		ORDER BY completed_at DESC
		LIMIT $2
	`

	rows, err := s.db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get mindfulness sessions: %w", err)
	}
	defer rows.Close()

	var sessions []types.MindfulnessSession
	for rows.Next() {
		var session types.MindfulnessSession
		err := rows.Scan(
			&session.SessionID,
			&session.UserID,
			&session.SessionType,
			&session.Duration,
			&session.CompletedAt,
			&session.Notes,
			&session.MoodBefore,
			&session.MoodAfter,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan mindfulness session: %w", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (s *Store) GetMindfulnessStats(ctx context.Context, userID string) (*types.MindfulnessStats, error) {
	stats := &types.MindfulnessStats{
		SessionsByType: make(map[string]int),
	}

	query := `
		SELECT 
			COUNT(*) as total_sessions,
			COALESCE(SUM(duration_seconds), 0) as total_seconds
		FROM mindfulness_sessions
		WHERE user_id = $1
	`
	var totalSeconds int
	err := s.db.QueryRow(ctx, query, userID).Scan(&stats.TotalSessions, &totalSeconds)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to get session stats: %w", err)
	}
	stats.TotalMinutes = totalSeconds / 60

	typeQuery := `
		SELECT session_type, COUNT(*)
		FROM mindfulness_sessions
		WHERE user_id = $1
		GROUP BY session_type
	`
	rows, err := s.db.Query(ctx, typeQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions by type: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var sessionType string
		var count int
		if err := rows.Scan(&sessionType, &count); err != nil {
			continue
		}
		stats.SessionsByType[sessionType] = count
	}

	streak, err := s.GetOrCreateStreak(ctx, userID)
	if err == nil {
		stats.CurrentStreak = streak.CurrentStreak
		stats.LongestStreak = streak.LongestStreak
	}

	stats.RecentSessions, _ = s.GetMindfulnessSessions(ctx, userID, 10)

	return stats, nil
}

func (s *Store) CreateBreathingExercise(ctx context.Context, userID string, req *types.CreateBreathingExerciseRequest) (*types.BreathingExercise, error) {
	query := `
		INSERT INTO breathing_exercises (user_id, breathing_type, duration_seconds, cycles_completed, heart_rate_before, heart_rate_after)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING exercise_id, user_id, breathing_type, duration_seconds, cycles_completed, completed_at, heart_rate_before, heart_rate_after
	`

	var exercise types.BreathingExercise
	err := s.db.QueryRow(ctx, query, userID, req.BreathingType, req.Duration, req.CyclesCompleted, req.HeartRateBefore, req.HeartRateAfter).Scan(
		&exercise.ExerciseID,
		&exercise.UserID,
		&exercise.BreathingType,
		&exercise.Duration,
		&exercise.CyclesCompleted,
		&exercise.CompletedAt,
		&exercise.HeartRateBefore,
		&exercise.HeartRateAfter,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create breathing exercise: %w", err)
	}

	go s.UpdateStreak(context.Background(), userID)

	return &exercise, nil
}

func (s *Store) GetBreathingExercises(ctx context.Context, userID string, limit int) ([]types.BreathingExercise, error) {
	query := `
		SELECT exercise_id, user_id, breathing_type, duration_seconds, cycles_completed, completed_at, heart_rate_before, heart_rate_after
		FROM breathing_exercises
		WHERE user_id = $1
		ORDER BY completed_at DESC
		LIMIT $2
	`

	rows, err := s.db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get breathing exercises: %w", err)
	}
	defer rows.Close()

	var exercises []types.BreathingExercise
	for rows.Next() {
		var exercise types.BreathingExercise
		err := rows.Scan(
			&exercise.ExerciseID,
			&exercise.UserID,
			&exercise.BreathingType,
			&exercise.Duration,
			&exercise.CyclesCompleted,
			&exercise.CompletedAt,
			&exercise.HeartRateBefore,
			&exercise.HeartRateAfter,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan breathing exercise: %w", err)
		}
		exercises = append(exercises, exercise)
	}

	return exercises, nil
}

func (s *Store) GetBreathingStats(ctx context.Context, userID string) (*types.BreathingStats, error) {
	stats := &types.BreathingStats{
		ExercisesByType: make(map[string]int),
	}

	query := `
		SELECT 
			COUNT(*) as total_exercises,
			COALESCE(SUM(duration_seconds), 0) as total_seconds,
			COALESCE(SUM(cycles_completed), 0) as total_cycles
		FROM breathing_exercises
		WHERE user_id = $1
	`
	var totalSeconds int
	err := s.db.QueryRow(ctx, query, userID).Scan(&stats.TotalExercises, &totalSeconds, &stats.TotalCycles)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to get breathing stats: %w", err)
	}
	stats.TotalMinutes = totalSeconds / 60

	typeQuery := `
		SELECT breathing_type, COUNT(*)
		FROM breathing_exercises
		WHERE user_id = $1
		GROUP BY breathing_type
	`
	rows, err := s.db.Query(ctx, typeQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exercises by type: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var breathingType string
		var count int
		if err := rows.Scan(&breathingType, &count); err != nil {
			continue
		}
		stats.ExercisesByType[breathingType] = count
	}

	stats.RecentExercises, _ = s.GetBreathingExercises(ctx, userID, 10)

	return stats, nil
}

func (s *Store) CreateGratitudeEntry(ctx context.Context, userID string, req *types.CreateGratitudeEntryRequest) (*types.GratitudeEntry, error) {
	query := `
		INSERT INTO gratitude_entries (user_id, entry_text, tags, mood, workout_session_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING entry_id, user_id, entry_text, tags, mood, created_at, workout_session_id
	`

	var entry types.GratitudeEntry
	err := s.db.QueryRow(ctx, query, userID, req.EntryText, req.Tags, req.Mood, req.WorkoutSessionID).Scan(
		&entry.EntryID,
		&entry.UserID,
		&entry.EntryText,
		&entry.Tags,
		&entry.Mood,
		&entry.CreatedAt,
		&entry.WorkoutSessionID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gratitude entry: %w", err)
	}

	return &entry, nil
}

func (s *Store) GetGratitudeEntries(ctx context.Context, userID string, limit int) ([]types.GratitudeEntry, error) {
	query := `
		SELECT entry_id, user_id, entry_text, COALESCE(tags, '{}'), mood, created_at, workout_session_id
		FROM gratitude_entries
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := s.db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get gratitude entries: %w", err)
	}
	defer rows.Close()

	var entries []types.GratitudeEntry
	for rows.Next() {
		var entry types.GratitudeEntry
		err := rows.Scan(
			&entry.EntryID,
			&entry.UserID,
			&entry.EntryText,
			&entry.Tags,
			&entry.Mood,
			&entry.CreatedAt,
			&entry.WorkoutSessionID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan gratitude entry: %w", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (s *Store) DeleteGratitudeEntry(ctx context.Context, userID string, entryID int) error {
	query := `DELETE FROM gratitude_entries WHERE entry_id = $1 AND user_id = $2`
	result, err := s.db.Exec(ctx, query, entryID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete gratitude entry: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("gratitude entry not found")
	}
	return nil
}

func (s *Store) GetReflectionPrompts(ctx context.Context, category *string) ([]types.ReflectionPrompt, error) {
	query := `
		SELECT prompt_id, prompt_text, category, is_active
		FROM reflection_prompts
		WHERE is_active = true
	`
	args := []interface{}{}

	if category != nil && *category != "" {
		query += ` AND category = $1`
		args = append(args, *category)
	}

	query += ` ORDER BY RANDOM() LIMIT 5`

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get reflection prompts: %w", err)
	}
	defer rows.Close()

	var prompts []types.ReflectionPrompt
	for rows.Next() {
		var prompt types.ReflectionPrompt
		err := rows.Scan(&prompt.PromptID, &prompt.PromptText, &prompt.Category, &prompt.IsActive)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reflection prompt: %w", err)
		}
		prompts = append(prompts, prompt)
	}

	return prompts, nil
}

func (s *Store) CreateReflectionResponse(ctx context.Context, userID string, req *types.CreateReflectionResponseRequest) (*types.ReflectionResponse, error) {
	query := `
		INSERT INTO reflection_responses (user_id, prompt_id, response_text)
		VALUES ($1, $2, $3)
		RETURNING response_id, user_id, prompt_id, response_text, created_at
	`

	var response types.ReflectionResponse
	err := s.db.QueryRow(ctx, query, userID, req.PromptID, req.ResponseText).Scan(
		&response.ResponseID,
		&response.UserID,
		&response.PromptID,
		&response.ResponseText,
		&response.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create reflection response: %w", err)
	}

	return &response, nil
}

func (s *Store) GetReflectionResponses(ctx context.Context, userID string, limit int) ([]types.ReflectionResponse, error) {
	query := `
		SELECT response_id, user_id, prompt_id, response_text, created_at
		FROM reflection_responses
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := s.db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get reflection responses: %w", err)
	}
	defer rows.Close()

	var responses []types.ReflectionResponse
	for rows.Next() {
		var response types.ReflectionResponse
		err := rows.Scan(
			&response.ResponseID,
			&response.UserID,
			&response.PromptID,
			&response.ResponseText,
			&response.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reflection response: %w", err)
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (s *Store) GetOrCreateStreak(ctx context.Context, userID string) (*types.MindfulnessStreak, error) {
	query := `
		INSERT INTO mindfulness_streaks (user_id, current_streak, longest_streak, last_activity_date, total_sessions)
		VALUES ($1, 0, 0, CURRENT_DATE, 0)
		ON CONFLICT (user_id) DO UPDATE SET user_id = EXCLUDED.user_id
		RETURNING streak_id, user_id, current_streak, longest_streak, last_activity_date, total_sessions
	`

	var streak types.MindfulnessStreak
	err := s.db.QueryRow(ctx, query, userID).Scan(
		&streak.StreakID,
		&streak.UserID,
		&streak.CurrentStreak,
		&streak.LongestStreak,
		&streak.LastActivityDate,
		&streak.TotalSessions,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create streak: %w", err)
	}

	return &streak, nil
}

func (s *Store) UpdateStreak(ctx context.Context, userID string) error {
	query := `
		INSERT INTO mindfulness_streaks (user_id, current_streak, longest_streak, last_activity_date, total_sessions)
		VALUES ($1, 1, 1, CURRENT_DATE, 1)
		ON CONFLICT (user_id) DO UPDATE SET
			current_streak = CASE
				WHEN mindfulness_streaks.last_activity_date = CURRENT_DATE THEN mindfulness_streaks.current_streak
				WHEN mindfulness_streaks.last_activity_date = CURRENT_DATE - INTERVAL '1 day' THEN mindfulness_streaks.current_streak + 1
				ELSE 1
			END,
			longest_streak = GREATEST(
				mindfulness_streaks.longest_streak,
				CASE
					WHEN mindfulness_streaks.last_activity_date = CURRENT_DATE THEN mindfulness_streaks.current_streak
					WHEN mindfulness_streaks.last_activity_date = CURRENT_DATE - INTERVAL '1 day' THEN mindfulness_streaks.current_streak + 1
					ELSE 1
				END
			),
			last_activity_date = CURRENT_DATE,
			total_sessions = mindfulness_streaks.total_sessions + 1
	`

	_, err := s.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to update streak: %w", err)
	}

	return nil
}
