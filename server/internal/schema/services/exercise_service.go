package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/tdmdh/fit-up-server/internal/schema/repository"
	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

type exerciseService struct {
	repo repository.SchemaRepo
}

func NewExerciseService(repo repository.SchemaRepo) ExerciseService {
	return &exerciseService{
		repo: repo,
	}
}

func (s *exerciseService) GetExerciseByID(ctx context.Context, exerciseID int) (*types.Exercise, error) {
	return s.repo.Exercises().GetExerciseByID(ctx, exerciseID)
}

func (s *exerciseService) ListExercises(ctx context.Context, pagination types.PaginationParams) (*types.PaginatedResponse[types.Exercise], error) {
	return s.repo.Exercises().ListExercises(ctx, pagination)
}

func (s *exerciseService) FilterExercises(ctx context.Context, filter types.ExerciseFilter, pagination types.PaginationParams) (*types.PaginatedResponse[types.Exercise], error) {
	return s.repo.Exercises().FilterExercises(ctx, filter, pagination)
}

func (s *exerciseService) SearchExercises(ctx context.Context, query string, pagination types.PaginationParams) (*types.PaginatedResponse[types.Exercise], error) {
	if err := validator.New().Var(query, "required,min=1"); err != nil {
		return nil, fmt.Errorf("invalid search query: %w", err)
	}

	return s.repo.Exercises().SearchExercises(ctx, query, pagination)
}

func (s *exerciseService) GetExercisesByMuscleGroup(ctx context.Context, muscleGroup string) ([]types.Exercise, error) {
	if strings.TrimSpace(muscleGroup) == "" {
		return nil, fmt.Errorf("muscle group cannot be empty")
	}

	validMuscleGroups := []string{"chest", "back", "shoulders", "biceps", "triceps", "legs", "glutes", "core", "calves"}
	if !s.isValidMuscleGroup(muscleGroup, validMuscleGroups) {
		return nil, fmt.Errorf("invalid muscle group: %s", muscleGroup)
	}

	return s.repo.Exercises().GetExercisesByMuscleGroup(ctx, muscleGroup)
}

func (s *exerciseService) GetExercisesByEquipment(ctx context.Context, equipment types.EquipmentType) ([]types.Exercise, error) {
	return s.repo.Exercises().GetExercisesByEquipment(ctx, equipment)
}

func (s *exerciseService) GetExercisesByDifficulty(ctx context.Context, difficulty types.FitnessLevel) ([]types.Exercise, error) {
	return s.repo.Exercises().GetExercisesByDifficulty(ctx, difficulty)
}

func (s *exerciseService) GetRecommendedExercises(ctx context.Context, userID int, count int) ([]types.Exercise, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	if count <= 0 || count > 50 {
		count = 10 
	}
	userProfile, err := s.repo.FitnessProfiles().GetUserFitnessProfile(ctx, userID)
	if err != nil {
		fmt.Printf("Warning: Could not get user profile: %v\n", err)
	}

	allExercises, err := s.repo.Exercises().ListExercises(ctx, types.PaginationParams{Page: 1, Limit: 100})
	if err != nil {
		return nil, fmt.Errorf("failed to get exercises: %w", err)
	}

	recommendations := s.generateBasicExerciseRecommendations(allExercises.Data, userProfile, count)

	return recommendations, nil
}


func (s *exerciseService) isValidMuscleGroup(muscleGroup string, validGroups []string) bool {
	mgLower := strings.ToLower(muscleGroup)
	for _, valid := range validGroups {
		if mgLower == valid {
			return true
		}
	}
	return false
}

func (s *exerciseService) generateBasicExerciseRecommendations(exercises []types.Exercise, userProfile *types.FitnessProfile, count int) []types.Exercise {
	var recommendations []types.Exercise

	exerciseCount := 0
	for _, exercise := range exercises {
		if exerciseCount >= count {
			break
		}

		if userProfile != nil {
			if s.isExerciseAppropriate(exercise, userProfile) {
				recommendations = append(recommendations, exercise)
				exerciseCount++
			}
		} else {
			recommendations = append(recommendations, exercise)
			exerciseCount++
		}
	}

	return recommendations
}

func (s *exerciseService) isExerciseAppropriate(exercise types.Exercise, profile *types.FitnessProfile) bool {
	switch string(profile.CurrentLevel) {
	case "beginner":
		return strings.Contains(strings.ToLower(exercise.Name), "beginner") ||
			!strings.Contains(strings.ToLower(exercise.Name), "advanced")
	case "intermediate":
		return !strings.Contains(strings.ToLower(exercise.Name), "advanced")
	case "advanced":
		return true
	}
	return true
}


func (s *exerciseService) CreateWorkoutExercise(ctx context.Context, exercise *types.WorkoutExerciseRequest) (*types.WorkoutExercise, error) {
	return s.repo.WorkoutExercises().CreateWorkoutExercise(ctx, exercise)
}

func (s *exerciseService) GetWorkoutExerciseByID(ctx context.Context, exerciseID int) (*types.WorkoutExercise, error) {
	return s.repo.WorkoutExercises().GetWorkoutExerciseByID(ctx, exerciseID)
}

func (s *exerciseService) UpdateWorkoutExercise(ctx context.Context, exerciseID int, exercise *types.WorkoutExerciseRequest) (*types.WorkoutExercise, error) {
	return s.repo.WorkoutExercises().UpdateWorkoutExercise(ctx, exerciseID, exercise)
}

func (s *exerciseService) DeleteWorkoutExercise(ctx context.Context, exerciseID int) error {
	return s.repo.WorkoutExercises().DeleteWorkoutExercise(ctx, exerciseID)
}

func (s *exerciseService) GetWorkoutExercisesByWorkoutID(ctx context.Context, workoutID int) ([]types.WorkoutExercise, error) {
	return s.repo.WorkoutExercises().GetWorkoutExercisesByWorkoutID(ctx, workoutID)
}

func (s *exerciseService) GetWorkoutExercisesByUserID(ctx context.Context, userID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutExercise], error) {
	return &types.PaginatedResponse[types.WorkoutExercise]{
		Data:       []types.WorkoutExercise{},
		TotalCount: 0,
		Page:       pagination.Page,
		PageSize:   pagination.Limit,
		TotalPages: 0,
	}, nil
}

func (s *exerciseService) CompleteWorkoutExercise(ctx context.Context, exerciseID int, actualReps int, actualWeight float64) error {
	return fmt.Errorf("method not implemented in repository")
}

func (s *exerciseService) GetCompletedWorkoutExercises(ctx context.Context, userID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutExercise], error) {

	return &types.PaginatedResponse[types.WorkoutExercise]{
		Data:       []types.WorkoutExercise{},
		TotalCount: 0,
		Page:       pagination.Page,
		PageSize:   pagination.Limit,
		TotalPages: 0,
	}, nil
}

func (s *exerciseService) GetExerciseUsageStats(ctx context.Context, exerciseID int) (map[string]interface{}, error) {
	return s.repo.WorkoutExercises().GetExerciseUsageStats(ctx, exerciseID)
}

func (s *exerciseService) GetMostUsedExercises(ctx context.Context, limit int) ([]types.Exercise, error) {
	return s.repo.WorkoutExercises().GetMostUsedExercises(ctx, limit)
}
