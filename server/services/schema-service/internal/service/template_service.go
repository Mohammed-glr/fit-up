package service

import (
	"context"

	"github.com/tdmdh/fit-up-server/services/schema-service/internal/repository"
	"github.com/tdmdh/fit-up-server/services/schema-service/internal/types"
)

type workoutTemplateService struct {
	repo repository.SchemaRepo
}

// NewWorkoutTemplateService creates a new workout template service instance
func NewWorkoutTemplateService(repo repository.SchemaRepo) WorkoutTemplateService {
	return &workoutTemplateService{
		repo: repo,
	}
}

func (s *workoutTemplateService) CreateTemplate(ctx context.Context, template *types.WorkoutTemplateRequest) (*types.WorkoutTemplate, error) {
	return s.repo.Templates().CreateTemplate(ctx, template)
}

func (s *workoutTemplateService) GetTemplateByID(ctx context.Context, templateID int) (*types.WorkoutTemplate, error) {
	return s.repo.Templates().GetTemplateByID(ctx, templateID)
}

func (s *workoutTemplateService) UpdateTemplate(ctx context.Context, templateID int, template *types.WorkoutTemplateRequest) (*types.WorkoutTemplate, error) {
	return s.repo.Templates().UpdateTemplate(ctx, templateID, template)
}

func (s *workoutTemplateService) DeleteTemplate(ctx context.Context, templateID int) error {
	return s.repo.Templates().DeleteTemplate(ctx, templateID)
}

func (s *workoutTemplateService) ListTemplates(ctx context.Context, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutTemplate], error) {
	return s.repo.Templates().ListTemplates(ctx, pagination)
}

func (s *workoutTemplateService) FilterTemplates(ctx context.Context, filter types.TemplateFilter, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutTemplate], error) {
	return s.repo.Templates().FilterTemplates(ctx, filter, pagination)
}

func (s *workoutTemplateService) SearchTemplates(ctx context.Context, query string, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutTemplate], error) {
	return s.repo.Templates().SearchTemplates(ctx, query, pagination)
}

func (s *workoutTemplateService) GetRecommendedTemplates(ctx context.Context, userID int, count int) ([]types.WorkoutTemplate, error) {
	return s.repo.Templates().GetRecommendedTemplates(ctx, userID, count)
}

func (s *workoutTemplateService) GetTemplatesByGoal(ctx context.Context, goal types.FitnessGoal) ([]types.WorkoutTemplate, error) {
	return s.repo.Templates().GetTemplatesByGoal(ctx, goal)
}

func (s *workoutTemplateService) GetTemplatesByLevel(ctx context.Context, level types.FitnessLevel) ([]types.WorkoutTemplate, error) {
	return s.repo.Templates().GetTemplatesByLevel(ctx, level)
}

func (s *workoutTemplateService) GetPopularTemplates(ctx context.Context, count int) ([]types.WorkoutTemplate, error) {
	return s.repo.Templates().GetPopularTemplates(ctx, count)
}
