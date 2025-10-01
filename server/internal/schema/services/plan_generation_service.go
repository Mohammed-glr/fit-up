package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jung-kurt/gofpdf"
	"github.com/tdmdh/fit-up-server/internal/schema/data"
	"github.com/tdmdh/fit-up-server/internal/schema/repository"
	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

type planGenerationServiceImpl struct {
	repo repository.SchemaRepo
}

func NewPlanGenerationService(repo repository.SchemaRepo) PlanGenerationService {
	return &planGenerationServiceImpl{
		repo: repo,
	}
}

// =============================================================================
// PLAN GENERATION METHODS
// =============================================================================

func (s *planGenerationServiceImpl) CreatePlanGeneration(ctx context.Context, userID int, metadata *types.PlanGenerationMetadata) (*types.GeneratedPlan, error) {
	if err := validator.New().Struct(metadata); err != nil {
		return nil, err
	}
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	activePlan, err := s.repo.PlanGeneration().GetActivePlanForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if activePlan != nil {
		return nil, types.ErrActivePlanExists
	}

	planMetadata, err := s.generateAdaptivePlan(ctx, userID, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to generate adaptive plan: %w", err)
	}

	return s.repo.PlanGeneration().CreatePlanGeneration(ctx, userID, planMetadata)
}

func (s *planGenerationServiceImpl) generateAdaptivePlan(_ context.Context, _ int, metadata *types.PlanGenerationMetadata) (*types.PlanGenerationMetadata, error) {
	if len(metadata.UserGoals) == 0 {
		return nil, fmt.Errorf("at least one fitness goal is required")
	}
	if len(metadata.AvailableEquipment) == 0 {
		return nil, fmt.Errorf("at least one equipment type is required")
	}
	if metadata.WeeklyFrequency <= 0 || metadata.WeeklyFrequency > 7 {
		return nil, fmt.Errorf("weekly frequency must be between 1 and 7")
	}

	fitupData, err := data.LoadFitUpData()
	if err != nil {
		return nil, fmt.Errorf("failed to load fitness data: %w", err)
	}

	userLevel := string(metadata.FitnessLevel)
	primaryGoal := metadata.UserGoals[0]

	levelData, exists := fitupData.Levels[userLevel]
	if !exists {
		return nil, fmt.Errorf("invalid fitness level: %s", userLevel)
	}

	goalData, exists := fitupData.Goals[string(primaryGoal)]
	if !exists {
		return nil, fmt.Errorf("invalid fitness goal: %s", primaryGoal)
	}

	template, err := s.selectOptimalTemplate(fitupData, userLevel, string(primaryGoal), metadata.WeeklyFrequency)
	if err != nil {
		return nil, fmt.Errorf("failed to select workout template: %w", err)
	}

	exerciseSelection, err := s.generateExerciseSelection(fitupData, template, metadata.AvailableEquipment, userLevel, string(primaryGoal))
	if err != nil {
		return nil, fmt.Errorf("failed to generate exercise selection: %w", err)
	}

	adaptedTemplate := s.applyProgressiveOverload(template, exerciseSelection, levelData, goalData, metadata.TimePerWorkout)

	balancedPlan := s.optimizeMuscleGroupBalance(adaptedTemplate, exerciseSelection)
	enhancedMetadata := &types.PlanGenerationMetadata{
		UserGoals:          metadata.UserGoals,
		AvailableEquipment: metadata.AvailableEquipment,
		FitnessLevel:       metadata.FitnessLevel,
		WeeklyFrequency:    metadata.WeeklyFrequency,
		TimePerWorkout:     metadata.TimePerWorkout,
		Algorithm:          "fitup_adaptive_v1",
		Parameters: map[string]any{
			"template_used":          template.ID,
			"total_exercises":        len(exerciseSelection),
			"muscle_groups_targeted": s.extractMuscleGroups(exerciseSelection),
			"equipment_utilized":     s.extractEquipmentTypes(exerciseSelection),
			"estimated_volume":       s.calculateWeeklyVolume(balancedPlan),
			"progression_method":     goalData.ProgressionMethods[0],
			"intensity_guidelines":   levelData.IntensityGuidelines,
			"generated_plan":         balancedPlan,
		},
	}

	return enhancedMetadata, nil
}

func (s *planGenerationServiceImpl) selectOptimalTemplate(fitupData *data.FitUpData, level, goal string, frequency int) (*data.WorkoutTemplate, error) {
	suitableTemplates := fitupData.GetWorkoutTemplateByGoalAndLevel(goal, level)
	if len(suitableTemplates) == 0 {
		return nil, fmt.Errorf("no suitable templates found for level %s and goal %s", level, goal)
	}

	for _, template := range suitableTemplates {
		if template.DaysPerWeek == frequency {
			return &template, nil
		}
	}

	var bestTemplate *data.WorkoutTemplate
	minDiff := 10
	for _, template := range suitableTemplates {
		diff := abs(template.DaysPerWeek - frequency)
		if diff < minDiff {
			minDiff = diff
			bestTemplate = &template
		}
	}

	if bestTemplate == nil {
		return &suitableTemplates[0], nil
	}

	return bestTemplate, nil
}

func (s *planGenerationServiceImpl) generateExerciseSelection(fitupData *data.FitUpData, template *data.WorkoutTemplate, availableEquipment []types.EquipmentType, level, _ string) ([]data.Exercise, error) {
	var selectedExercises []data.Exercise
	exerciseMap := make(map[int]data.Exercise)

	for _, ex := range fitupData.Exercises {
		exerciseMap[ex.ID] = ex
	}

	availableExercises := s.filterExercisesByEquipment(fitupData.Exercises, availableEquipment)

	bodyweightExercises := fitupData.GetExercisesByEquipment("bodyweight")
	availableExercises = append(availableExercises, bodyweightExercises...)

	levelAppropriateExercises := s.filterExercisesByLevel(availableExercises, level)

	for _, day := range template.Structure {
		for _, exerciseSpec := range day.Exercises {
			if exercise, exists := exerciseMap[exerciseSpec.ExerciseID]; exists {
				if s.isExerciseAvailable(exercise, availableEquipment) {
					selectedExercises = append(selectedExercises, exercise)
				} else {
					substitute := s.findExerciseSubstitute(exercise, levelAppropriateExercises)
					if substitute != nil {
						selectedExercises = append(selectedExercises, *substitute)
					}
				}
			}
		}
	}

	return selectedExercises, nil
}

func (s *planGenerationServiceImpl) applyProgressiveOverload(template *data.WorkoutTemplate, exercises []data.Exercise, level data.Level, goal data.Goal, timePerWorkout int) *data.WorkoutTemplate {
	adaptedTemplate := *template

	volumeMultiplier := s.getVolumeMultiplier(level.ID)

	repRange := goal.RepRanges.Primary

	for dayKey, day := range adaptedTemplate.Structure {
		adaptedDay := day
		for i, exerciseSpec := range day.Exercises {
			var targetExercise *data.Exercise
			for _, ex := range exercises {
				if ex.ID == exerciseSpec.ExerciseID {
					targetExercise = &ex
					break
				}
			}

			if targetExercise != nil {
				adaptedSpec := exerciseSpec
				adaptedSpec.Sets = int(float64(targetExercise.DefaultSets) * volumeMultiplier)
				adaptedSpec.Reps = s.adaptRepsForGoal(targetExercise.DefaultReps, repRange)
				adaptedSpec.Rest = s.adaptRestForGoal(targetExercise.RestSeconds, goal.RestPeriods, targetExercise.Type)

				adaptedDay.Exercises[i] = adaptedSpec
			}
		}
		adaptedTemplate.Structure[dayKey] = adaptedDay
	}

	return &adaptedTemplate
}

func (s *planGenerationServiceImpl) optimizeMuscleGroupBalance(template *data.WorkoutTemplate, exercises []data.Exercise) *data.WorkoutTemplate {
	muscleGroupCount := make(map[string]int)

	for _, day := range template.Structure {
		for _, exerciseSpec := range day.Exercises {
			for _, ex := range exercises {
				if ex.ID == exerciseSpec.ExerciseID {
					for _, mg := range ex.MuscleGroups {
						muscleGroupCount[mg]++
					}
					break
				}
			}
		}
	}

	// For now, return template as-is
	// In a more sophisticated implementation, we would:
	// - Identify underrepresented muscle groups
	// - Add corrective exercises
	// - Ensure push/pull balance
	// - Optimize for recovery patterns

	return template
}

func (s *planGenerationServiceImpl) filterExercisesByEquipment(exercises []data.Exercise, availableEquipment []types.EquipmentType) []data.Exercise {
	var filtered []data.Exercise
	equipmentSet := make(map[string]bool)

	for _, eq := range availableEquipment {
		equipmentSet[string(eq)] = true
	}

	for _, ex := range exercises {
		if equipmentSet[ex.Equipment] {
			filtered = append(filtered, ex)
		}
	}

	return filtered
}

func (s *planGenerationServiceImpl) filterExercisesByLevel(exercises []data.Exercise, level string) []data.Exercise {
	var filtered []data.Exercise
	levelPriority := map[string]int{"beginner": 1, "intermediate": 2, "advanced": 3}
	userLevelPriority := levelPriority[level]

	for _, ex := range exercises {
		exLevelPriority := levelPriority[ex.Difficulty]
		if exLevelPriority <= userLevelPriority {
			filtered = append(filtered, ex)
		}
	}

	return filtered
}

func (s *planGenerationServiceImpl) isExerciseAvailable(exercise data.Exercise, availableEquipment []types.EquipmentType) bool {
	for _, eq := range availableEquipment {
		if string(eq) == exercise.Equipment {
			return true
		}
	}
	return exercise.Equipment == "bodyweight"
}

func (s *planGenerationServiceImpl) findExerciseSubstitute(targetExercise data.Exercise, availableExercises []data.Exercise) *data.Exercise {
	for _, ex := range availableExercises {
		if ex.MovementPattern == targetExercise.MovementPattern {
			overlap := s.calculateMuscleGroupOverlap(targetExercise.MuscleGroups, ex.MuscleGroups)
			if overlap >= 0.5 {
				return &ex
			}
		}
	}

	if len(targetExercise.MuscleGroups) > 0 {
		primaryMuscle := targetExercise.MuscleGroups[0]
		for _, ex := range availableExercises {
			for _, mg := range ex.MuscleGroups {
				if mg == primaryMuscle {
					return &ex
				}
			}
		}
	}

	return nil
}

func (s *planGenerationServiceImpl) getVolumeMultiplier(level string) float64 {
	switch level {
	case "beginner":
		return 0.8
	case "intermediate":
		return 1.0
	case "advanced":
		return 1.2
	default:
		return 1.0
	}
}

func (s *planGenerationServiceImpl) adaptRepsForGoal(defaultReps, goalReps string) string {
	if goalReps != "" {
		return goalReps
	}
	return defaultReps
}

func (s *planGenerationServiceImpl) adaptRestForGoal(defaultRest int, goalRest data.RestPeriods, exerciseType string) int {
	if exerciseType == "strength" {
		if parsed := s.parseRestTime(goalRest.Compound); parsed > 0 {
			return parsed
		}
	}
	return defaultRest
}

func (s *planGenerationServiceImpl) extractMuscleGroups(exercises []data.Exercise) []string {
	muscleSet := make(map[string]bool)
	for _, ex := range exercises {
		for _, mg := range ex.MuscleGroups {
			muscleSet[mg] = true
		}
	}

	var muscles []string
	for mg := range muscleSet {
		muscles = append(muscles, mg)
	}
	return muscles
}

func (s *planGenerationServiceImpl) extractEquipmentTypes(exercises []data.Exercise) []string {
	equipSet := make(map[string]bool)
	for _, ex := range exercises {
		equipSet[ex.Equipment] = true
	}

	var equipment []string
	for eq := range equipSet {
		equipment = append(equipment, eq)
	}
	return equipment
}

func (s *planGenerationServiceImpl) calculateWeeklyVolume(template *data.WorkoutTemplate) int {
	totalSets := 0
	for _, day := range template.Structure {
		for _, ex := range day.Exercises {
			totalSets += ex.Sets
		}
	}
	return totalSets
}

func (s *planGenerationServiceImpl) calculateMuscleGroupOverlap(muscles1, muscles2 []string) float64 {
	if len(muscles1) == 0 || len(muscles2) == 0 {
		return 0.0
	}

	set1 := make(map[string]bool)
	for _, m := range muscles1 {
		set1[m] = true
	}

	overlap := 0
	for _, m := range muscles2 {
		if set1[m] {
			overlap++
		}
	}

	return float64(overlap) / float64(len(muscles1))
}

func (s *planGenerationServiceImpl) parseRestTime(restStr string) int {
	switch restStr {
	case "180-300":
		return 240
	case "120-180":
		return 150
	case "60-120":
		return 90
	case "45-90":
		return 75
	case "30-60":
		return 45
	default:
		return 0
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (s *planGenerationServiceImpl) GetActivePlanForUser(ctx context.Context, userID int) (*types.GeneratedPlan, error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	activePlan, err := s.repo.PlanGeneration().GetActivePlanForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if activePlan == nil {
		return nil, fmt.Errorf("no active plan found for user %d", userID)
	}

	if err := s.enrichPlanWithProgress(ctx, activePlan); err != nil {
		return nil, fmt.Errorf("failed to enrich plan with progress: %w", err)
	}
	return activePlan, nil
}

func (s *planGenerationServiceImpl) enrichPlanWithProgress(ctx context.Context, plan *types.GeneratedPlan) error {
	logs, err := s.repo.Progress().GetProgressLogsByUserID(ctx, plan.UserID, types.PaginationParams{Limit: 100, Offset: 0})
	if err != nil {
		return fmt.Errorf("failed to get progress logs: %w", err)
	}

	exerciseProgress := make(map[int][]types.ProgressLog)
	for _, log := range logs.Data {
		exerciseProgress[log.ExerciseID] = append(exerciseProgress[log.ExerciseID], log)
	}

	if plan.Metadata == nil {
		return nil
	}

	var metadata map[string]any
	if err := json.Unmarshal(plan.Metadata, &metadata); err != nil {
		return fmt.Errorf("failed to unmarshal plan metadata: %w", err)
	}

	parametersIface, ok := metadata["parameters"]
	if !ok {
		return nil
	}

	parameters, ok := parametersIface.(map[string]any)
	if !ok {
		return nil
	}

	generated, ok := parameters["generated_plan"].([]any)
	if !ok {
		return nil
	}
	for _, dayIface := range generated {
		dayMap, ok := dayIface.(map[string]any)
		if !ok {
			continue
		}

		exercisesIface, ok := dayMap["exercises"].([]any)
		if !ok {
			continue
		}

		for _, exIface := range exercisesIface {
			exMap, ok := exIface.(map[string]any)
			if !ok {
				continue
			}

			exIDFloat, ok := exMap["exercise_id"].(float64)
			if !ok {
				continue
			}
			exID := int(exIDFloat)

			if logs, exists := exerciseProgress[exID]; exists {
				var progressSummaries []map[string]any
				for _, log := range logs {
					progressSummaries = append(progressSummaries, map[string]any{
						"date":        log.Date,
						"sets":        log.SetsCompleted,
						"reps":        log.RepsCompleted,
						"weight":      log.WeightUsed,
						"duration":    log.DurationSeconds,
						"exercise_id": log.ExerciseID,
						"user_id":     log.UserID,
					})
				}
				exMap["progress_logs"] = progressSummaries
			}
		}

	}

	parameters["generated_plan"] = generated
	metadata["parameters"] = parameters

	updateMetadata, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal updated metadata: %w", err)
	}

	plan.Metadata = updateMetadata
	return nil
}

func (s *planGenerationServiceImpl) GetPlanGenerationHistory(ctx context.Context, userID int, limit int) ([]types.GeneratedPlan, error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	planHistory, err := s.repo.PlanGeneration().GetPlanGenerationHistory(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan generation history: %w", err)
	}

	for i := range planHistory {
		if err := s.enrichPlanWithProgress(ctx, &planHistory[i]); err != nil {
			fmt.Printf("Warning: Failed to enrich plan %d with progress: %v\n", planHistory[i].PlanID, err)
		}
	}

	if len(planHistory) == 0 {
		return []types.GeneratedPlan{}, nil
	}

	if len(planHistory) > 1 {
		evolution, err := s.analyzePlanEvolution(planHistory)
		if err != nil {
			fmt.Printf("Warning: Failed to analyze plan evolution: %v\n", err)
		} else {
			fmt.Printf("Plan evolution insights for user %d: %+v\n", userID, evolution)
		}
	}

	return planHistory, nil
}

func (s *planGenerationServiceImpl) analyzePlanEvolution(plans []types.GeneratedPlan) (map[string]any, error) {
	if len(plans) < 2 {
		return nil, nil
	}

	var evolutionInsights []string
	for i := 1; i < len(plans); i++ {
		var prevMetadata, currMetadata map[string]any
		if err := json.Unmarshal(plans[i-1].Metadata, &prevMetadata); err != nil {
			continue
		}
		if err := json.Unmarshal(plans[i].Metadata, &currMetadata); err != nil {
			continue
		}

		prevParams, ok1 := prevMetadata["parameters"].(map[string]any)
		currParams, ok2 := currMetadata["parameters"].(map[string]any)
		if !ok1 || !ok2 {
			continue
		}

		prevTemplate, _ := prevParams["template_used"].(string)
		currTemplate, _ := currParams["template_used"].(string)
		if prevTemplate != currTemplate {
			evolutionInsights = append(evolutionInsights, fmt.Sprintf("Changed template from %s to %s", prevTemplate, currTemplate))
		}

		prevVolume, _ := prevParams["estimated_volume"].(float64)
		currVolume, _ := currParams["estimated_volume"].(float64)
		if currVolume > prevVolume {
			evolutionInsights = append(evolutionInsights, fmt.Sprintf("Increased weekly volume from %.0f to %.0f sets", prevVolume, currVolume))
		} else if currVolume < prevVolume {
			evolutionInsights = append(evolutionInsights, fmt.Sprintf("Decreased weekly volume from %.0f to %.0f sets", prevVolume, currVolume))
		}

		prevMuscles, _ := prevParams["muscle_groups_targeted"].([]string)
		currMuscles, _ := currParams["muscle_groups_targeted"].([]string)
		if len(prevMuscles) != len(currMuscles) {
			evolutionInsights = append(evolutionInsights, fmt.Sprintf("Changed muscle groups targeted from %v to %v", prevMuscles, currMuscles))
		}
	}

	insights := map[string]any{
		"total_plans":     len(plans),
		"evolution_notes": evolutionInsights,
	}
	return insights, nil
}

// =============================================================================
// PLAN PERFORMANCE TRACKING METHODS
// =============================================================================

func (s *planGenerationServiceImpl) TrackPlanPerformance(ctx context.Context, planID int, performance *types.PlanPerformanceData) error {
	if planID <= 0 {
		return fmt.Errorf("invalid plan ID")
	}

	if performance == nil {
		return fmt.Errorf("performance data cannot be nil")
	}

	if err := validator.New().Struct(performance); err != nil {
		return fmt.Errorf("invalid performance data: %w", err)
	}

	if err := s.repo.PlanGeneration().TrackPlanPerformance(ctx, planID, performance); err != nil {
		return fmt.Errorf("failed to track plan performance: %w", err)
	}

	if err := s.analyzeAndAdaptPlan(ctx, planID, performance); err != nil {
		// Log error but don't fail the main operation
		fmt.Printf("Warning: Failed to analyze plan for adaptations: %v\n", err)
	}

	return nil
}

func (s *planGenerationServiceImpl) analyzeAndAdaptPlan(ctx context.Context, planID int, performance *types.PlanPerformanceData) error {
	if performance.CompletionRate < 0.6 {
		changes, _ := json.Marshal(map[string]any{
			"type":        "volume_reduction",
			"description": "Reduced workout intensity and volume due to low completion rate",
			"adjustments": []string{"reduced_sets", "reduced_intensity"},
		})

		adaptation := &types.PlanAdaptation{
			PlanID:         planID,
			Reason:         "low_completion_rate",
			Changes:        changes,
			Trigger:        "automatic_analysis",
			AdaptationDate: time.Now(),
		}
		return s.LogPlanAdaptation(ctx, planID, adaptation)
	}

	if performance.AverageRPE > 8.5 && performance.CompletionRate < 0.8 {
		changes, _ := json.Marshal(map[string]any{
			"type":        "recovery_focus",
			"description": "Added rest days and reduced intensity due to high RPE and low completion",
			"adjustments": []string{"additional_rest_days", "reduced_intensity"},
		})

		adaptation := &types.PlanAdaptation{
			PlanID:         planID,
			Reason:         "potential_overtraining",
			Changes:        changes,
			Trigger:        "automatic_analysis",
			AdaptationDate: time.Now(),
		}
		return s.LogPlanAdaptation(ctx, planID, adaptation)
	}

	if performance.CompletionRate > 0.9 && performance.AverageRPE < 6.0 {
		changes, _ := json.Marshal(map[string]any{
			"type":        "progression",
			"description": "Increased volume and intensity due to high completion rate and low RPE",
			"adjustments": []string{"increased_volume", "increased_intensity"},
		})

		adaptation := &types.PlanAdaptation{
			PlanID:         planID,
			Reason:         "ready_for_progression",
			Changes:        changes,
			Trigger:        "automatic_analysis",
			AdaptationDate: time.Now(),
		}
		return s.LogPlanAdaptation(ctx, planID, adaptation)
	}

	return nil
}

func (s *planGenerationServiceImpl) GetPlanEffectivenessScore(ctx context.Context, planID int) (float64, error) {
	if planID <= 0 {
		return 0, fmt.Errorf("invalid plan ID")
	}

	score, err := s.repo.PlanGeneration().GetPlanEffectivenessScore(ctx, planID)
	if err != nil {
		return 0, fmt.Errorf("failed to get plan effectiveness score: %w", err)
	}

	if score < 0 || score > 100 {
		score = 75.0
	}

	return score, nil
}

func (s *planGenerationServiceImpl) MarkPlanForRegeneration(ctx context.Context, planID int, reason string) error {
	if planID <= 0 {
		return fmt.Errorf("invalid plan ID")
	}

	if reason == "" {
		return fmt.Errorf("reason for regeneration cannot be empty")
	}

	changes, _ := json.Marshal(map[string]any{
		"type":        "plan_regeneration",
		"description": fmt.Sprintf("Plan marked for regeneration: %s", reason),
		"status":      "pending_regeneration",
	})

	adaptation := &types.PlanAdaptation{
		PlanID:         planID,
		Reason:         "plan_regeneration_request",
		Changes:        changes,
		Trigger:        reason,
		AdaptationDate: time.Now(),
	}

	if err := s.LogPlanAdaptation(ctx, planID, adaptation); err != nil {
		return fmt.Errorf("failed to log regeneration adaptation: %w", err)
	}

	return s.repo.PlanGeneration().MarkPlanForRegeneration(ctx, planID, reason)
}

// =============================================================================
// PLAN ADAPTATION METHODS
// =============================================================================

func (s *planGenerationServiceImpl) LogPlanAdaptation(ctx context.Context, planID int, adaptation *types.PlanAdaptation) error {
	if planID <= 0 {
		return fmt.Errorf("invalid plan ID")
	}

	if adaptation == nil {
		return fmt.Errorf("adaptation data cannot be nil")
	}

	if adaptation.Reason == "" {
		return fmt.Errorf("adaptation reason cannot be empty")
	}

	adaptation.PlanID = planID
	if adaptation.AdaptationDate.IsZero() {
		adaptation.AdaptationDate = time.Now()
	}
	if adaptation.Trigger == "" {
		adaptation.Trigger = "manual"
	}

	return s.repo.PlanGeneration().LogPlanAdaptation(ctx, planID, adaptation)
}

func (s *planGenerationServiceImpl) GetAdaptationHistory(ctx context.Context, userID int) ([]types.PlanAdaptation, error) {
	if userID <= 0 {
		return nil, types.ErrInvalidUserID
	}

	adaptations, err := s.repo.PlanGeneration().GetAdaptationHistory(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get adaptation history: %w", err)
	}

	if len(adaptations) > 0 {
		insights := s.analyzeAdaptationPatterns(adaptations)
		fmt.Printf("Adaptation insights for user %d: %+v\n", userID, insights)
	}

	return adaptations, nil
}

// =============================================================================
// TEMPLATE MANAGEMENT METHODS (merged from WorkoutTemplateService)
// =============================================================================

func (s *planGenerationServiceImpl) GetTemplateByID(ctx context.Context, templateID int) (*types.WorkoutTemplate, error) {
	return s.repo.Templates().GetTemplateByID(ctx, templateID)
}

func (s *planGenerationServiceImpl) ListTemplates(ctx context.Context, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutTemplate], error) {
	return s.repo.Templates().ListTemplates(ctx, pagination)
}

func (s *planGenerationServiceImpl) FilterTemplates(ctx context.Context, filter types.TemplateFilter, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutTemplate], error) {
	return s.repo.Templates().FilterTemplates(ctx, filter, pagination)
}

func (s *planGenerationServiceImpl) SearchTemplates(ctx context.Context, query string, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutTemplate], error) {
	return s.repo.Templates().SearchTemplates(ctx, query, pagination)
}

func (s *planGenerationServiceImpl) GetTemplatesByLevel(ctx context.Context, level types.FitnessLevel) ([]types.WorkoutTemplate, error) {
	return s.repo.Templates().GetTemplatesByLevel(ctx, level)
}

func (s *planGenerationServiceImpl) GetTemplatesByGoal(ctx context.Context, goal types.FitnessGoal) ([]types.WorkoutTemplate, error) {
	return s.repo.Templates().GetTemplatesByGoal(ctx, goal)
}

func (s *planGenerationServiceImpl) GetRecommendedTemplates(ctx context.Context, userID int, count int) ([]types.WorkoutTemplate, error) {
	return s.repo.Templates().GetRecommendedTemplates(ctx, userID, count)
}

func (s *planGenerationServiceImpl) GetPopularTemplates(ctx context.Context, count int) ([]types.WorkoutTemplate, error) {
	return s.repo.Templates().GetPopularTemplates(ctx, count)
}

// =============================================================================
// WEEKLY SCHEMA MANAGEMENT METHODS (merged from WeeklySchemaService)
// =============================================================================

func (s *planGenerationServiceImpl) GetWeeklySchemaByID(ctx context.Context, schemaID int) (*types.WeeklySchema, error) {
	return s.repo.Schemas().GetWeeklySchemaByID(ctx, schemaID)
}

func (s *planGenerationServiceImpl) GetWeeklySchemasByUserID(ctx context.Context, userID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.WeeklySchema], error) {
	return s.repo.Schemas().GetWeeklySchemasByUserID(ctx, userID, pagination)
}

func (s *planGenerationServiceImpl) GetActiveWeeklySchemaByUserID(ctx context.Context, userID int) (*types.WeeklySchema, error) {
	return s.repo.Schemas().GetActiveWeeklySchemaByUserID(ctx, userID)
}

func (s *planGenerationServiceImpl) GetWeeklySchemaByUserAndWeek(ctx context.Context, userID int, weekStart time.Time) (*types.WeeklySchema, error) {
	return s.repo.Schemas().GetWeeklySchemaByUserAndWeek(ctx, userID, weekStart)
}

func (s *planGenerationServiceImpl) GetCurrentWeekSchema(ctx context.Context, userID int) (*types.WeeklySchema, error) {
	return s.repo.Schemas().GetCurrentWeekSchema(ctx, userID)
}

func (s *planGenerationServiceImpl) GetWeeklySchemaHistory(ctx context.Context, userID int, limit int) ([]types.WeeklySchema, error) {
	return s.repo.Schemas().GetWeeklySchemaHistory(ctx, userID, limit)
}

func (s *planGenerationServiceImpl) CreateWeeklySchemaFromTemplate(ctx context.Context, userID, templateID int, weekStart time.Time) (*types.WeeklySchemaWithWorkouts, error) {
	// Create weekly schema request
	schemaRequest := &types.WeeklySchemaRequest{
		UserID:    userID,
		WeekStart: weekStart,
	}

	// Create the weekly schema
	schema, err := s.repo.Schemas().CreateWeeklySchema(ctx, schemaRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create weekly schema: %w", err)
	}

	return &types.WeeklySchemaWithWorkouts{
		SchemaID:  schema.SchemaID,
		UserID:    schema.UserID,
		WeekStart: schema.WeekStart,
		Active:    schema.Active,
		Workouts:  []types.WorkoutWithExercises{},
	}, nil
}

func (s *planGenerationServiceImpl) analyzeAdaptationPatterns(adaptations []types.PlanAdaptation) map[string]any {
	if len(adaptations) == 0 {
		return nil
	}

	reasonCounts := make(map[string]int)
	triggerCounts := make(map[string]int)

	var recentAdaptations []types.PlanAdaptation
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	for _, adaptation := range adaptations {
		reasonCounts[adaptation.Reason]++
		triggerCounts[adaptation.Trigger]++

		if adaptation.AdaptationDate.After(thirtyDaysAgo) {
			recentAdaptations = append(recentAdaptations, adaptation)
		}
	}

	mostCommonReason := ""
	maxCount := 0
	for reason, count := range reasonCounts {
		if count > maxCount {
			maxCount = count
			mostCommonReason = reason
		}
	}

	insights := map[string]any{
		"total_adaptations":       len(adaptations),
		"recent_adaptations":      len(recentAdaptations),
		"most_common_reason":      mostCommonReason,
		"reason_frequency":        reasonCounts,
		"trigger_frequency":       triggerCounts,
		"adaptation_rate_30_days": float64(len(recentAdaptations)) / 30.0,
	}

	return insights
}

// =============================================================================
// ADDITIONAL HELPER METHODS
// =============================================================================

func (s *planGenerationServiceImpl) validatePlanGenerationRequest(metadata *types.PlanGenerationMetadata) error {
	if metadata == nil {
		return fmt.Errorf("plan generation metadata cannot be nil")
	}

	if len(metadata.UserGoals) == 0 {
		return fmt.Errorf("at least one fitness goal is required")
	}

	if len(metadata.AvailableEquipment) == 0 {
		return fmt.Errorf("at least one equipment type is required")
	}

	if metadata.WeeklyFrequency <= 0 || metadata.WeeklyFrequency > 7 {
		return fmt.Errorf("weekly frequency must be between 1 and 7")
	}

	if metadata.TimePerWorkout <= 0 || metadata.TimePerWorkout > 300 {
		return fmt.Errorf("time per workout must be between 1 and 300 minutes")
	}

	validLevels := map[types.FitnessLevel]bool{
		types.LevelBeginner:     true,
		types.LevelIntermediate: true,
		types.LevelAdvanced:     true,
	}

	if !validLevels[metadata.FitnessLevel] {
		return fmt.Errorf("invalid fitness level: %s", metadata.FitnessLevel)
	}

	return nil
}

func (s *planGenerationServiceImpl) calculatePlanComplexity(exercises []data.Exercise, frequency int) float64 {
	if len(exercises) == 0 {
		return 0.0
	}

	complexityScore := 0.0
	difficultyWeights := map[string]float64{
		"beginner":     1.0,
		"intermediate": 2.0,
		"advanced":     3.0,
	}

	for _, exercise := range exercises {
		if weight, exists := difficultyWeights[exercise.Difficulty]; exists {
			complexityScore += weight
		}
	}

	averageComplexity := complexityScore / float64(len(exercises))
	frequencyMultiplier := 1.0 + (float64(frequency-1) * 0.1)

	return averageComplexity * frequencyMultiplier
}

func (s *planGenerationServiceImpl) estimatePlanDuration(level types.FitnessLevel, goals []types.FitnessGoal) int {
	baseDuration := map[types.FitnessLevel]int{
		types.LevelBeginner:     4, // 4 weeks
		types.LevelIntermediate: 6, // 6 weeks
		types.LevelAdvanced:     8, // 8 weeks
	}

	duration := baseDuration[level]

	if len(goals) > 0 {
		switch goals[0] {
		case types.GoalMuscleGain:
			duration += 2 // Muscle building needs longer phases
		case types.GoalFatLoss:
			duration += 1 // Weight loss benefits from longer consistency
		case types.GoalEndurance:
			duration += 1 // Endurance building is gradual
		case types.GoalStrength:
			duration += 2 // Strength building needs progressive phases
		}
	}

	return duration
}

func (s *planGenerationServiceImpl) optimizeForUserPreferences(plan *data.WorkoutTemplate, userID int, metadata *types.PlanGenerationMetadata) *data.WorkoutTemplate {
	optimizedPlan := *plan

	if metadata.TimePerWorkout > 0 {
		optimizedPlan = *s.adjustForTimeConstraints(&optimizedPlan, metadata.TimePerWorkout)
	}

	if len(metadata.AvailableEquipment) > 0 {
		optimizedPlan = *s.adaptForEquipment(&optimizedPlan, metadata.AvailableEquipment)
	}

	return &optimizedPlan
}

// adjustForTimeConstraints modifies the workout to fit within time constraints
func (s *planGenerationServiceImpl) adjustForTimeConstraints(plan *data.WorkoutTemplate, maxMinutes int) *data.WorkoutTemplate {
	if maxMinutes >= 60 {
		return plan
	}

	adjustedPlan := *plan

	restReduction := 1.0
	if maxMinutes < 30 {
		restReduction = 0.6 // 40% reduction in rest time
	} else if maxMinutes < 45 {
		restReduction = 0.8 // 20% reduction in rest time
	}

	for dayKey, day := range adjustedPlan.Structure {
		adjustedDay := day
		for i, exercise := range day.Exercises {
			exercise.Rest = int(float64(exercise.Rest) * restReduction)
			adjustedDay.Exercises[i] = exercise
		}

		if maxMinutes < 30 && len(adjustedDay.Exercises) > 4 {
			adjustedDay.Exercises = adjustedDay.Exercises[:4]
		}

		adjustedPlan.Structure[dayKey] = adjustedDay
	}

	return &adjustedPlan
}

// adaptForEquipment ensures the plan only uses available equipment
func (s *planGenerationServiceImpl) adaptForEquipment(plan *data.WorkoutTemplate, availableEquipment []types.EquipmentType) *data.WorkoutTemplate {
	// This is a simplified version - in practice, you'd need to:
	// 1. Check each exercise in the plan
	// 2. Replace exercises that require unavailable equipment
	// 3. Maintain the balance and effectiveness of the plan

	return plan // For now, return as-is
}

func (s *planGenerationServiceImpl) validatePlanGeneration(plan *types.GeneratedPlan, metadata *types.PlanGenerationMetadata) error {
	if plan == nil {
		return fmt.Errorf("plan cannot be nil")
	}

	if plan.UserID <= 0 {
		return fmt.Errorf("invalid user ID in plan")
	}

	if metadata == nil {
		return fmt.Errorf("plan metadata cannot be nil")
	}

	if len(metadata.UserGoals) == 0 {
		return fmt.Errorf("plan must have at least one goal")
	}

	if metadata.WeeklyFrequency <= 0 {
		return fmt.Errorf("plan must have positive weekly frequency")
	}

	return nil
}

func (s *planGenerationServiceImpl) ExportPlanToPDF(ctx context.Context, planID int) ([]byte, error) {
	if planID <= 0 {
		return nil, fmt.Errorf("invalid plan ID")
	}

	plan, err := s.repo.PlanGeneration().GetPlanID(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan: %w", err)
	}

	if plan == nil {
		return nil, fmt.Errorf("no plan found with ID %d", planID)
	}

	var metadata map[string]any
	if err := json.Unmarshal(plan.Metadata, &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plan metadata: %w", err)
	}

	parametersIface, ok := metadata["parameters"]
	if !ok {
		return nil, fmt.Errorf("plan metadata missing parameters")
	}

	parameters, ok := parametersIface.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid parameters format in metadata")
	}

	generatedIface, ok := parameters["generated_plan"]
	if !ok {
		return nil, fmt.Errorf("generated_plan not found in parameters")
	}

	generated, ok := generatedIface.([]any)
	if !ok {
		return nil, fmt.Errorf("invalid generated_plan format")
	}

	// Initialize PDF with better margins
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.SetAutoPageBreak(true, 20)
	pdf.AddPage()

	// Helper function for drawing header line
	drawHeaderLine := func() {
		pdf.SetDrawColor(41, 128, 185) // Blue color
		pdf.SetLineWidth(0.5)
		pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
		pdf.Ln(3)
	}

	// Header Section with Logo/Title
	pdf.SetFillColor(41, 128, 185)  // Blue background
	pdf.SetTextColor(255, 255, 255) // White text
	pdf.SetFont("Arial", "B", 24)
	pdf.CellFormat(0, 15, "FIT-UP", "", 0, "C", true, 0, "")
	pdf.Ln(15)

	pdf.SetFillColor(236, 240, 241) // Light gray background
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 18)
	pdf.CellFormat(0, 12, "PERSONALIZED WORKOUT PLAN", "", 0, "C", true, 0, "")
	pdf.Ln(15)

	// Plan Metadata Section
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(52, 73, 94) // Dark blue-gray
	pdf.Cell(0, 8, "PLAN INFORMATION")
	pdf.Ln(8)
	drawHeaderLine()

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(0, 0, 0)

	// Plan details in a table-like format
	leftCol := 50.0
	pdf.Cell(leftCol, 6, "Plan ID:")
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, fmt.Sprintf("#%d", plan.PlanID))
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(leftCol, 6, "Generated:")
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, plan.GeneratedAt.Format("January 2, 2006 at 3:04 PM"))
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(leftCol, 6, "Week Start:")
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, plan.WeekStart.Format("Monday, January 2, 2006"))
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(leftCol, 6, "Algorithm:")
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, plan.Algorithm)
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(leftCol, 6, "Status:")
	pdf.SetFont("Arial", "B", 10)
	statusText := "Active"
	if !plan.IsActive {
		statusText = "Inactive"
	}
	pdf.Cell(0, 6, statusText)
	pdf.Ln(12)

	// Training Parameters Section
	if templateUsed, ok := parameters["template_used"].(string); ok {
		pdf.SetFont("Arial", "B", 12)
		pdf.SetTextColor(52, 73, 94)
		pdf.Cell(0, 8, "TRAINING PARAMETERS")
		pdf.Ln(8)
		drawHeaderLine()

		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(0, 0, 0)

		pdf.Cell(leftCol, 6, "Template:")
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(0, 6, templateUsed)
		pdf.Ln(6)

		if totalEx, ok := parameters["total_exercises"].(float64); ok {
			pdf.SetFont("Arial", "", 10)
			pdf.Cell(leftCol, 6, "Total Exercises:")
			pdf.SetFont("Arial", "B", 10)
			pdf.Cell(0, 6, fmt.Sprintf("%.0f", totalEx))
			pdf.Ln(6)
		}

		if muscleGroups, ok := parameters["muscle_groups_targeted"].([]interface{}); ok {
			pdf.SetFont("Arial", "", 10)
			pdf.Cell(leftCol, 6, "Muscle Groups:")
			pdf.SetFont("Arial", "B", 10)
			muscleList := make([]string, len(muscleGroups))
			for i, mg := range muscleGroups {
				if mgStr, ok := mg.(string); ok {
					muscleList[i] = mgStr
				}
			}
			pdf.MultiCell(0, 6, fmt.Sprintf("%s", muscleList), "", "", false)
		}

		if equipment, ok := parameters["equipment_utilized"].([]interface{}); ok && len(equipment) > 0 {
			pdf.SetFont("Arial", "", 10)
			pdf.Cell(leftCol, 6, "Equipment:")
			pdf.SetFont("Arial", "B", 10)
			equipList := make([]string, len(equipment))
			for i, eq := range equipment {
				if eqStr, ok := eq.(string); ok {
					equipList[i] = eqStr
				}
			}
			pdf.MultiCell(0, 6, fmt.Sprintf("%s", equipList), "", "", false)
		}

		pdf.Ln(6)
	}

	// Weekly Schedule Section
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(52, 73, 94)
	pdf.Cell(0, 10, "WEEKLY TRAINING SCHEDULE")
	pdf.Ln(10)
	drawHeaderLine()
	pdf.Ln(2)

	// Workout days with enhanced formatting
	totalWorkoutTime := 0
	totalExercises := 0

	for dayIdx, dayIface := range generated {
		dayMap, ok := dayIface.(map[string]any)
		if !ok {
			continue
		}

		dayTitle, _ := dayMap["day_title"].(string)
		focus, _ := dayMap["focus"].(string)

		// Day header with background
		pdf.SetFillColor(52, 152, 219) // Bright blue
		pdf.SetTextColor(255, 255, 255)
		pdf.SetFont("Arial", "B", 13)
		pdf.CellFormat(0, 10, fmt.Sprintf("DAY %d: %s", dayIdx+1, dayTitle), "", 0, "L", true, 0, "")
		pdf.Ln(10)

		// Focus area
		if focus != "" && focus != "recovery" {
			pdf.SetTextColor(0, 0, 0)
			pdf.SetFont("Arial", "I", 10)
			pdf.SetFillColor(236, 240, 241)
			pdf.CellFormat(0, 6, fmt.Sprintf("Focus: %s", focus), "", 0, "L", true, 0, "")
			pdf.Ln(8)
		}

		exercisesIface, ok := dayMap["exercises"].([]any)
		if !ok || len(exercisesIface) == 0 {
			pdf.SetTextColor(127, 140, 141)
			pdf.SetFont("Arial", "I", 10)
			pdf.Cell(0, 8, "Rest day - Focus on recovery, stretching, and mobility work")
			pdf.Ln(12)
			continue
		}

		// Exercise table headers
		pdf.SetFillColor(189, 195, 199) // Light gray
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(10, 8, "#", "1", 0, "C", true, 0, "")
		pdf.CellFormat(70, 8, "Exercise", "1", 0, "L", true, 0, "")
		pdf.CellFormat(25, 8, "Sets", "1", 0, "C", true, 0, "")
		pdf.CellFormat(30, 8, "Reps", "1", 0, "C", true, 0, "")
		pdf.CellFormat(35, 8, "Rest", "1", 0, "C", true, 0, "")
		pdf.Ln(8)

		// Exercise rows
		pdf.SetFont("Arial", "", 9)
		dayEstimatedTime := 0

		for exIdx, exIface := range exercisesIface {
			exMap, ok := exIface.(map[string]any)
			if !ok {
				continue
			}

			exName, _ := exMap["name"].(string)
			setsFloat, _ := exMap["sets"].(float64)
			reps, _ := exMap["reps"].(string)
			restFloat, _ := exMap["rest"].(float64)
			notes, _ := exMap["notes"].(string)

			sets := int(setsFloat)
			rest := int(restFloat)

			// Alternating row colors
			if exIdx%2 == 0 {
				pdf.SetFillColor(245, 245, 245)
			} else {
				pdf.SetFillColor(255, 255, 255)
			}

			pdf.CellFormat(10, 7, fmt.Sprintf("%d", exIdx+1), "1", 0, "C", true, 0, "")
			pdf.CellFormat(70, 7, exName, "1", 0, "L", true, 0, "")
			pdf.CellFormat(25, 7, fmt.Sprintf("%d", sets), "1", 0, "C", true, 0, "")
			pdf.CellFormat(30, 7, reps, "1", 0, "C", true, 0, "")
			pdf.CellFormat(35, 7, fmt.Sprintf("%d sec", rest), "1", 0, "C", true, 0, "")
			pdf.Ln(7)

			// Exercise notes if available
			if notes != "" {
				pdf.SetFont("Arial", "I", 8)
				pdf.SetTextColor(95, 95, 95)
				pdf.CellFormat(10, 5, "", "", 0, "L", false, 0, "")
				pdf.MultiCell(160, 5, fmt.Sprintf("Note: %s", notes), "", "L", false)
				pdf.SetFont("Arial", "", 9)
				pdf.SetTextColor(0, 0, 0)
			}

			// Estimate time (assume 3 seconds per rep, sets * reps * 3 + rest time between sets)
			repsNum := 10 // default estimation
			if reps != "" && reps != "AMRAP" && reps != "To failure" {
				fmt.Sscanf(reps, "%d", &repsNum)
			}
			exerciseTime := (sets * repsNum * 3) + (sets-1)*rest
			dayEstimatedTime += exerciseTime
			totalExercises++
		}

		totalWorkoutTime += dayEstimatedTime

		// Day summary
		pdf.SetFont("Arial", "I", 9)
		pdf.SetTextColor(52, 73, 94)
		pdf.Ln(2)
		pdf.Cell(0, 6, fmt.Sprintf("Estimated Duration: %d minutes | Exercises: %d", dayEstimatedTime/60, len(exercisesIface)))
		pdf.Ln(12)
	}

	// Add a new page for training notes and tips
	pdf.AddPage()

	// Training Guidelines Section
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(52, 73, 94)
	pdf.Cell(0, 10, "TRAINING GUIDELINES & TIPS")
	pdf.Ln(10)
	drawHeaderLine()
	pdf.Ln(4)

	guidelines := []struct {
		title string
		text  string
	}{
		{
			title: "Warm-Up Protocol",
			text:  "Always begin each session with 5-10 minutes of light cardio and dynamic stretching. This increases blood flow, raises body temperature, and prepares your muscles and joints for the workout ahead.",
		},
		{
			title: "Proper Form",
			text:  "Quality over quantity - focus on controlled movements and proper form rather than lifting heavy weights. Poor form can lead to injuries and reduces the effectiveness of exercises.",
		},
		{
			title: "Rest Between Sets",
			text:  "Follow the prescribed rest periods. Shorter rests (30-60s) are better for endurance and fat loss, while longer rests (90-180s) support strength and power development.",
		},
		{
			title: "Progressive Overload",
			text:  "Gradually increase the difficulty of your workouts by adding weight, increasing reps, or reducing rest time. Track your progress to ensure continuous improvement.",
		},
		{
			title: "Recovery",
			text:  "Rest days are crucial for muscle growth and recovery. Get adequate sleep (7-9 hours), stay hydrated, and consider active recovery like walking or yoga on rest days.",
		},
		{
			title: "Nutrition",
			text:  "Support your training with proper nutrition. Consume adequate protein (1.6-2.2g per kg body weight), stay hydrated, and time your meals around your workouts for optimal results.",
		},
		{
			title: "Listen to Your Body",
			text:  "If you experience sharp pain, excessive fatigue, or signs of overtraining, take extra rest. It's better to miss one workout than to be sidelined with an injury.",
		},
	}

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(0, 0, 0)

	for _, guide := range guidelines {
		pdf.SetFont("Arial", "B", 11)
		pdf.SetTextColor(41, 128, 185)
		pdf.Write(6, "• "+guide.title+": ")

		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(0, 0, 0)
		pdf.MultiCell(0, 5, guide.text, "", "L", false)
		pdf.Ln(3)
	}

	pdf.Ln(5)

	// Exercise Modifications
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(52, 73, 94)
	pdf.Cell(0, 8, "EXERCISE MODIFICATIONS")
	pdf.Ln(8)
	drawHeaderLine()
	pdf.Ln(2)

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(0, 0, 0)
	pdf.MultiCell(0, 5, "If an exercise feels uncomfortable or you don't have the required equipment:", "", "L", false)
	pdf.Ln(2)

	modifications := []string{
		"Push-ups → Knee push-ups, incline push-ups, or wall push-ups",
		"Pull-ups → Assisted pull-ups, inverted rows, or resistance band pull-downs",
		"Squats → Chair squats, wall sits, or split squats",
		"Deadlifts → Romanian deadlifts, single-leg deadlifts, or good mornings",
		"Bench Press → Floor press, dumbbell press, or resistance band press",
	}

	pdf.SetFont("Arial", "", 9)
	for _, mod := range modifications {
		pdf.Cell(5, 5, "")
		pdf.Cell(0, 5, mod)
		pdf.Ln(5)
	}

	pdf.Ln(8)

	// Weekly Summary Section
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(52, 73, 94)
	pdf.Cell(0, 8, "WEEKLY SUMMARY")
	pdf.Ln(8)
	drawHeaderLine()
	pdf.Ln(2)

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(0, 0, 0)

	workoutDays := 0
	for _, dayIface := range generated {
		if dayMap, ok := dayIface.(map[string]any); ok {
			if exercises, ok := dayMap["exercises"].([]any); ok && len(exercises) > 0 {
				workoutDays++
			}
		}
	}

	pdf.Cell(0, 6, fmt.Sprintf("Total Workout Days: %d", workoutDays))
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Total Exercises: %d", totalExercises))
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Estimated Weekly Training Time: %d minutes (%d hours)", totalWorkoutTime/60, totalWorkoutTime/3600))
	pdf.Ln(12)

	// Footer
	pdf.SetY(-30)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(127, 140, 141)
	pdf.CellFormat(0, 4, "Generated by FIT-UP Adaptive Training System", "", 0, "C", false, 0, "")
	pdf.Ln(4)
	pdf.CellFormat(0, 4, fmt.Sprintf("Plan ID: #%d | Generated: %s", plan.PlanID, plan.GeneratedAt.Format("2006-01-02")), "", 0, "C", false, 0, "")
	pdf.Ln(4)
	pdf.CellFormat(0, 4, "For questions or support, contact: support@fitup.com", "", 0, "C", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}
