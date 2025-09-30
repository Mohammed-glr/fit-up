package service

import (
	"context"

	"github.com/tdmdh/fit-up-server/services/schema-service/internal/repository"
	"github.com/tdmdh/fit-up-server/services/schema-service/internal/types"
)

type userService struct {
	repo repository.SchemaRepo
}

func NewUserService(repo repository.SchemaRepo) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *types.WorkoutUserRequest) (*types.WorkoutUser, error) {
	return s.repo.Users().CreateUser(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, userID int) (*types.WorkoutUser, error) {
	return s.repo.Users().GetUserByID(ctx, userID)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*types.WorkoutUser, error) {
	return s.repo.Users().GetUserByEmail(ctx, email)
}

func (s *userService) UpdateUser(ctx context.Context, userID int, user *types.WorkoutUserRequest) (*types.WorkoutUser, error) {
	return s.repo.Users().UpdateUser(ctx, userID, user)
}

func (s *userService) DeleteUser(ctx context.Context, userID int) error {
	return s.repo.Users().DeleteUser(ctx, userID)
}

func (s *userService) ListUsers(ctx context.Context, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutUser], error) {
	return s.repo.Users().ListUsers(ctx, pagination)
}

func (s *userService) SearchUsers(ctx context.Context, query string, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutUser], error) {
	return s.repo.Users().SearchUsers(ctx, query, pagination)
}

func (s *userService) GetUsersByLevel(ctx context.Context, level types.FitnessLevel) ([]types.WorkoutUser, error) {
	return s.repo.Users().GetUsersByLevel(ctx, level)
}

func (s *userService) GetUsersByGoal(ctx context.Context, goal types.FitnessGoal) ([]types.WorkoutUser, error) {
	return s.repo.Users().GetUsersByGoal(ctx, goal)
}
