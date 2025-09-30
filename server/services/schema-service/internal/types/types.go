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

type WorkoutProfile struct {
	WorkoutProfileID int             `json:"workout_profile_id" db:"workout_profile_id"`
	AuthUserID       string          `json:"auth_user_id" db:"auth_user_id"` // References auth service user ID
	Level            FitnessLevel    `json:"level" db:"level"`
	Goal             FitnessGoal     `json:"goal" db:"goal"`
	Frequency        int             `json:"frequency" db:"frequency"` // workouts per week
	Equipment        json.RawMessage `json:"equipment" db:"equipment"` // JSONB array of equipment types
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
}

type WorkoutProfileRequest struct {
	Level     FitnessLevel `json:"level" validate:"required"`
	Goal      FitnessGoal  `json:"goal" validate:"required"`
	Frequency int          `json:"frequency" validate:"required,min=1,max=7"`
	Equipment []string     `json:"equipment" validate:"required,min=1"`
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
	TemplateID    int          `json:"template_id" db:"template_id"`
	Name          string       `json:"name" db:"name"`
	Description   string       `json:"description" db:"description"`
	MinLevel      FitnessLevel `json:"min_level" db:"min_level"`
	MaxLevel      FitnessLevel `json:"max_level" db:"max_level"`
	SuitableGoals string       `json:"suitable_goals" db:"suitable_goals"` // comma-separated goals
	DaysPerWeek   int          `json:"days_per_week" db:"days_per_week"`
}

type WorkoutTemplateRequest struct {
	Name          string        `json:"name" validate:"required,min=2,max=50"`
	Description   string        `json:"description" validate:"max=500"`
	MinLevel      FitnessLevel  `json:"min_level" validate:"required"`
	MaxLevel      FitnessLevel  `json:"max_level" validate:"required"`
	SuitableGoals []FitnessGoal `json:"suitable_goals" validate:"required,min=1"`
	DaysPerWeek   int           `json:"days_per_week" validate:"required,min=1,max=7"`
}

type WorkoutTemplateResponse struct {
	TemplateID    int           `json:"template_id"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	MinLevel      FitnessLevel  `json:"min_level"`
	MaxLevel      FitnessLevel  `json:"max_level"`
	SuitableGoals []FitnessGoal `json:"suitable_goals"`
	DaysPerWeek   int           `json:"days_per_week"`
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
	WorkoutID int    `json:"workout_id" db:"workout_id"`
	SchemaID  int    `json:"schema_id" db:"schema_id"`
	DayOfWeek int    `json:"day_of_week" db:"day_of_week"` // 1=Monday ... 7=Sunday
	Focus     string `json:"focus" db:"focus"`             // e.g., "upper", "lower", "cardio"
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
	LogID           int       `json:"log_id" db:"log_id"`
	UserID          int       `json:"user_id" db:"user_id"`
	ExerciseID      int       `json:"exercise_id" db:"exercise_id"`
	Date            time.Time `json:"date" db:"date"`
	SetsCompleted   *int      `json:"sets_completed" db:"sets_completed"`
	RepsCompleted   *int      `json:"reps_completed" db:"reps_completed"`
	WeightUsed      *float64  `json:"weight_used" db:"weight_used"`
	DurationSeconds *int      `json:"duration_seconds" db:"duration_seconds"`
}

type ProgressLogRequest struct {
	UserID          int       `json:"user_id" validate:"required"`
	ExerciseID      int       `json:"exercise_id" validate:"required"`
	Date            time.Time `json:"date" validate:"required"`
	SetsCompleted   *int      `json:"sets_completed" validate:"omitempty,min=0,max=20"`
	RepsCompleted   *int      `json:"reps_completed" validate:"omitempty,min=0,max=1000"`
	WeightUsed      *float64  `json:"weight_used" validate:"omitempty,min=0"`
	DurationSeconds *int      `json:"duration_seconds" validate:"omitempty,min=0"`
}

// =============================================================================
// COMPLEX RESPONSE TYPES (WITH JOINS)
// =============================================================================

type WorkoutWithExercises struct {
	WorkoutID int                     `json:"workout_id"`
	SchemaID  int                     `json:"schema_id"`
	DayOfWeek int                     `json:"day_of_week"`
	Focus     string                  `json:"focus"`
	Exercises []WorkoutExerciseDetail `json:"exercises"`
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
	UserID        int            `json:"user_id"`
	TotalWorkouts int            `json:"total_workouts"`
	CurrentStreak int            `json:"current_streak"`
	LastWorkout   *time.Time     `json:"last_workout"`
	PersonalBests []PersonalBest `json:"personal_bests"`
}

type PersonalBest struct {
	ExerciseID   int       `json:"exercise_id"`
	ExerciseName string    `json:"exercise_name"`
	BestWeight   *float64  `json:"best_weight"`
	BestReps     *int      `json:"best_reps"`
	BestVolume   *float64  `json:"best_volume"` // weight * reps * sets
	AchievedAt   time.Time `json:"achieved_at"`
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
	Level       *FitnessLevel `json:"level"`
	Goals       []FitnessGoal `json:"goals"`
	DaysPerWeek *int          `json:"days_per_week"`
	Search      string        `json:"search"`
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
	Offset   int `json:"offset" validate:"min=0"`
	Limit    int `json:"limit" validate:"min=1,max=100"`
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}

// =============================================================================
// FITUP SMART LOGIC TYPES
// =============================================================================

// Fitness Assessment Types
type FitnessAssessment struct {
	AssessmentID     int             `json:"assessment_id" db:"assessment_id"`
	UserID           int             `json:"user_id" db:"user_id"`
	AssessmentDate   time.Time       `json:"assessment_date" db:"assessment_date"`
	OverallLevel     FitnessLevel    `json:"overall_level" db:"overall_level"`
	StrengthLevel    FitnessLevel    `json:"strength_level" db:"strength_level"`
	CardioLevel      FitnessLevel    `json:"cardio_level" db:"cardio_level"`
	FlexibilityLevel FitnessLevel    `json:"flexibility_level" db:"flexibility_level"`
	AssessmentData   json.RawMessage `json:"assessment_data" db:"assessment_data"`
}

type FitnessAssessmentRequest struct {
	UserID           int                    `json:"user_id" validate:"required"`
	OverallLevel     FitnessLevel           `json:"overall_level" validate:"required"`
	StrengthLevel    FitnessLevel           `json:"strength_level" validate:"required"`
	CardioLevel      FitnessLevel           `json:"cardio_level" validate:"required"`
	FlexibilityLevel FitnessLevel           `json:"flexibility_level" validate:"required"`
	AssessmentData   map[string]interface{} `json:"assessment_data"`
}

type FitnessProfile struct {
	UserID           int                 `json:"user_id"`
	CurrentLevel     FitnessLevel        `json:"current_level"`
	StrengthLevel    FitnessLevel        `json:"strength_level"`
	CardioLevel      FitnessLevel        `json:"cardio_level"`
	FlexibilityLevel FitnessLevel        `json:"flexibility_level"`
	Goals            []FitnessGoalTarget `json:"goals"`
	Equipment        []EquipmentType     `json:"equipment"`
	LastAssessment   *time.Time          `json:"last_assessment"`
	TrainingHistory  *TrainingHistory    `json:"training_history"`
}

type FitnessGoalTarget struct {
	GoalID       int             `json:"goal_id" db:"goal_id"`
	UserID       int             `json:"user_id" db:"user_id"`
	GoalType     FitnessGoal     `json:"goal_type" db:"goal_type"`
	TargetValue  float64         `json:"target_value" db:"target_value"`
	CurrentValue float64         `json:"current_value" db:"current_value"`
	TargetDate   time.Time       `json:"target_date" db:"target_date"`
	IsActive     bool            `json:"is_active" db:"is_active"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	Metadata     json.RawMessage `json:"metadata" db:"metadata"`
}

type FitnessGoalRequest struct {
	GoalType    FitnessGoal            `json:"goal_type" validate:"required"`
	TargetValue float64                `json:"target_value" validate:"required,min=0"`
	TargetDate  time.Time              `json:"target_date" validate:"required"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// Workout Session Types
type WorkoutSession struct {
	SessionID          int           `json:"session_id" db:"session_id"`
	UserID             int           `json:"user_id" db:"user_id"`
	WorkoutID          int           `json:"workout_id" db:"workout_id"`
	StartTime          time.Time     `json:"start_time" db:"start_time"`
	EndTime            *time.Time    `json:"end_time" db:"end_time"`
	Status             SessionStatus `json:"status" db:"status"`
	TotalExercises     int           `json:"total_exercises" db:"total_exercises"`
	CompletedExercises int           `json:"completed_exercises" db:"completed_exercises"`
	TotalVolume        float64       `json:"total_volume" db:"total_volume"`
	Notes              string        `json:"notes" db:"notes"`
}

type SessionStatus string

const (
	SessionActive    SessionStatus = "active"
	SessionCompleted SessionStatus = "completed"
	SessionSkipped   SessionStatus = "skipped"
	SessionAbandoned SessionStatus = "abandoned"
)

type SessionSummary struct {
	TotalDuration      int                   `json:"total_duration_seconds"`
	ExercisesCompleted int                   `json:"exercises_completed"`
	TotalVolume        float64               `json:"total_volume"`
	AverageRPE         float64               `json:"average_rpe"`
	Notes              string                `json:"notes"`
	Exercises          []ExercisePerformance `json:"exercises"`
}

type ExercisePerformance struct {
	ExerciseID    int            `json:"exercise_id"`
	SetsCompleted int            `json:"sets_completed"`
	BestSet       SetPerformance `json:"best_set"`
	TotalVolume   float64        `json:"total_volume"`
	RPE           float64        `json:"rpe"`
	Notes         string         `json:"notes"`
}

type SetPerformance struct {
	Reps   int     `json:"reps"`
	Weight float64 `json:"weight"`
	RPE    float64 `json:"rpe"`
	Rest   int     `json:"rest_seconds"`
}

type SkippedWorkout struct {
	SkipID    int       `json:"skip_id" db:"skip_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	WorkoutID int       `json:"workout_id" db:"workout_id"`
	SkipDate  time.Time `json:"skip_date" db:"skip_date"`
	Reason    string    `json:"reason" db:"reason"`
}

// Performance Analytics Types
type OneRepMaxEstimate struct {
	EstimateID   int       `json:"estimate_id" db:"estimate_id"`
	UserID       int       `json:"user_id" db:"user_id"`
	ExerciseID   int       `json:"exercise_id" db:"exercise_id"`
	EstimatedMax float64   `json:"estimated_max" db:"estimated_max"`
	EstimateDate time.Time `json:"estimate_date" db:"estimate_date"`
	Method       string    `json:"method" db:"method"`
	Confidence   float64   `json:"confidence" db:"confidence"`
}

type PerformanceData struct {
	Weight   float64 `json:"weight"`
	Reps     int     `json:"reps"`
	Sets     int     `json:"sets"`
	RPE      float64 `json:"rpe"`
	Duration int     `json:"duration_seconds"`
}

type MovementAssessment struct {
	AssessmentID   int                  `json:"assessment_id" db:"assessment_id"`
	UserID         int                  `json:"user_id" db:"user_id"`
	AssessmentDate time.Time            `json:"assessment_date" db:"assessment_date"`
	MovementData   json.RawMessage      `json:"movement_data" db:"movement_data"`
	Limitations    []MovementLimitation `json:"limitations"`
}

type MovementAssessmentRequest struct {
	UserID       int                    `json:"user_id" validate:"required"`
	MovementData map[string]interface{} `json:"movement_data" validate:"required"`
	Limitations  []string               `json:"limitations"`
}

type MovementLimitation struct {
	LimitationID int    `json:"limitation_id" db:"limitation_id"`
	UserID       int    `json:"user_id" db:"user_id"`
	MovementType string `json:"movement_type" db:"movement_type"`
	Severity     string `json:"severity" db:"severity"`
	Description  string `json:"description" db:"description"`
}

// Plan Generation Types
type GeneratedPlan struct {
	PlanID        int             `json:"plan_id" db:"plan_id"`
	UserID        int             `json:"user_id" db:"user_id"`
	WeekStart     time.Time       `json:"week_start" db:"week_start"`
	GeneratedAt   time.Time       `json:"generated_at" db:"generated_at"`
	Algorithm     string          `json:"algorithm" db:"algorithm"`
	Effectiveness float64         `json:"effectiveness" db:"effectiveness"`
	IsActive      bool            `json:"is_active" db:"is_active"`
	Metadata      json.RawMessage `json:"metadata" db:"metadata"`
}

type PlanGenerationMetadata struct {
	UserGoals          []FitnessGoal          `json:"user_goals"`
	AvailableEquipment []EquipmentType        `json:"available_equipment"`
	FitnessLevel       FitnessLevel           `json:"fitness_level"`
	WeeklyFrequency    int                    `json:"weekly_frequency"`
	TimePerWorkout     int                    `json:"time_per_workout"`
	Algorithm          string                 `json:"algorithm"`
	Parameters         map[string]interface{} `json:"parameters"`
}

type PlanPerformanceData struct {
	CompletionRate   float64 `json:"completion_rate"`
	AverageRPE       float64 `json:"average_rpe"`
	ProgressRate     float64 `json:"progress_rate"`
	UserSatisfaction float64 `json:"user_satisfaction"`
	InjuryRate       float64 `json:"injury_rate"`
}

type PlanAdaptation struct {
	AdaptationID   int             `json:"adaptation_id" db:"adaptation_id"`
	PlanID         int             `json:"plan_id" db:"plan_id"`
	AdaptationDate time.Time       `json:"adaptation_date" db:"adaptation_date"`
	Reason         string          `json:"reason" db:"reason"`
	Changes        json.RawMessage `json:"changes" db:"changes"`
	Trigger        string          `json:"trigger" db:"trigger"`
}

// Recovery and Analytics Types
type RecoveryMetrics struct {
	MetricID     int       `json:"metric_id" db:"metric_id"`
	UserID       int       `json:"user_id" db:"user_id"`
	Date         time.Time `json:"date" db:"date"`
	SleepHours   float64   `json:"sleep_hours" db:"sleep_hours"`
	SleepQuality float64   `json:"sleep_quality" db:"sleep_quality"`
	StressLevel  float64   `json:"stress_level" db:"stress_level"`
	EnergyLevel  float64   `json:"energy_level" db:"energy_level"`
	Soreness     float64   `json:"soreness" db:"soreness"`
}

type RecoveryStatus struct {
	UserID               int     `json:"user_id"`
	RecoveryScore        float64 `json:"recovery_score"`
	Recommendation       string  `json:"recommendation"`
	RecommendedIntensity float64 `json:"recommended_intensity"`
	RestDayRecommended   bool    `json:"rest_day_recommended"`
}

type SleepQuality struct {
	Hours      float64 `json:"hours"`
	Quality    float64 `json:"quality"`
	Efficiency float64 `json:"efficiency"`
}

type RestDayRecommendation struct {
	Recommended bool     `json:"recommended"`
	Reason      string   `json:"reason"`
	Duration    int      `json:"duration_days"`
	Activities  []string `json:"recommended_activities"`
}

// Additional Analytics Types
type SessionMetrics struct {
	SessionID        int     `json:"session_id"`
	Duration         int     `json:"duration_seconds"`
	TotalVolume      float64 `json:"total_volume"`
	AverageIntensity float64 `json:"average_intensity"`
	CompletionRate   float64 `json:"completion_rate"`
	RPE              float64 `json:"rpe"`
	CaloriesBurned   int     `json:"calories_burned"`
}

type WeeklySessionStats struct {
	WeekStart         time.Time `json:"week_start"`
	SessionsPlanned   int       `json:"sessions_planned"`
	SessionsCompleted int       `json:"sessions_completed"`
	TotalVolume       float64   `json:"total_volume"`
	AverageRPE        float64   `json:"average_rpe"`
	CompletionRate    float64   `json:"completion_rate"`
}

type StrengthProgression struct {
	ExerciseID      int     `json:"exercise_id"`
	StartingMax     float64 `json:"starting_max"`
	CurrentMax      float64 `json:"current_max"`
	ProgressionRate float64 `json:"progression_rate"`
	Trend           string  `json:"trend"`
}

type PlateauDetection struct {
	ExerciseID      int    `json:"exercise_id"`
	PlateauDetected bool   `json:"plateau_detected"`
	PlateauDuration int    `json:"plateau_duration_days"`
	Recommendation  string `json:"recommendation"`
}

type GoalProgress struct {
	GoalID              int        `json:"goal_id"`
	ProgressPercent     float64    `json:"progress_percent"`
	OnTrack             bool       `json:"on_track"`
	EstimatedCompletion *time.Time `json:"estimated_completion"`
}

type GoalPrediction struct {
	GoalID               int     `json:"goal_id"`
	ProbabilityOfSuccess float64 `json:"probability_of_success"`
	EstimatedTime        int     `json:"estimated_days"`
	Confidence           float64 `json:"confidence"`
}

type TrainingVolume struct {
	WeekStart    time.Time `json:"week_start"`
	TotalSets    int       `json:"total_sets"`
	TotalReps    int       `json:"total_reps"`
	TotalWeight  float64   `json:"total_weight"`
	VolumeLoad   float64   `json:"volume_load"`
	IntensityAvg float64   `json:"intensity_average"`
}

type IntensityProgression struct {
	ExerciseID        int     `json:"exercise_id"`
	BaselineIntensity float64 `json:"baseline_intensity"`
	CurrentIntensity  float64 `json:"current_intensity"`
	ProgressionRate   float64 `json:"progression_rate"`
	RecommendedNext   float64 `json:"recommended_next"`
}

type OptimalLoad struct {
	UserID          int     `json:"user_id"`
	RecommendedSets int     `json:"recommended_sets"`
	RecommendedReps int     `json:"recommended_reps"`
	IntensityRange  string  `json:"intensity_range"`
	VolumeTarget    float64 `json:"volume_target"`
}

type TimeToGoalEstimate struct {
	GoalID        int      `json:"goal_id"`
	EstimatedDays int      `json:"estimated_days"`
	Confidence    float64  `json:"confidence"`
	Assumptions   []string `json:"assumptions"`
}

type GoalAdjustment struct {
	GoalID             int    `json:"goal_id"`
	RecommendationType string `json:"recommendation_type"`
	Adjustment         string `json:"adjustment"`
	Reason             string `json:"reason"`
}

type TrainingHistory struct {
	TotalWorkouts    int     `json:"total_workouts"`
	WeeksActive      int     `json:"weeks_active"`
	AverageFrequency float64 `json:"average_frequency"`
	ConsistencyScore float64 `json:"consistency_score"`
}
