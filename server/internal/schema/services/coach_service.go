package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/tdmdh/fit-up-server/internal/schema/repository"
	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

type coachService struct {
	repo repository.SchemaRepo
	validator *validator.Validate
}

func NewCoachService(repo repository.SchemaRepo) CoachService {
	return &coachService{
		repo: repo,
		validator: validator.New(),
	}
}

func (s *coachService) CreateManualSchemaForClient(ctx context.Context, coachID string, req *types.ManualSchemaRequest) (*types.WeeklySchemaExtended, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	
	if err := s.ValidateCoachPermission(ctx, coachID, req.UserID); err != nil {
		return nil, err
	}
	
	schemaReq := &types.WeeklySchemaRequest{
		UserID:    req.UserID,
		WeekStart: req.StartDate,
	}
	
	schema, err := s.repo.Schemas().CreateWeeklySchema(ctx, schemaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}
	
	var workouts []types.Workout
	for _, workoutReq := range req.Workouts {
		workout := &types.WorkoutRequest{
			SchemaID:  schema.SchemaID,
			DayOfWeek: workoutReq.DayOfWeek,
			Focus:     workoutReq.Focus,
		}
		
		createdWorkout, err := s.repo.Workouts().CreateWorkout(ctx, workout)
		if err != nil {
			return nil, fmt.Errorf("failed to create workout: %w", err)
		}
		
		for _, exReq := range workoutReq.Exercises {
			exerciseReq := &types.WorkoutExerciseRequest{
				WorkoutID:   createdWorkout.WorkoutID,
				ExerciseID:  exReq.ExerciseID,
				Sets:        exReq.Sets,
				Reps:        exReq.Reps,
				RestSeconds: exReq.RestSeconds,
			}
			
			_, err := s.repo.WorkoutExercises().CreateWorkoutExercise(ctx, exerciseReq)
			if err != nil {
				return nil, fmt.Errorf("failed to add exercise: %w", err)
			}
		}
		
		workouts = append(workouts, *createdWorkout)
	}
	
	activity := &types.CoachActivity{
		ActivityType: "schema_created",
		UserID:       req.UserID,
		Description:  fmt.Sprintf("Created manual schema: %s", req.Name),
		Timestamp:    time.Now(),
	}
	_ = s.repo.CoachAssignments().LogCoachActivity(ctx, activity)
	

	return &types.WeeklySchemaExtended{
		WeeklySchema: *schema,
		CoachID:      &coachID,

	}, nil
}

func (s *coachService) UpdateManualSchema(ctx context.Context, coachID string, schemaID int, req *types.ManualSchemaRequest) (*types.WeeklySchemaExtended, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	isCoach, err := s.repo.CoachAssignments().IsCoachForUser(ctx, coachID, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check coach permission: %w", err)
	}
	if !isCoach {
		return nil, fmt.Errorf("coach %s is not authorized for user %d", coachID, req.UserID)
	}

	schema, err := s.repo.Schemas().GetWeeklySchemaByID(ctx, schemaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get schema: %w", err)
	}
	if schema.UserID != req.UserID {
		return nil, fmt.Errorf("schema does not belong to user %d", req.UserID)
	}
	if !schema.WeekStart.Equal(req.StartDate) {
		schema.WeekStart = req.StartDate
		updatedSchema, err := s.repo.Schemas().UpdateWeeklySchema(ctx, schema.SchemaID, true)
		if err != nil {
			return nil, fmt.Errorf("failed to update schema: %w", err)
		}
		schema = updatedSchema
	}

	workoutsToDelete, err := s.repo.Workouts().GetWorkoutsBySchemaID(ctx, schema.SchemaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workouts: %w", err)
	}
	for _, workout := range workoutsToDelete {
		if err := s.repo.Workouts().DeleteWorkout(ctx, workout.WorkoutID); err != nil {
			return nil, fmt.Errorf("failed to delete existing workouts: %w", err)
		}
	}
	var workouts []types.Workout
	for _, workoutReq := range req.Workouts {
		workout := &types.WorkoutRequest{
			SchemaID:  schema.SchemaID,
			DayOfWeek: workoutReq.DayOfWeek,
			Focus:     workoutReq.Focus,
		}
		createdWorkout, err := s.repo.Workouts().CreateWorkout(ctx, workout)
		if err != nil {
			return nil, fmt.Errorf("failed to create workout: %w", err)
		}
		for _, exReq := range workoutReq.Exercises {
			exerciseReq := &types.WorkoutExerciseRequest{
				WorkoutID:   createdWorkout.WorkoutID,
				ExerciseID:  exReq.ExerciseID,
				Sets:        exReq.Sets,
				Reps:        exReq.Reps,
				RestSeconds: exReq.RestSeconds,
			}
			_, err := s.repo.WorkoutExercises().CreateWorkoutExercise(ctx, exerciseReq)
			if err != nil {
				return nil, fmt.Errorf("failed to add exercise: %w", err)
			}
		}
		workouts = append(workouts, *createdWorkout)
	}
	activity := &types.CoachActivity{
		ActivityType: "schema_updated",
		UserID:       req.UserID,
		Description:  fmt.Sprintf("Updated manual schema: %s", req.Name),
		Timestamp:    time.Now(),
	}
	_ = s.repo.CoachAssignments().LogCoachActivity(ctx, activity)
	return &types.WeeklySchemaExtended{
		WeeklySchema: *schema,
		CoachID:      &coachID,
	}, nil
}

func (s *coachService) DeleteSchema(ctx context.Context, coachID string, schemaID int) error {
	schema, err := s.repo.Schemas().GetWeeklySchemaByID(ctx, schemaID)
	if err != nil {
		return fmt.Errorf("failed to get schema: %w", err)
	}
	isCoach, err := s.repo.CoachAssignments().IsCoachForUser(ctx, coachID, schema.UserID)
	if err != nil {
		return fmt.Errorf("failed to check coach permission: %w", err)
	}
	if !isCoach {
		return fmt.Errorf("coach %s is not authorized for schema %d", coachID, schemaID)
	}
	if err := s.repo.Schemas().DeleteWeeklySchema(ctx, schemaID); err != nil {
		return fmt.Errorf("failed to delete schema: %w", err)
	}
	return nil
}

func (s *coachService) CloneSchemaToClient(ctx context.Context, coachID string, sourceSchemaID int, targetUserID int) (*types.WeeklySchemaExtended, error) {
	sourceSchema, err := s.repo.Schemas().GetWeeklySchemaByID(ctx, sourceSchemaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source schema: %w", err)
	}
	isCoach, err := s.repo.CoachAssignments().IsCoachForUser(ctx, coachID, sourceSchema.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check coach permission: %w", err)
	}
	if !isCoach {
		return nil, fmt.Errorf("coach %s is not authorized for schema %d", coachID, sourceSchemaID)
	}
	clonedSchemaReq := &types.WeeklySchemaRequest{
		UserID:    targetUserID,
		WeekStart: sourceSchema.WeekStart,
		
	}
	createdSchema, err := s.repo.Schemas().CreateWeeklySchema(ctx, clonedSchemaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloned schema: %w", err)
	}
	return &types.WeeklySchemaExtended{
		WeeklySchema: *createdSchema,
		CoachID:      &coachID,
	}, nil
}

func (s *coachService) SaveSchemaAsTemplate(ctx context.Context, coachID string, schemaID int, templateName string) error {
	schema, err := s.repo.Schemas().GetWeeklySchemaByID(ctx, schemaID)
	if err != nil {
		return fmt.Errorf("failed to get schema: %w", err)
	}
	isCoach, err := s.repo.CoachAssignments().IsCoachForUser(ctx, coachID, schema.UserID)
	if err != nil {
		return fmt.Errorf("failed to check coach permission: %w", err)
	}
	if !isCoach {
		return fmt.Errorf("coach %s is not authorized for schema %d", coachID, schemaID)
	}
	
	if err := s.repo.Schemas().SaveSchemaAsTemplate(ctx, schemaID, templateName); err != nil {
		return fmt.Errorf("failed to save schema as template: %w", err)
	}
	return nil

}

func (s *coachService) ValidateCoachPermission(ctx context.Context, coachID string, userID int) error {
	isCoach, err := s.repo.CoachAssignments().IsCoachForUser(ctx, coachID, userID)
	if err != nil {
		return fmt.Errorf("failed to check coach permission: %w", err)
	}
	
	if !isCoach {
		return fmt.Errorf("coach %s is not authorized for user %d", coachID, userID)
	}
	
	return nil
}

func (s *coachService) GetCoachDashboard(ctx context.Context, coachID string) (*types.CoachDashboard, error) {
	dashboard, err := s.repo.CoachAssignments().GetCoachDashboard(ctx, coachID)
	if err != nil {
		return nil, fmt.Errorf("failed to get coach dashboard: %w", err)
	}
	return dashboard, nil
}
func (s *coachService) GetCoachClients(ctx context.Context, coachID string) ([]types.ClientSummary, error) {
	return s.repo.CoachAssignments().GetClientsByCoachID(ctx, coachID)
}

func (s *coachService) AssignClientToCoach(ctx context.Context, req *types.CoachAssignmentRequest) (*types.CoachAssignment, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	assignment, err := s.repo.CoachAssignments().CreateCoachAssignment(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to assign client to coach: %w", err)
	}
	return assignment, nil
}

func (s *coachService) RemoveClientFromCoach(ctx context.Context, assignmentID int) error {
	return s.repo.CoachAssignments().DeactivateAssignment(ctx, assignmentID)
}

func (s *coachService) GetCoachTemplates(ctx context.Context, coachID string) ([]types.WorkoutTemplate, error) {
	return s.repo.Schemas().GetTemplatesByCoachID(ctx, coachID)
}

func (s *coachService) CreateSchemaFromCoachTemplate(ctx context.Context, coachID string, templateID int, userID int) (*types.WeeklySchemaExtended, error) {
	_, err := s.repo.Templates().GetTemplateByID(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}
	isCoach, err := s.repo.CoachAssignments().IsCoachForUser(ctx, coachID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check coach permission: %w", err)
	}
	if !isCoach {
		return nil, fmt.Errorf("coach %s is not authorized for user %d", coachID, userID)
	}
	schemaReq := &types.WeeklySchemaRequest{
		UserID:    userID,
		WeekStart: time.Now().AddDate(0, 0, -int(time.Now().Weekday())),
	}
	createdSchema, err := s.repo.Schemas().CreateWeeklySchema(ctx, schemaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create schema from template: %w", err)
	}
	return &types.WeeklySchemaExtended{
		WeeklySchema: *createdSchema,
		CoachID:      &coachID,
	}, nil
}

func (s *coachService) GetClientProgress(ctx context.Context, coachID string, userID int) (*types.UserProgressSummary, error) {
	isCoach, err := s.repo.CoachAssignments().IsCoachForUser(ctx, coachID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check coach permission: %w", err)
	}
	if !isCoach {
		return nil, fmt.Errorf("coach %s is not authorized for user %d", coachID, userID)
	}
	

	pagination := types.PaginationParams{
		Page:     1,
		PageSize: 100,
	}
	progress, err := s.repo.Progress().GetProgressLogsByUserID(ctx, userID, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to get user progress: %w", err)
	}

	
	return &types.UserProgressSummary{
		UserID:       userID,
		TotalWorkouts: progress.TotalCount,
		LastWorkout:  nil,
	}, nil
}
