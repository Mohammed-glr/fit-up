package service

import (
	"context"

	"github.com/tdmdh/fit-up-server/services/schema-service/internal/repository"
	"github.com/tdmdh/fit-up-server/services/schema-service/internal/types"
)

type exerciseService struct {
	repo repository.SchemaRepo
}

func NewExerciseService(repo repository.SchemaRepo) ExerciseService {
	return &exerciseService{
		repo: repo,
	}
}

func (s *exerciseService) CreateExercise(ctx context.Context, exercise *types.ExerciseRequest) (*types.Exercise, error) {
	return s.repo.Exercises().CreateExercise(ctx, exercise)
}

func (s *exerciseService) GetExerciseByID(ctx context.Context, exerciseID int) (*types.Exercise, error) {
	return s.repo.Exercises().GetExerciseByID(ctx, exerciseID)
}

func (s *exerciseService) UpdateExercise(ctx context.Context, exerciseID int, exercise *types.ExerciseRequest) (*types.Exercise, error) {
	return s.repo.Exercises().UpdateExercise(ctx, exerciseID, exercise)
}

func (s *exerciseService) DeleteExercise(ctx context.Context, exerciseID int) error {
	return s.repo.Exercises().DeleteExercise(ctx, exerciseID)
}

func (s *exerciseService) ListExercises(ctx context.Context, pagination types.PaginationParams) (*types.PaginatedResponse[types.Exercise], error) {
	return s.repo.Exercises().ListExercises(ctx, pagination)
}

func (s *exerciseService) FilterExercises(ctx context.Context, filter types.ExerciseFilter, pagination types.PaginationParams) (*types.PaginatedResponse[types.Exercise], error) {
	return s.repo.Exercises().FilterExercises(ctx, filter, pagination)
}

func (s *exerciseService) SearchExercises(ctx context.Context, query string, pagination types.PaginationParams) (*types.PaginatedResponse[types.Exercise], error) {
	return s.repo.Exercises().SearchExercises(ctx, query, pagination)
}

func (s *exerciseService) GetExercisesByMuscleGroup(ctx context.Context, muscleGroup string) ([]types.Exercise, error) {
	return s.repo.Exercises().GetExercisesByMuscleGroup(ctx, muscleGroup)
}

func (s *exerciseService) GetExercisesByEquipment(ctx context.Context, equipment types.EquipmentType) ([]types.Exercise, error) {
	return s.repo.Exercises().GetExercisesByEquipment(ctx, equipment)
}

func (s *exerciseService) GetExercisesByDifficulty(ctx context.Context, difficulty types.FitnessLevel) ([]types.Exercise, error) {
	return s.repo.Exercises().GetExercisesByDifficulty(ctx, difficulty)
}

func (s *exerciseService) GetRecommendedExercises(ctx context.Context, userID int, count int) ([]types.Exercise, error) {
	return s.repo.Exercises().GetRecommendedExercises(ctx, userID, count)
}

func (s *exerciseService) BulkCreateExercises(ctx context.Context, exercises []types.ExerciseRequest) ([]types.Exercise, error) {
	return s.repo.Exercises().BulkCreateExercises(ctx, exercises)
}
