package types

import "time"

type MindfulnessSession struct {
	SessionID   int       `json:"session_id"`
	UserID      string    `json:"user_id"`
	SessionType string    `json:"session_type"` // pre_workout, post_workout, breathing, meditation, gratitude
	Duration    int       `json:"duration_seconds"`
	CompletedAt time.Time `json:"completed_at"`
	Notes       *string   `json:"notes,omitempty"`
	MoodBefore  *int      `json:"mood_before,omitempty"`
	MoodAfter   *int      `json:"mood_after,omitempty"`
}

type BreathingExercise struct {
	ExerciseID      int       `json:"exercise_id"`
	UserID          string    `json:"user_id"`
	BreathingType   string    `json:"breathing_type"`
	Duration        int       `json:"duration_seconds"`
	CyclesCompleted int       `json:"cycles_completed"`
	CompletedAt     time.Time `json:"completed_at"`
	HeartRateBefore *int      `json:"heart_rate_before,omitempty"`
	HeartRateAfter  *int      `json:"heart_rate_after,omitempty"`
}

type GratitudeEntry struct {
	EntryID          int       `json:"entry_id"`
	UserID           string    `json:"user_id"`
	EntryText        string    `json:"entry_text"`
	Tags             []string  `json:"tags,omitempty"`
	Mood             *int      `json:"mood,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	WorkoutSessionID *int      `json:"workout_session_id,omitempty"`
}

type ReflectionPrompt struct {
	PromptID   int    `json:"prompt_id"`
	PromptText string `json:"prompt_text"`
	Category   string `json:"category"`
	IsActive   bool   `json:"is_active"`
}

type ReflectionResponse struct {
	ResponseID   int       `json:"response_id"`
	UserID       string    `json:"user_id"`
	PromptID     *int      `json:"prompt_id,omitempty"`
	ResponseText string    `json:"response_text"`
	CreatedAt    time.Time `json:"created_at"`
}

type MindfulnessStreak struct {
	StreakID         int       `json:"streak_id"`
	UserID           string    `json:"user_id"`
	CurrentStreak    int       `json:"current_streak"`
	LongestStreak    int       `json:"longest_streak"`
	LastActivityDate time.Time `json:"last_activity_date"`
	TotalSessions    int       `json:"total_sessions"`
}

type CreateMindfulnessSessionRequest struct {
	SessionType string  `json:"session_type"`
	Duration    int     `json:"duration_seconds"`
	Notes       *string `json:"notes,omitempty"`
	MoodBefore  *int    `json:"mood_before,omitempty"`
	MoodAfter   *int    `json:"mood_after,omitempty"`
}

type CreateBreathingExerciseRequest struct {
	BreathingType   string `json:"breathing_type"`
	Duration        int    `json:"duration_seconds"`
	CyclesCompleted int    `json:"cycles_completed"`
	HeartRateBefore *int   `json:"heart_rate_before,omitempty"`
	HeartRateAfter  *int   `json:"heart_rate_after,omitempty"`
}

type CreateGratitudeEntryRequest struct {
	EntryText        string   `json:"entry_text"`
	Tags             []string `json:"tags,omitempty"`
	Mood             *int     `json:"mood,omitempty"`
	WorkoutSessionID *int     `json:"workout_session_id,omitempty"`
}

type CreateReflectionResponseRequest struct {
	PromptID     *int   `json:"prompt_id,omitempty"`
	ResponseText string `json:"response_text"`
}

type MindfulnessStats struct {
	TotalSessions  int                  `json:"total_sessions"`
	TotalMinutes   int                  `json:"total_minutes"`
	CurrentStreak  int                  `json:"current_streak"`
	LongestStreak  int                  `json:"longest_streak"`
	SessionsByType map[string]int       `json:"sessions_by_type"`
	RecentSessions []MindfulnessSession `json:"recent_sessions"`
	MoodTrend      []MoodDataPoint      `json:"mood_trend,omitempty"`
}

type MoodDataPoint struct {
	Date      string  `json:"date"`
	AvgBefore float64 `json:"avg_before"`
	AvgAfter  float64 `json:"avg_after"`
}

type BreathingStats struct {
	TotalExercises  int                 `json:"total_exercises"`
	TotalMinutes    int                 `json:"total_minutes"`
	TotalCycles     int                 `json:"total_cycles"`
	ExercisesByType map[string]int      `json:"exercises_by_type"`
	RecentExercises []BreathingExercise `json:"recent_exercises"`
}
