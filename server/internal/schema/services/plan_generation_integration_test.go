package service

// import (
// 	"context"
// 	"encoding/json"
// 	"os"
// 	"testing"

// 	"github.com/tdmdh/fit-up-server/internal/schema/types"
// )

// func TestCompleteWorkflow(t *testing.T) {
// 	mockPlanRepo := &MockPlanGenerationRepo{
// 		plans:     make(map[int]*types.GeneratedPlan),
// 		plansByID: make(map[int]*types.GeneratedPlan),
// 	}
// 	mockRepo := &MockSchemaRepo{
// 		planGenRepo: mockPlanRepo,
// 	}

// 	service := NewPlanGenerationService(mockRepo)
// 	ctx := context.Background()

// 	t.Run("Intermediate_MuscleGain_Plan", func(t *testing.T) {
// 		userID := 100
// 		metadata := &types.PlanGenerationMetadata{
// 			UserGoals: []types.FitnessGoal{
// 				types.GoalMuscleGain,
// 			},
// 			AvailableEquipment: []types.EquipmentType{
// 				types.EquipmentDumbbell,
// 				types.EquipmentBodyweight,
// 			},
// 			FitnessLevel:    types.LevelIntermediate,
// 			WeeklyFrequency: 4,
// 			TimePerWorkout:  60,
// 			Algorithm:       "",
// 			Parameters:      make(map[string]interface{}),
// 		}

// 		plan, err := service.CreatePlanGeneration(ctx, userID, metadata)
// 		if err != nil {
// 			t.Fatalf("Failed to create plan: %v", err)
// 		}

// 		t.Logf("✓ Plan generated successfully - Plan ID: %d", plan.PlanID)

// 		if plan.UserID != userID {
// 			t.Errorf("Expected UserID %d, got %d", userID, plan.UserID)
// 		}
// 		if !plan.IsActive {
// 			t.Error("Plan should be active")
// 		}
// 		if plan.Algorithm != "fitup_adaptive_v1" {
// 			t.Errorf("Expected algorithm 'fitup_adaptive_v1', got '%s'", plan.Algorithm)
// 		}

// 		var metadataMap map[string]interface{}
// 		err = json.Unmarshal(plan.Metadata, &metadataMap)
// 		if err != nil {
// 			t.Fatalf("Failed to unmarshal plan metadata: %v", err)
// 		}

// 		parameters, ok := metadataMap["parameters"].(map[string]interface{})
// 		if !ok {
// 			t.Fatal("Metadata should contain parameters")
// 		}

// 		t.Logf("✓ Plan metadata structure is valid")

// 		requiredParams := []string{
// 			"template_used",
// 			"total_exercises",
// 			"muscle_groups_targeted",
// 			"equipment_utilized",
// 			"estimated_volume",
// 			"progression_method",
// 			"generated_plan",
// 		}

// 		for _, param := range requiredParams {
// 			if _, exists := parameters[param]; !exists {
// 				t.Errorf("Missing required parameter: %s", param)
// 			}
// 		}

// 		t.Logf("✓ All required parameters present")

// 		generatedPlan, ok := parameters["generated_plan"].([]interface{})
// 		if !ok {
// 			t.Fatal("Generated plan should be an array")
// 		}

// 		if len(generatedPlan) == 0 {
// 			t.Fatal("Generated plan should have at least one day")
// 		}

// 		t.Logf("✓ Generated plan has %d days", len(generatedPlan))

// 		workoutDays := 0
// 		restDays := 0
// 		totalExercises := 0

// 		for dayIdx, dayIface := range generatedPlan {
// 			dayMap, ok := dayIface.(map[string]interface{})
// 			if !ok {
// 				t.Errorf("Day %d is not a valid map", dayIdx)
// 				continue
// 			}

// 			dayTitle, _ := dayMap["day_title"].(string)
// 			focus, _ := dayMap["focus"].(string)
// 			exercises, _ := dayMap["exercises"].([]interface{})

// 			if dayTitle == "" {
// 				t.Errorf("Day %d missing title", dayIdx)
// 			}
// 			if focus == "" {
// 				t.Errorf("Day %d missing focus", dayIdx)
// 			}

// 			if len(exercises) > 0 {
// 				workoutDays++
// 				totalExercises += len(exercises)

// 				for exIdx, exIface := range exercises {
// 					exMap, ok := exIface.(map[string]interface{})
// 					if !ok {
// 						t.Errorf("Day %d, Exercise %d is not a valid map", dayIdx, exIdx)
// 						continue
// 					}

// 					if _, ok := exMap["name"].(string); !ok {
// 						t.Errorf("Day %d, Exercise %d missing name", dayIdx, exIdx)
// 					}
// 					if _, ok := exMap["sets"].(float64); !ok {
// 						t.Errorf("Day %d, Exercise %d missing sets", dayIdx, exIdx)
// 					}
// 					if _, ok := exMap["reps"].(string); !ok {
// 						t.Errorf("Day %d, Exercise %d missing reps", dayIdx, exIdx)
// 					}
// 					if _, ok := exMap["rest"].(float64); !ok {
// 						t.Errorf("Day %d, Exercise %d missing rest", dayIdx, exIdx)
// 					}
// 				}
// 			} else {
// 				restDays++
// 			}
// 		}

// 		t.Logf("✓ Plan structure valid: %d workout days, %d rest days, %d total exercises",
// 			workoutDays, restDays, totalExercises)

// 		pdfBytes, err := service.ExportPlanToPDF(ctx, plan.PlanID)
// 		if err != nil {
// 			t.Fatalf("Failed to generate PDF: %v", err)
// 		}

// 		if len(pdfBytes) == 0 {
// 			t.Fatal("PDF bytes are empty")
// 		}

// 		t.Logf("✓ PDF generated successfully - Size: %d bytes", len(pdfBytes))

// 		outputPath := "test_integration_workout_plan.pdf"
// 		err = os.WriteFile(outputPath, pdfBytes, 0644)
// 		if err != nil {
// 			t.Logf("Warning: Could not save PDF to file: %v", err)
// 		} else {
// 			t.Logf("✓ PDF saved to: %s", outputPath)
// 		}

// 		_, err = service.CreatePlanGeneration(ctx, userID, metadata)
// 		if err != types.ErrActivePlanExists {
// 			t.Errorf("Expected ErrActivePlanExists, got %v", err)
// 		}

// 		t.Logf("✓ Active plan validation works correctly")
// 	})

// 	t.Run("Beginner_Bodyweight_Plan", func(t *testing.T) {
// 		userID := 200
// 		metadata := &types.PlanGenerationMetadata{
// 			UserGoals: []types.FitnessGoal{
// 				types.GoalGeneralFitness,
// 			},
// 			AvailableEquipment: []types.EquipmentType{
// 				types.EquipmentBodyweight,
// 			},
// 			FitnessLevel:    types.LevelBeginner,
// 			WeeklyFrequency: 3,
// 			TimePerWorkout:  30,
// 			Algorithm:       "",
// 			Parameters:      make(map[string]interface{}),
// 		}

// 		plan, err := service.CreatePlanGeneration(ctx, userID, metadata)
// 		if err != nil {
// 			t.Fatalf("Failed to create beginner plan: %v", err)
// 		}

// 		t.Logf("✓ Beginner plan generated - Plan ID: %d", plan.PlanID)

// 		pdfBytes, err := service.ExportPlanToPDF(ctx, plan.PlanID)
// 		if err != nil {
// 			t.Fatalf("Failed to generate PDF: %v", err)
// 		}

// 		outputPath := "test_beginner_workout_plan.pdf"
// 		err = os.WriteFile(outputPath, pdfBytes, 0644)
// 		if err == nil {
// 			t.Logf("✓ Beginner PDF saved to: %s", outputPath)
// 		}
// 	})

// 	t.Run("Advanced_Strength_Plan", func(t *testing.T) {
// 		userID := 300
// 		metadata := &types.PlanGenerationMetadata{
// 			UserGoals: []types.FitnessGoal{
// 				types.GoalStrength,
// 			},
// 			AvailableEquipment: []types.EquipmentType{
// 				types.EquipmentBarbell,
// 				types.EquipmentDumbbell,
// 				types.EquipmentBodyweight,
// 			},
// 			FitnessLevel:    types.LevelAdvanced,
// 			WeeklyFrequency: 5,
// 			TimePerWorkout:  90,
// 			Algorithm:       "",
// 			Parameters:      make(map[string]interface{}),
// 		}

// 		plan, err := service.CreatePlanGeneration(ctx, userID, metadata)
// 		if err != nil {
// 			t.Fatalf("Failed to create advanced plan: %v", err)
// 		}

// 		t.Logf("✓ Advanced plan generated - Plan ID: %d", plan.PlanID)

// 		pdfBytes, err := service.ExportPlanToPDF(ctx, plan.PlanID)
// 		if err != nil {
// 			t.Fatalf("Failed to generate PDF: %v", err)
// 		}

// 		outputPath := "test_advanced_workout_plan.pdf"
// 		err = os.WriteFile(outputPath, pdfBytes, 0644)
// 		if err == nil {
// 			t.Logf("✓ Advanced PDF saved to: %s", outputPath)
// 		}
// 	})

// 	t.Log("\n========================================")
// 	t.Log("✓ ALL INTEGRATION TESTS PASSED!")
// 	t.Log("========================================")
// 	t.Log("Plans generated successfully for all user levels")
// 	t.Log("PDFs created with proper structure and formatting")
// 	t.Log("Data flow verified from input to output")
// 	t.Log("Check the generated PDF files for visual confirmation")
// }

// func TestPlanDataIntegrity(t *testing.T) {
// 	mockPlanRepo := &MockPlanGenerationRepo{
// 		plans:     make(map[int]*types.GeneratedPlan),
// 		plansByID: make(map[int]*types.GeneratedPlan),
// 	}
// 	mockRepo := &MockSchemaRepo{
// 		planGenRepo: mockPlanRepo,
// 	}

// 	service := NewPlanGenerationService(mockRepo)
// 	ctx := context.Background()

// 	userID := 400
// 	metadata := &types.PlanGenerationMetadata{
// 		UserGoals:          []types.FitnessGoal{types.GoalMuscleGain},
// 		AvailableEquipment: []types.EquipmentType{types.EquipmentDumbbell},
// 		FitnessLevel:       types.LevelIntermediate,
// 		WeeklyFrequency:    4,
// 		TimePerWorkout:     60,
// 		Algorithm:          "",
// 		Parameters:         make(map[string]interface{}),
// 	}

// 	plan, err := service.CreatePlanGeneration(ctx, userID, metadata)
// 	if err != nil {
// 		t.Fatalf("Failed to create plan: %v", err)
// 	}

// 	retrievedPlan, err := service.GetActivePlanForUser(ctx, userID)
// 	if err != nil {
// 		t.Fatalf("Failed to retrieve plan: %v", err)
// 	}

// 	if retrievedPlan.PlanID != plan.PlanID {
// 		t.Errorf("Plan ID mismatch: expected %d, got %d", plan.PlanID, retrievedPlan.PlanID)
// 	}

// 	if retrievedPlan.UserID != plan.UserID {
// 		t.Errorf("User ID mismatch: expected %d, got %d", plan.UserID, retrievedPlan.UserID)
// 	}

// 	var originalMeta, retrievedMeta map[string]interface{}
// 	json.Unmarshal(plan.Metadata, &originalMeta)
// 	json.Unmarshal(retrievedPlan.Metadata, &retrievedMeta)

// 	originalParams := originalMeta["parameters"].(map[string]interface{})
// 	retrievedParams := retrievedMeta["parameters"].(map[string]interface{})

// 	if originalParams["template_used"] != retrievedParams["template_used"] {
// 		t.Error("Template information was corrupted during storage/retrieval")
// 	}

// 	t.Log("✓ Plan data integrity verified - all data preserved correctly")
// }
