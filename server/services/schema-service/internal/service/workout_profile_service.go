package service

import (
	"context"

	"github.com/tdmdh/fit-up-server/services/schema-service/internal/repository"
	"github.com/tdmdh/fit-up-server/services/schema-service/internal/types"
)

type workoutProfileService struct {
	repo repository.SchemaRepo
}

func NewWorkoutProfileService(repo repository.SchemaRepo) WorkoutProfileService {
	return &workoutProfileService{
		repo: repo,
	}
}

func (s *workoutProfileService) CreateWorkoutProfile(ctx context.Context, authUserID string, profile *types.WorkoutProfileRequest) (*types.WorkoutProfile, error) {
	return s.repo.WorkoutProfiles().CreateWorkoutProfile(ctx, authUserID, profile)
}

func (s *workoutProfileService) GetWorkoutProfileByAuthID(ctx context.Context, authUserID string) (*types.WorkoutProfile, error) {
	return s.repo.WorkoutProfiles().GetWorkoutProfileByAuthID(ctx, authUserID)
}

func (s *workoutProfileService) GetWorkoutProfileByID(ctx context.Context, workoutProfileID int) (*types.WorkoutProfile, error) {
	return s.repo.WorkoutProfiles().GetWorkoutProfileByID(ctx, workoutProfileID)
}

func (s *workoutProfileService) UpdateWorkoutProfile(ctx context.Context, authUserID string, profile *types.WorkoutProfileRequest) (*types.WorkoutProfile, error) {
	return s.repo.WorkoutProfiles().UpdateWorkoutProfile(ctx, authUserID, profile)
}

func (s *workoutProfileService) DeleteWorkoutProfile(ctx context.Context, authUserID string) error {
	return s.repo.WorkoutProfiles().DeleteWorkoutProfile(ctx, authUserID)
}

func (s *workoutProfileService) ListWorkoutProfiles(ctx context.Context, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutProfile], error) {
	return s.repo.WorkoutProfiles().ListWorkoutProfiles(ctx, pagination)
}

func (s *workoutProfileService) SearchWorkoutProfiles(ctx context.Context, query string, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutProfile], error) {
	return s.repo.WorkoutProfiles().SearchWorkoutProfiles(ctx, query, pagination)
}

func (s *workoutProfileService) GetProfilesByLevel(ctx context.Context, level types.FitnessLevel) ([]types.WorkoutProfile, error) {
	return s.repo.WorkoutProfiles().GetProfilesByLevel(ctx, level)
}

func (s *workoutProfileService) GetProfilesByGoal(ctx context.Context, goal types.FitnessGoal) ([]types.WorkoutProfile, error) {
	return s.repo.WorkoutProfiles().GetProfilesByGoal(ctx, goal)
}

func (s *workoutProfileService) CountActiveProfiles(ctx context.Context) (int, error) {
	return s.repo.WorkoutProfiles().CountActiveProfiles(ctx)
}
