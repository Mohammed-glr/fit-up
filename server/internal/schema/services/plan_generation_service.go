package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jung-kurt/gofpdf"
	"github.com/tdmdh/fit-up-server/internal/schema/data"
	"github.com/tdmdh/fit-up-server/internal/schema/repository"
	"github.com/tdmdh/fit-up-server/internal/schema/types"
	"github.com/tdmdh/fit-up-server/shared/middleware"
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
	if metadata == nil {
		slog.Error("plan generation metadata missing", slog.Int("user_id", userID))
		return nil, fmt.Errorf("plan generation metadata cannot be nil")
	}

	if err := validator.New().Struct(metadata); err != nil {
		slog.Warn("invalid plan generation metadata", slog.Int("user_id", userID), slog.Any("error", err))
		return nil, err
	}

	resolvedUserID, authUserID, err := s.resolveUserIdentity(ctx, userID, metadata, true)
	if err != nil {
		slog.Warn("failed to resolve user identity for plan generation", slog.Int("requested_user_id", userID), slog.Any("error", err))
		return nil, err
	}

	userID = resolvedUserID

	slog.Info("starting plan generation", slog.Int("user_id", userID), slog.Int("weekly_frequency", metadata.WeeklyFrequency), slog.String("fitness_level", string(metadata.FitnessLevel)), slog.Any("goals", metadata.UserGoals))

	activeCount, err := s.repo.PlanGeneration().CountActivePlans(ctx, userID)
	if err != nil {
		slog.Error("failed to count active plans", slog.Int("user_id", userID), slog.Any("error", err))
		return nil, err
	}
	if activeCount >= 3 {
		slog.Warn("plan generation blocked due to active plan limit", slog.Int("user_id", userID), slog.Int("active_plans", activeCount))
		return nil, types.ErrPlanLimitReached
	}

	planMetadata, err := s.generateAdaptivePlan(ctx, userID, metadata)
	if err != nil {
		slog.Error("adaptive plan generation failed", slog.Int("user_id", userID), slog.Any("error", err))
		return nil, fmt.Errorf("failed to generate adaptive plan: %w", err)
	}

	plan, err := s.repo.PlanGeneration().CreatePlanGeneration(ctx, userID, authUserID, planMetadata)
	if err != nil {
		slog.Error("failed to persist generated plan", slog.Int("user_id", userID), slog.Any("error", err))
		return nil, err
	}
	if plan == nil {
		slog.Error("repository returned nil plan", slog.Int("user_id", userID))
		return nil, fmt.Errorf("failed to create plan generation")
	}

	if err := s.persistGeneratedStructure(ctx, plan.PlanID, planMetadata); err != nil {
		slog.Warn("failed to persist generated structure", slog.Int("plan_id", plan.PlanID), slog.Any("error", err))
	}

	slog.Info("plan generation completed", slog.Int("user_id", userID), slog.Int("plan_id", plan.PlanID))

	return plan, nil
}

func (s *planGenerationServiceImpl) persistGeneratedStructure(ctx context.Context, planID int, metadata *types.PlanGenerationMetadata) error {
	if metadata == nil || metadata.Parameters == nil {
		return nil
	}

	generatedValue, exists := metadata.Parameters["generated_plan"]
	if !exists {
		return nil
	}

	generated, ok := generatedValue.([]any)
	if !ok {
		switch value := generatedValue.(type) {
		case []interface{}:
			generated = interfaceSliceToAnySlice(value)
		case *data.WorkoutTemplate:
			fitupData, err := data.LoadFitUpData()
			if err != nil {
				slog.Warn("unable to load fitup data for structure persistence", slog.Any("error", err))
				return nil
			}
			generated = s.serializePlanStructure(value, fitupData)
			metadata.Parameters["generated_plan"] = generated
		case data.WorkoutTemplate:
			fitupData, err := data.LoadFitUpData()
			if err != nil {
				slog.Warn("unable to load fitup data for structure persistence", slog.Any("error", err))
				return nil
			}
			tmpl := value
			generated = s.serializePlanStructure(&tmpl, fitupData)
			metadata.Parameters["generated_plan"] = generated
		default:
			return nil
		}
	}

	if len(generated) == 0 {
		return nil
	}

	structure := planStructureInputsFromGeneratedDays(generated)
	if len(structure) == 0 {
		return nil
	}

	return s.repo.PlanGeneration().SaveGeneratedPlanStructure(ctx, planID, structure)
}

func (s *planGenerationServiceImpl) buildPlanStructureInputs(template *data.WorkoutTemplate, fitupData *data.FitUpData) []types.PlanStructureDayInput {
	if template == nil {
		return nil
	}

	var lookup map[int]data.Exercise
	if fitupData != nil {
		lookup = make(map[int]data.Exercise, len(fitupData.Exercises))
		for _, exercise := range fitupData.Exercises {
			lookup[exercise.ID] = exercise
		}
	}

	dayKeys := make([]string, 0, len(template.Structure))
	for key := range template.Structure {
		dayKeys = append(dayKeys, key)
	}

	sort.Slice(dayKeys, func(i, j int) bool {
		return dayKeyOrder(dayKeys[i]) < dayKeyOrder(dayKeys[j])
	})

	structure := make([]types.PlanStructureDayInput, 0, len(dayKeys))
	for idx, key := range dayKeys {
		day := template.Structure[key]

		dayTitle := fmt.Sprintf("Day %d", idx+1)
		if idx < len(template.Schedule) {
			titled := prettifyDayLabel(template.Schedule[idx])
			if titled != "" {
				dayTitle = titled
			}
		} else {
			titled := prettifyDayLabel(key)
			if titled != "" {
				dayTitle = titled
			}
		}

		focus := day.Focus
		if focus != "" {
			focus = prettifyDayLabel(focus)
		}

		exercises := make([]types.PlanStructureExerciseInput, 0, len(day.Exercises))
		for _, spec := range day.Exercises {
			var exerciseID *int
			if spec.ExerciseID > 0 {
				id := spec.ExerciseID
				exerciseID = &id
			}

			name := ""
			if lookup != nil {
				if exercise, ok := lookup[spec.ExerciseID]; ok {
					name = exercise.Name
				}
			}

			exercises = append(exercises, types.PlanStructureExerciseInput{
				ExerciseID:  exerciseID,
				Name:        name,
				Sets:        spec.Sets,
				Reps:        spec.Reps,
				RestSeconds: spec.Rest,
			})
		}

		structure = append(structure, types.PlanStructureDayInput{
			DayIndex:  idx + 1,
			DayTitle:  dayTitle,
			Focus:     focus,
			IsRest:    len(day.Exercises) == 0,
			Exercises: exercises,
		})
	}

	if len(template.Schedule) > len(dayKeys) {
		for idx := len(dayKeys); idx < len(template.Schedule); idx++ {
			dayTitle := prettifyDayLabel(template.Schedule[idx])
			if dayTitle == "" {
				dayTitle = fmt.Sprintf("Day %d", idx+1)
			}

			structure = append(structure, types.PlanStructureDayInput{
				DayIndex:  idx + 1,
				DayTitle:  dayTitle,
				Focus:     "Recovery",
				IsRest:    true,
				Exercises: []types.PlanStructureExerciseInput{},
			})
		}
	}

	return structure
}

func (s *planGenerationServiceImpl) serializePlanStructure(template *data.WorkoutTemplate, fitupData *data.FitUpData) []any {
	if template == nil {
		return []any{}
	}

	structure := s.buildPlanStructureInputs(template, fitupData)
	return generatedDaysFromInputs(structure)
}

func interfaceSliceToAnySlice(src []interface{}) []any {
	if len(src) == 0 {
		return []any{}
	}

	dst := make([]any, len(src))
	copy(dst, src)
	return dst
}

func (s *planGenerationServiceImpl) structureInputsFromMetadata(raw json.RawMessage) ([]types.PlanStructureDayInput, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	var metadata map[string]any
	if err := json.Unmarshal(raw, &metadata); err != nil {
		return nil, err
	}

	parameters, ok := metadata["parameters"].(map[string]any)
	if !ok {
		return nil, nil
	}

	generatedValue, ok := parameters["generated_plan"]
	if !ok {
		return nil, nil
	}

	if days, ok := generatedValue.([]any); ok {
		return planStructureInputsFromGeneratedDays(days), nil
	}

	marshaled, err := json.Marshal(generatedValue)
	if err != nil {
		return nil, err
	}

	var decoded any
	if err := json.Unmarshal(marshaled, &decoded); err != nil {
		return nil, err
	}

	switch value := decoded.(type) {
	case []any:
		return planStructureInputsFromGeneratedDays(value), nil
	case map[string]any:
		if structureMap, ok := value["structure"].(map[string]any); ok {
			days := structureMapToGeneratedDays(structureMap)
			return planStructureInputsFromGeneratedDays(days), nil
		}
		if days, ok := value["days"].([]any); ok {
			return planStructureInputsFromGeneratedDays(days), nil
		}
	}

	return nil, nil
}

func structureMapToGeneratedDays(structure map[string]any) []any {
	if len(structure) == 0 {
		return []any{}
	}

	keys := make([]string, 0, len(structure))
	for key := range structure {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return dayKeyOrder(keys[i]) < dayKeyOrder(keys[j])
	})

	days := make([]any, 0, len(keys))
	for idx, key := range keys {
		dayValue, ok := structure[key].(map[string]any)
		if !ok {
			continue
		}

		rawExercises, ok := dayValue["exercises"].([]any)
		if !ok {
			if generic, ok := dayValue["exercises"].([]interface{}); ok {
				rawExercises = interfaceSliceToAnySlice(generic)
			}
		}

		dayTitle := prettifyDayLabel(key)
		if dayTitle == "" {
			dayTitle = fmt.Sprintf("Day %d", idx+1)
		}

		focus, _ := dayValue["focus"].(string)

		days = append(days, map[string]any{
			"day_index": idx + 1,
			"day_title": dayTitle,
			"focus":     focus,
			"is_rest":   len(rawExercises) == 0,
			"exercises": rawExercises,
		})
	}

	return days
}

func planInputsToGeneratedWorkouts(planID int, inputs []types.PlanStructureDayInput) []types.GeneratedPlanWorkout {
	workouts := make([]types.GeneratedPlanWorkout, 0, len(inputs))
	for _, day := range inputs {
		exercises := make([]types.GeneratedPlanExerciseDetail, 0, len(day.Exercises))
		for idx, exercise := range day.Exercises {
			var exerciseID *int
			if exercise.ExerciseID != nil {
				idCopy := *exercise.ExerciseID
				exerciseID = &idCopy
			}

			exercises = append(exercises, types.GeneratedPlanExerciseDetail{
				PlanExerciseID: 0,
				PlanDayID:      0,
				ExerciseOrder:  idx + 1,
				ExerciseID:     exerciseID,
				Name:           exercise.Name,
				Sets:           exercise.Sets,
				Reps:           exercise.Reps,
				RestSeconds:    exercise.RestSeconds,
				Notes:          exercise.Notes,
			})
		}

		workouts = append(workouts, types.GeneratedPlanWorkout{
			WorkoutID: 0,
			PlanID:    planID,
			DayIndex:  day.DayIndex,
			DayTitle:  day.DayTitle,
			Focus:     day.Focus,
			IsRest:    day.IsRest,
			Exercises: exercises,
		})
	}

	return workouts
}

func planStructureInputsFromGeneratedDays(generated []any) []types.PlanStructureDayInput {
	var structure []types.PlanStructureDayInput

	for idx, dayIface := range generated {
		dayMap, ok := dayIface.(map[string]any)
		if !ok {
			continue
		}

		title, _ := dayMap["day_title"].(string)
		focus, _ := dayMap["focus"].(string)
		rest := false
		if val, ok := dayMap["is_rest"].(bool); ok {
			rest = val
		}

		rawExercises, ok := dayMap["exercises"].([]any)
		if !ok {
			if generic, ok := dayMap["exercises"].([]interface{}); ok {
				rawExercises = interfaceSliceToAnySlice(generic)
			}
		}

		var exercises []types.PlanStructureExerciseInput
		for _, exIface := range rawExercises {
			exMap, ok := exIface.(map[string]any)
			if !ok {
				continue
			}

			var exerciseID *int
			switch idVal := exMap["exercise_id"].(type) {
			case float64:
				parsed := int(idVal)
				exerciseID = &parsed
			case int:
				parsed := idVal
				exerciseID = &parsed
			case int32:
				parsed := int(idVal)
				exerciseID = &parsed
			case int64:
				parsed := int(idVal)
				exerciseID = &parsed
			}

			name, _ := exMap["name"].(string)
			reps, _ := exMap["reps"].(string)
			notes, _ := exMap["notes"].(string)

			exercises = append(exercises, types.PlanStructureExerciseInput{
				ExerciseID:  exerciseID,
				Name:        name,
				Sets:        intFromAny(exMap["sets"], 0),
				Reps:        reps,
				RestSeconds: intFromAny(exMap["rest"], 0),
				Notes:       notes,
			})
		}

		dayIndex := idx + 1
		if override := intFromAny(dayMap["day_index"], 0); override > 0 {
			dayIndex = override
		}

		structure = append(structure, types.PlanStructureDayInput{
			DayIndex:  dayIndex,
			DayTitle:  title,
			Focus:     focus,
			IsRest:    rest,
			Exercises: exercises,
		})
	}

	return structure
}

func generatedDaysFromInputs(inputs []types.PlanStructureDayInput) []any {
	generated := make([]any, 0, len(inputs))
	for _, day := range inputs {
		exercises := make([]any, 0, len(day.Exercises))
		for _, exercise := range day.Exercises {
			exerciseMap := map[string]any{
				"sets": exercise.Sets,
				"reps": exercise.Reps,
				"rest": exercise.RestSeconds,
				"name": exercise.Name,
			}
			if exercise.ExerciseID != nil {
				exerciseMap["exercise_id"] = *exercise.ExerciseID
			}
			if exercise.Notes != "" {
				exerciseMap["notes"] = exercise.Notes
			}
			exercises = append(exercises, exerciseMap)
		}

		generated = append(generated, map[string]any{
			"day_index": day.DayIndex,
			"day_title": day.DayTitle,
			"focus":     day.Focus,
			"is_rest":   day.IsRest,
			"exercises": exercises,
		})
	}

	return generated
}

func dayKeyOrder(key string) int {
	digits := 0
	for _, r := range key {
		if r >= '0' && r <= '9' {
			digits = digits*10 + int(r-'0')
		}
	}
	if digits == 0 {
		return 1 << 30
	}
	return digits
}

func prettifyDayLabel(input string) string {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return ""
	}

	parts := strings.FieldsFunc(trimmed, func(r rune) bool {
		switch r {
		case '_', '-', ' ':
			return true
		default:
			return false
		}
	})

	for idx, part := range parts {
		lower := strings.ToLower(part)
		if len(lower) == 0 {
			continue
		}
		parts[idx] = strings.ToUpper(lower[:1]) + lower[1:]
	}

	return strings.Join(parts, " ")
}

func intFromAny(value any, fallback int) int {
	switch v := value.(type) {
	case float64:
		return int(v)
	case float32:
		return int(v)
	case int:
		return v
	case int64:
		return int(v)
	case json.Number:
		if parsed, err := v.Int64(); err == nil {
			return int(parsed)
		}
	}
	return fallback
}

func stringListFromAny(value any) []string {
	switch v := value.(type) {
	case []string:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if trimmed := strings.TrimSpace(item); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	case []interface{}:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if str, ok := item.(string); ok {
				if trimmed := strings.TrimSpace(str); trimmed != "" {
					result = append(result, trimmed)
				}
			}
		}
		return result
	case string:
		trimmed := strings.TrimSpace(v)
		if trimmed == "" {
			return nil
		}
		parts := strings.Split(trimmed, ",")
		if len(parts) == 1 {
			return []string{trimmed}
		}
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			piece := strings.TrimSpace(part)
			if piece != "" {
				result = append(result, piece)
			}
		}
		return result
	default:
		return nil
	}
}

func extractLeadingInt(input string) int {
	var builder strings.Builder
	for _, r := range input {
		if r >= '0' && r <= '9' {
			builder.WriteRune(r)
		} else if builder.Len() > 0 {
			break
		}
	}
	if builder.Len() == 0 {
		return 0
	}
	if value, err := strconv.Atoi(builder.String()); err == nil {
		return value
	}
	return 0
}

func parseRepsToInt(reps string) int {
	trimmed := strings.TrimSpace(strings.ToLower(reps))
	if trimmed == "" {
		return 10
	}

	if strings.Contains(trimmed, "-") {
		parts := strings.Split(trimmed, "-")
		total := 0
		count := 0
		for _, part := range parts {
			if value := extractLeadingInt(part); value > 0 {
				total += value
				count++
			}
		}
		if count > 0 {
			return total / count
		}
	}

	if value := extractLeadingInt(trimmed); value > 0 {
		return value
	}

	if strings.Contains(trimmed, "x") {
		segments := strings.Split(trimmed, "x")
		if len(segments) > 0 {
			if value := extractLeadingInt(segments[0]); value > 0 {
				return value
			}
		}
	}

	switch trimmed {
	case "amrap", "to failure":
		return 12
	}

	return 10
}

func estimateWorkoutDuration(workout types.GeneratedPlanWorkout) (int, int) {
	totalSeconds := 0
	totalSets := 0

	for _, exercise := range workout.Exercises {
		sets := exercise.Sets
		if sets <= 0 {
			continue
		}

		reps := parseRepsToInt(exercise.Reps)
		rest := exercise.RestSeconds
		totalSets += sets
		totalSeconds += sets * reps * 3
		if sets > 1 && rest > 0 {
			totalSeconds += (sets - 1) * rest
		}
	}

	return totalSeconds, totalSets
}

func (s *planGenerationServiceImpl) generateAdaptivePlan(_ context.Context, userID int, metadata *types.PlanGenerationMetadata) (*types.PlanGenerationMetadata, error) {
	if len(metadata.UserGoals) == 0 {
		slog.Warn("plan generation missing goals", slog.Int("user_id", userID))
		return nil, fmt.Errorf("at least one fitness goal is required")
	}
	if len(metadata.AvailableEquipment) == 0 {
		slog.Warn("plan generation missing equipment", slog.Int("user_id", userID))
		return nil, fmt.Errorf("at least one equipment type is required")
	}
	if metadata.WeeklyFrequency <= 0 || metadata.WeeklyFrequency > 7 {
		slog.Warn("plan generation invalid frequency", slog.Int("user_id", userID), slog.Int("weekly_frequency", metadata.WeeklyFrequency))
		return nil, fmt.Errorf("weekly frequency must be between 1 and 7")
	}

	fitupData, err := data.LoadFitUpData()
	if err != nil {
		slog.Error("failed to load fitup data", slog.Int("user_id", userID), slog.Any("error", err))
		return nil, fmt.Errorf("failed to load fitness data: %w", err)
	}

	userLevel := string(metadata.FitnessLevel)
	primaryGoal := metadata.UserGoals[0]

	levelData, exists := fitupData.Levels[userLevel]
	if !exists {
		slog.Warn("invalid fitness level", slog.Int("user_id", userID), slog.String("fitness_level", userLevel))
		return nil, fmt.Errorf("invalid fitness level: %s", userLevel)
	}

	goalData, exists := fitupData.Goals[string(primaryGoal)]
	if !exists {
		slog.Warn("invalid fitness goal", slog.Int("user_id", userID), slog.String("goal", string(primaryGoal)))
		return nil, fmt.Errorf("invalid fitness goal: %s", primaryGoal)
	}

	template, err := s.selectOptimalTemplate(fitupData, userLevel, string(primaryGoal), metadata.WeeklyFrequency)
	if err != nil {
		slog.Error("failed to select template", slog.Int("user_id", userID), slog.Any("error", err))
		return nil, fmt.Errorf("failed to select workout template: %w", err)
	}

	exerciseSelection, err := s.generateExerciseSelection(fitupData, template, metadata.AvailableEquipment, userLevel, string(primaryGoal))
	if err != nil {
		slog.Error("failed to select exercises", slog.Int("user_id", userID), slog.Any("error", err))
		return nil, fmt.Errorf("failed to generate exercise selection: %w", err)
	}

	adaptedTemplate := s.applyProgressiveOverload(template, exerciseSelection, levelData, goalData, metadata.TimePerWorkout)

	balancedPlan := s.optimizeMuscleGroupBalance(adaptedTemplate, exerciseSelection)
	generatedPlan := s.serializePlanStructure(balancedPlan, fitupData)

	weekStart := startOfWeek(time.Now().UTC())
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
			"generated_plan":         generatedPlan,
			"week_start":             weekStart.Format("2006-01-02"),
		},
	}

	slog.Info("adaptive plan generated", slog.Int("user_id", userID), slog.String("template_id", template.ID), slog.Int("exercise_count", len(exerciseSelection)))

	return enhancedMetadata, nil
}

func (s *planGenerationServiceImpl) resolveUserIdentity(ctx context.Context, schemaUserID int, metadata *types.PlanGenerationMetadata, createIfMissing bool) (int, string, error) {
	authUserID, hasAuth := middleware.GetAuthUserIDFromContext(ctx)

	if schemaUserID > 0 {
		profile, err := s.repo.WorkoutProfiles().GetWorkoutProfileByID(ctx, schemaUserID)
		if err == nil {
			if hasAuth && authUserID != "" && profile.AuthUserID != authUserID {
				return 0, "", types.ErrInvalidUserID
			}
			return profile.WorkoutProfileID, profile.AuthUserID, nil
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return 0, "", fmt.Errorf("failed to get workout profile: %w", err)
		}
	}

	if !hasAuth || authUserID == "" {
		return 0, "", types.ErrInvalidUserID
	}

	profile, err := s.repo.WorkoutProfiles().GetWorkoutProfileByAuthID(ctx, authUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			if !createIfMissing {
				return 0, "", types.ErrInvalidUserID
			}

			req := s.buildWorkoutProfileRequest(metadata)
			if req == nil {
				return 0, "", types.ErrInvalidUserID
			}

			profile, err = s.repo.WorkoutProfiles().CreateWorkoutProfile(ctx, authUserID, req)
			if err != nil {
				return 0, "", fmt.Errorf("failed to create workout profile: %w", err)
			}

			slog.Info("created workout profile for user", slog.String("auth_user_id", authUserID), slog.Int("user_id", profile.WorkoutProfileID))
		} else {
			return 0, "", fmt.Errorf("failed to lookup workout profile: %w", err)
		}
	}

	return profile.WorkoutProfileID, profile.AuthUserID, nil
}

func (s *planGenerationServiceImpl) buildWorkoutProfileRequest(metadata *types.PlanGenerationMetadata) *types.WorkoutProfileRequest {
	if metadata == nil {
		return nil
	}

	goal := types.GoalGeneralFitness
	if len(metadata.UserGoals) > 0 {
		goal = metadata.UserGoals[0]
	}

	equipment := make([]string, 0, len(metadata.AvailableEquipment))
	for _, eq := range metadata.AvailableEquipment {
		equipment = append(equipment, string(eq))
	}

	if metadata.WeeklyFrequency <= 0 {
		return nil
	}

	return &types.WorkoutProfileRequest{
		Level:     metadata.FitnessLevel,
		Goal:      goal,
		Frequency: metadata.WeeklyFrequency,
		Equipment: equipment,
	}
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

func startOfWeek(t time.Time) time.Time {
	loc := t.Location()
	base := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
	daysSinceMonday := (int(base.Weekday()) + 6) % 7
	return base.AddDate(0, 0, -daysSinceMonday)
}

func (s *planGenerationServiceImpl) GetActivePlanForUser(ctx context.Context, userID int) (*types.GeneratedPlan, error) {
	resolvedID, _, err := s.resolveUserIdentity(ctx, userID, nil, false)
	if err != nil {
		return nil, err
	}

	activePlan, err := s.repo.PlanGeneration().GetActivePlanForUser(ctx, resolvedID)
	if err != nil {
		return nil, err
	}
	if activePlan == nil {
		return nil, fmt.Errorf("no active plan found for user %d", resolvedID)
	}

	s.populatePlanWorkouts(ctx, activePlan)

	if err := s.enrichPlanWithProgress(ctx, activePlan); err != nil {
		return nil, fmt.Errorf("failed to enrich plan with progress: %w", err)
	}
	return activePlan, nil
}

func (s *planGenerationServiceImpl) enrichPlanWithProgress(ctx context.Context, plan *types.GeneratedPlan) error {
	logs, err := s.repo.Progress().GetProgressLogsByUserID(ctx, plan.UserID, types.PaginationParams{Limit: 100, Offset: 0})
	if err != nil {
		slog.Warn("unable to load progress logs for plan", slog.Int("plan_id", plan.PlanID), slog.Int("profile_id", plan.UserID), slog.Any("error", err))
		return nil
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

func mapStructureToWorkouts(days []types.GeneratedPlanDay) []types.GeneratedPlanWorkout {
	workouts := make([]types.GeneratedPlanWorkout, 0, len(days))
	for _, day := range days {
		var exercises []types.GeneratedPlanExerciseDetail
		for _, ex := range day.Exercises {
			exDetail := types.GeneratedPlanExerciseDetail{
				PlanExerciseID: ex.PlanExerciseID,
				PlanDayID:      ex.PlanDayID,
				ExerciseOrder:  ex.ExerciseOrder,
				Name:           ex.Name,
				Sets:           ex.Sets,
				Reps:           ex.Reps,
				RestSeconds:    ex.RestSeconds,
				Notes:          ex.Notes,
			}
			if ex.ExerciseID != nil {
				copyID := *ex.ExerciseID
				exDetail.ExerciseID = &copyID
			}
			exercises = append(exercises, exDetail)
		}

		workouts = append(workouts, types.GeneratedPlanWorkout{
			WorkoutID: day.PlanDayID,
			PlanID:    day.PlanID,
			DayIndex:  day.DayIndex,
			DayTitle:  day.DayTitle,
			Focus:     day.Focus,
			IsRest:    day.IsRest,
			Exercises: exercises,
		})
	}
	return workouts
}

func (s *planGenerationServiceImpl) populatePlanWorkouts(ctx context.Context, plan *types.GeneratedPlan) {
	if plan == nil {
		return
	}

	structure, err := s.repo.PlanGeneration().GetGeneratedPlanStructure(ctx, plan.PlanID)
	if err == nil && len(structure) > 0 {
		plan.Workouts = mapStructureToWorkouts(structure)
		return
	}

	inputs, extractErr := s.structureInputsFromMetadata(plan.Metadata)
	if extractErr != nil || len(inputs) == 0 {
		if extractErr != nil {
			slog.Warn("failed to extract plan structure from metadata", slog.Int("plan_id", plan.PlanID), slog.Any("error", extractErr))
		}
		return
	}

	if saveErr := s.repo.PlanGeneration().SaveGeneratedPlanStructure(ctx, plan.PlanID, inputs); saveErr == nil {
		if rebuilt, err := s.repo.PlanGeneration().GetGeneratedPlanStructure(ctx, plan.PlanID); err == nil && len(rebuilt) > 0 {
			plan.Workouts = mapStructureToWorkouts(rebuilt)
			return
		}
	} else {
		slog.Warn("failed to persist reconstructed plan structure", slog.Int("plan_id", plan.PlanID), slog.Any("error", saveErr))
	}

	plan.Workouts = planInputsToGeneratedWorkouts(plan.PlanID, inputs)
}

func (s *planGenerationServiceImpl) GetPlanGenerationHistory(ctx context.Context, userID int, limit int) ([]types.GeneratedPlan, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	resolvedID, _, err := s.resolveUserIdentity(ctx, userID, nil, false)
	if err != nil {
		if errors.Is(err, types.ErrInvalidUserID) {
			return []types.GeneratedPlan{}, nil
		}
		return nil, err
	}

	planHistory, err := s.repo.PlanGeneration().GetPlanGenerationHistory(ctx, resolvedID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan generation history: %w", err)
	}

	for i := range planHistory {
		s.populatePlanWorkouts(ctx, &planHistory[i])
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
			fmt.Printf("Plan evolution insights for user %d: %+v\n", resolvedID, evolution)
		}
	}

	return planHistory, nil
}

func (s *planGenerationServiceImpl) DeletePlan(ctx context.Context, userID int, planID int) error {
	if planID <= 0 {
		slog.Warn("invalid plan id for deletion", slog.Int("plan_id", planID), slog.Int("requested_user_id", userID))
		return types.ErrPlanNotFound
	}

	resolvedID, authUserID, err := s.resolveUserIdentity(ctx, userID, nil, false)
	if err != nil {
		return err
	}

	plan, err := s.repo.PlanGeneration().GetPlanID(ctx, planID)
	if err != nil {
		if errors.Is(err, types.ErrPlanNotFound) {
			return err
		}
		return fmt.Errorf("failed to load plan for deletion: %w", err)
	}

	if plan == nil {
		return types.ErrPlanNotFound
	}

	if plan.UserID != resolvedID {
		slog.Warn("plan deletion denied", slog.Int("plan_id", planID), slog.Int("plan_owner_id", plan.UserID), slog.Int("requester_id", resolvedID))
		return types.ErrPlanDeleteDenied
	}

	if err := s.repo.PlanGeneration().DeletePlanForUser(ctx, planID, authUserID); err != nil {
		if errors.Is(err, types.ErrPlanNotFound) {
			return err
		}
		return fmt.Errorf("failed to delete plan: %w", err)
	}

	slog.Info("plan deleted", slog.Int("plan_id", planID), slog.Int("user_id", resolvedID))
	return nil
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

	normalizedPerformance := normalizePlanPerformanceData(performance)

	if err := validator.New().Struct(normalizedPerformance); err != nil {
		return fmt.Errorf("invalid performance data: %w", err)
	}

	if err := s.repo.PlanGeneration().TrackPlanPerformance(ctx, planID, normalizedPerformance); err != nil {
		return fmt.Errorf("failed to track plan performance: %w", err)
	}

	if err := s.analyzeAndAdaptPlan(ctx, planID, normalizedPerformance); err != nil {
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
	resolvedID, _, err := s.resolveUserIdentity(ctx, userID, nil, false)
	if err != nil {
		if errors.Is(err, types.ErrInvalidUserID) {
			return []types.PlanAdaptation{}, nil
		}
		return nil, err
	}

	adaptations, err := s.repo.PlanGeneration().GetAdaptationHistory(ctx, resolvedID)
	if err != nil {
		return nil, fmt.Errorf("failed to get adaptation history: %w", err)
	}

	if len(adaptations) > 0 {
		insights := s.analyzeAdaptationPatterns(adaptations)
		fmt.Printf("Adaptation insights for user %d: %+v\n", resolvedID, insights)
	}

	return adaptations, nil
}

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

func normalizePlanPerformanceData(performance *types.PlanPerformanceData) *types.PlanPerformanceData {
	if performance == nil {
		return nil
	}

	normalized := *performance
	normalized.CompletionRate = clampFloat(normalized.CompletionRate, 0, 1)
	normalized.ProgressRate = clampFloat(normalized.ProgressRate, 0, 1)
	normalized.InjuryRate = clampFloat(normalized.InjuryRate, 0, 1)
	normalized.AverageRPE = clampFloat(normalized.AverageRPE, 1, 10)
	if normalized.UserSatisfaction > 1 {
		normalized.UserSatisfaction = normalized.UserSatisfaction / 10
	}
	normalized.UserSatisfaction = clampFloat(normalized.UserSatisfaction, 0, 1)

	return &normalized
}

func clampFloat(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
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

	s.populatePlanWorkouts(ctx, plan)

	metadata := map[string]any{}
	if len(plan.Metadata) > 0 {
		if err := json.Unmarshal(plan.Metadata, &metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal plan metadata: %w", err)
		}
	}

	parameters := map[string]any{}
	if rawParams, ok := metadata["parameters"]; ok {
		if typed, ok := rawParams.(map[string]any); ok {
			parameters = typed
		}
	}

	workouts := plan.Workouts
	if len(workouts) == 0 {
		if rawGenerated, ok := parameters["generated_plan"]; ok {
			if generatedDays, ok := rawGenerated.([]any); ok {
				structure := planStructureInputsFromGeneratedDays(generatedDays)
				workouts = planInputsToGeneratedWorkouts(plan.PlanID, structure)
			} else if genericDays, ok := rawGenerated.([]interface{}); ok {
				structure := planStructureInputsFromGeneratedDays(interfaceSliceToAnySlice(genericDays))
				workouts = planInputsToGeneratedWorkouts(plan.PlanID, structure)
			}
		}
	}

	if len(workouts) == 0 {
		if inputs, err := s.structureInputsFromMetadata(plan.Metadata); err == nil && len(inputs) > 0 {
			workouts = planInputsToGeneratedWorkouts(plan.PlanID, inputs)
		}
	}

	if len(workouts) == 0 {
		return nil, fmt.Errorf("plan contains no workout data to export")
	}

	return s.renderPlanPDF(plan, workouts, parameters)
}

func (s *planGenerationServiceImpl) renderPlanPDF(plan *types.GeneratedPlan, workouts []types.GeneratedPlanWorkout, parameters map[string]any) ([]byte, error) {
	const (
		brandPrimaryR     = 143
		brandPrimaryG     = 229
		brandPrimaryB     = 7
		brandPrimaryDarkR = 106
		brandPrimaryDarkG = 176
		brandPrimaryDarkB = 0
		brandCanvasR      = 249
		brandCanvasG      = 250
		brandCanvasB      = 251
		brandTextDarkR    = 28
		brandTextDarkG    = 28
		brandTextDarkB    = 30
		brandTextMutedR   = 107
		brandTextMutedG   = 114
		brandTextMutedB   = 128
		brandCanvasDarkR  = 10
		brandCanvasDarkG  = 10
		brandCanvasDarkB  = 10
	)

	leftMargin := 18.0
	rightMargin := 210.0 - leftMargin
	contentWidth := rightMargin - leftMargin

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(leftMargin, 24, leftMargin)
	pdf.SetAutoPageBreak(true, 20)
	pdf.AddPage()

	// Header
	pdf.SetFillColor(brandCanvasDarkR, brandCanvasDarkG, brandCanvasDarkB)
	pdf.Rect(0, 0, 210, 54, "F")
	pdf.SetXY(leftMargin, 18)
	pdf.SetFont("Arial", "B", 26)
	pdf.SetTextColor(brandPrimaryR, brandPrimaryG, brandPrimaryB)
	pdf.Cell(0, 12, "FitUp Training Plan")
	pdf.Ln(11)
	pdf.SetFont("Arial", "", 11)
	pdf.SetTextColor(255, 255, 255)
	pdf.MultiCell(0, 5, "Personalized coaching blueprint generated by the FitUp adaptive engine.", "", "L", false)
	pdf.Ln(2)
	pdf.SetY(58)

	drawSectionHeader := func(title string) {
		pdf.SetX(leftMargin)
		pdf.SetFont("Arial", "B", 12)
		pdf.SetTextColor(brandTextDarkR, brandTextDarkG, brandTextDarkB)
		pdf.Cell(0, 8, strings.ToUpper(title))
		pdf.Ln(6)
		y := pdf.GetY()
		pdf.SetDrawColor(brandPrimaryDarkR, brandPrimaryDarkG, brandPrimaryDarkB)
		pdf.SetLineWidth(0.6)
		pdf.Line(leftMargin, y, rightMargin, y)
		pdf.Ln(5)
		pdf.SetLineWidth(0.2)
		pdf.SetTextColor(brandTextDarkR, brandTextDarkG, brandTextDarkB)
	}

	drawKeyValue := func(label, value string) {
		if strings.TrimSpace(value) == "" {
			return
		}
		pdf.SetX(leftMargin)
		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(brandTextMutedR, brandTextMutedG, brandTextMutedB)
		pdf.Cell(42, 6, label)
		pdf.SetFont("Arial", "B", 10)
		pdf.SetTextColor(brandTextDarkR, brandTextDarkG, brandTextDarkB)
		pdf.Cell(0, 6, value)
		pdf.Ln(6)
	}

	// Plan overview
	drawSectionHeader("Plan Overview")
	drawKeyValue("Plan ID", fmt.Sprintf("#%d", plan.PlanID))
	drawKeyValue("Generated", plan.GeneratedAt.Format("January 2, 2006 at 3:04 PM"))
	drawKeyValue("Week Start", plan.WeekStart.Format("Monday, January 2, 2006"))
	drawKeyValue("Algorithm", plan.Algorithm)
	status := "Active"
	if !plan.IsActive {
		status = "Inactive"
	}
	drawKeyValue("Status", status)
	if plan.Effectiveness > 0 {
		drawKeyValue("Effectiveness Score", fmt.Sprintf("%.0f%%", plan.Effectiveness))
	}
	pdf.Ln(2)

	// Training parameters
	if parameters != nil && len(parameters) > 0 {
		templateUsed, _ := parameters["template_used"].(string)
		totalExercisesParam := ""
		if totalEx, ok := parameters["total_exercises"].(float64); ok && totalEx > 0 {
			totalExercisesParam = fmt.Sprintf("%.0f", totalEx)
		}
		targetedMuscles := stringListFromAny(parameters["muscle_groups_targeted"])
		equipment := stringListFromAny(parameters["equipment_utilized"])
		weeklyFrequency := ""
		if freq := intFromAny(parameters["weekly_frequency"], 0); freq > 0 {
			weeklyFrequency = fmt.Sprintf("%d sessions", freq)
		}

		if templateUsed != "" || totalExercisesParam != "" || len(targetedMuscles) > 0 || len(equipment) > 0 || weeklyFrequency != "" {
			drawSectionHeader("Training Parameters")
			if templateUsed != "" {
				drawKeyValue("Template", templateUsed)
			}
			if weeklyFrequency != "" {
				drawKeyValue("Weekly Frequency", weeklyFrequency)
			}
			if totalExercisesParam != "" {
				drawKeyValue("Target Volume", totalExercisesParam+" exercises")
			}
			if len(targetedMuscles) > 0 {
				drawKeyValue("Muscle Groups", strings.Join(targetedMuscles, ", "))
			}
			if len(equipment) > 0 {
				drawKeyValue("Equipment", strings.Join(equipment, ", "))
			}
			if intensity, ok := parameters["intensity_focus"].(string); ok && strings.TrimSpace(intensity) != "" {
				drawKeyValue("Intensity Focus", intensity)
			}
			pdf.Ln(2)
		}
	}

	// Weekly schedule
	drawSectionHeader("Weekly Schedule")

	totalWorkoutTime := 0
	totalExercises := 0
	totalSets := 0
	workoutDays := 0
	focusTags := map[string]struct{}{}

	for idx, workout := range workouts {
		dayTitle := strings.TrimSpace(workout.DayTitle)
		if dayTitle == "" {
			dayTitle = fmt.Sprintf("Day %d", idx+1)
		}

		focusLabel := ""
		if focus := strings.TrimSpace(workout.Focus); focus != "" {
			focusLabel = prettifyDayLabel(focus)
			if strings.ToLower(focus) != "recovery" {
				focusTags[focusLabel] = struct{}{}
			}
		}

		isRestDay := workout.IsRest || len(workout.Exercises) == 0
		if !isRestDay {
			workoutDays++
		}

		pdf.SetX(leftMargin)
		pdf.SetFillColor(brandPrimaryR, brandPrimaryG, brandPrimaryB)
		pdf.SetTextColor(brandCanvasDarkR, brandCanvasDarkG, brandCanvasDarkB)
		pdf.SetFont("Arial", "B", 13)
		pdf.CellFormat(0, 9, fmt.Sprintf("DAY %d · %s", idx+1, strings.ToUpper(dayTitle)), "", 0, "L", true, 0, "")
		pdf.Ln(8)

		if focusLabel != "" {
			pdf.SetX(leftMargin)
			pdf.SetFont("Arial", "I", 10)
			pdf.SetFillColor(232, 255, 209)
			pdf.SetTextColor(brandPrimaryDarkR, brandPrimaryDarkG, brandPrimaryDarkB)
			pdf.CellFormat(0, 6, fmt.Sprintf("Focus: %s", focusLabel), "", 0, "L", true, 0, "")
			pdf.Ln(7)
		}

		if isRestDay {
			pdf.SetX(leftMargin)
			pdf.SetFont("Arial", "I", 10)
			pdf.SetTextColor(brandTextMutedR, brandTextMutedG, brandTextMutedB)
			pdf.Cell(0, 8, "Recovery emphasis • mobility, stretching, and active rest")
			pdf.Ln(12)
			continue
		}

		pdf.SetX(leftMargin)
		pdf.SetFillColor(232, 255, 209)
		pdf.SetTextColor(brandCanvasDarkR, brandCanvasDarkG, brandCanvasDarkB)
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(10, 8, "#", "1", 0, "C", true, 0, "")
		pdf.CellFormat(70, 8, "Exercise", "1", 0, "L", true, 0, "")
		pdf.CellFormat(25, 8, "Sets", "1", 0, "C", true, 0, "")
		pdf.CellFormat(30, 8, "Reps", "1", 0, "C", true, 0, "")
		pdf.CellFormat(35, 8, "Rest", "1", 0, "C", true, 0, "")
		pdf.Ln(8)

		pdf.SetFont("Arial", "", 9)
		dayDuration, daySets := estimateWorkoutDuration(workout)
		totalWorkoutTime += dayDuration
		totalSets += daySets
		totalExercises += len(workout.Exercises)

		for exIdx, exercise := range workout.Exercises {
			if exIdx%2 == 0 {
				pdf.SetFillColor(248, 253, 237)
			} else {
				pdf.SetFillColor(255, 255, 255)
			}

			restDisplay := "-"
			if exercise.RestSeconds > 0 {
				restDisplay = fmt.Sprintf("%d sec", exercise.RestSeconds)
			}

			repDisplay := strings.TrimSpace(exercise.Reps)
			if repDisplay == "" {
				repDisplay = "-"
			}

			pdf.SetX(leftMargin)
			pdf.CellFormat(10, 7, fmt.Sprintf("%d", exIdx+1), "1", 0, "C", true, 0, "")
			pdf.CellFormat(70, 7, exercise.Name, "1", 0, "L", true, 0, "")
			pdf.CellFormat(25, 7, fmt.Sprintf("%d", exercise.Sets), "1", 0, "C", true, 0, "")
			pdf.CellFormat(30, 7, repDisplay, "1", 0, "C", true, 0, "")
			pdf.CellFormat(35, 7, restDisplay, "1", 0, "C", true, 0, "")
			pdf.Ln(7)

			if strings.TrimSpace(exercise.Notes) != "" {
				pdf.SetFont("Arial", "I", 8)
				pdf.SetTextColor(brandTextMutedR, brandTextMutedG, brandTextMutedB)
				pdf.SetX(leftMargin + 10)
				pdf.MultiCell(contentWidth-10, 5, fmt.Sprintf("Coach note: %s", exercise.Notes), "", "L", false)
				pdf.SetFont("Arial", "", 9)
				pdf.SetTextColor(brandTextDarkR, brandTextDarkG, brandTextDarkB)
			}
		}

		pdf.SetFont("Arial", "I", 9)
		pdf.SetTextColor(brandTextMutedR, brandTextMutedG, brandTextMutedB)
		pdf.SetX(leftMargin)
		pdf.Cell(0, 6, fmt.Sprintf("Estimated duration: %d minutes • Exercises: %d", dayDuration/60, len(workout.Exercises)))
		pdf.Ln(10)
	}

	// Training guidelines page
	pdf.AddPage()
	drawSectionHeader("Training Guidelines")

	guidelines := []struct {
		title string
		text  string
	}{
		{"Prime & Mobilize", "Begin every session with 5-8 minutes of dynamic movement and joint prep to unlock full range of motion."},
		{"Own The Form", "Controlled reps with full range beat heavier loads with compromised technique. Film difficult lifts once a week."},
		{"Dial Rest Windows", "Short rest (45-75s) drives conditioning, longer rest (90-150s) supports strength progression."},
		{"Progressive Stress", "Track weights, reps, or rest. Nudge one variable upwards weekly for steady progress."},
		{"Recovery Rituals", "Sleep 7-9 hours, hydrate, and embrace mobility or low-intensity cardio on off days."},
		{"Fuel The Work", "Center meals around lean protein, smart carbs, and electrolytes to stay ready for the next session."},
		{"Listen To Signals", "Differentiate soreness from pain. Deload or modify when joints protest or fatigue lingers."},
	}

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(brandTextDarkR, brandTextDarkG, brandTextDarkB)

	for _, guide := range guidelines {
		pdf.SetX(leftMargin)
		pdf.SetFont("Arial", "B", 10)
		pdf.SetTextColor(brandPrimaryDarkR, brandPrimaryDarkG, brandPrimaryDarkB)
		pdf.Write(5, "• "+guide.title+": ")
		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(brandTextDarkR, brandTextDarkG, brandTextDarkB)
		pdf.MultiCell(0, 5, guide.text, "", "L", false)
		pdf.Ln(2)
	}

	pdf.Ln(4)
	drawSectionHeader("Exercise Swaps")
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(brandTextDarkR, brandTextDarkG, brandTextDarkB)
	pdf.SetX(leftMargin)
	pdf.MultiCell(0, 5, "If equipment is limited or a movement bothers your joints, rotate to one of the FitUp-approved substitutions:", "", "L", false)
	pdf.Ln(2)

	modifications := []string{
		"Push-up → Incline push-up • Dumbbell floor press",
		"Pull-up → Assisted pull-up • Band lat pull-down",
		"Barbell squat → Goblet squat • Split squat",
		"Deadlift → Romanian deadlift • Single-leg hinge",
		"Bench press → Dumbbell press • Band chest press",
	}

	for _, mod := range modifications {
		pdf.SetFont("Arial", "", 9)
		pdf.SetTextColor(brandTextMutedR, brandTextMutedG, brandTextMutedB)
		pdf.SetX(leftMargin + 5)
		pdf.Cell(0, 5, mod)
		pdf.Ln(5)
	}

	pdf.Ln(6)
	drawSectionHeader("Weekly Summary")
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(brandTextDarkR, brandTextDarkG, brandTextDarkB)

	focusList := make([]string, 0, len(focusTags))
	for focus := range focusTags {
		focusList = append(focusList, focus)
	}
	sort.Strings(focusList)

	pdf.SetX(leftMargin)
	pdf.Cell(0, 6, fmt.Sprintf("Workout Days: %d of %d", workoutDays, len(workouts)))
	pdf.Ln(6)
	pdf.SetX(leftMargin)
	pdf.Cell(0, 6, fmt.Sprintf("Total Exercises: %d", totalExercises))
	pdf.Ln(6)
	pdf.SetX(leftMargin)
	pdf.Cell(0, 6, fmt.Sprintf("Total Sets Logged: %d", totalSets))
	pdf.Ln(6)
	if totalWorkoutTime > 0 {
		pdf.SetX(leftMargin)
		pdf.Cell(0, 6, fmt.Sprintf("Estimated Weekly Training Time: %d minutes (~%.1f hours)", totalWorkoutTime/60, float64(totalWorkoutTime)/3600))
		pdf.Ln(6)
	}
	if len(focusList) > 0 {
		pdf.SetX(leftMargin)
		pdf.Cell(0, 6, fmt.Sprintf("Primary Focus Areas: %s", strings.Join(focusList, ", ")))
		pdf.Ln(6)
	}

	// Footer
	pdf.SetY(-28)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(brandTextMutedR, brandTextMutedG, brandTextMutedB)
	pdf.CellFormat(0, 4, "FitUp • Adaptive Training Intelligence", "", 0, "C", false, 0, "")
	pdf.Ln(4)
	pdf.CellFormat(0, 4, fmt.Sprintf("Plan #%d • Generated %s", plan.PlanID, plan.GeneratedAt.Format("2006-01-02")), "", 0, "C", false, 0, "")
	pdf.Ln(4)
	pdf.CellFormat(0, 4, "Need support? Reach us at support@fitup.com", "", 0, "C", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}
