package service

import (
	"context"
	"testing"
	"time"

	gofpdf "github.com/jung-kurt/gofpdf"
	"github.com/tdmdh/fit-up-server/services/schema-service/internal/repository"
	"github.com/tdmdh/fit-up-server/services/schema-service/internal/types"
)

type MockPlanGenerationRepo struct {
	plans map[int]*types.GeneratedPlan
}

func (m *MockPlanGenerationRepo) GetActivePlanForUser(ctx context.Context, userID int) (*types.GeneratedPlan, error) {
	plan, exists := m.plans[userID]
	if !exists {
		return nil, nil // No active plan
	}
	return plan, nil
}

func (m *MockPlanGenerationRepo) CreatePlanGeneration(ctx context.Context, userID int, metadata *types.PlanGenerationMetadata) (*types.GeneratedPlan, error) {
	plan := &types.GeneratedPlan{
		PlanID:      len(m.plans) + 1,
		UserID:      userID,
		WeekStart:   time.Now(),
		GeneratedAt: time.Now(),
		Algorithm:   metadata.Algorithm,
		IsActive:    true,
		Metadata:    nil, 
	}
	m.plans[userID] = plan
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

// MockSchemaRepo for testing
type MockSchemaRepo struct {
	planGenRepo *MockPlanGenerationRepo
}

func (m *MockSchemaRepo) PlanGeneration() repository.PlanGenerationRepo {
	return m.planGenRepo
}

func (m *MockSchemaRepo) WorkoutProfiles() repository.WorkoutProfileRepo            { return nil }
func (m *MockSchemaRepo) Exercises() repository.ExerciseRepo                        { return nil }
func (m *MockSchemaRepo) Templates() repository.WorkoutTemplateRepo                 { return nil }
func (m *MockSchemaRepo) Schemas() repository.WeeklySchemaRepo                      { return nil }
func (m *MockSchemaRepo) Workouts() repository.WorkoutRepo                          { return nil }
func (m *MockSchemaRepo) WorkoutExercises() repository.WorkoutExerciseRepo          { return nil }
func (m *MockSchemaRepo) Progress() repository.ProgressRepo                         { return nil }
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
		plans: make(map[int]*types.GeneratedPlan),
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
		plans: make(map[int]*types.GeneratedPlan),
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
		plans: make(map[int]*types.GeneratedPlan),
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
		plans: make(map[int]*types.GeneratedPlan),
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

	plans, err := service.CreatePlanGeneration(ctx, 123, metaDate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if plans == nil {
		t.Fatal("Expected plan to be created, got nil")
	}

	pdfBytes, err := service.ExportPlanToPDF(ctx, plans.PlanID)
	if err != nil {
		t.Fatalf("Expected no error creating PDF, got %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Fatal("Expected PDF bytes, got empty")
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	if pdf == nil {
		t.Fatal("Failed to create PDF instance")
	}
}
