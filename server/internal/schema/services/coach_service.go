package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/tdmdh/fit-up-server/internal/schema/repository"
	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

type coachService struct {
	repo      repository.SchemaRepo
	validator *validator.Validate
}

func NewCoachService(repo repository.SchemaRepo) CoachService {
	return &coachService{
		repo:      repo,
		validator: validator.New(),
	}
}

func (s *coachService) CreateManualSchemaForClient(ctx context.Context, coachID string, req *types.ManualSchemaRequest) (*types.WeeklySchemaExtended, error) {
	reqJSON, _ := json.Marshal(req)
	log.Printf("CreateManualSchemaForClient - Request data: %s", string(reqJSON))
	log.Printf("CreateManualSchemaForClient - CoachID: %s, UserID: %d", coachID, req.UserID)

	if err := s.validator.Struct(req); err != nil {
		log.Printf("CreateManualSchemaForClient - Validation error: %v", err)
		return nil, fmt.Errorf("validation error: %w", err)
	}

	if err := s.ValidateCoachPermission(ctx, coachID, req.UserID); err != nil {
		log.Printf("CreateManualSchemaForClient - Permission validation failed: %v", err)
		return nil, err
	}

	authUserID, err := s.repo.WorkoutProfiles().GetWorkoutProfileByID(ctx, req.UserID)
	if err != nil {
		log.Printf("CreateManualSchemaForClient - Failed to lookup auth_user_id for workout_profile_id %d: %v", req.UserID, err)
		return nil, fmt.Errorf("failed to find user profile: %w", err)
	}

	log.Printf("CreateManualSchemaForClient - Creating schema for workout_profile_id=%d, auth_user_id=%s", req.UserID, authUserID.AuthUserID)

	schemaReq := &types.WeeklySchemaRequest{
		UserID:    authUserID.AuthUserID,
		WeekStart: req.StartDate,
	}

	schema, err := s.repo.Schemas().CreateWeeklySchema(ctx, schemaReq)
	if err != nil {
		log.Printf("CreateManualSchemaForClient - Failed to create schema: %v", err)
		log.Printf("CreateManualSchemaForClient - SchemaRequest was: UserID=%s, WeekStart=%v",
			schemaReq.UserID, schemaReq.WeekStart)
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

	// Verify schema belongs to user by checking auth_user_id
	profile, err := s.repo.WorkoutProfiles().GetWorkoutProfileByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	if schema.UserID != profile.AuthUserID {
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

	profile, err := s.repo.WorkoutProfiles().GetWorkoutProfileByAuthID(ctx, schema.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user profile: %w", err)
	}

	isCoach, err := s.repo.CoachAssignments().IsCoachForUser(ctx, coachID, profile.WorkoutProfileID)
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

	sourceProfile, err := s.repo.WorkoutProfiles().GetWorkoutProfileByAuthID(ctx, sourceSchema.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source user profile: %w", err)
	}

	isCoach, err := s.repo.CoachAssignments().IsCoachForUser(ctx, coachID, sourceProfile.WorkoutProfileID)
	if err != nil {
		return nil, fmt.Errorf("failed to check coach permission: %w", err)
	}
	if !isCoach {
		return nil, fmt.Errorf("coach %s is not authorized for schema %d", coachID, sourceSchemaID)
	}

	targetProfile, err := s.repo.WorkoutProfiles().GetWorkoutProfileByID(ctx, targetUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get target user profile: %w", err)
	}

	clonedSchemaReq := &types.WeeklySchemaRequest{
		UserID:    targetProfile.AuthUserID,
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

	// Get workout_profile_id from auth_user_id to check coach permission
	profile, err := s.repo.WorkoutProfiles().GetWorkoutProfileByAuthID(ctx, schema.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user profile: %w", err)
	}

	isCoach, err := s.repo.CoachAssignments().IsCoachForUser(ctx, coachID, profile.WorkoutProfileID)
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

func (s *coachService) SearchUsers(ctx context.Context, query string, coachID string, limit int) ([]types.UserSearchResult, error) {
	if query == "" {
		return []types.UserSearchResult{}, nil
	}

	if limit <= 0 || limit > 50 {
		limit = 20
	}

	return s.repo.WorkoutProfiles().SearchUsers(ctx, query, coachID, limit)
}

func (s *coachService) AssignClientToCoach(ctx context.Context, req *types.CoachAssignmentRequest) (*types.CoachAssignment, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	profile, err := s.repo.WorkoutProfiles().GetWorkoutProfileByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			authUserID, err := s.repo.UserRoles().GetUserIDByUsername(ctx, req.Username)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil, fmt.Errorf("user not found")
				}
				return nil, fmt.Errorf("failed to get user: %w", err)
			}

			defaultProfile := &types.WorkoutProfileRequest{
				Level:     types.LevelBeginner,
				Goal:      types.GoalGeneralFitness,
				Frequency: 3,
				Equipment: []string{"bodyweight"},
			}

			profile, err = s.repo.WorkoutProfiles().CreateWorkoutProfile(ctx, authUserID, defaultProfile)
			if err != nil {
				return nil, fmt.Errorf("failed to create workout profile: %w", err)
			}
			log.Printf("[AssignClient] Created default workout profile %d for user %s", profile.WorkoutProfileID, req.Username)
		} else {
			return nil, fmt.Errorf("failed to verify client profile: %w", err)
		}
	}

	profileID := profile.WorkoutProfileID
	req.WorkoutProfileID = profileID

	exists, err := s.repo.CoachAssignments().IsCoachForUser(ctx, req.CoachID, profileID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing assignment: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("client is already assigned to you")
	}

	if currentAssignment, err := s.repo.CoachAssignments().GetCoachByUserID(ctx, profileID); err == nil {
		if currentAssignment.IsActive && currentAssignment.CoachID != req.CoachID {
			return nil, fmt.Errorf("client is already assigned to another coach")
		}
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to verify existing assignment: %w", err)
	}

	assignment, err := s.repo.CoachAssignments().CreateCoachAssignment(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to assign client to coach: %w", err)
	}

	_ = s.repo.CoachAssignments().LogCoachActivity(ctx, &types.CoachActivity{
		CoachID:      req.CoachID,
		UserID:       profileID,
		ActivityType: "client_assigned",
		Description:  fmt.Sprintf("Assigned client %s", req.Username),
		Timestamp:    time.Now(),
	})

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

	// Get user's auth_user_id for schema creation
	profile, err := s.repo.WorkoutProfiles().GetWorkoutProfileByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	schemaReq := &types.WeeklySchemaRequest{
		UserID:    profile.AuthUserID,
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
		UserID:        userID,
		TotalWorkouts: progress.TotalCount,
		LastWorkout:   nil,
	}, nil
}
