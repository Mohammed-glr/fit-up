package service

import (
	"context"

	"github.com/tdmdh/fit-up-server/services/schema-service/internal/repository"
	"github.com/tdmdh/fit-up-server/services/schema-service/internal/types"
)

type workoutService struct {
	repo repository.SchemaRepo
}

// NewWorkoutService creates a new workout service instance
func NewWorkoutService(repo repository.SchemaRepo) WorkoutService {
	return &workoutService{
		repo: repo,
	}
}

func (s *workoutService) CreateWorkout(ctx context.Context, workout *types.WorkoutRequest) (*types.Workout, error) {
	return s.repo.Workouts().CreateWorkout(ctx, workout)
}

func (s *workoutService) GetWorkoutByID(ctx context.Context, workoutID int) (*types.Workout, error) {
	return s.repo.Workouts().GetWorkoutByID(ctx, workoutID)
}

func (s *workoutService) UpdateWorkout(ctx context.Context, workoutID int, workout *types.WorkoutRequest) (*types.Workout, error) {
	return s.repo.Workouts().UpdateWorkout(ctx, workoutID, workout)
}

func (s *workoutService) DeleteWorkout(ctx context.Context, workoutID int) error {
	return s.repo.Workouts().DeleteWorkout(ctx, workoutID)
}

func (s *workoutService) GetWorkoutsBySchemaID(ctx context.Context, schemaID int) ([]types.Workout, error) {
	return s.repo.Workouts().GetWorkoutsBySchemaID(ctx, schemaID)
}

func (s *workoutService) GetWorkoutBySchemaAndDay(ctx context.Context, schemaID int, dayOfWeek int) (*types.Workout, error) {
	return s.repo.Workouts().GetWorkoutBySchemaAndDay(ctx, schemaID, dayOfWeek)
}

func (s *workoutService) GetWorkoutWithExercises(ctx context.Context, workoutID int) (*types.WorkoutWithExercises, error) {
	return s.repo.Workouts().GetWorkoutWithExercises(ctx, workoutID)
}

func (s *workoutService) GetSchemaWithAllWorkouts(ctx context.Context, schemaID int) (*types.WeeklySchemaWithWorkouts, error) {
	return s.repo.Workouts().GetSchemaWithAllWorkouts(ctx, schemaID)
}

func (s *workoutService) BulkCreateWorkoutsForSchema(ctx context.Context, schemaID int, workouts []types.WorkoutRequest) ([]types.Workout, error) {
	return s.repo.Workouts().BulkCreateWorkoutsForSchema(ctx, schemaID, workouts)
}
