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
	PlanGeneration() PlanGenerationService
}

type Service struct {
	repo repository.SchemaRepo

	exerciseService ExerciseService
	workoutService  WorkoutService
	coachService    CoachService
	planGenerationService PlanGenerationService
}

func NewService(repo repository.SchemaRepo) SchemaService {
	return &Service{
		repo:                  repo,
		exerciseService:       NewExerciseService(repo),
		workoutService:        NewWorkoutService(repo),
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



func (s *Service) PlanGeneration() PlanGenerationService {
	return s.planGenerationService
}

func (s *Service) Coaches() CoachService {
	return s.coachService
}
