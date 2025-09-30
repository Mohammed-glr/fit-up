package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"

	"github.com/tdmdh/fit-up-server/services/schema-service/internal/types"
)

// FitUpData represents the complete fitness data structure
type FitUpData struct {
	Meta                  Meta                            `json:"meta"`
	Levels                map[string]Level                `json:"levels"`
	Goals                 map[string]Goal                 `json:"goals"`
	EquipmentTypes        map[string]Equipment            `json:"equipment_types"`
	FocusAreas            map[string]FocusArea            `json:"focus_areas"`
	ExerciseTypes         map[string]ExerciseType         `json:"exercise_types"`
	Exercises             []Exercise                      `json:"exercises"`
	WorkoutTemplates      map[string]WorkoutTemplate      `json:"workout_templates"`
	WeeklySchemaExample   WeeklySchemaExample             `json:"weekly_schema_example"`
	ProgressionAlgorithms map[string]ProgressionAlgorithm `json:"progression_algorithms"`
	AdaptationTriggers    map[string]AdaptationTrigger    `json:"adaptation_triggers"`
}

// Meta contains metadata about the data structure
type Meta struct {
	Version     string `json:"version"`
	GeneratedAt string `json:"generated_at"`
	Description string `json:"description"`
	LastUpdated string `json:"last_updated"`
}

// Level represents fitness experience levels
type Level struct {
	ID                  string              `json:"id"`
	Name                string              `json:"name"`
	Description         string              `json:"description"`
	ExperienceMonths    string              `json:"experience_months"`
	WeeklyVolume        WeeklyVolumeGuide   `json:"weekly_volume"`
	IntensityGuidelines IntensityGuidelines `json:"intensity_guidelines"`
}

type WeeklyVolumeGuide struct {
	MinDays                int    `json:"min_days"`
	MaxDays                int    `json:"max_days"`
	SessionDurationMinutes string `json:"session_duration_minutes"`
	TotalWeeklySets        string `json:"total_weekly_sets"`
}

type IntensityGuidelines struct {
	StrengthRPE     string `json:"strength_rpe"`
	CardioIntensity string `json:"cardio_intensity"`
	ProgressionRate string `json:"progression_rate"`
}

// Goal represents fitness goals
type Goal struct {
	ID                 string         `json:"id"`
	Name               string         `json:"name"`
	Description        string         `json:"description"`
	PrimaryAdaptations []string       `json:"primary_adaptations"`
	RepRanges          RepRanges      `json:"rep_ranges"`
	RestPeriods        RestPeriods    `json:"rest_periods"`
	WeeklyFrequency    map[string]int `json:"weekly_frequency"`
	ProgressionMethods []string       `json:"progression_methods"`
}

type RepRanges struct {
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
}

type RestPeriods struct {
	Compound  string `json:"compound"`
	Isolation string `json:"isolation"`
}

// Equipment represents available equipment types
type Equipment struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	SpaceRequired      string   `json:"space_required"`
	Cost               string   `json:"cost"`
	Accessibility      string   `json:"accessibility"`
	ProgressionMethods []string `json:"progression_methods"`
}

// FocusArea represents workout focus areas
type FocusArea struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	Description      string         `json:"description"`
	MuscleGroups     []string       `json:"muscle_groups,omitempty"`
	Frequency        map[string]int `json:"frequency,omitempty"`
	MovementPatterns []string       `json:"movement_patterns,omitempty"`
	Types            []string       `json:"types,omitempty"`
	IntensityZones   []string       `json:"intensity_zones,omitempty"`
	WorkRestRatios   []string       `json:"work_rest_ratios,omitempty"`
	Intensity        string         `json:"intensity,omitempty"`
	Focus            []string       `json:"focus,omitempty"`
}

// ExerciseType represents different exercise categories
type ExerciseType struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Characteristics []string `json:"characteristics"`
	Adaptations     []string `json:"adaptations"`
}

// Exercise represents individual exercises
type Exercise struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Difficulty      string   `json:"difficulty"`
	MuscleGroups    []string `json:"muscle_groups"`
	Equipment       string   `json:"equipment"`
	Type            string   `json:"type"`
	MovementPattern string   `json:"movement_pattern"`
	DefaultSets     int      `json:"default_sets"`
	DefaultReps     string   `json:"default_reps"`
	RestSeconds     int      `json:"rest_seconds"`
	Instructions    []string `json:"instructions"`
	Progressions    []string `json:"progressions"`
	Regressions     []string `json:"regressions"`
	Benefits        []string `json:"benefits"`
}

// WorkoutTemplate represents complete workout structures
type WorkoutTemplate struct {
	ID             string                `json:"id"`
	Name           string                `json:"name"`
	Description    string                `json:"description"`
	DaysPerWeek    int                   `json:"days_per_week"`
	SuitableLevels []string              `json:"suitable_levels"`
	SuitableGoals  []string              `json:"suitable_goals"`
	Schedule       []string              `json:"schedule"`
	Structure      map[string]WorkoutDay `json:"structure"`
}

type WorkoutDay struct {
	Focus     string                `json:"focus"`
	Exercises []WorkoutExerciseSpec `json:"exercises"`
}

type WorkoutExerciseSpec struct {
	ExerciseID int    `json:"exercise_id"`
	Sets       int    `json:"sets"`
	Reps       string `json:"reps"`
	Rest       int    `json:"rest"`
}

// WeeklySchemaExample represents a complete generated workout plan
type WeeklySchemaExample struct {
	UserProfile   UserProfile   `json:"user_profile"`
	GeneratedPlan GeneratedPlan `json:"generated_plan"`
}

type UserProfile struct {
	UserID             int      `json:"user_id"`
	Level              string   `json:"level"`
	Goals              []string `json:"goals"`
	AvailableEquipment []string `json:"available_equipment"`
	WeeklyFrequency    int      `json:"weekly_frequency"`
	TimePerSession     int      `json:"time_per_session"`
	FocusPreference    string   `json:"focus_preference"`
}

type GeneratedPlan struct {
	WeekStart             string             `json:"week_start"`
	TemplateUsed          string             `json:"template_used"`
	TotalSessions         int                `json:"total_sessions"`
	EstimatedWeeklyVolume int                `json:"estimated_weekly_volume"`
	Workouts              []GeneratedWorkout `json:"workouts"`
	WeeklySummary         WeeklySummary      `json:"weekly_summary"`
}

type GeneratedWorkout struct {
	Day               string              `json:"day"`
	DayOfWeek         int                 `json:"day_of_week"`
	Focus             string              `json:"focus"`
	EstimatedDuration int                 `json:"estimated_duration"`
	Exercises         []GeneratedExercise `json:"exercises"`
}

type GeneratedExercise struct {
	ExerciseID        int      `json:"exercise_id"`
	Name              string   `json:"name"`
	Sets              int      `json:"sets"`
	Reps              string   `json:"reps"`
	RestSeconds       int      `json:"rest_seconds"`
	EstimatedDuration int      `json:"estimated_duration"`
	MuscleGroups      []string `json:"muscle_groups"`
	Notes             string   `json:"notes"`
}

type WeeklySummary struct {
	TotalExercises          int      `json:"total_exercises"`
	TotalSets               int      `json:"total_sets"`
	EstimatedTotalTime      int      `json:"estimated_total_time"`
	MuscleGroupsTrained     []string `json:"muscle_groups_trained"`
	MovementPatterns        []string `json:"movement_patterns"`
	EquipmentUsed           []string `json:"equipment_used"`
	ProgressionNotes        string   `json:"progression_notes"`
	RecoveryRecommendations []string `json:"recovery_recommendations"`
}

type ProgressionAlgorithm struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	SuitableFor []string               `json:"suitable_for"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type AdaptationTrigger struct {
	Metric    string   `json:"metric"`
	Threshold string   `json:"threshold"`
	Actions   []string `json:"actions"`
}

func LoadFitUpData() (*FitUpData, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get current file path")
	}

	dataDir := filepath.Dir(currentFile)
	jsonPath := filepath.Join(dataDir, "fitup_data.json")

	jsonData, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read fitup_data.json: %w", err)
	}

	var fitupData FitUpData
	if err := json.Unmarshal(jsonData, &fitupData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &fitupData, nil
}

func (f *FitUpData) GetExerciseByID(id int) (*Exercise, error) {
	for _, exercise := range f.Exercises {
		if exercise.ID == id {
			return &exercise, nil
		}
	}
	return nil, fmt.Errorf("exercise with ID %d not found", id)
}

func (f *FitUpData) GetExercisesByEquipment(equipment string) []Exercise {
	var filtered []Exercise
	for _, exercise := range f.Exercises {
		if exercise.Equipment == equipment {
			filtered = append(filtered, exercise)
		}
	}
	return filtered
}

func (f *FitUpData) GetExercisesByMuscleGroup(muscleGroup string) []Exercise {
	var filtered []Exercise
	for _, exercise := range f.Exercises {
		for _, mg := range exercise.MuscleGroups {
			if mg == muscleGroup {
				filtered = append(filtered, exercise)
				break
			}
		}
	}
	return filtered
}

func (f *FitUpData) GetExercisesByDifficulty(difficulty string) []Exercise {
	var filtered []Exercise
	for _, exercise := range f.Exercises {
		if exercise.Difficulty == difficulty {
			filtered = append(filtered, exercise)
		}
	}
	return filtered
}

func (f *FitUpData) GetWorkoutTemplateByGoalAndLevel(goal, level string) []WorkoutTemplate {
	var suitable []WorkoutTemplate
	for _, template := range f.WorkoutTemplates {
		goalMatch := false
		levelMatch := false

		for _, g := range template.SuitableGoals {
			if g == goal {
				goalMatch = true
				break
			}
		}

		for _, l := range template.SuitableLevels {
			if l == level {
				levelMatch = true
				break
			}
		}

		if goalMatch && levelMatch {
			suitable = append(suitable, template)
		}
	}
	return suitable
}

func (f *FitUpData) ConvertToGoTypes() ([]types.Exercise, []types.WorkoutTemplate, error) {
	var exercises []types.Exercise
	for _, ex := range f.Exercises {
		muscleGroupsStr := ""
		if len(ex.MuscleGroups) > 0 {
			muscleGroupsJSON, _ := json.Marshal(ex.MuscleGroups)
			muscleGroupsStr = string(muscleGroupsJSON)
		}

		exercise := types.Exercise{
			ExerciseID:   ex.ID,
			Name:         ex.Name,
			MuscleGroups: muscleGroupsStr,
			Difficulty:   types.FitnessLevel(ex.Difficulty),
			Equipment:    types.EquipmentType(ex.Equipment),
			Type:         types.ExerciseType(ex.Type),
			DefaultSets:  ex.DefaultSets,
			DefaultReps:  ex.DefaultReps,
			RestSeconds:  ex.RestSeconds,
		}
		exercises = append(exercises, exercise)
	}

	var templates []types.WorkoutTemplate
	for _, tmpl := range f.WorkoutTemplates {
		suitableGoalsJSON, _ := json.Marshal(tmpl.SuitableGoals)

		template := types.WorkoutTemplate{
			Name:          tmpl.Name,
			Description:   tmpl.Description,
			MinLevel:      types.FitnessLevel(tmpl.SuitableLevels[0]),
			MaxLevel:      types.FitnessLevel(tmpl.SuitableLevels[len(tmpl.SuitableLevels)-1]),
			SuitableGoals: string(suitableGoalsJSON),
			DaysPerWeek:   tmpl.DaysPerWeek,
		}
		templates = append(templates, template)
	}

	return exercises, templates, nil
}

func (f *FitUpData) ValidateData() error {
	exerciseIDs := make(map[int]bool)
	for _, exercise := range f.Exercises {
		if exerciseIDs[exercise.ID] {
			return fmt.Errorf("duplicate exercise ID found: %d", exercise.ID)
		}
		exerciseIDs[exercise.ID] = true
	}

	for templateID, template := range f.WorkoutTemplates {
		for dayName, day := range template.Structure {
			for _, ex := range day.Exercises {
				if !exerciseIDs[ex.ExerciseID] {
					return fmt.Errorf("template %s, day %s references invalid exercise ID: %d",
						templateID, dayName, ex.ExerciseID)
				}
			}
		}
	}

	return nil
}
