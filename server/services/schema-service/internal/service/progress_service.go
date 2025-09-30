package service

import (
	"context"
	"time"

	"github.com/tdmdh/fit-up-server/services/schema-service/internal/repository"
	"github.com/tdmdh/fit-up-server/services/schema-service/internal/types"
)

type progressService struct {
	repo repository.SchemaRepo
}

// NewProgressService creates a new progress service instance
func NewProgressService(repo repository.SchemaRepo) ProgressService {
	return &progressService{
		repo: repo,
	}
}

func (s *progressService) CreateProgressLog(ctx context.Context, progress *types.ProgressLogRequest) (*types.ProgressLog, error) {
	return s.repo.Progress().CreateProgressLog(ctx, progress)
}

func (s *progressService) GetProgressLogByID(ctx context.Context, logID int) (*types.ProgressLog, error) {
	return s.repo.Progress().GetProgressLogByID(ctx, logID)
}

func (s *progressService) UpdateProgressLog(ctx context.Context, logID int, progress *types.ProgressLogRequest) (*types.ProgressLog, error) {
	return s.repo.Progress().UpdateProgressLog(ctx, logID, progress)
}

func (s *progressService) DeleteProgressLog(ctx context.Context, logID int) error {
	return s.repo.Progress().DeleteProgressLog(ctx, logID)
}

func (s *progressService) GetProgressLogsByUserID(ctx context.Context, userID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.ProgressLog], error) {
	return s.repo.Progress().GetProgressLogsByUserID(ctx, userID, pagination)
}

func (s *progressService) GetLatestProgressLogsForUser(ctx context.Context, userID int) ([]types.ProgressLog, error) {
	return s.repo.Progress().GetLatestProgressLogsForUser(ctx, userID)
}

func (s *progressService) GetProgressLogsByUserAndExercise(ctx context.Context, userID int, exerciseID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.ProgressLog], error) {
	return s.repo.Progress().GetProgressLogsByUserAndExercise(ctx, userID, exerciseID, pagination)
}

func (s *progressService) GetProgressLogsByUserAndDateRange(ctx context.Context, userID int, startDate, endDate time.Time) ([]types.ProgressLog, error) {
	return s.repo.Progress().GetProgressLogsByUserAndDateRange(ctx, userID, startDate, endDate)
}

func (s *progressService) FilterProgressLogs(ctx context.Context, filter types.ProgressFilter, pagination types.PaginationParams) (*types.PaginatedResponse[types.ProgressLog], error) {
	return s.repo.Progress().FilterProgressLogs(ctx, filter, pagination)
}

func (s *progressService) GetUserProgressSummary(ctx context.Context, userID int) (*types.UserProgressSummary, error) {
	return s.repo.Progress().GetUserProgressSummary(ctx, userID)
}

func (s *progressService) GetPersonalBests(ctx context.Context, userID int) ([]types.PersonalBest, error) {
	return s.repo.Progress().GetPersonalBests(ctx, userID)
}

func (s *progressService) GetProgressTrend(ctx context.Context, userID int, exerciseID int, days int) ([]types.ProgressLog, error) {
	return s.repo.Progress().GetProgressTrend(ctx, userID, exerciseID, days)
}

func (s *progressService) GetWorkoutStreak(ctx context.Context, userID int) (int, error) {
	return s.repo.Progress().GetWorkoutStreak(ctx, userID)
}

func (s *progressService) BulkCreateProgressLogs(ctx context.Context, logs []types.ProgressLogRequest) ([]types.ProgressLog, error) {
	return s.repo.Progress().BulkCreateProgressLogs(ctx, logs)
}
