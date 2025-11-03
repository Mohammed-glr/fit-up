package repository

import (
	"context"
	"time"

	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

type WorkoutProfileRepo interface {
	CreateWorkoutProfile(ctx context.Context, authUserID string, profile *types.WorkoutProfileRequest) (*types.WorkoutProfile, error)
	GetWorkoutProfileByAuthID(ctx context.Context, authUserID string) (*types.WorkoutProfile, error)
	GetWorkoutProfileByID(ctx context.Context, workoutProfileID int) (*types.WorkoutProfile, error)
	UpdateWorkoutProfile(ctx context.Context, authUserID string, profile *types.WorkoutProfileRequest) (*types.WorkoutProfile, error)
	DeleteWorkoutProfile(ctx context.Context, authUserID string) error

	ListWorkoutProfiles(ctx context.Context, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutProfile], error)
	SearchWorkoutProfiles(ctx context.Context, query string, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutProfile], error)

	GetProfilesByLevel(ctx context.Context, level types.FitnessLevel) ([]types.WorkoutProfile, error)
	GetProfilesByGoal(ctx context.Context, goal types.FitnessGoal) ([]types.WorkoutProfile, error)
	CountActiveProfiles(ctx context.Context) (int, error)
}

type ExerciseRepo interface {
	CreateExercise(ctx context.Context, exercise *types.ExerciseRequest) (*types.Exercise, error)
	GetExerciseByID(ctx context.Context, exerciseID int) (*types.Exercise, error)
	UpdateExercise(ctx context.Context, exerciseID int, exercise *types.ExerciseRequest) (*types.Exercise, error)
	DeleteExercise(ctx context.Context, exerciseID int) error

	ListExercises(ctx context.Context, pagination types.PaginationParams) (*types.PaginatedResponse[types.Exercise], error)
	FilterExercises(ctx context.Context, filter types.ExerciseFilter, pagination types.PaginationParams) (*types.PaginatedResponse[types.Exercise], error)
	SearchExercises(ctx context.Context, query string, pagination types.PaginationParams) (*types.PaginatedResponse[types.Exercise], error)

	GetExercisesByMuscleGroup(ctx context.Context, muscleGroup string) ([]types.Exercise, error)
	GetExercisesByEquipment(ctx context.Context, equipment types.EquipmentType) ([]types.Exercise, error)
	GetExercisesByDifficulty(ctx context.Context, difficulty types.FitnessLevel) ([]types.Exercise, error)
	GetRecommendedExercises(ctx context.Context, userID int, count int) ([]types.Exercise, error)
	BulkCreateExercises(ctx context.Context, exercises []types.ExerciseRequest) ([]types.Exercise, error)
}

type WorkoutTemplateRepo interface {
	CreateTemplate(ctx context.Context, template *types.WorkoutTemplateRequest) (*types.WorkoutTemplate, error)
	GetTemplateByID(ctx context.Context, templateID int) (*types.WorkoutTemplate, error)
	UpdateTemplate(ctx context.Context, templateID int, template *types.WorkoutTemplateRequest) (*types.WorkoutTemplate, error)
	DeleteTemplate(ctx context.Context, templateID int) error

	ListTemplates(ctx context.Context, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutTemplate], error)
	FilterTemplates(ctx context.Context, filter types.TemplateFilter, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutTemplate], error)
	SearchTemplates(ctx context.Context, query string, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutTemplate], error)

	GetTemplatesByLevel(ctx context.Context, level types.FitnessLevel) ([]types.WorkoutTemplate, error)
	GetTemplatesByGoal(ctx context.Context, goal types.FitnessGoal) ([]types.WorkoutTemplate, error)
	GetRecommendedTemplates(ctx context.Context, userID int, count int) ([]types.WorkoutTemplate, error)
	GetPopularTemplates(ctx context.Context, count int) ([]types.WorkoutTemplate, error)
}

type WeeklySchemaRepo interface {
	GetTemplatesByCoachID(ctx context.Context, coachID string) ([]types.WorkoutTemplate, error)
	CreateWeeklySchema(ctx context.Context, schema *types.WeeklySchemaRequest) (*types.WeeklySchema, error)
	GetWeeklySchemaByID(ctx context.Context, schemaID int) (*types.WeeklySchema, error)
	UpdateWeeklySchema(ctx context.Context, schemaID int, active bool) (*types.WeeklySchema, error)
	DeleteWeeklySchema(ctx context.Context, schemaID int) error
	SaveSchemaAsTemplate(ctx context.Context, schemaID int, templateName string) error
	GetWeeklySchemasByUserID(ctx context.Context, userID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.WeeklySchema], error)
	GetActiveWeeklySchemaByUserID(ctx context.Context, userID int) (*types.WeeklySchema, error)
	GetWeeklySchemaByUserAndWeek(ctx context.Context, userID int, weekStart time.Time) (*types.WeeklySchema, error)

	DeactivateAllWeeklySchemasForUser(ctx context.Context, userID int) error
	GetCurrentWeekSchema(ctx context.Context, userID int) (*types.WeeklySchema, error)
	GetWeeklySchemaHistory(ctx context.Context, userID int, limit int) ([]types.WeeklySchema, error)
}

type WorkoutRepo interface {
	CreateWorkout(ctx context.Context, workout *types.WorkoutRequest) (*types.Workout, error)
	GetWorkoutByID(ctx context.Context, workoutID int) (*types.Workout, error)
	UpdateWorkout(ctx context.Context, workoutID int, workout *types.WorkoutRequest) (*types.Workout, error)
	DeleteWorkout(ctx context.Context, workoutID int) error

	GetWorkoutsBySchemaID(ctx context.Context, schemaID int) ([]types.Workout, error)
	GetWorkoutBySchemaAndDay(ctx context.Context, schemaID int, dayOfWeek int) (*types.Workout, error)
	BulkCreateWorkoutsForSchema(ctx context.Context, schemaID int, workouts []types.WorkoutRequest) ([]types.Workout, error)

	GetWorkoutWithExercises(ctx context.Context, workoutID int) (*types.WorkoutWithExercises, error)
	GetSchemaWithAllWorkouts(ctx context.Context, schemaID int) (*types.WeeklySchemaWithWorkouts, error)
}

type WorkoutExerciseRepo interface {
	CreateWorkoutExercise(ctx context.Context, workoutExercise *types.WorkoutExerciseRequest) (*types.WorkoutExercise, error)
	GetWorkoutExerciseByID(ctx context.Context, weID int) (*types.WorkoutExercise, error)
	UpdateWorkoutExercise(ctx context.Context, weID int, workoutExercise *types.WorkoutExerciseRequest) (*types.WorkoutExercise, error)
	DeleteWorkoutExercise(ctx context.Context, weID int) error

	GetWorkoutExercisesByWorkoutID(ctx context.Context, workoutID int) ([]types.WorkoutExercise, error)
	BulkCreateWorkoutExercisesForWorkout(ctx context.Context, workoutID int, exercises []types.WorkoutExerciseRequest) ([]types.WorkoutExercise, error)
	BulkUpdateWorkoutExercisesForWorkout(ctx context.Context, workoutID int, exercises []types.WorkoutExercise) error
	DeleteAllWorkoutExercisesForWorkout(ctx context.Context, workoutID int) error

	GetMostUsedExercises(ctx context.Context, limit int) ([]types.Exercise, error)
	GetExerciseUsageStats(ctx context.Context, exerciseID int) (map[string]interface{}, error)
}

type ProgressRepo interface {
	CreateProgressLog(ctx context.Context, progress *types.ProgressLogRequest) (*types.ProgressLog, error)
	GetProgressLogByID(ctx context.Context, logID int) (*types.ProgressLog, error)
	UpdateProgressLog(ctx context.Context, logID int, progress *types.ProgressLogRequest) (*types.ProgressLog, error)
	DeleteProgressLog(ctx context.Context, logID int) error

	GetProgressLogsByUserID(ctx context.Context, userID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.ProgressLog], error)
	GetProgressLogsByUserAndExercise(ctx context.Context, userID int, exerciseID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.ProgressLog], error)
	GetProgressLogsByUserAndDateRange(ctx context.Context, userID int, startDate, endDate time.Time) ([]types.ProgressLog, error)

	FilterProgressLogs(ctx context.Context, filter types.ProgressFilter, pagination types.PaginationParams) (*types.PaginatedResponse[types.ProgressLog], error)

	GetUserProgressSummary(ctx context.Context, userID int) (*types.UserProgressSummary, error)
	GetPersonalBests(ctx context.Context, userID int) ([]types.PersonalBest, error)
	GetProgressTrend(ctx context.Context, userID int, exerciseID int, days int) ([]types.ProgressLog, error)
	GetWorkoutStreak(ctx context.Context, userID int) (int, error)

	BulkCreateProgressLogs(ctx context.Context, logs []types.ProgressLogRequest) ([]types.ProgressLog, error)
	GetLatestProgressLogsForUser(ctx context.Context, userID int) ([]types.ProgressLog, error)
}
type FitnessProfileRepo interface {
	CreateFitnessAssessment(ctx context.Context, userID int, assessment *types.FitnessAssessmentRequest) (*types.FitnessAssessment, error)
	GetUserFitnessProfile(ctx context.Context, userID int) (*types.FitnessProfile, error)
	UpdateFitnessLevel(ctx context.Context, userID int, level types.FitnessLevel) error
	UpdateFitnessGoals(ctx context.Context, userID int, goals []types.FitnessGoalTarget) error

	EstimateOneRepMax(ctx context.Context, userID int, exerciseID int, performance *types.PerformanceData) (*types.OneRepMaxEstimate, error)
	GetOneRepMaxHistory(ctx context.Context, userID int, exerciseID int) ([]types.OneRepMaxEstimate, error)
	UpdateOneRepMax(ctx context.Context, userID int, exerciseID int, estimate float64) error

	CreateMovementAssessment(ctx context.Context, userID int, assessment *types.MovementAssessmentRequest) (*types.MovementAssessment, error)
	GetMovementLimitations(ctx context.Context, userID int) ([]types.MovementLimitation, error)
}

type WorkoutSessionRepo interface {
	StartWorkoutSession(ctx context.Context, userID int, workoutID int) (*types.WorkoutSession, error)
	CompleteWorkoutSession(ctx context.Context, sessionID int, summary *types.SessionSummary) (*types.WorkoutSession, error)
	SkipWorkout(ctx context.Context, userID int, workoutID int, reason string) (*types.SkippedWorkout, error)

	LogExercisePerformance(ctx context.Context, sessionID int, exerciseID int, performance *types.ExercisePerformance) error
	GetActiveSession(ctx context.Context, userID int) (*types.WorkoutSession, error)
	GetSessionHistory(ctx context.Context, userID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutSession], error)

	GetSessionMetrics(ctx context.Context, sessionID int) (*types.SessionMetrics, error)
	GetWeeklySessionStats(ctx context.Context, userID int, weekStart time.Time) (*types.WeeklySessionStats, error)
}

type PlanGenerationRepo interface {
	CreatePlanGeneration(ctx context.Context, userID int, authUserID string, metadata *types.PlanGenerationMetadata) (*types.GeneratedPlan, error)
	GetActivePlanForUser(ctx context.Context, userID int) (*types.GeneratedPlan, error)
	GetPlanGenerationHistory(ctx context.Context, userID int, limit int) ([]types.GeneratedPlan, error)
	GetPlanID(ctx context.Context, planID int) (*types.GeneratedPlan, error)

	TrackPlanPerformance(ctx context.Context, planID int, performance *types.PlanPerformanceData) error
	GetPlanEffectivenessScore(ctx context.Context, planID int) (float64, error)
	MarkPlanForRegeneration(ctx context.Context, planID int, reason string) error

	LogPlanAdaptation(ctx context.Context, planID int, adaptation *types.PlanAdaptation) error
	GetAdaptationHistory(ctx context.Context, userID int) ([]types.PlanAdaptation, error)

	CountActivePlans(ctx context.Context, userID int) (int, error)
	SaveGeneratedPlanStructure(ctx context.Context, planID int, structure []types.PlanStructureDayInput) error
	GetGeneratedPlanStructure(ctx context.Context, planID int) ([]types.GeneratedPlanDay, error)
	DeletePlanForUser(ctx context.Context, planID int, authUserID string) error
}

type RecoveryMetricsRepo interface {
	LogRecoveryMetrics(ctx context.Context, userID int, metrics *types.RecoveryMetrics) error
	GetRecoveryStatus(ctx context.Context, userID int) (*types.RecoveryStatus, error)
	GetRecoveryTrend(ctx context.Context, userID int, days int) ([]types.RecoveryMetrics, error)

	CalculateFatigueScore(ctx context.Context, userID int) (float64, error)
	RecommendRestDay(ctx context.Context, userID int) (*types.RestDayRecommendation, error)
	TrackSleepQuality(ctx context.Context, userID int, quality *types.SleepQuality) error
}

type GoalTrackingRepo interface {
	CreateFitnessGoal(ctx context.Context, userID int, goal *types.FitnessGoalRequest) (*types.FitnessGoalTarget, error)
	UpdateGoalProgress(ctx context.Context, goalID int, progress float64) error
	GetActiveGoals(ctx context.Context, userID int) ([]types.FitnessGoalTarget, error)
	CompleteGoal(ctx context.Context, goalID int) error

	CalculateGoalProgress(ctx context.Context, goalID int) (*types.GoalProgress, error)
	EstimateTimeToGoal(ctx context.Context, goalID int) (*types.TimeToGoalEstimate, error)
	SuggestGoalAdjustments(ctx context.Context, userID int) ([]types.GoalAdjustment, error)
}

type CoachAssignmentRepo interface {
	CreateCoachAssignment(ctx context.Context, assignment *types.CoachAssignmentRequest) (*types.CoachAssignment, error)
	GetCoachAssignment(ctx context.Context, assignmentID int) (*types.CoachAssignment, error)
	GetClientsByCoachID(ctx context.Context, coachID string) ([]types.ClientSummary, error)
	GetCoachByUserID(ctx context.Context, userID int) (*types.CoachAssignment, error)
	DeactivateAssignment(ctx context.Context, assignmentID int) error
	IsCoachForUser(ctx context.Context, coachID string, userID int) (bool, error)
	GetCoachDashboard(ctx context.Context, coachID string) (*types.CoachDashboard, error)

	LogCoachActivity(ctx context.Context, activity *types.CoachActivity) error
	GetCoachActivityLog(ctx context.Context, coachID string, limit int) ([]types.CoachActivity, error)
}

type UserRoleRepo interface {
	GetUserRole(ctx context.Context, authUserID string) (types.UserRole, error)
	UpsertUserRole(ctx context.Context, authUserID string, role types.UserRole) error
	BatchUpsertUserRoles(ctx context.Context, roles map[string]types.UserRole) error
	DeleteUserRole(ctx context.Context, authUserID string) error
	GetStaleRoles(ctx context.Context, staleDuration time.Duration) ([]types.UserRoleCache, error)
}

type SchemaRepo interface {
	WorkoutProfiles() WorkoutProfileRepo
	Exercises() ExerciseRepo
	Templates() WorkoutTemplateRepo
	Schemas() WeeklySchemaRepo
	Workouts() WorkoutRepo
	WorkoutExercises() WorkoutExerciseRepo
	Progress() ProgressRepo
	PlanGeneration() PlanGenerationRepo
	RecoveryMetrics() RecoveryMetricsRepo
	GoalTracking() GoalTrackingRepo
	CoachAssignments() CoachAssignmentRepo
	UserRoles() UserRoleRepo

	WithTransaction(ctx context.Context, fn func(context.Context) error) error
}
