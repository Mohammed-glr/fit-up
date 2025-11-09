package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tdmdh/fit-up-server/internal/auth/types"
)

func (s *Store) GetUserTemplates(ctx context.Context, userID string, page, pageSize int) (*types.TemplateListResponse, error) {
	offset := (page - 1) * pageSize

	var totalCount int
	countQuery := `SELECT COUNT(*) FROM workout_templates WHERE user_id = $1`
	err := s.db.QueryRow(ctx, countQuery, userID).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count templates: %w", err)
	}

	query := `
		SELECT template_id, user_id, name, description, is_public, exercises, created_at, updated_at
		FROM workout_templates
		WHERE user_id = $1
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.Query(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query templates: %w", err)
	}
	defer rows.Close()

	var templates []types.WorkoutTemplate
	for rows.Next() {
		var template types.WorkoutTemplate
		var exercisesJSON []byte

		err := rows.Scan(
			&template.TemplateID,
			&template.UserID,
			&template.Name,
			&template.Description,
			&template.IsPublic,
			&exercisesJSON,
			&template.CreatedAt,
			&template.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}

		if err := json.Unmarshal(exercisesJSON, &template.Exercises); err != nil {
			return nil, fmt.Errorf("failed to unmarshal exercises: %w", err)
		}

		templates = append(templates, template)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating templates: %w", err)
	}

	hasMore := offset+len(templates) < totalCount

	return &types.TemplateListResponse{
		Templates:  templates,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		HasMore:    hasMore,
	}, nil
}

func (s *Store) GetPublicTemplates(ctx context.Context, page, pageSize int) (*types.TemplateListResponse, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM workout_templates WHERE is_public = TRUE`
	err := s.db.QueryRow(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count public templates: %w", err)
	}

	query := `
		SELECT template_id, user_id, name, description, is_public, exercises, created_at, updated_at
		FROM workout_templates
		WHERE is_public = TRUE
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query public templates: %w", err)
	}
	defer rows.Close()

	var templates []types.WorkoutTemplate
	for rows.Next() {
		var template types.WorkoutTemplate
		var exercisesJSON []byte

		err := rows.Scan(
			&template.TemplateID,
			&template.UserID,
			&template.Name,
			&template.Description,
			&template.IsPublic,
			&exercisesJSON,
			&template.CreatedAt,
			&template.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}

		if err := json.Unmarshal(exercisesJSON, &template.Exercises); err != nil {
			return nil, fmt.Errorf("failed to unmarshal exercises: %w", err)
		}

		templates = append(templates, template)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating templates: %w", err)
	}

	hasMore := offset+len(templates) < totalCount

	return &types.TemplateListResponse{
		Templates:  templates,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		HasMore:    hasMore,
	}, nil
}

// GetTemplateByID retrieves a template by ID
func (s *Store) GetTemplateByID(ctx context.Context, templateID int) (*types.WorkoutTemplate, error) {
	query := `
		SELECT template_id, user_id, name, description, is_public, exercises, created_at, updated_at
		FROM workout_templates
		WHERE template_id = $1
	`

	var template types.WorkoutTemplate
	var exercisesJSON []byte

	err := s.db.QueryRow(ctx, query, templateID).Scan(
		&template.TemplateID,
		&template.UserID,
		&template.Name,
		&template.Description,
		&template.IsPublic,
		&exercisesJSON,
		&template.CreatedAt,
		&template.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	if err := json.Unmarshal(exercisesJSON, &template.Exercises); err != nil {
		return nil, fmt.Errorf("failed to unmarshal exercises: %w", err)
	}

	return &template, nil
}

// CreateTemplate creates a new workout template
func (s *Store) CreateTemplate(ctx context.Context, userID string, req *types.CreateTemplateRequest) (*types.WorkoutTemplate, error) {
	// Marshal exercises to JSON
	exercisesJSON, err := json.Marshal(req.Exercises)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal exercises: %w", err)
	}

	query := `
		INSERT INTO workout_templates (user_id, name, description, is_public, exercises)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING template_id, user_id, name, description, is_public, exercises, created_at, updated_at
	`

	var template types.WorkoutTemplate
	var returnedExercisesJSON []byte

	err = s.db.QueryRow(ctx, query,
		userID,
		req.Name,
		req.Description,
		req.IsPublic,
		exercisesJSON,
	).Scan(
		&template.TemplateID,
		&template.UserID,
		&template.Name,
		&template.Description,
		&template.IsPublic,
		&returnedExercisesJSON,
		&template.CreatedAt,
		&template.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	if err := json.Unmarshal(returnedExercisesJSON, &template.Exercises); err != nil {
		return nil, fmt.Errorf("failed to unmarshal exercises: %w", err)
	}

	return &template, nil
}

// UpdateTemplate updates an existing template
func (s *Store) UpdateTemplate(ctx context.Context, templateID int, userID string, req *types.UpdateTemplateRequest) (*types.WorkoutTemplate, error) {
	// First, check if template exists and belongs to user
	existingTemplate, err := s.GetTemplateByID(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}

	if existingTemplate.UserID != userID {
		return nil, fmt.Errorf("unauthorized: template does not belong to user")
	}

	// Build dynamic update query
	updateFields := []string{}
	args := []interface{}{}
	argCounter := 1

	if req.Name != nil {
		updateFields = append(updateFields, fmt.Sprintf("name = $%d", argCounter))
		args = append(args, *req.Name)
		argCounter++
	}

	if req.Description != nil {
		updateFields = append(updateFields, fmt.Sprintf("description = $%d", argCounter))
		args = append(args, *req.Description)
		argCounter++
	}

	if req.IsPublic != nil {
		updateFields = append(updateFields, fmt.Sprintf("is_public = $%d", argCounter))
		args = append(args, *req.IsPublic)
		argCounter++
	}

	if req.Exercises != nil {
		exercisesJSON, err := json.Marshal(req.Exercises)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal exercises: %w", err)
		}
		updateFields = append(updateFields, fmt.Sprintf("exercises = $%d", argCounter))
		args = append(args, exercisesJSON)
		argCounter++
	}

	if len(updateFields) == 0 {
		return existingTemplate, nil // No updates to perform
	}

	// Add template_id to args
	args = append(args, templateID)

	// Build SET clause
	setClause := ""
	for i, field := range updateFields {
		if i > 0 {
			setClause += ", "
		}
		setClause += field
	}

	query := fmt.Sprintf(`
		UPDATE workout_templates
		SET %s
		WHERE template_id = $%d
		RETURNING template_id, user_id, name, description, is_public, exercises, created_at, updated_at
	`, setClause, argCounter)

	var template types.WorkoutTemplate
	var exercisesJSON []byte

	err = s.db.QueryRow(ctx, query, args...).Scan(
		&template.TemplateID,
		&template.UserID,
		&template.Name,
		&template.Description,
		&template.IsPublic,
		&exercisesJSON,
		&template.CreatedAt,
		&template.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update template: %w", err)
	}

	// Parse exercises JSON
	if err := json.Unmarshal(exercisesJSON, &template.Exercises); err != nil {
		return nil, fmt.Errorf("failed to unmarshal exercises: %w", err)
	}

	return &template, nil
}

// DeleteTemplate deletes a template
func (s *Store) DeleteTemplate(ctx context.Context, templateID int, userID string) error {
	// First, check if template exists and belongs to user
	existingTemplate, err := s.GetTemplateByID(ctx, templateID)
	if err != nil {
		return fmt.Errorf("template not found: %w", err)
	}

	if existingTemplate.UserID != userID {
		return fmt.Errorf("unauthorized: template does not belong to user")
	}

	query := `DELETE FROM workout_templates WHERE template_id = $1`
	_, err = s.db.Exec(ctx, query, templateID)
	if err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	return nil
}
