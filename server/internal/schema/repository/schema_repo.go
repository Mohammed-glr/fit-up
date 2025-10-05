package repository

import (
	"context"
	"time"

	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

func (s *Store) CreateWeeklySchema(ctx context.Context, schema *types.WeeklySchemaRequest) (*types.WeeklySchema, error) {
	q := `
		INSERT INTO weekly_schemas (user_id, week_start, active)
		VALUES ($1, $2, $3)
		RETURNING schema_id, user_id, week_start, active
	`
	row := s.db.QueryRow(ctx, q,
		schema.UserID,
		schema.WeekStart,
		true, 
	)

	var ws types.WeeklySchema
	err := row.Scan(
		&ws.SchemaID,
		&ws.UserID,
		&ws.WeekStart,
		&ws.Active,
	)
	if err != nil {
		return nil, err
	}

	return &ws, nil
}

func (s *Store) GetWeeklySchemaByID(ctx context.Context, schemaID int) (*types.WeeklySchema, error) {
	q := `
		SELECT schema_id, user_id, week_start, active
		FROM weekly_schemas
		WHERE schema_id = $1
	`
	row := s.db.QueryRow(ctx, q, schemaID)
	var ws types.WeeklySchema
	err := row.Scan(
		&ws.SchemaID,
		&ws.UserID,
		&ws.WeekStart,
		&ws.Active,
	)
	if err != nil {
		return nil, err
	}
	return &ws, nil
}

func (s *Store) UpdateWeeklySchema(ctx context.Context, schemaID int, active bool) (*types.WeeklySchema, error) {
	q := `
		UPDATE weekly_schemas
		SET active = $1
		WHERE schema_id = $2
		RETURNING schema_id, user_id, week_start, active
	`

	row := s.db.QueryRow(ctx, q, active, schemaID)
	var ws types.WeeklySchema
	err := row.Scan(
		&ws.SchemaID,
		&ws.UserID,
		&ws.WeekStart,
		&ws.Active,
	)
	if err != nil {
		return nil, err
	}
	return &ws, nil
}

func (s *Store) DeleteWeeklySchema(ctx context.Context, schemaID int) error {
	q := `
		DELETE FROM weekly_schemas
		WHERE schema_id = $1
	`
	_, err := s.db.Exec(ctx, q, schemaID)
	return err
}

func (s *Store) GetWeeklySchemasByUserID(ctx context.Context, userID int, pagination types.PaginationParams) (*types.PaginatedResponse[types.WeeklySchema], error) {
	q := `
		SELECT schema_id, user_id, week_start, active
		FROM weekly_schemas
		WHERE user_id = $1
		ORDER BY week_start DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := s.db.Query(ctx, q, userID, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemas []types.WeeklySchema
	for rows.Next() {
		var ws types.WeeklySchema
		err := rows.Scan(
			&ws.SchemaID,
			&ws.UserID,
			&ws.WeekStart,
			&ws.Active,
		)
		if err != nil {
			return nil, err
		}
		schemas = append(schemas, ws)
	}

	countQuery := `SELECT COUNT(*) FROM weekly_schemas WHERE user_id = $1`
	var totalCount int
	err = s.db.QueryRow(ctx, countQuery, userID).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	totalPages := (totalCount + pagination.PageSize - 1) / pagination.PageSize

	return &types.PaginatedResponse[types.WeeklySchema]{
		Data:       schemas,
		TotalCount: totalCount,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *Store) GetActiveWeeklySchemaByUserID(ctx context.Context, userID int) (*types.WeeklySchema, error) {
	q := `
		SELECT schema_id, user_id, week_start, active
		FROM weekly_schemas
		WHERE user_id = $1 AND active = true
	`
	row := s.db.QueryRow(ctx, q, userID)
	var ws types.WeeklySchema
	err := row.Scan(
		&ws.SchemaID,
		&ws.UserID,
		&ws.WeekStart,
		&ws.Active,
	)
	if err != nil {
		return nil, err
	}
	return &ws, nil
}

func (s *Store) GetWeeklySchemaByUserAndWeek(ctx context.Context, userID int, weekStart time.Time) (*types.WeeklySchema, error) {
	q := `
		SELECT schema_id, user_id, week_start, active
		FROM weekly_schemas
		WHERE user_id = $1 AND week_start = $2
	`
	row := s.db.QueryRow(ctx, q, userID, weekStart)
	var ws types.WeeklySchema
	err := row.Scan(
		&ws.SchemaID,
		&ws.UserID,
		&ws.WeekStart,
		&ws.Active,
	)
	if err != nil {
		return nil, err
	}
	return &ws, nil
}

func (s *Store) DeactivateAllWeeklySchemasForUser(ctx context.Context, userID int) error {
	q := `
		UPDATE weekly_schemas
		SET active = false, updated_at = NOW()
		WHERE user_id = $1
	`
	_, err := s.db.Exec(ctx, q, userID)
	return err
}

func (s *Store) GetCurrentWeekSchema(ctx context.Context, userID int) (*types.WeeklySchema, error) {
	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))
	q := `
		SELECT id, user_id, template_id, week_start, week_end, created_at, updated_at
		FROM weekly_schemas
		WHERE user_id = $1 AND week_start = $2
	`
	row := s.db.QueryRow(ctx, q, userID, weekStart)
	var ws types.WeeklySchema
	err := row.Scan(
		&ws.SchemaID,
		&ws.UserID,
		&ws.WeekStart,
		&ws.Active,
	)
	if err != nil {
		return nil, err
	}
	return &ws, nil
}

func (s *Store) GetWeeklySchemaHistory(ctx context.Context, userID int, limit int) ([]types.WeeklySchema, error) {
	q := `
		SELECT id, user_id, template_id, week_start, week_end, created_at, updated_at
		FROM weekly_schemas
		WHERE user_id = $1
		ORDER BY week_start DESC
		LIMIT $2
	`
	rows, err := s.db.Query(ctx, q, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemas []types.WeeklySchema
	for rows.Next() {
		var ws types.WeeklySchema
		err := rows.Scan(
			&ws.SchemaID,
			&ws.UserID,
			&ws.WeekStart,
		)
		if err != nil {
			return nil, err
		}
		schemas = append(schemas, ws)
	}
	return schemas, nil
}

func (s *Store) SaveSchemaAsTemplate(ctx context.Context, schemaID int, templateName string) error {
	q := `
		INSERT INTO schema_templates (schema_id, template_name, created_at)
		VALUES ($1, $2, NOW())
	`
	_, err := s.db.Exec(ctx, q, schemaID, templateName)
	return err
}