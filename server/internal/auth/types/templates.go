package types

import "time"

type TemplateExercise struct {
	ExerciseName string  `json:"exercise_name"`
	Sets         int     `json:"sets"`
	TargetReps   int     `json:"target_reps"`
	TargetWeight float64 `json:"target_weight"`
	RestSeconds  int     `json:"rest_seconds"`
}

type WorkoutTemplate struct {
	TemplateID  int                `json:"template_id"`
	UserID      string             `json:"user_id"`
	Name        string             `json:"name"`
	Description *string            `json:"description,omitempty"`
	IsPublic    bool               `json:"is_public"`
	Exercises   []TemplateExercise `json:"exercises"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type CreateTemplateRequest struct {
	Name        string             `json:"name" validate:"required,min=1,max=255"`
	Description *string            `json:"description,omitempty"`
	IsPublic    bool               `json:"is_public"`
	Exercises   []TemplateExercise `json:"exercises" validate:"required,min=1"`
}

type UpdateTemplateRequest struct {
	Name        *string            `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string            `json:"description,omitempty"`
	IsPublic    *bool              `json:"is_public,omitempty"`
	Exercises   []TemplateExercise `json:"exercises,omitempty" validate:"omitempty,min=1"`
}

type TemplateListResponse struct {
	Templates  []WorkoutTemplate `json:"templates"`
	TotalCount int               `json:"total_count"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	HasMore    bool              `json:"has_more"`
}
