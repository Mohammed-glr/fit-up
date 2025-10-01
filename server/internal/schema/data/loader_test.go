package data

import (
	"encoding/json"
	"testing"
)

func TestLoadFitUpData(t *testing.T) {
	// Test loading the JSON data
	data, err := LoadFitUpData()
	if err != nil {
		t.Fatalf("Failed to load FitUp data: %v", err)
	}

	// Basic validation tests
	if data.Meta.Version == "" {
		t.Error("Meta version should not be empty")
	}

	if len(data.Exercises) == 0 {
		t.Error("Should have at least one exercise")
	}

	if len(data.Levels) == 0 {
		t.Error("Should have fitness levels defined")
	}

	if len(data.Goals) == 0 {
		t.Error("Should have fitness goals defined")
	}

	// Test validation
	err = data.ValidateData()
	if err != nil {
		t.Errorf("Data validation failed: %v", err)
	}
}

func TestConvertToGoTypes(t *testing.T) {
	// Test conversion to existing Go types
	data, err := LoadFitUpData()
	if err != nil {
		t.Fatalf("Failed to load FitUp data: %v", err)
	}

	exercises, templates, err := data.ConvertToGoTypes()
	if err != nil {
		t.Fatalf("Failed to convert to Go types: %v", err)
	}

	if len(exercises) == 0 {
		t.Error("Should have converted exercises")
	}

	if len(templates) == 0 {
		t.Error("Should have converted templates")
	}

	// Test that first exercise has valid fields
	if len(exercises) > 0 {
		ex := exercises[0]
		if ex.Name == "" {
			t.Error("Exercise name should not be empty")
		}
		if ex.DefaultSets <= 0 {
			t.Error("Exercise should have positive default sets")
		}
		if ex.RestSeconds < 0 {
			t.Error("Exercise rest seconds should not be negative")
		}
	}
}

func TestGetExerciseByID(t *testing.T) {
	data, err := LoadFitUpData()
	if err != nil {
		t.Fatalf("Failed to load FitUp data: %v", err)
	}

	// Test getting existing exercise
	exercise, err := data.GetExerciseByID(1)
	if err != nil {
		t.Errorf("Failed to get exercise by ID: %v", err)
	}
	if exercise == nil {
		t.Error("Exercise should not be nil")
	}

	// Test getting non-existing exercise
	_, err = data.GetExerciseByID(9999)
	if err == nil {
		t.Error("Should return error for non-existing exercise")
	}
}

func TestGetExercisesByEquipment(t *testing.T) {
	data, err := LoadFitUpData()
	if err != nil {
		t.Fatalf("Failed to load FitUp data: %v", err)
	}

	// Test filtering by bodyweight equipment
	bodyweightExercises := data.GetExercisesByEquipment("bodyweight")
	if len(bodyweightExercises) == 0 {
		t.Error("Should have bodyweight exercises")
	}

	// Verify all returned exercises are bodyweight
	for _, ex := range bodyweightExercises {
		if ex.Equipment != "bodyweight" {
			t.Errorf("Expected bodyweight exercise, got %s", ex.Equipment)
		}
	}
}

func TestGetWorkoutTemplateByGoalAndLevel(t *testing.T) {
	data, err := LoadFitUpData()
	if err != nil {
		t.Fatalf("Failed to load FitUp data: %v", err)
	}

	// Test finding templates for muscle gain and intermediate level
	templates := data.GetWorkoutTemplateByGoalAndLevel("muscle_gain", "intermediate")
	if len(templates) == 0 {
		t.Error("Should find templates for muscle gain and intermediate level")
	}

	// Verify returned templates match criteria
	for _, tmpl := range templates {
		goalMatch := false
		levelMatch := false

		for _, goal := range tmpl.SuitableGoals {
			if goal == "muscle_gain" {
				goalMatch = true
				break
			}
		}

		for _, level := range tmpl.SuitableLevels {
			if level == "intermediate" {
				levelMatch = true
				break
			}
		}

		if !goalMatch {
			t.Errorf("Template %s should support muscle_gain goal", tmpl.Name)
		}
		if !levelMatch {
			t.Errorf("Template %s should support intermediate level", tmpl.Name)
		}
	}
}

func TestJSONStructure(t *testing.T) {
	// Test that the JSON can be marshaled and unmarshaled correctly
	data, err := LoadFitUpData()
	if err != nil {
		t.Fatalf("Failed to load FitUp data: %v", err)
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal data to JSON: %v", err)
	}

	// Unmarshal back
	var newData FitUpData
	err = json.Unmarshal(jsonData, &newData)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	// Basic comparison
	if newData.Meta.Version != data.Meta.Version {
		t.Error("Version mismatch after JSON round-trip")
	}

	if len(newData.Exercises) != len(data.Exercises) {
		t.Error("Exercise count mismatch after JSON round-trip")
	}
}
