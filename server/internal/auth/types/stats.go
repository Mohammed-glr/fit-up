package types

import "time"

// UserStats represents comprehensive user statistics
type UserStats struct {
	UserID           string     `json:"user_id"`
	TotalWorkouts    int        `json:"total_workouts"`
	ActivePrograms   int        `json:"active_programs"`
	DaysActive       int        `json:"days_active"`
	CurrentStreak    int        `json:"current_streak"`
	LongestStreak    int        `json:"longest_streak"`
	TotalWeeks       int        `json:"total_weeks"`
	CompletionRate   float64    `json:"completion_rate"`
	LastWorkoutDate  *time.Time `json:"last_workout_date"`
	FirstWorkoutDate *time.Time `json:"first_workout_date"`
	AssignedCoach    *CoachInfo `json:"assigned_coach,omitempty"`
}

// CoachInfo represents coach information for user stats
type CoachInfo struct {
	CoachID       string     `json:"coach_id"`
	Name          string     `json:"name"`
	Image         *string    `json:"image"`
	Specialty     string     `json:"specialty,omitempty"`
	AssignedAt    time.Time  `json:"assigned_at"`
	TotalMessages int        `json:"total_messages"`
	LastMessageAt *time.Time `json:"last_message_at,omitempty"`
}

// TodayWorkout represents the current day's workout
type TodayWorkout struct {
	PlanID           int             `json:"plan_id"`
	PlanName         string          `json:"plan_name"`
	DayIndex         int             `json:"day_index"`
	DayTitle         string          `json:"day_title"`
	Focus            string          `json:"focus"`
	IsRest           bool            `json:"is_rest"`
	TotalExercises   int             `json:"total_exercises"`
	EstimatedMinutes int             `json:"estimated_minutes"`
	IsCompleted      bool            `json:"is_completed"`
	CompletedAt      *time.Time      `json:"completed_at,omitempty"`
	Exercises        []TodayExercise `json:"exercises"`
}

// TodayExercise represents an exercise in today's workout
type TodayExercise struct {
	ExerciseID  *int   `json:"exercise_id,omitempty"`
	Name        string `json:"name"`
	Sets        int    `json:"sets"`
	Reps        string `json:"reps"`
	RestSeconds int    `json:"rest_seconds"`
	Notes       string `json:"notes,omitempty"`
}

// ActivityFeedItem represents an item in the user's activity feed
type ActivityFeedItem struct {
	ID          string                 `json:"id"`
	Type        ActivityType           `json:"type"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Timestamp   time.Time              `json:"timestamp"`
	Icon        string                 `json:"icon"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ActivityType represents the type of activity
type ActivityType string

const (
	ActivityWorkoutCompleted ActivityType = "workout_completed"
	ActivityPRChieved        ActivityType = "pr_achieved"
	ActivityStreakMilestone  ActivityType = "streak_milestone"
	ActivityCoachMessage     ActivityType = "coach_message"
	ActivityNewPlan          ActivityType = "new_plan"
	ActivityGoalAchieved     ActivityType = "goal_achieved"
	ActivityPlanCompleted    ActivityType = "plan_completed"
)

// WorkoutCompletionRequest represents the data sent when completing a workout
type WorkoutCompletionRequest struct {
	PlanID          int              `json:"plan_id"`
	PlanDayID       *int             `json:"plan_day_id,omitempty"`
	DayIndex        int              `json:"day_index"`
	DurationSeconds int              `json:"duration_seconds"`
	CompletedAt     time.Time        `json:"completed_at"`
	Exercises       []ExerciseSetLog `json:"exercises"`
	Notes           string           `json:"notes,omitempty"`
}

// ExerciseSetLog represents a single set performed during a workout
type ExerciseSetLog struct {
	ExerciseID   *int    `json:"exercise_id,omitempty"`
	ExerciseName string  `json:"exercise_name"`
	SetNumber    int     `json:"set_number"`
	Reps         int     `json:"reps"`
	Weight       float64 `json:"weight"`
	Completed    bool    `json:"completed"`
	Notes        *string `json:"notes,omitempty"`
}

// WorkoutCompletionResponse represents the response after saving a workout
type WorkoutCompletionResponse struct {
	Success                 bool              `json:"success"`
	Message                 string            `json:"message"`
	WorkoutDate             time.Time         `json:"workout_date"`
	TotalSets               int               `json:"total_sets"`
	CompletedSets           int               `json:"completed_sets"`
	CompletionRate          float64           `json:"completion_rate"`
	TotalVolume             float64           `json:"total_volume"`
	DurationMinutes         int               `json:"duration_minutes"`
	NewStreak               int               `json:"new_streak"`
	IsPersonalBest          bool              `json:"is_personal_best"`
	NewlyEarnedAchievements []UserAchievement `json:"newly_earned_achievements,omitempty"`
}

// WorkoutHistoryItem represents a single workout in history
type WorkoutHistoryItem struct {
	Date            time.Time `json:"date"`
	PlanID          *int      `json:"plan_id,omitempty"`
	PlanName        string    `json:"plan_name,omitempty"`
	DayTitle        string    `json:"day_title,omitempty"`
	TotalExercises  int       `json:"total_exercises"`
	CompletedSets   int       `json:"completed_sets"`
	TotalVolume     float64   `json:"total_volume"`
	DurationMinutes int       `json:"duration_minutes"`
	Exercises       []string  `json:"exercises"`
}

// WorkoutHistoryResponse represents paginated workout history
type WorkoutHistoryResponse struct {
	Workouts   []WorkoutHistoryItem `json:"workouts"`
	TotalCount int                  `json:"total_count"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	HasMore    bool                 `json:"has_more"`
}

// ExerciseProgressData represents performance data for a specific exercise
type ExerciseProgressData struct {
	ExerciseID   *int                        `json:"exercise_id,omitempty"`
	ExerciseName string                      `json:"exercise_name"`
	DataPoints   []ExerciseProgressDataPoint `json:"data_points"`
	MaxWeight    float64                     `json:"max_weight"`
	MaxVolume    float64                     `json:"max_volume"`
	TotalSets    int                         `json:"total_sets"`
}

// ExerciseProgressDataPoint represents a single data point in exercise progress
type ExerciseProgressDataPoint struct {
	Date             time.Time `json:"date"`
	Weight           float64   `json:"weight"`
	Reps             int       `json:"reps"`
	Sets             int       `json:"sets"`
	Volume           float64   `json:"volume"`
	IsPersonalRecord bool      `json:"is_personal_record"`
}
