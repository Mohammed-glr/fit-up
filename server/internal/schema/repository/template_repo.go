package repository

import (
	"context"

	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

func (s *Store) CreateTemplate(ctx context.Context, template *types.WorkoutTemplateRequest) (*types.WorkoutTemplate, error) {
	q := `
		INSERT INTO workout_templates (name, description, min_level, max_level, suitable_goals, days_per_week)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING template_id, name, description, min_level, max_level, suitable_goals, days_per_week
	`

	// Convert suitable goals slice to comma-separated string
	suitableGoalsStr := ""
	if len(template.SuitableGoals) > 0 {
		for i, goal := range template.SuitableGoals {
			if i > 0 {
				suitableGoalsStr += ","
			}
			suitableGoalsStr += string(goal)
		}
	}

	row := s.db.QueryRow(ctx, q,
		template.Name,
		template.Description,
		template.MinLevel,
		template.MaxLevel,
		suitableGoalsStr,
		template.DaysPerWeek,
	)

	var wt types.WorkoutTemplate
	err := row.Scan(
		&wt.TemplateID,
		&wt.Name,
		&wt.Description,
		&wt.MinLevel,
		&wt.MaxLevel,
		&wt.SuitableGoals,
		&wt.DaysPerWeek,
	)
	if err != nil {
		return nil, err
	}

	return &wt, nil
}

func (s *Store) GetTemplateByID(ctx context.Context, templateID int) (*types.WorkoutTemplate, error) {
	q := `
		SELECT template_id, name, description, min_level, max_level, suitable_goals, days_per_week
		FROM workout_templates
		WHERE template_id = $1
	`
	row := s.db.QueryRow(ctx, q, templateID)
	var wt types.WorkoutTemplate
	err := row.Scan(
		&wt.TemplateID,
		&wt.Name,
		&wt.Description,
		&wt.MinLevel,
		&wt.MaxLevel,
		&wt.SuitableGoals,
		&wt.DaysPerWeek,
	)
	if err != nil {
		return nil, err
	}

	return &wt, nil
}

func (s *Store) UpdateTemplate(ctx context.Context, templateID int, template *types.WorkoutTemplateRequest) (*types.WorkoutTemplate, error) {
	q := `
		UPDATE workout_templates
		SET name = $1, description = $2, level = $3, goal = $4, duration_minutes = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING id, name, description, level, goal, duration_minutes, created_at, updated_at
	`
	row := s.db.QueryRow(ctx, q,
		template.Name,
		template.DaysPerWeek,
		template.Description,
		template.MaxLevel,
		template.MinLevel,
		template.SuitableGoals,
		templateID,
	)
	var wt types.WorkoutTemplate
	err := row.Scan(
		&wt.TemplateID,
		&wt.SuitableGoals,
		&wt.Name,
		&wt.MinLevel,
		&wt.MaxLevel,
		&wt.Description,
		&wt.DaysPerWeek,
	)
	if err != nil {
		return nil, err
	}

	return &wt, nil
}

func (s *Store) DeleteTemplate(ctx context.Context, templateID int) error {
	q := `
		DELETE FROM workout_templates
		WHERE id = $1
	`
	_, err := s.db.Exec(ctx, q, templateID)
	return err
}

func (s *Store) ListTemplates(ctx context.Context, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutTemplate], error) {
	q := `
		SELECT id, name, description, level, goal, duration_minutes, created_at, updated_at
		FROM workout_templates
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := s.db.Query(ctx, q, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var templates []types.WorkoutTemplate
	for rows.Next() {
		var wt types.WorkoutTemplate
		err := rows.Scan(
			&wt.TemplateID,
			&wt.SuitableGoals,
			&wt.Name,
			&wt.MinLevel,
			&wt.MaxLevel,
			&wt.Description,
			&wt.DaysPerWeek,
		)
		if err != nil {
			return nil, err
		}
		templates = append(templates, wt)
	}

	return &types.PaginatedResponse[types.WorkoutTemplate]{
		Data:       templates,
		TotalCount: len(templates),
		TotalPages: 1,
		Page:       1,
		PageSize:   pagination.Limit,
	}, nil
}

func (s *Store) FilterTemplates(ctx context.Context, filter types.TemplateFilter, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutTemplate], error) {

	q := `
		SELECT id, name, description, level, goal, duration_minutes, created_at, updated_at
		FROM workout_templates
		WHERE ($1::INT[] IS NULL OR level && $1::INT[])
		AND ($2::TEXT[] IS NULL OR goal && $2::TEXT[])
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := s.db.Query(ctx, q,
		filter.Search,
		filter.Level,
		filter.Goals,
		filter.DaysPerWeek,
		pagination.Limit,
		pagination.Offset,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var templates []types.WorkoutTemplate
	for rows.Next() {
		var wt types.WorkoutTemplate
		err := rows.Scan(
			&wt.TemplateID,
			&wt.SuitableGoals,
			&wt.Name,
			&wt.MinLevel,
			&wt.MaxLevel,
			&wt.Description,
			&wt.DaysPerWeek,
		)
		if err != nil {
			return nil, err
		}
		templates = append(templates, wt)
	}

	return &types.PaginatedResponse[types.WorkoutTemplate]{
		Data:       templates,
		TotalCount: len(templates),
		TotalPages: 1,
		Page:       1,
		PageSize:   pagination.Limit,
	}, nil
}

func (s *Store) SearchTemplates(ctx context.Context, query string, pagination types.PaginationParams) (*types.PaginatedResponse[types.WorkoutTemplate], error) {
	q := `
		SELECT id, name, description, level, goal, duration_minutes, created_at, updated_at
		FROM workout_templates
		WHERE name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := s.db.Query(ctx, q, query, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var templates []types.WorkoutTemplate
	for rows.Next() {
		var wt types.WorkoutTemplate
		err := rows.Scan(
			&wt.TemplateID,
			&wt.SuitableGoals,
			&wt.Name,
			&wt.MinLevel,
			&wt.MaxLevel,
			&wt.Description,
			&wt.DaysPerWeek,
		)
		if err != nil {
			return nil, err
		}
		templates = append(templates, wt)
	}

	return &types.PaginatedResponse[types.WorkoutTemplate]{
		Data:       templates,
		TotalCount: len(templates),
		TotalPages: 1,
		Page:       1,
		PageSize:   pagination.Limit,
	}, nil
}

func (s *Store) GetTemplatesByLevel(ctx context.Context, level types.FitnessLevel) ([]types.WorkoutTemplate, error) {
	q := `
		SELECT id, name, description, level, goal, duration_minutes, created_at, updated_at
		FROM workout_templates
		WHERE $1 = ANY(level)
		ORDER BY created_at DESC
	`
	rows, err := s.db.Query(ctx, q, level)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []types.WorkoutTemplate
	for rows.Next() {
		var wt types.WorkoutTemplate
		err := rows.Scan(
			&wt.TemplateID,
			&wt.SuitableGoals,
			&wt.Name,
			&wt.MinLevel,
			&wt.MaxLevel,
			&wt.Description,
			&wt.DaysPerWeek,
		)
		if err != nil {
			return nil, err
		}
		templates = append(templates, wt)
	}

	return templates, nil
}

func (s *Store) GetTemplatesByGoal(ctx context.Context, goal types.FitnessGoal) ([]types.WorkoutTemplate, error) {
	q := `
		SELECT id, name, description, level, goal, duration_minutes, created_at, updated_at
		FROM workout_templates
		WHERE $1 = ANY(goal)
		ORDER BY created_at DESC
	`
	rows, err := s.db.Query(ctx, q, goal)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []types.WorkoutTemplate
	for rows.Next() {
		var wt types.WorkoutTemplate
		err := rows.Scan(
			&wt.TemplateID,
			&wt.SuitableGoals,
			&wt.Name,
			&wt.MinLevel,
			&wt.MaxLevel,
			&wt.Description,
			&wt.DaysPerWeek,
		)
		if err != nil {
			return nil, err
		}
		templates = append(templates, wt)
	}

	return templates, nil
}

func (s *Store) GetRecommendedTemplates(ctx context.Context, userID int, count int) ([]types.WorkoutTemplate, error) {
	q := `
		SELECT wt.id, wt.name, wt.description, wt.level, wt.goal, wt.duration_minutes, wt.created_at, wt.updated_at
		FROM workout_templates wt
		JOIN users u ON u.fitness_level = ANY(wt.level) AND u.fitness_goal = ANY(wt.goal)
		WHERE u.id = $1
		ORDER BY wt.created_at DESC
		LIMIT $2
	`
	rows, err := s.db.Query(ctx, q, userID, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []types.WorkoutTemplate
	for rows.Next() {
		var wt types.WorkoutTemplate
		err := rows.Scan(
			&wt.TemplateID,
			&wt.SuitableGoals,
			&wt.Name,
			&wt.MinLevel,
			&wt.MaxLevel,
			&wt.Description,
			&wt.DaysPerWeek,
		)
		if err != nil {
			return nil, err
		}
		templates = append(templates, wt)
	}

	return templates, nil
}

func (s *Store) GetPopularTemplates(ctx context.Context, count int) ([]types.WorkoutTemplate, error) {
	q := `
		SELECT wt.id, wt.name, wt.description, wt.level, wt.goal, wt.duration_minutes, wt.created_at, wt.updated_at
		FROM workout_templates wt
		LEFT JOIN workouts w ON w.template_id = wt.id
		GROUP BY wt.id
		ORDER BY COUNT(w.id) DESC
		LIMIT $1
	`
	rows, err := s.db.Query(ctx, q, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []types.WorkoutTemplate
	for rows.Next() {
		var wt types.WorkoutTemplate
		err := rows.Scan(
			&wt.TemplateID,
			&wt.SuitableGoals,
			&wt.Name,
			&wt.MinLevel,
			&wt.MaxLevel,
			&wt.Description,
			&wt.DaysPerWeek,
		)
		if err != nil {
			return nil, err
		}
		templates = append(templates, wt)
	}

	return templates, nil
}
