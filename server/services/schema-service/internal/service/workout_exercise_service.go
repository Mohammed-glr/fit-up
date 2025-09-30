package service

import (
	"context"

	"github.com/tdmdh/fit-up-server/services/schema-service/internal/repository"
	"github.com/tdmdh/fit-up-server/services/schema-service/internal/types"
)

type workoutExerciseService struct {
	repo repository.SchemaRepo
}

// NewWorkoutExerciseService creates a new workout exercise service instance
func NewWorkoutExerciseService(repo repository.SchemaRepo) WorkoutExerciseService {
	return &workoutExerciseService{
		repo: repo,
	}
}

func (s *workoutExerciseService) CreateWorkoutExercise(ctx context.Context, workoutExercise *types.WorkoutExerciseRequest) (*types.WorkoutExercise, error) {
	return s.repo.WorkoutExercises().CreateWorkoutExercise(ctx, workoutExercise)
}

func (s *workoutExerciseService) GetWorkoutExerciseByID(ctx context.Context, weID int) (*types.WorkoutExercise, error) {
	return s.repo.WorkoutExercises().GetWorkoutExerciseByID(ctx, weID)
}

func (s *workoutExerciseService) UpdateWorkoutExercise(ctx context.Context, weID int, workoutExercise *types.WorkoutExerciseRequest) (*types.WorkoutExercise, error) {
	return s.repo.WorkoutExercises().UpdateWorkoutExercise(ctx, weID, workoutExercise)
}

func (s *workoutExerciseService) DeleteWorkoutExercise(ctx context.Context, weID int) error {
	return s.repo.WorkoutExercises().DeleteWorkoutExercise(ctx, weID)
}

func (s *workoutExerciseService) GetWorkoutExercisesByWorkoutID(ctx context.Context, workoutID int) ([]types.WorkoutExercise, error) {
	return s.repo.WorkoutExercises().GetWorkoutExercisesByWorkoutID(ctx, workoutID)
}

func (s *workoutExerciseService) BulkCreateWorkoutExercisesForWorkout(ctx context.Context, workoutID int, exercises []types.WorkoutExerciseRequest) ([]types.WorkoutExercise, error) {
	return s.repo.WorkoutExercises().BulkCreateWorkoutExercisesForWorkout(ctx, workoutID, exercises)
}

func (s *workoutExerciseService) BulkUpdateWorkoutExercisesForWorkout(ctx context.Context, workoutID int, exercises []types.WorkoutExercise) error {
	return s.repo.WorkoutExercises().BulkUpdateWorkoutExercisesForWorkout(ctx, workoutID, exercises)
}

func (s *workoutExerciseService) DeleteAllWorkoutExercisesForWorkout(ctx context.Context, workoutID int) error {
	return s.repo.WorkoutExercises().DeleteAllWorkoutExercisesForWorkout(ctx, workoutID)
}

func (s *workoutExerciseService) GetExerciseUsageStats(ctx context.Context, exerciseID int) (map[string]interface{}, error) {
	return s.repo.WorkoutExercises().GetExerciseUsageStats(ctx, exerciseID)
}

func (s *workoutExerciseService) GetMostUsedExercises(ctx context.Context, limit int) ([]types.Exercise, error) {
	return s.repo.WorkoutExercises().GetMostUsedExercises(ctx, limit)
}
