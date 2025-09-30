package service

import (
	"context"
	"fmt"
	"time"

	"github.com/tdmdh/fit-up-server/services/schema-service/internal/repository"
	"github.com/tdmdh/fit-up-server/services/schema-service/internal/types"
)

type weeklySchemaService struct {
	repo repository.SchemaRepo
}

// NewWeeklySchemaService creates a new weekly schema service instance
func NewWeeklySchemaService(repo repository.SchemaRepo) WeeklySchemaService {
	return &weeklySchemaService{
		repo: repo,
	}
}

func (s *weeklySchemaService) CreateWeeklySchema(ctx context.Context, schema *types.WeeklySchemaRequest) (*types.WeeklySchema, error) {
	return s.repo.Schemas().CreateWeeklySchema(ctx, schema)
}

func (s *weeklySchemaService) GetWeeklySchemaByID(ctx context.Context, schemaID int) (*types.WeeklySchema, error) {
	return s.repo.Schemas().GetWeeklySchemaByID(ctx, schemaID)
}

func (s *weeklySchemaService) UpdateWeeklySchema(ctx context.Context, schemaID int, active bool) (*types.WeeklySchema, error) {
	return s.repo.Schemas().UpdateWeeklySchema(ctx, schemaID, active)
}

func (s *weeklySchemaService) DeleteWeeklySchema(ctx context.Context, schemaID int) error {
	return s.repo.Schemas().DeleteWeeklySchema(ctx, schemaID)
}

func (s *weeklySchemaService) GetWeeklySchemasByUserID(ctx context.Context, userID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.WeeklySchema], error) {
	return s.repo.Schemas().GetWeeklySchemasByUserID(ctx, userID, pagination)
}

func (s *weeklySchemaService) GetActiveWeeklySchemaByUserID(ctx context.Context, userID int) (*types.WeeklySchema, error) {
	return s.repo.Schemas().GetActiveWeeklySchemaByUserID(ctx, userID)
}

func (s *weeklySchemaService) GetWeeklySchemaByUserAndWeek(ctx context.Context, userID int, weekStart time.Time) (*types.WeeklySchema, error) {
	return s.repo.Schemas().GetWeeklySchemaByUserAndWeek(ctx, userID, weekStart)
}

func (s *weeklySchemaService) GetCurrentWeekSchema(ctx context.Context, userID int) (*types.WeeklySchema, error) {
	return s.repo.Schemas().GetCurrentWeekSchema(ctx, userID)
}

func (s *weeklySchemaService) CreateWeeklySchemaFromTemplate(ctx context.Context, userID, templateID int, weekStart time.Time) (*types.WeeklySchemaWithWorkouts, error) {
	// Start transaction for creating schema from template
	var result *types.WeeklySchemaWithWorkouts
	err := s.repo.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. Get the template
		template, err := s.repo.Templates().GetTemplateByID(txCtx, templateID)
		if err != nil {
			return fmt.Errorf("failed to get template: %w", err)
		}

		// 2. Create weekly schema
		schemaReq := &types.WeeklySchemaRequest{
			UserID:    userID,
			WeekStart: weekStart,
		}

		schema, err := s.repo.Schemas().CreateWeeklySchema(txCtx, schemaReq)
		if err != nil {
			return fmt.Errorf("failed to create weekly schema: %w", err)
		}

		// 3. Get recommended exercises for the user
		exercises, err := s.repo.Exercises().GetRecommendedExercises(txCtx, userID, 50)
		if err != nil {
			return fmt.Errorf("failed to get recommended exercises: %w", err)
		}

		if len(exercises) == 0 {
			return fmt.Errorf("no recommended exercises found for user")
		}

		// 4. Create workouts for each day based on template
		var workouts []types.Workout
		for day := 1; day <= template.DaysPerWeek; day++ {
			workoutReq := &types.WorkoutRequest{
				SchemaID:  schema.SchemaID,
				DayOfWeek: day,
				Focus:     fmt.Sprintf("Day %d", day),
			}

			workout, err := s.repo.Workouts().CreateWorkout(txCtx, workoutReq)
			if err != nil {
				return fmt.Errorf("failed to create workout for day %d: %w", day, err)
			}
			workouts = append(workouts, *workout)

			// 5. Add exercises to each workout (simplified logic)
			exerciseCount := min(5, len(exercises)) // Use up to 5 exercises per workout
			for i := 0; i < exerciseCount; i++ {
				exerciseIdx := (day-1)*exerciseCount + i%len(exercises)
				if exerciseIdx >= len(exercises) {
					exerciseIdx = exerciseIdx % len(exercises)
				}

				workoutExerciseReq := &types.WorkoutExerciseRequest{
					WorkoutID:   workout.WorkoutID,
					ExerciseID:  exercises[exerciseIdx].ExerciseID,
					Sets:        exercises[exerciseIdx].DefaultSets,
					Reps:        exercises[exerciseIdx].DefaultReps,
					RestSeconds: exercises[exerciseIdx].RestSeconds,
				}

				_, err := s.repo.WorkoutExercises().CreateWorkoutExercise(txCtx, workoutExerciseReq)
				if err != nil {
					return fmt.Errorf("failed to create workout exercise: %w", err)
				}
			}
		}

		// 6. Get the complete schema with workouts
		result, err = s.repo.Workouts().GetSchemaWithAllWorkouts(txCtx, schema.SchemaID)
		if err != nil {
			return fmt.Errorf("failed to get complete schema: %w", err)
		}

		return nil
	})

	return result, err
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
