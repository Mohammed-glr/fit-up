package types

import "time"

// ShareWorkoutRequest represents a request to share a completed workout
type ShareWorkoutRequest struct {
	SessionID int    `json:"session_id" validate:"required,gt=0"`
	ShareType string `json:"share_type" validate:"required,oneof=coach image text social"`
	Message   string `json:"message,omitempty"`
}

// WorkoutShareSummary contains all the data for a shareable workout summary
type WorkoutShareSummary struct {
	SessionID       int                    `json:"session_id"`
	WorkoutTitle    string                 `json:"workout_title"`
	CompletedAt     time.Time              `json:"completed_at"`
	DurationMinutes int                    `json:"duration_minutes"`
	TotalExercises  int                    `json:"total_exercises"`
	TotalSets       int                    `json:"total_sets"`
	TotalReps       int                    `json:"total_reps"`
	TotalVolumeLbs  float64                `json:"total_volume_lbs"`
	Exercises       []WorkoutShareExercise `json:"exercises"`
	PRsAchieved     int                    `json:"prs_achieved"`
	UserName        string                 `json:"user_name,omitempty"`
	UserPhotoURL    string                 `json:"user_photo_url,omitempty"`
}

// WorkoutShareExercise contains exercise data for sharing
type WorkoutShareExercise struct {
	ExerciseName   string   `json:"exercise_name"`
	SetsCompleted  int      `json:"sets_completed"`
	TotalReps      int      `json:"total_reps"`
	TotalVolumeLbs float64  `json:"total_volume_lbs"`
	PRAchieved     bool     `json:"pr_achieved"`
	BestSet        *BestSet `json:"best_set,omitempty"`
}

// BestSet represents the best set in an exercise
type BestSet struct {
	Weight float64 `json:"weight"`
	Reps   int     `json:"reps"`
}

// ShareWorkoutResponse is the response after sharing a workout
type ShareWorkoutResponse struct {
	Success   bool   `json:"success"`
	ShareURL  string `json:"share_url,omitempty"`
	ShareText string `json:"share_text,omitempty"`
	Message   string `json:"message"`
}

// ShareToCoachRequest represents a request to share workout with coach
type ShareToCoachRequest struct {
	SessionID      int    `json:"session_id" validate:"required,gt=0"`
	CoachID        string `json:"coach_id" validate:"required"`
	Message        string `json:"message,omitempty"`
	IncludeSummary bool   `json:"include_summary"`
}

// ShareToCoachResponse is the response after sharing with coach
type ShareToCoachResponse struct {
	Success   bool   `json:"success"`
	MessageID int    `json:"message_id,omitempty"`
	Message   string `json:"message"`
}
