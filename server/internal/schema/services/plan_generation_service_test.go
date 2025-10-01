package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	gofpdf "github.com/jung-kurt/gofpdf"
	"github.com/tdmdh/fit-up-server/internal/schema/repository"
	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

func countExercisesInPlan(plan []interface{}) int {
	count := 0
	for _, dayIface := range plan {
		if dayMap, ok := dayIface.(map[string]interface{}); ok {
			if exercises, ok := dayMap["exercises"].([]interface{}); ok {
				count += len(exercises)
			}
		}
	}
	return count
}

func extractMuscleGroups(level types.FitnessLevel) []string {
	switch level {
	case types.LevelBeginner:
		return []string{"full_body", "core", "legs"}
	case types.LevelIntermediate:
		return []string{"chest", "back", "shoulders", "legs", "arms"}
	case types.LevelAdvanced:
		return []string{"chest", "back", "shoulders", "legs", "arms", "core", "glutes"}
	default:
		return []string{"general"}
	}
}

func convertEquipmentTypes(equipment []types.EquipmentType) []string {
	result := make([]string, len(equipment))
	for i, eq := range equipment {
		result[i] = string(eq)
	}
	return result
}

type MockPlanGenerationRepo struct {
	plans     map[int]*types.GeneratedPlan
	plansByID map[int]*types.GeneratedPlan
}

// GetPlanID implements repository.PlanGenerationRepo.
func (m *MockPlanGenerationRepo) GetPlanID(ctx context.Context, planID int) (*types.GeneratedPlan, error) {
	if plan, exists := m.plansByID[planID]; exists {
		return plan, nil
	}
	return nil, nil
}

func (m *MockPlanGenerationRepo) GetActivePlanForUser(ctx context.Context, userID int) (*types.GeneratedPlan, error) {
	plan, exists := m.plans[userID]
	if !exists {
		return nil, nil
	}
	return plan, nil
}

func (m *MockPlanGenerationRepo) CreatePlanGeneration(ctx context.Context, userID int, metadata *types.PlanGenerationMetadata) (*types.GeneratedPlan, error) {
	planID := len(m.plansByID) + 1

	// Generate different plans based on fitness level and goals
	var mockGeneratedPlan []interface{}

	switch metadata.FitnessLevel {
	case types.LevelBeginner:
		mockGeneratedPlan = []interface{}{
			map[string]interface{}{
				"day_title": "Full Body Basics",
				"focus":     "foundational_strength",
				"exercises": []interface{}{
					map[string]interface{}{
						"exercise_id": 1,
						"name":        "Bodyweight Squats",
						"sets":        2.0,
						"reps":        "12-15",
						"rest":        60.0,
						"notes":       "Focus on proper form, keep knees aligned with toes",
					},
					map[string]interface{}{
						"exercise_id": 2,
						"name":        "Wall Push-ups",
						"sets":        2.0,
						"reps":        "10-12",
						"rest":        60.0,
						"notes":       "Keep core tight, controlled movement",
					},
					map[string]interface{}{
						"exercise_id": 3,
						"name":        "Assisted Lunges",
						"sets":        2.0,
						"reps":        "8-10 per leg",
						"rest":        60.0,
						"notes":       "Use chair for balance if needed",
					},
					map[string]interface{}{
						"exercise_id": 4,
						"name":        "Plank Hold",
						"sets":        2.0,
						"reps":        "20-30 sec",
						"rest":        60.0,
						"notes":       "Maintain straight line from head to heels",
					},
				},
			},
			map[string]interface{}{
				"day_title": "Rest & Recovery",
				"focus":     "recovery",
				"exercises": []interface{}{},
			},
			map[string]interface{}{
				"day_title": "Light Cardio & Mobility",
				"focus":     "cardiovascular_health",
				"exercises": []interface{}{
					map[string]interface{}{
						"exercise_id": 5,
						"name":        "Walking",
						"sets":        1.0,
						"reps":        "20 minutes",
						"rest":        0.0,
						"notes":       "Maintain steady pace, focus on breathing",
					},
					map[string]interface{}{
						"exercise_id": 6,
						"name":        "Basic Stretching",
						"sets":        1.0,
						"reps":        "10 minutes",
						"rest":        0.0,
						"notes":       "Hold each stretch for 20-30 seconds",
					},
				},
			},
		}
	case types.LevelIntermediate:
		if len(metadata.UserGoals) > 0 && metadata.UserGoals[0] == types.GoalMuscleGain {
			mockGeneratedPlan = []interface{}{
				map[string]interface{}{
					"day_title": "Upper Body Push",
					"focus":     "chest_shoulders_triceps",
					"exercises": []interface{}{
						map[string]interface{}{
							"exercise_id": 1,
							"name":        "Dumbbell Bench Press",
							"sets":        4.0,
							"reps":        "8-12",
							"rest":        90.0,
							"notes":       "Control the weight, full range of motion",
						},
						map[string]interface{}{
							"exercise_id": 2,
							"name":        "Dumbbell Shoulder Press",
							"sets":        3.0,
							"reps":        "10-12",
							"rest":        75.0,
							"notes":       "Keep core engaged, press straight up",
						},
						map[string]interface{}{
							"exercise_id": 3,
							"name":        "Dips",
							"sets":        3.0,
							"reps":        "8-12",
							"rest":        75.0,
							"notes":       "Lean forward for chest emphasis",
						},
						map[string]interface{}{
							"exercise_id": 4,
							"name":        "Lateral Raises",
							"sets":        3.0,
							"reps":        "12-15",
							"rest":        60.0,
							"notes":       "Control the movement, slight bend in elbows",
						},
					},
				},
				map[string]interface{}{
					"day_title": "Lower Body",
					"focus":     "legs_glutes",
					"exercises": []interface{}{
						map[string]interface{}{
							"exercise_id": 5,
							"name":        "Dumbbell Goblet Squats",
							"sets":        4.0,
							"reps":        "10-15",
							"rest":        90.0,
							"notes":       "Keep chest up, squat to parallel or below",
						},
						map[string]interface{}{
							"exercise_id": 6,
							"name":        "Romanian Deadlifts",
							"sets":        3.0,
							"reps":        "10-12",
							"rest":        90.0,
							"notes":       "Feel stretch in hamstrings, keep back straight",
						},
						map[string]interface{}{
							"exercise_id": 7,
							"name":        "Bulgarian Split Squats",
							"sets":        3.0,
							"reps":        "10-12 per leg",
							"rest":        75.0,
							"notes":       "Front leg does most of the work",
						},
					},
				},
				map[string]interface{}{
					"day_title": "Rest Day",
					"focus":     "recovery",
					"exercises": []interface{}{},
				},
				map[string]interface{}{
					"day_title": "Upper Body Pull",
					"focus":     "back_biceps",
					"exercises": []interface{}{
						map[string]interface{}{
							"exercise_id": 8,
							"name":        "Pull-ups",
							"sets":        4.0,
							"reps":        "6-10",
							"rest":        90.0,
							"notes":       "Full extension at bottom, chin over bar at top",
						},
						map[string]interface{}{
							"exercise_id": 9,
							"name":        "Dumbbell Rows",
							"sets":        4.0,
							"reps":        "10-12",
							"rest":        75.0,
							"notes":       "Pull elbow back, squeeze shoulder blade",
						},
						map[string]interface{}{
							"exercise_id": 10,
							"name":        "Bicep Curls",
							"sets":        3.0,
							"reps":        "12-15",
							"rest":        60.0,
							"notes":       "Keep elbows stationary, control the weight",
						},
					},
				},
			}
		} else {
			mockGeneratedPlan = []interface{}{
				map[string]interface{}{
					"day_title": "Full Body Workout A",
					"focus":     "compound_movements",
					"exercises": []interface{}{
						map[string]interface{}{
							"exercise_id": 1,
							"name":        "Push-ups",
							"sets":        3.0,
							"reps":        "12-15",
							"rest":        60.0,
						},
						map[string]interface{}{
							"exercise_id": 2,
							"name":        "Dumbbell Squats",
							"sets":        3.0,
							"reps":        "12-15",
							"rest":        75.0,
						},
						map[string]interface{}{
							"exercise_id": 3,
							"name":        "Rows",
							"sets":        3.0,
							"reps":        "10-12",
							"rest":        75.0,
						},
					},
				},
				map[string]interface{}{
					"day_title": "Rest Day",
					"focus":     "recovery",
					"exercises": []interface{}{},
				},
			}
		}
	case types.LevelAdvanced:
		mockGeneratedPlan = []interface{}{
			map[string]interface{}{
				"day_title": "Heavy Lower Body - Strength Focus",
				"focus":     "maximum_strength",
				"exercises": []interface{}{
					map[string]interface{}{
						"exercise_id": 1,
						"name":        "Barbell Back Squats",
						"sets":        5.0,
						"reps":        "3-5",
						"rest":        180.0,
						"notes":       "Heavy weight, full depth, maximum effort",
					},
					map[string]interface{}{
						"exercise_id": 2,
						"name":        "Barbell Deadlifts",
						"sets":        5.0,
						"reps":        "3-5",
						"rest":        180.0,
						"notes":       "Maintain neutral spine, explosive pull",
					},
					map[string]interface{}{
						"exercise_id": 3,
						"name":        "Front Squats",
						"sets":        4.0,
						"reps":        "6-8",
						"rest":        120.0,
						"notes":       "Keep chest up, elbows high",
					},
					map[string]interface{}{
						"exercise_id": 4,
						"name":        "Walking Lunges",
						"sets":        3.0,
						"reps":        "12 per leg",
						"rest":        90.0,
						"notes":       "Add weight for increased difficulty",
					},
				},
			},
			map[string]interface{}{
				"day_title": "Upper Body Power",
				"focus":     "explosive_strength",
				"exercises": []interface{}{
					map[string]interface{}{
						"exercise_id": 5,
						"name":        "Barbell Bench Press",
						"sets":        5.0,
						"reps":        "3-5",
						"rest":        180.0,
						"notes":       "Heavy weight, controlled descent, explosive push",
					},
					map[string]interface{}{
						"exercise_id": 6,
						"name":        "Weighted Pull-ups",
						"sets":        4.0,
						"reps":        "5-8",
						"rest":        150.0,
						"notes":       "Add weight with belt or vest",
					},
					map[string]interface{}{
						"exercise_id": 7,
						"name":        "Barbell Overhead Press",
						"sets":        4.0,
						"reps":        "5-8",
						"rest":        120.0,
						"notes":       "Strict form, no leg drive",
					},
					map[string]interface{}{
						"exercise_id": 8,
						"name":        "Barbell Rows",
						"sets":        4.0,
						"reps":        "8-10",
						"rest":        90.0,
						"notes":       "Pull to lower chest, squeeze back",
					},
				},
			},
			map[string]interface{}{
				"day_title": "Active Recovery",
				"focus":     "recovery",
				"exercises": []interface{}{},
			},
			map[string]interface{}{
				"day_title": "Hypertrophy Upper",
				"focus":     "muscle_building",
				"exercises": []interface{}{
					map[string]interface{}{
						"exercise_id": 9,
						"name":        "Incline Dumbbell Press",
						"sets":        4.0,
						"reps":        "10-12",
						"rest":        75.0,
						"notes":       "Focus on muscle contraction",
					},
					map[string]interface{}{
						"exercise_id": 10,
						"name":        "Dumbbell Flyes",
						"sets":        3.0,
						"reps":        "12-15",
						"rest":        60.0,
						"notes":       "Stretch at bottom, squeeze at top",
					},
					map[string]interface{}{
						"exercise_id": 11,
						"name":        "Cable Rows",
						"sets":        4.0,
						"reps":        "12-15",
						"rest":        60.0,
						"notes":       "Slow and controlled",
					},
				},
			},
			map[string]interface{}{
				"day_title": "Hypertrophy Lower",
				"focus":     "leg_development",
				"exercises": []interface{}{
					map[string]interface{}{
						"exercise_id": 12,
						"name":        "Leg Press",
						"sets":        4.0,
						"reps":        "12-15",
						"rest":        90.0,
						"notes":       "Deep range of motion",
					},
					map[string]interface{}{
						"exercise_id": 13,
						"name":        "Leg Curls",
						"sets":        4.0,
						"reps":        "12-15",
						"rest":        60.0,
						"notes":       "Squeeze hamstrings at peak",
					},
					map[string]interface{}{
						"exercise_id": 14,
						"name":        "Calf Raises",
						"sets":        4.0,
						"reps":        "15-20",
						"rest":        60.0,
						"notes":       "Full stretch and contraction",
					},
				},
			},
		}
	default:
		mockGeneratedPlan = []interface{}{
			map[string]interface{}{
				"day_title": "General Fitness",
				"focus":     "overall_fitness",
				"exercises": []interface{}{
					map[string]interface{}{
						"exercise_id": 1,
						"name":        "Push-ups",
						"sets":        3.0,
						"reps":        "10-12",
						"rest":        60.0,
					},
				},
			},
		}
	}

	if metadata.Parameters == nil {
		metadata.Parameters = make(map[string]interface{})
	}
	metadata.Parameters["generated_plan"] = mockGeneratedPlan
	metadata.Parameters["template_used"] = fmt.Sprintf("%s_%s_%dday", metadata.FitnessLevel, metadata.UserGoals[0], metadata.WeeklyFrequency)
	metadata.Parameters["total_exercises"] = float64(countExercisesInPlan(mockGeneratedPlan))
	metadata.Parameters["muscle_groups_targeted"] = extractMuscleGroups(metadata.FitnessLevel)
	metadata.Parameters["equipment_utilized"] = convertEquipmentTypes(metadata.AvailableEquipment)

	// Marshal metadata to JSON
	metadataJSON, _ := json.Marshal(map[string]interface{}{
		"parameters": metadata.Parameters,
	})

	plan := &types.GeneratedPlan{
		PlanID:      planID,
		UserID:      userID,
		WeekStart:   time.Now(),
		GeneratedAt: time.Now(),
		Algorithm:   metadata.Algorithm,
		IsActive:    true,
		Metadata:    metadataJSON,
	}
	m.plans[userID] = plan
	m.plansByID[planID] = plan
	return plan, nil
}

func (m *MockPlanGenerationRepo) GetPlanGenerationHistory(ctx context.Context, userID int, limit int) ([]types.GeneratedPlan, error) {
	return nil, nil
}

func (m *MockPlanGenerationRepo) TrackPlanPerformance(ctx context.Context, planID int, performance *types.PlanPerformanceData) error {
	return nil
}

func (m *MockPlanGenerationRepo) GetPlanEffectivenessScore(ctx context.Context, planID int) (float64, error) {
	return 0.85, nil
}

func (m *MockPlanGenerationRepo) MarkPlanForRegeneration(ctx context.Context, planID int, reason string) error {
	return nil
}

func (m *MockPlanGenerationRepo) LogPlanAdaptation(ctx context.Context, planID int, adaptation *types.PlanAdaptation) error {
	return nil
}

func (m *MockPlanGenerationRepo) GetAdaptationHistory(ctx context.Context, userID int) ([]types.PlanAdaptation, error) {
	return nil, nil
}

// MockProgressRepo for testing
type MockProgressRepo struct{}

func (m *MockProgressRepo) GetProgressLogsByUserID(ctx context.Context, userID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.ProgressLog], error) {
	return &types.PaginatedResponse[types.ProgressLog]{
		Data:       []types.ProgressLog{},
		TotalCount: 0,
		Page:       1,
		PageSize:   pagination.Limit,
	}, nil
}

func (m *MockProgressRepo) CreateProgressLog(ctx context.Context, log *types.ProgressLogRequest) (*types.ProgressLog, error) {
	return &types.ProgressLog{}, nil
}

func (m *MockProgressRepo) GetProgressLogByID(ctx context.Context, logID int) (*types.ProgressLog, error) {
	return nil, nil
}

func (m *MockProgressRepo) UpdateProgressLog(ctx context.Context, logID int, log *types.ProgressLogRequest) (*types.ProgressLog, error) {
	return &types.ProgressLog{}, nil
}

func (m *MockProgressRepo) DeleteProgressLog(ctx context.Context, logID int) error {
	return nil
}

func (m *MockProgressRepo) GetProgressByExercise(ctx context.Context, userID int, exerciseID int) ([]types.ProgressLog, error) {
	return []types.ProgressLog{}, nil
}

func (m *MockProgressRepo) BulkCreateProgressLogs(ctx context.Context, logs []types.ProgressLogRequest) ([]types.ProgressLog, error) {
	return []types.ProgressLog{}, nil
}

func (m *MockProgressRepo) FilterProgressLogs(ctx context.Context, filter types.ProgressFilter, pagination types.PaginationParams) (*types.PaginatedResponse[types.ProgressLog], error) {
	return &types.PaginatedResponse[types.ProgressLog]{
		Data:       []types.ProgressLog{},
		TotalCount: 0,
		Page:       1,
		PageSize:   pagination.Limit,
	}, nil
}

func (m *MockProgressRepo) GetLatestProgressLogsForUser(ctx context.Context, userID int) ([]types.ProgressLog, error) {
	return []types.ProgressLog{}, nil
}

func (m *MockProgressRepo) GetPersonalBests(ctx context.Context, userID int) ([]types.PersonalBest, error) {
	return []types.PersonalBest{}, nil
}

func (m *MockProgressRepo) GetProgressLogsByUserAndDateRange(ctx context.Context, userID int, startDate, endDate time.Time) ([]types.ProgressLog, error) {
	return []types.ProgressLog{}, nil
}

func (m *MockProgressRepo) GetProgressLogsByUserAndExercise(ctx context.Context, userID int, exerciseID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.ProgressLog], error) {
	return &types.PaginatedResponse[types.ProgressLog]{
		Data:       []types.ProgressLog{},
		TotalCount: 0,
		Page:       1,
		PageSize:   pagination.Limit,
	}, nil
}

func (m *MockProgressRepo) GetProgressTrend(ctx context.Context, userID int, exerciseID int, days int) ([]types.ProgressLog, error) {
	return []types.ProgressLog{}, nil
}

func (m *MockProgressRepo) GetUserProgressSummary(ctx context.Context, userID int) (*types.UserProgressSummary, error) {
	return &types.UserProgressSummary{}, nil
}

func (m *MockProgressRepo) GetWorkoutStreak(ctx context.Context, userID int) (int, error) {
	return 0, nil
}

// MockSchemaRepo for testing
type MockSchemaRepo struct {
	planGenRepo  *MockPlanGenerationRepo
	progressRepo *MockProgressRepo
}

func (m *MockSchemaRepo) PlanGeneration() repository.PlanGenerationRepo {
	return m.planGenRepo
}
func (m *MockSchemaRepo) WorkoutProfiles() repository.WorkoutProfileRepo { return nil }
func (m *MockSchemaRepo) Exercises() repository.ExerciseRepo             { return nil }
func (m *MockSchemaRepo) Templates() repository.WorkoutTemplateRepo      { return nil }
func (m *MockSchemaRepo) Schemas() repository.WeeklySchemaRepo           { return nil }
func (m *MockSchemaRepo) Workouts() repository.WorkoutRepo               { return nil }
func (m *MockSchemaRepo) WorkoutExercises() repository.WorkoutExerciseRepo {
	return nil
}
func (m *MockSchemaRepo) Progress() repository.ProgressRepo {
	if m.progressRepo == nil {
		return &MockProgressRepo{}
	}
	return m.progressRepo
}
func (m *MockSchemaRepo) FitnessProfiles() repository.FitnessProfileRepo            { return nil }
func (m *MockSchemaRepo) WorkoutSessions() repository.WorkoutSessionRepo            { return nil }
func (m *MockSchemaRepo) RecoveryMetrics() repository.RecoveryMetricsRepo           { return nil }
func (m *MockSchemaRepo) PerformanceAnalytics() repository.PerformanceAnalyticsRepo { return nil }
func (m *MockSchemaRepo) GoalTracking() repository.GoalTrackingRepo                 { return nil }
func (m *MockSchemaRepo) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	return nil
}

func TestCreatePlanGeneration(t *testing.T) {
	// Setup mock repository
	mockPlanRepo := &MockPlanGenerationRepo{
		plans:     make(map[int]*types.GeneratedPlan),
		plansByID: make(map[int]*types.GeneratedPlan),
	}
	mockRepo := &MockSchemaRepo{
		planGenRepo: mockPlanRepo,
	}

	// Create service
	service := NewPlanGenerationService(mockRepo)

	// Test data
	userID := 123
	metadata := &types.PlanGenerationMetadata{
		UserGoals: []types.FitnessGoal{
			types.GoalMuscleGain,
		},
		AvailableEquipment: []types.EquipmentType{
			types.EquipmentDumbbell,
			types.EquipmentBodyweight,
		},
		FitnessLevel:    types.LevelIntermediate,
		WeeklyFrequency: 4,
		TimePerWorkout:  60,
		Algorithm:       "",
		Parameters:      make(map[string]interface{}),
	}

	// Test successful plan creation
	ctx := context.Background()
	plan, err := service.CreatePlanGeneration(ctx, userID, metadata)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if plan == nil {
		t.Fatal("Expected plan to be created, got nil")
	}

	if plan.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, plan.UserID)
	}

	if plan.Algorithm != "fitup_adaptive_v1" {
		t.Errorf("Expected algorithm 'fitup_adaptive_v1', got '%s'", plan.Algorithm)
	}

	if !plan.IsActive {
		t.Error("Expected plan to be active")
	}

	// Test error for existing active plan
	_, err = service.CreatePlanGeneration(ctx, userID, metadata)
	if err != types.ErrActivePlanExists {
		t.Errorf("Expected ErrActivePlanExists, got %v", err)
	}
}

func TestCreatePlanGeneration_InvalidInput(t *testing.T) {
	mockPlanRepo := &MockPlanGenerationRepo{
		plans:     make(map[int]*types.GeneratedPlan),
		plansByID: make(map[int]*types.GeneratedPlan),
	}
	mockRepo := &MockSchemaRepo{
		planGenRepo: mockPlanRepo,
	}

	service := NewPlanGenerationService(mockRepo)
	ctx := context.Background()

	// Test invalid userID
	validMetadata := &types.PlanGenerationMetadata{
		UserGoals:          []types.FitnessGoal{types.GoalMuscleGain},
		AvailableEquipment: []types.EquipmentType{types.EquipmentBodyweight},
		FitnessLevel:       types.LevelBeginner,
		WeeklyFrequency:    3,
		TimePerWorkout:     45,
		Algorithm:          "",
		Parameters:         make(map[string]interface{}),
	}

	_, err := service.CreatePlanGeneration(ctx, 0, validMetadata)
	if err != types.ErrInvalidUserID {
		t.Errorf("Expected ErrInvalidUserID, got %v", err)
	}

	// Test invalid metadata (missing goals)
	invalidMetadata := &types.PlanGenerationMetadata{
		UserGoals:          []types.FitnessGoal{}, // Empty goals
		AvailableEquipment: []types.EquipmentType{types.EquipmentBodyweight},
		FitnessLevel:       types.LevelBeginner,
		WeeklyFrequency:    3,
		TimePerWorkout:     45,
		Algorithm:          "",
		Parameters:         make(map[string]interface{}),
	}
	_, err = service.CreatePlanGeneration(ctx, 123, invalidMetadata)
	if err == nil {
		t.Error("Expected error for empty goals")
	}

	// Test invalid metadata (missing equipment)
	invalidMetadata2 := &types.PlanGenerationMetadata{
		UserGoals:          []types.FitnessGoal{types.GoalMuscleGain},
		AvailableEquipment: []types.EquipmentType{}, // Empty equipment
		FitnessLevel:       types.LevelBeginner,
		WeeklyFrequency:    3,
		TimePerWorkout:     45,
		Algorithm:          "",
		Parameters:         make(map[string]interface{}),
	}
	_, err = service.CreatePlanGeneration(ctx, 123, invalidMetadata2)
	if err == nil {
		t.Error("Expected error for empty equipment")
	}
}

func TestGenerateAdaptivePlan_Integration(t *testing.T) {
	// This test verifies that the adaptive plan generation works with real fitness data
	mockPlanRepo := &MockPlanGenerationRepo{
		plans:     make(map[int]*types.GeneratedPlan),
		plansByID: make(map[int]*types.GeneratedPlan),
	}
	mockRepo := &MockSchemaRepo{
		planGenRepo: mockPlanRepo,
	}

	service := NewPlanGenerationService(mockRepo).(*planGenerationServiceImpl)
	ctx := context.Background()

	// Test with various user profiles
	testCases := []struct {
		name           string
		userID         int
		level          types.FitnessLevel
		goal           types.FitnessGoal
		equipment      []types.EquipmentType
		frequency      int
		timePerWorkout int
	}{
		{
			name:           "Beginner Bodyweight",
			userID:         1,
			level:          types.LevelBeginner,
			goal:           types.GoalGeneralFitness,
			equipment:      []types.EquipmentType{types.EquipmentBodyweight},
			frequency:      3,
			timePerWorkout: 30,
		},
		{
			name:           "Intermediate Muscle Gain",
			userID:         2,
			level:          types.LevelIntermediate,
			goal:           types.GoalMuscleGain,
			equipment:      []types.EquipmentType{types.EquipmentDumbbell, types.EquipmentBodyweight},
			frequency:      4,
			timePerWorkout: 60,
		},
		{
			name:           "Advanced Strength",
			userID:         3,
			level:          types.LevelAdvanced,
			goal:           types.GoalStrength,
			equipment:      []types.EquipmentType{types.EquipmentBarbell, types.EquipmentDumbbell, types.EquipmentBodyweight},
			frequency:      5,
			timePerWorkout: 90,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			metadata := &types.PlanGenerationMetadata{
				UserGoals:          []types.FitnessGoal{tc.goal},
				AvailableEquipment: tc.equipment,
				FitnessLevel:       tc.level,
				WeeklyFrequency:    tc.frequency,
				TimePerWorkout:     tc.timePerWorkout,
				Algorithm:          "",
				Parameters:         make(map[string]interface{}),
			}

			// Test the adaptive plan generation
			enhancedMetadata, err := service.generateAdaptivePlan(ctx, tc.userID, metadata)
			if err != nil {
				t.Fatalf("Failed to generate adaptive plan: %v", err)
			}

			if enhancedMetadata.Algorithm != "fitup_adaptive_v1" {
				t.Errorf("Expected algorithm 'fitup_adaptive_v1', got '%s'", enhancedMetadata.Algorithm)
			}

			// Verify parameters were populated
			if len(enhancedMetadata.Parameters) == 0 {
				t.Error("Expected parameters to be populated")
			}

			// Check for required parameters
			requiredParams := []string{
				"template_used",
				"total_exercises",
				"muscle_groups_targeted",
				"equipment_utilized",
				"estimated_volume",
				"progression_method",
			}

			for _, param := range requiredParams {
				if _, exists := enhancedMetadata.Parameters[param]; !exists {
					t.Errorf("Missing required parameter: %s", param)
				}
			}
		})
	}
}

func TestCreatePDFPlan(t *testing.T) {
	mockPlanRepo := &MockPlanGenerationRepo{
		plans:     make(map[int]*types.GeneratedPlan),
		plansByID: make(map[int]*types.GeneratedPlan),
	}

	mockRepo := &MockSchemaRepo{
		planGenRepo: mockPlanRepo,
	}

	service := NewPlanGenerationService(mockRepo).(*planGenerationServiceImpl)
	ctx := context.Background()

	metaDate := &types.PlanGenerationMetadata{
		UserGoals:          []types.FitnessGoal{types.GoalMuscleGain},
		AvailableEquipment: []types.EquipmentType{types.EquipmentBodyweight},
		FitnessLevel:       types.LevelBeginner,
		WeeklyFrequency:    3,
		TimePerWorkout:     45,
		Algorithm:          "",
		Parameters:         make(map[string]interface{}),
	}

	plan, err := service.CreatePlanGeneration(ctx, 123, metaDate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if plan == nil {
		t.Fatal("Expected plan to be created, got nil")
	}

	pdfBytes, err := service.ExportPlanToPDF(ctx, plan.PlanID)
	if err != nil {
		t.Fatalf("Expected no error creating PDF, got %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Fatal("Expected PDF bytes, got empty")
	}

	outputPath := "test_workout_plan.pdf"
	err = os.WriteFile(outputPath, pdfBytes, 0644)
	if err != nil {
		t.Logf("Warning: Could not save PDF to file: %v", err)
	} else {
		t.Logf("PDF saved successfully to: %s", outputPath)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	if pdf == nil {
		t.Fatal("Failed to create PDF instance")
	}
}
