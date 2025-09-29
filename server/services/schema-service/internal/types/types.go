package types

import (
	"encoding/json"
	"time"
)

// =============================================================================
// ENUMS AND CONSTANTS
// =============================================================================

type FitnessLevel string

const (
	LevelBeginner     FitnessLevel = "beginner"
	LevelIntermediate FitnessLevel = "intermediate"
	LevelAdvanced     FitnessLevel = "advanced"
)

type FitnessGoal string

const (
	GoalStrength   FitnessGoal = "strength"
	GoalMuscleGain FitnessGoal = "muscle_gain"
	GoalFatLoss    FitnessGoal = "fat_loss"
	GoalEndurance  FitnessGoal = "endurance"
)

type ExerciseType string

const (
	TypeStrength ExerciseType = "strength"
	TypeCardio   ExerciseType = "cardio"
	TypeMobility ExerciseType = "mobility"
	TypeHIIT     ExerciseType = "hiit"
)

type EquipmentType string

const (
	EquipmentBarbell    EquipmentType = "barbell"
	EquipmentDumbbell   EquipmentType = "dumbbell"
	EquipmentBodyweight EquipmentType = "bodyweight"
	EquipmentMachine    EquipmentType = "machine"
	EquipmentKettlebell EquipmentType = "kettlebell"
	EquipmentBands      EquipmentType = "bands"
)

// =============================================================================
// CORE WORKOUT TYPES
// =============================================================================

type WorkoutUser struct {
	UserID    int               `json:"user_id" db:"user_id"`
	Name      string            `json:"name" db:"name"`
	Email     string            `json:"email" db:"email"`
	Level     FitnessLevel      `json:"level" db:"level"`
	Goal      FitnessGoal       `json:"goal" db:"goal"`
	Frequency int               `json:"frequency" db:"frequency"` // workouts per week
	Equipment json.RawMessage   `json:"equipment" db:"equipment"` // JSONB array of equipment types
	CreatedAt time.Time         `json:"created_at" db:"created_at"`
}

type WorkoutUserRequest struct {
	Name      string        `json:"name" validate:"required,min=2,max=50"`
	Email     string        `json:"email" validate:"required,email,max=100"`
	Level     FitnessLevel  `json:"level" validate:"required"`
	Goal      FitnessGoal   `json:"goal" validate:"required"`
	Frequency int           `json:"frequency" validate:"required,min=1,max=7"`
	Equipment []string      `json:"equipment" validate:"required,min=1"`
}

type Exercise struct {
	ExerciseID   int           `json:"exercise_id" db:"exercise_id"`
	Name         string        `json:"name" db:"name"`
	MuscleGroups string        `json:"muscle_groups" db:"muscle_groups"` // comma-separated
	Difficulty   FitnessLevel  `json:"difficulty" db:"difficulty"`
	Equipment    EquipmentType `json:"equipment" db:"equipment"`
	Type         ExerciseType  `json:"type" db:"type"`
	DefaultSets  int           `json:"default_sets" db:"default_sets"`
	DefaultReps  string        `json:"default_reps" db:"default_reps"` // e.g., "8-12"
	RestSeconds  int           `json:"rest_seconds" db:"rest_seconds"`
}

type ExerciseRequest struct {
	Name         string        `json:"name" validate:"required,min=2,max=100"`
	MuscleGroups []string      `json:"muscle_groups" validate:"required,min=1"`
	Difficulty   FitnessLevel  `json:"difficulty" validate:"required"`
	Equipment    EquipmentType `json:"equipment" validate:"required"`
	Type         ExerciseType  `json:"type" validate:"required"`
	DefaultSets  int           `json:"default_sets" validate:"required,min=1,max=10"`
	DefaultReps  string        `json:"default_reps" validate:"required"`
	RestSeconds  int           `json:"rest_seconds" validate:"required,min=0,max=600"`
}

type ExerciseResponse struct {
	ExerciseID   int           `json:"exercise_id"`
	Name         string        `json:"name"`
	MuscleGroups []string      `json:"muscle_groups"`
	Difficulty   FitnessLevel  `json:"difficulty"`
	Equipment    EquipmentType `json:"equipment"`
	Type         ExerciseType  `json:"type"`
	DefaultSets  int           `json:"default_sets"`
	DefaultReps  string        `json:"default_reps"`
	RestSeconds  int           `json:"rest_seconds"`
}

type WorkoutTemplate struct {
	TemplateID    int      `json:"template_id" db:"template_id"`
	Name          string   `json:"name" db:"name"`
	Description   string   `json:"description" db:"description"`
	MinLevel      FitnessLevel `json:"min_level" db:"min_level"`
	MaxLevel      FitnessLevel `json:"max_level" db:"max_level"`
	SuitableGoals string   `json:"suitable_goals" db:"suitable_goals"` // comma-separated goals
	DaysPerWeek   int      `json:"days_per_week" db:"days_per_week"`
}

type WorkoutTemplateRequest struct {
	Name          string         `json:"name" validate:"required,min=2,max=50"`
	Description   string         `json:"description" validate:"max=500"`
	MinLevel      FitnessLevel   `json:"min_level" validate:"required"`
	MaxLevel      FitnessLevel   `json:"max_level" validate:"required"`
	SuitableGoals []FitnessGoal  `json:"suitable_goals" validate:"required,min=1"`
	DaysPerWeek   int            `json:"days_per_week" validate:"required,min=1,max=7"`
}

type WorkoutTemplateResponse struct {
	TemplateID    int            `json:"template_id"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	MinLevel      FitnessLevel   `json:"min_level"`
	MaxLevel      FitnessLevel   `json:"max_level"`
	SuitableGoals []FitnessGoal  `json:"suitable_goals"`
	DaysPerWeek   int            `json:"days_per_week"`
}

// =============================================================================
// WEEKLY SCHEMA AND WORKOUT TYPES
// =============================================================================

type WeeklySchema struct {
	SchemaID  int       `json:"schema_id" db:"schema_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	WeekStart time.Time `json:"week_start" db:"week_start"` // Monday of that week
	Active    bool      `json:"active" db:"active"`
}

type WeeklySchemaRequest struct {
	UserID    int       `json:"user_id" validate:"required"`
	WeekStart time.Time `json:"week_start" validate:"required"`
}

type Workout struct {
	WorkoutID   int    `json:"workout_id" db:"workout_id"`
	SchemaID    int    `json:"schema_id" db:"schema_id"`
	DayOfWeek   int    `json:"day_of_week" db:"day_of_week"` // 1=Monday ... 7=Sunday
	Focus       string `json:"focus" db:"focus"`             // e.g., "upper", "lower", "cardio"
}

type WorkoutRequest struct {
	SchemaID  int    `json:"schema_id" validate:"required"`
	DayOfWeek int    `json:"day_of_week" validate:"required,min=1,max=7"`
	Focus     string `json:"focus" validate:"required,min=2,max=50"`
}

type WorkoutExercise struct {
	WeID        int    `json:"we_id" db:"we_id"`
	WorkoutID   int    `json:"workout_id" db:"workout_id"`
	ExerciseID  int    `json:"exercise_id" db:"exercise_id"`
	Sets        int    `json:"sets" db:"sets"`
	Reps        string `json:"reps" db:"reps"`
	RestSeconds int    `json:"rest_seconds" db:"rest_seconds"`
}

type WorkoutExerciseRequest struct {
	WorkoutID   int    `json:"workout_id" validate:"required"`
	ExerciseID  int    `json:"exercise_id" validate:"required"`
	Sets        int    `json:"sets" validate:"required,min=1,max=10"`
	Reps        string `json:"reps" validate:"required"`
	RestSeconds int    `json:"rest_seconds" validate:"required,min=0,max=600"`
}

// =============================================================================
// PROGRESS TRACKING TYPES
// =============================================================================

type ProgressLog struct {
	LogID          int       `json:"log_id" db:"log_id"`
	UserID         int       `json:"user_id" db:"user_id"`
	ExerciseID     int       `json:"exercise_id" db:"exercise_id"`
	Date           time.Time `json:"date" db:"date"`
	SetsCompleted  *int      `json:"sets_completed" db:"sets_completed"`
	RepsCompleted  *int      `json:"reps_completed" db:"reps_completed"`
	WeightUsed     *float64  `json:"weight_used" db:"weight_used"`
	DurationSeconds *int     `json:"duration_seconds" db:"duration_seconds"`
}

type ProgressLogRequest struct {
	UserID          int      `json:"user_id" validate:"required"`
	ExerciseID      int      `json:"exercise_id" validate:"required"`
	Date            time.Time `json:"date" validate:"required"`
	SetsCompleted   *int     `json:"sets_completed" validate:"omitempty,min=0,max=20"`
	RepsCompleted   *int     `json:"reps_completed" validate:"omitempty,min=0,max=1000"`
	WeightUsed      *float64 `json:"weight_used" validate:"omitempty,min=0"`
	DurationSeconds *int     `json:"duration_seconds" validate:"omitempty,min=0"`
}

// =============================================================================
// COMPLEX RESPONSE TYPES (WITH JOINS)
// =============================================================================

type WorkoutWithExercises struct {
	WorkoutID   int                       `json:"workout_id"`
	SchemaID    int                       `json:"schema_id"`
	DayOfWeek   int                       `json:"day_of_week"`
	Focus       string                    `json:"focus"`
	Exercises   []WorkoutExerciseDetail   `json:"exercises"`
}

type WorkoutExerciseDetail struct {
	WeID        int              `json:"we_id"`
	Sets        int              `json:"sets"`
	Reps        string           `json:"reps"`
	RestSeconds int              `json:"rest_seconds"`
	Exercise    ExerciseResponse `json:"exercise"`
}

type WeeklySchemaWithWorkouts struct {
	SchemaID  int                    `json:"schema_id"`
	UserID    int                    `json:"user_id"`
	WeekStart time.Time              `json:"week_start"`
	Active    bool                   `json:"active"`
	Workouts  []WorkoutWithExercises `json:"workouts"`
}

type UserProgressSummary struct {
	UserID        int                    `json:"user_id"`
	TotalWorkouts int                    `json:"total_workouts"`
	CurrentStreak int                    `json:"current_streak"`
	LastWorkout   *time.Time             `json:"last_workout"`
	PersonalBests []PersonalBest         `json:"personal_bests"`
}

type PersonalBest struct {
	ExerciseID   int              `json:"exercise_id"`
	ExerciseName string           `json:"exercise_name"`
	BestWeight   *float64         `json:"best_weight"`
	BestReps     *int             `json:"best_reps"`
	BestVolume   *float64         `json:"best_volume"` // weight * reps * sets
	AchievedAt   time.Time        `json:"achieved_at"`
}

// =============================================================================
// FILTER AND SEARCH TYPES
// =============================================================================

type ExerciseFilter struct {
	MuscleGroups []string        `json:"muscle_groups"`
	Difficulty   *FitnessLevel   `json:"difficulty"`
	Equipment    []EquipmentType `json:"equipment"`
	Type         []ExerciseType  `json:"type"`
	Search       string          `json:"search"` // for name search
}

type TemplateFilter struct {
	Level         *FitnessLevel  `json:"level"`
	Goals         []FitnessGoal  `json:"goals"`
	DaysPerWeek   *int           `json:"days_per_week"`
	Search        string         `json:"search"`
}

type ProgressFilter struct {
	UserID     int        `json:"user_id"`
	ExerciseID *int       `json:"exercise_id"`
	DateFrom   *time.Time `json:"date_from"`
	DateTo     *time.Time `json:"date_to"`
}

// =============================================================================
// API RESPONSE WRAPPERS
// =============================================================================

type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

type PaginatedResponse[T any] struct {
	Data       []T `json:"data"`
	TotalCount int `json:"total_count"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
}

type PaginationParams struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}
