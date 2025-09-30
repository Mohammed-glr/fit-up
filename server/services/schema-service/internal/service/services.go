package service

import (
	"context"
	"time"

	"github.com/tdmdh/fit-up-server/services/schema-service/internal/repository"
	"github.com/tdmdh/fit-up-server/services/schema-service/internal/types"
)

// =============================================================================
// SERVICE INTERFACES
// =============================================================================

type UserService interface {
	CreateUser(ctx context.Context, user *types.WorkoutUserRequest) (*types.WorkoutUser, error)
	GetUserByID(ctx context.Context, userID int) (*types.WorkoutUser, error)
	GetUserByEmail(ctx context.Context, email string) (*types.WorkoutUser, error)
	UpdateUser(ctx context.Context, userID int, user *types.WorkoutUserRequest) (*types.WorkoutUser, error)
	DeleteUser(ctx context.Context, userID int) error
	ListUsers(ctx context.Context, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutUser], error)
	SearchUsers(ctx context.Context, query string, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutUser], error)
	GetUsersByLevel(ctx context.Context, level types.FitnessLevel) ([]types.WorkoutUser, error)
	GetUsersByGoal(ctx context.Context, goal types.FitnessGoal) ([]types.WorkoutUser, error)
}

type ExerciseService interface {
	CreateExercise(ctx context.Context, exercise *types.ExerciseRequest) (*types.Exercise, error)
	GetExerciseByID(ctx context.Context, exerciseID int) (*types.Exercise, error)
	UpdateExercise(ctx context.Context, exerciseID int, exercise *types.ExerciseRequest) (*types.Exercise, error)
	DeleteExercise(ctx context.Context, exerciseID int) error
	ListExercises(ctx context.Context, pagination types.PaginationParams) (*types.PaginatedResponse[types.Exercise], error)
	FilterExercises(ctx context.Context, filter types.ExerciseFilter, pagination types.PaginationParams) (*types.PaginatedResponse[types.Exercise], error)
	SearchExercises(ctx context.Context, query string, pagination types.PaginationParams) (*types.PaginatedResponse[types.Exercise], error)
	GetExercisesByMuscleGroup(ctx context.Context, muscleGroup string) ([]types.Exercise, error)
	GetExercisesByEquipment(ctx context.Context, equipment types.EquipmentType) ([]types.Exercise, error)
	GetRecommendedExercises(ctx context.Context, userID int, count int) ([]types.Exercise, error)
	BulkCreateExercises(ctx context.Context, exercises []types.ExerciseRequest) ([]types.Exercise, error)
}

type WorkoutTemplateService interface {
	CreateTemplate(ctx context.Context, template *types.WorkoutTemplateRequest) (*types.WorkoutTemplate, error)
	GetTemplateByID(ctx context.Context, templateID int) (*types.WorkoutTemplate, error)
	UpdateTemplate(ctx context.Context, templateID int, template *types.WorkoutTemplateRequest) (*types.WorkoutTemplate, error)
	DeleteTemplate(ctx context.Context, templateID int) error
	ListTemplates(ctx context.Context, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutTemplate], error)
	FilterTemplates(ctx context.Context, filter types.TemplateFilter, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutTemplate], error)
	SearchTemplates(ctx context.Context, query string, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutTemplate], error)
	GetRecommendedTemplates(ctx context.Context, userID int, count int) ([]types.WorkoutTemplate, error)
	GetPopularTemplates(ctx context.Context, count int) ([]types.WorkoutTemplate, error)
}

type WeeklySchemaService interface {
	CreateWeeklySchema(ctx context.Context, schema *types.WeeklySchemaRequest) (*types.WeeklySchema, error)
	GetWeeklySchemaByID(ctx context.Context, schemaID int) (*types.WeeklySchema, error)
	UpdateWeeklySchema(ctx context.Context, schemaID int, active bool) (*types.WeeklySchema, error)
	DeleteWeeklySchema(ctx context.Context, schemaID int) error
	GetWeeklySchemasByUserID(ctx context.Context, userID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.WeeklySchema], error)
	GetActiveWeeklySchemaByUserID(ctx context.Context, userID int) (*types.WeeklySchema, error)
	GetWeeklySchemaByUserAndWeek(ctx context.Context, userID int, weekStart time.Time) (*types.WeeklySchema, error)
	CreateWeeklySchemaFromTemplate(ctx context.Context, userID, templateID int, weekStart time.Time) (*types.WeeklySchemaWithWorkouts, error)
	GetCurrentWeekSchema(ctx context.Context, userID int) (*types.WeeklySchema, error)
}

type WorkoutService interface {
	CreateWorkout(ctx context.Context, workout *types.WorkoutRequest) (*types.Workout, error)
	GetWorkoutByID(ctx context.Context, workoutID int) (*types.Workout, error)
	UpdateWorkout(ctx context.Context, workoutID int, workout *types.WorkoutRequest) (*types.Workout, error)
	DeleteWorkout(ctx context.Context, workoutID int) error
	GetWorkoutsBySchemaID(ctx context.Context, schemaID int) ([]types.Workout, error)
	GetWorkoutBySchemaAndDay(ctx context.Context, schemaID int, dayOfWeek int) (*types.Workout, error)
	GetWorkoutWithExercises(ctx context.Context, workoutID int) (*types.WorkoutWithExercises, error)
	GetSchemaWithAllWorkouts(ctx context.Context, schemaID int) (*types.WeeklySchemaWithWorkouts, error)
	BulkCreateWorkoutsForSchema(ctx context.Context, schemaID int, workouts []types.WorkoutRequest) ([]types.Workout, error)
}

type WorkoutExerciseService interface {
	CreateWorkoutExercise(ctx context.Context, workoutExercise *types.WorkoutExerciseRequest) (*types.WorkoutExercise, error)
	GetWorkoutExerciseByID(ctx context.Context, weID int) (*types.WorkoutExercise, error)
	UpdateWorkoutExercise(ctx context.Context, weID int, workoutExercise *types.WorkoutExerciseRequest) (*types.WorkoutExercise, error)
	DeleteWorkoutExercise(ctx context.Context, weID int) error
	GetWorkoutExercisesByWorkoutID(ctx context.Context, workoutID int) ([]types.WorkoutExercise, error)
	BulkCreateWorkoutExercisesForWorkout(ctx context.Context, workoutID int, exercises []types.WorkoutExerciseRequest) ([]types.WorkoutExercise, error)
	BulkUpdateWorkoutExercisesForWorkout(ctx context.Context, workoutID int, exercises []types.WorkoutExercise) error
	GetMostUsedExercises(ctx context.Context, limit int) ([]types.Exercise, error)
}

type ProgressService interface {
	CreateProgressLog(ctx context.Context, progress *types.ProgressLogRequest) (*types.ProgressLog, error)
	GetProgressLogByID(ctx context.Context, logID int) (*types.ProgressLog, error)
	UpdateProgressLog(ctx context.Context, logID int, progress *types.ProgressLogRequest) (*types.ProgressLog, error)
	DeleteProgressLog(ctx context.Context, logID int) error
	GetProgressLogsByUserID(ctx context.Context, userID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.ProgressLog], error)
	GetProgressLogsByUserAndExercise(ctx context.Context, userID int, exerciseID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.ProgressLog], error)
	FilterProgressLogs(ctx context.Context, filter types.ProgressFilter, pagination types.PaginationParams) (*types.PaginatedResponse[types.ProgressLog], error)
	GetUserProgressSummary(ctx context.Context, userID int) (*types.UserProgressSummary, error)
	GetPersonalBests(ctx context.Context, userID int) ([]types.PersonalBest, error)
	GetProgressTrend(ctx context.Context, userID int, exerciseID int, days int) ([]types.ProgressLog, error)
	GetWorkoutStreak(ctx context.Context, userID int) (int, error)
	BulkCreateProgressLogs(ctx context.Context, logs []types.ProgressLogRequest) ([]types.ProgressLog, error)
}

// =============================================================================
// AGGREGATED SERVICE INTERFACE
// =============================================================================

type SchemaService interface {
	Users() UserService
	Exercises() ExerciseService
	Templates() WorkoutTemplateService
	Schemas() WeeklySchemaService
	Workouts() WorkoutService
	WorkoutExercises() WorkoutExerciseService
	Progress() ProgressService
}

// =============================================================================
// SERVICE IMPLEMENTATION
// =============================================================================

type Service struct {
	repo repository.SchemaRepo

	userService            UserService
	exerciseService        ExerciseService
	templateService        WorkoutTemplateService
	weeklySchemaService    WeeklySchemaService
	workoutService         WorkoutService
	workoutExerciseService WorkoutExerciseService
	progressService        ProgressService
}

func NewService(repo repository.SchemaRepo) SchemaService {
	s := &Service{
		repo: repo,
	}

	s.userService = NewUserService(repo)
	s.exerciseService = NewExerciseService(repo)
	s.templateService = NewWorkoutTemplateService(repo)
	s.weeklySchemaService = NewWeeklySchemaService(repo)
	s.workoutService = NewWorkoutService(repo)
	s.workoutExerciseService = NewWorkoutExerciseService(repo)
	s.progressService = NewProgressService(repo)

	return s
}

func (s *Service) Users() UserService {
	return s.userService
}

func (s *Service) Exercises() ExerciseService {
	return s.exerciseService
}

func (s *Service) Templates() WorkoutTemplateService {
	return s.templateService
}

func (s *Service) Schemas() WeeklySchemaService {
	return s.weeklySchemaService
}

func (s *Service) Workouts() WorkoutService {
	return s.workoutService
}

func (s *Service) WorkoutExercises() WorkoutExerciseService {
	return s.workoutExerciseService
}

func (s *Service) Progress() ProgressService {
	return s.progressService
}
