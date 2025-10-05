package service

import (
	"context"
	"time"

	"github.com/tdmdh/fit-up-server/internal/schema/repository"
	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

type ExerciseService interface {
	GetExerciseByID(ctx context.Context, exerciseID int) (*types.Exercise, error)
	ListExercises(ctx context.Context, pagination types.PaginationParams) (*types.PaginatedResponse[types.Exercise], error)
	FilterExercises(ctx context.Context, filter types.ExerciseFilter, pagination types.PaginationParams) (*types.PaginatedResponse[types.Exercise], error)
	SearchExercises(ctx context.Context, query string, pagination types.PaginationParams) (*types.PaginatedResponse[types.Exercise], error)
	GetExercisesByMuscleGroup(ctx context.Context, muscleGroup string) ([]types.Exercise, error)
	GetExercisesByEquipment(ctx context.Context, equipment types.EquipmentType) ([]types.Exercise, error)
	GetExercisesByDifficulty(ctx context.Context, difficulty types.FitnessLevel) ([]types.Exercise, error)
	GetRecommendedExercises(ctx context.Context, userID int, count int) ([]types.Exercise, error)

	GetWorkoutExerciseByID(ctx context.Context, weID int) (*types.WorkoutExercise, error)
	GetWorkoutExercisesByWorkoutID(ctx context.Context, workoutID int) ([]types.WorkoutExercise, error)
	GetMostUsedExercises(ctx context.Context, limit int) ([]types.Exercise, error)
	GetExerciseUsageStats(ctx context.Context, exerciseID int) (map[string]interface{}, error)
}

type WorkoutService interface {
	GetWorkoutByID(ctx context.Context, workoutID int) (*types.Workout, error)
	GetWorkoutsBySchemaID(ctx context.Context, schemaID int) ([]types.Workout, error)
	GetWorkoutBySchemaAndDay(ctx context.Context, schemaID int, dayOfWeek int) (*types.Workout, error)
	GetWorkoutWithExercises(ctx context.Context, workoutID int) (*types.WorkoutWithExercises, error)
	GetSchemaWithAllWorkouts(ctx context.Context, schemaID int) (*types.WeeklySchemaWithWorkouts, error)
}

type FitnessProfileService interface {
	CreateFitnessAssessment(ctx context.Context, userID int, assessment *types.FitnessAssessmentRequest) (*types.FitnessAssessment, error)
	GetUserFitnessProfile(ctx context.Context, userID int) (*types.FitnessProfile, error)
	UpdateFitnessLevel(ctx context.Context, userID int, level types.FitnessLevel) error
	UpdateFitnessGoals(ctx context.Context, userID int, goals []types.FitnessGoalTarget) error
	EstimateOneRepMax(ctx context.Context, userID int, exerciseID int, performance *types.PerformanceData) (*types.OneRepMaxEstimate, error)
	GetOneRepMaxHistory(ctx context.Context, userID int, exerciseID int) ([]types.OneRepMaxEstimate, error)
	UpdateOneRepMax(ctx context.Context, userID int, exerciseID int, estimate float64) error
	CreateMovementAssessment(ctx context.Context, userID int, assessment *types.MovementAssessmentRequest) (*types.MovementAssessment, error)
	GetMovementLimitations(ctx context.Context, userID int) ([]types.MovementLimitation, error)

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

	CreateFitnessGoal(ctx context.Context, userID int, goal *types.FitnessGoalRequest) (*types.FitnessGoalTarget, error)
	UpdateGoalProgress(ctx context.Context, goalID int, progress float64) error
	GetActiveGoals(ctx context.Context, userID int) ([]types.FitnessGoalTarget, error)
	CompleteGoal(ctx context.Context, goalID int) error
	CalculateGoalProgress(ctx context.Context, goalID int) (*types.GoalProgress, error)
	EstimateTimeToGoal(ctx context.Context, goalID int) (*types.TimeToGoalEstimate, error)
	SuggestGoalAdjustments(ctx context.Context, userID int) ([]types.GoalAdjustment, error)
}

type WorkoutSessionService interface {
	StartWorkoutSession(ctx context.Context, userID int, workoutID int) (*types.WorkoutSession, error)
	CompleteWorkoutSession(ctx context.Context, sessionID int, summary *types.SessionSummary) (*types.WorkoutSession, error)
	SkipWorkout(ctx context.Context, userID int, workoutID int, reason string) (*types.SkippedWorkout, error)
	LogExercisePerformance(ctx context.Context, sessionID int, exerciseID int, performance *types.ExercisePerformance) error
	GetActiveSession(ctx context.Context, userID int) (*types.WorkoutSession, error)
	GetSessionHistory(ctx context.Context, userID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutSession], error)
	GetSessionMetrics(ctx context.Context, sessionID int) (*types.SessionMetrics, error)
	GetWeeklySessionStats(ctx context.Context, userID int, weekStart time.Time) (*types.WeeklySessionStats, error)
}

type PlanGenerationService interface {
	CreatePlanGeneration(ctx context.Context, userID int, metadata *types.PlanGenerationMetadata) (*types.GeneratedPlan, error)
	GetActivePlanForUser(ctx context.Context, userID int) (*types.GeneratedPlan, error)
	GetPlanGenerationHistory(ctx context.Context, userID int, limit int) ([]types.GeneratedPlan, error)
	TrackPlanPerformance(ctx context.Context, planID int, performance *types.PlanPerformanceData) error
	GetPlanEffectivenessScore(ctx context.Context, planID int) (float64, error)
	MarkPlanForRegeneration(ctx context.Context, planID int, reason string) error
	LogPlanAdaptation(ctx context.Context, planID int, adaptation *types.PlanAdaptation) error
	GetAdaptationHistory(ctx context.Context, userID int) ([]types.PlanAdaptation, error)

	ExportPlanToPDF(ctx context.Context, planID int) ([]byte, error)

	GetTemplateByID(ctx context.Context, templateID int) (*types.WorkoutTemplate, error)
	ListTemplates(ctx context.Context, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutTemplate], error)
	FilterTemplates(ctx context.Context, filter types.TemplateFilter, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutTemplate], error)
	SearchTemplates(ctx context.Context, query string, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutTemplate], error)
	GetTemplatesByLevel(ctx context.Context, level types.FitnessLevel) ([]types.WorkoutTemplate, error)
	GetTemplatesByGoal(ctx context.Context, goal types.FitnessGoal) ([]types.WorkoutTemplate, error)
	GetRecommendedTemplates(ctx context.Context, userID int, count int) ([]types.WorkoutTemplate, error)
	GetPopularTemplates(ctx context.Context, count int) ([]types.WorkoutTemplate, error)

	GetWeeklySchemaByID(ctx context.Context, schemaID int) (*types.WeeklySchema, error)
	GetWeeklySchemasByUserID(ctx context.Context, userID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.WeeklySchema], error)
	GetActiveWeeklySchemaByUserID(ctx context.Context, userID int) (*types.WeeklySchema, error)
	GetWeeklySchemaByUserAndWeek(ctx context.Context, userID int, weekStart time.Time) (*types.WeeklySchema, error)
	GetCurrentWeekSchema(ctx context.Context, userID int) (*types.WeeklySchema, error)
	GetWeeklySchemaHistory(ctx context.Context, userID int, limit int) ([]types.WeeklySchema, error)
	CreateWeeklySchemaFromTemplate(ctx context.Context, userID, templateID int, weekStart time.Time) (*types.WeeklySchemaWithWorkouts, error)
}

type CoachService interface {
	AssignClientToCoach(ctx context.Context, req *types.CoachAssignmentRequest) (*types.CoachAssignment, error)
	GetCoachClients(ctx context.Context, coachID string) ([]types.ClientSummary, error)
	GetCoachDashboard(ctx context.Context, coachID string) (*types.CoachDashboard, error)
	RemoveClientFromCoach(ctx context.Context, assignmentID int) error
	CreateManualSchemaForClient(ctx context.Context, coachID string, req *types.ManualSchemaRequest) (*types.WeeklySchemaExtended, error)
	UpdateManualSchema(ctx context.Context, coachID string, schemaID int, req *types.ManualSchemaRequest) (*types.WeeklySchemaExtended, error)
	DeleteSchema(ctx context.Context, coachID string, schemaID int) error
	CloneSchemaToClient(ctx context.Context, coachID string, sourceSchemaID int, targetUserID int) (*types.WeeklySchemaExtended, error)
	SaveSchemaAsTemplate(ctx context.Context, coachID string, schemaID int, templateName string) error
	GetCoachTemplates(ctx context.Context, coachID string) ([]types.WorkoutTemplate, error)
	CreateSchemaFromCoachTemplate(ctx context.Context, coachID string, templateID int, userID int) (*types.WeeklySchemaExtended, error)
	GetClientProgress(ctx context.Context, coachID string, userID int) (*types.UserProgressSummary, error)	
	ValidateCoachPermission(ctx context.Context, coachID string, userID int) error
}



type SchemaService interface {
	Exercises() ExerciseService
	Workouts() WorkoutService
	Coaches() CoachService
	FitnessProfiles() FitnessProfileService
	WorkoutSessions() WorkoutSessionService
	PlanGeneration() PlanGenerationService
}

type Service struct {
	repo repository.SchemaRepo

	exerciseService ExerciseService
	workoutService  WorkoutService
	coachService    CoachService
	fitnessProfileService FitnessProfileService
	workoutSessionService WorkoutSessionService
	planGenerationService PlanGenerationService
}

func NewService(repo repository.SchemaRepo) SchemaService {
	return &Service{
		repo:                  repo,
		exerciseService:       NewExerciseService(repo),
		workoutService:        NewWorkoutService(repo),
		fitnessProfileService: NewFitnessProfileService(repo),
		workoutSessionService: NewWorkoutSessionService(repo),
		planGenerationService: NewPlanGenerationService(repo),
		coachService:          NewCoachService(repo),
	}
}

func (s *Service) Exercises() ExerciseService {
	return s.exerciseService
}

func (s *Service) Workouts() WorkoutService {
	return s.workoutService
}

func (s *Service) FitnessProfiles() FitnessProfileService {
	return s.fitnessProfileService
}

func (s *Service) WorkoutSessions() WorkoutSessionService {
	return s.workoutSessionService
}

func (s *Service) PlanGeneration() PlanGenerationService {
	return s.planGenerationService
}

func (s *Service) Coaches() CoachService {
	return s.coachService
}
