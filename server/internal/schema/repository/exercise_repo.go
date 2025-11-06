package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

func (s *Store) CreateExercise(ctx context.Context, exercise *types.ExerciseRequest) (*types.Exercise, error) {
	q := `
		INSERT INTO exercises (name, muscle_groups, difficulty, equipment, type, default_sets, default_reps, rest_seconds)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING exercise_id, name, muscle_groups, difficulty, equipment, type, default_sets, default_reps, rest_seconds
	`

	muscleGroupsStr := ""
	if len(exercise.MuscleGroups) > 0 {
		for i, group := range exercise.MuscleGroups {
			if i > 0 {
				muscleGroupsStr += ","
			}
			muscleGroupsStr += group
		}
	}

	var createdExercise types.Exercise
	err := s.db.QueryRow(ctx, q,
		exercise.Name,
		muscleGroupsStr,
		exercise.Difficulty,
		exercise.Equipment,
		exercise.Type,
		exercise.DefaultSets,
		exercise.DefaultReps,
		exercise.RestSeconds,
	).Scan(
		&createdExercise.ExerciseID,
		&createdExercise.Name,
		&createdExercise.MuscleGroups,
		&createdExercise.Difficulty,
		&createdExercise.Equipment,
		&createdExercise.Type,
		&createdExercise.DefaultSets,
		&createdExercise.DefaultReps,
		&createdExercise.RestSeconds,
	)
	if err != nil {
		return nil, err
	}

	return &createdExercise, nil
}

func (s *Store) GetExerciseByID(ctx context.Context, exerciseID int) (*types.Exercise, error) {
	q := `
		SELECT exercise_id, name, muscle_groups, difficulty, equipment, type, default_sets, default_reps, rest_seconds
		FROM exercises
		WHERE exercise_id = $1
	`

	var exercise types.Exercise
	err := s.db.QueryRow(ctx, q, exerciseID).Scan(
		&exercise.ExerciseID,
		&exercise.Name,
		&exercise.MuscleGroups,
		&exercise.Difficulty,
		&exercise.Equipment,
		&exercise.Type,
		&exercise.DefaultSets,
		&exercise.DefaultReps,
		&exercise.RestSeconds,
	)
	if err != nil {
		return nil, err
	}

	return &exercise, nil
}

func (s *Store) UpdateExercise(ctx context.Context, exerciseID int, exercise *types.ExerciseRequest) (*types.Exercise, error) {
	q := `
		UPDATE exercises
		SET name = $1, muscle_groups = $2, difficulty = $3, equipment = $4, type = $5, default_sets = $6, default_reps = $7, rest_seconds = $8
		WHERE exercise_id = $9
		RETURNING exercise_id, name, muscle_groups, difficulty, equipment, type, default_sets, default_reps, rest_seconds
	`

	muscleGroupsStr := ""
	if len(exercise.MuscleGroups) > 0 {
		for i, group := range exercise.MuscleGroups {
			if i > 0 {
				muscleGroupsStr += ","
			}
			muscleGroupsStr += group
		}
	}

	var updatedExercise types.Exercise
	err := s.db.QueryRow(ctx, q,
		exercise.Name,
		muscleGroupsStr,
		exercise.Difficulty,
		exercise.Equipment,
		exercise.Type,
		exercise.DefaultSets,
		exercise.DefaultReps,
		exercise.RestSeconds,
		exerciseID,
	).Scan(
		&updatedExercise.ExerciseID,
		&updatedExercise.Name,
		&updatedExercise.MuscleGroups,
		&updatedExercise.Difficulty,
		&updatedExercise.Equipment,
		&updatedExercise.Type,
		&updatedExercise.DefaultSets,
		&updatedExercise.DefaultReps,
		&updatedExercise.RestSeconds,
	)
	if err != nil {
		return nil, err
	}

	return &updatedExercise, nil
}

func (s *Store) DeleteExercise(ctx context.Context, exerciseID int) error {
	q := `
		DELETE FROM exercises
		WHERE id = $1
	`

	_, err := s.db.Exec(ctx, q, exerciseID)
	return err
}

func (s *Store) ListExercises(ctx context.Context, pagination types.PaginationParams) (*types.PaginatedResponse[types.Exercise], error) {
	q := `
		SELECT created_at, name, description, muscle_group, equipment, difficulty, id AS exercise_id
		FROM exercises
		ORDER BY created_at DESC
		OFFSET $1 LIMIT $2
	`

	rows, err := s.db.Query(ctx, q, pagination.Offset, pagination.Limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var exercises []types.Exercise
	for rows.Next() {
		var exercise types.Exercise
		if err := rows.Scan(
			&exercise.DefaultReps,
			&exercise.DefaultSets,
			&exercise.Difficulty,
			&exercise.Equipment,
			&exercise.ExerciseID,
			&exercise.MuscleGroups,
			&exercise.Name,
			&exercise.RestSeconds,
			&exercise.Type,
		); err != nil {
			return nil, err
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	countQ := `SELECT COUNT(*) FROM exercises`
	var total int
	err = s.db.QueryRow(ctx, countQ).Scan(&total)
	if err != nil {
		return nil, err
	}

	totalPage := (total + pagination.Limit - 1) / pagination.Limit
	if pagination.Page > totalPage {
		exercises = []types.Exercise{}
	}

	response := &types.PaginatedResponse[types.Exercise]{
		Data:       exercises,
		TotalCount: total,
		TotalPages: totalPage,
		Page:       pagination.Page,
		PageSize:   pagination.Limit,
	}
	return response, nil
}

func (s *Store) FilterExercises(ctx context.Context, filter types.ExerciseFilter, pagination types.PaginationParams) (*types.PaginatedResponse[types.Exercise], error) {
	baseQuery := `FROM exercises WHERE 1=1`
	args := make([]interface{}, 0)

	if len(filter.MuscleGroups) > 0 {
		placeholders := make([]string, len(filter.MuscleGroups))
		startIndex := len(args)
		for i, group := range filter.MuscleGroups {
			placeholders[i] = fmt.Sprintf("muscle_groups ILIKE $%d", startIndex+i+1)
			args = append(args, "%"+group+"%")
		}
		baseQuery += " AND (" + strings.Join(placeholders, " OR ") + ")"
	}

	if filter.Difficulty != nil {
		baseQuery += fmt.Sprintf(" AND difficulty = $%d", len(args)+1)
		args = append(args, *filter.Difficulty)
	}

	if len(filter.Equipment) > 0 {
		equipmentValues := make([]string, len(filter.Equipment))
		for i, equipment := range filter.Equipment {
			equipmentValues[i] = string(equipment)
		}
		if len(equipmentValues) == 1 {
			baseQuery += fmt.Sprintf(" AND equipment = $%d", len(args)+1)
			args = append(args, equipmentValues[0])
		} else {
			baseQuery += fmt.Sprintf(" AND equipment = ANY($%d)", len(args)+1)
			args = append(args, equipmentValues)
		}
	}

	if len(filter.Type) > 0 {
		typeValues := make([]string, len(filter.Type))
		for i, exerciseType := range filter.Type {
			typeValues[i] = string(exerciseType)
		}
		if len(typeValues) == 1 {
			baseQuery += fmt.Sprintf(" AND type = $%d", len(args)+1)
			args = append(args, typeValues[0])
		} else {
			baseQuery += fmt.Sprintf(" AND type = ANY($%d)", len(args)+1)
			args = append(args, typeValues)
		}
	}

	if filter.Search != "" {
		nameParam := len(args) + 1
		groupParam := len(args) + 2
		baseQuery += fmt.Sprintf(" AND (name ILIKE $%d OR muscle_groups ILIKE $%d)", nameParam, groupParam)
		args = append(args, "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	limit := pagination.Limit
	if limit <= 0 {
		limit = 20
	}

	offset := pagination.Offset
	if offset < 0 {
		offset = 0
	}

	dataQuery := `SELECT exercise_id, name, muscle_groups, difficulty, equipment, type, default_sets, default_reps, rest_seconds ` + baseQuery +
		fmt.Sprintf(" ORDER BY name ASC OFFSET $%d LIMIT $%d", len(args)+1, len(args)+2)
	dataArgs := append(append([]interface{}{}, args...), offset, limit)

	rows, err := s.db.Query(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	exercises := make([]types.Exercise, 0)
	for rows.Next() {
		var exercise types.Exercise
		if err := rows.Scan(
			&exercise.ExerciseID,
			&exercise.Name,
			&exercise.MuscleGroups,
			&exercise.Difficulty,
			&exercise.Equipment,
			&exercise.Type,
			&exercise.DefaultSets,
			&exercise.DefaultReps,
			&exercise.RestSeconds,
		); err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	countQuery := `SELECT COUNT(*) ` + baseQuery
	countArgs := append([]interface{}{}, args...)

	var total int
	if err := s.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, err
	}

	totalPages := 0
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}

	response := &types.PaginatedResponse[types.Exercise]{
		Data:       exercises,
		TotalCount: total,
		TotalPages: totalPages,
		Page:       pagination.Page,
		PageSize:   limit,
	}
	return response, nil
}

func (s *Store) SearchExercises(ctx context.Context, query string, pagination types.PaginationParams) (*types.PaginatedResponse[types.Exercise], error) {
	q := `
		SELECT created_at, name, description, muscle_group, equipment, difficulty, id AS exercise_id
		FROM exercises
		WHERE name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%'
		ORDER BY created_at DESC
		OFFSET $2 LIMIT $3
	`

	rows, err := s.db.Query(ctx, q, query, pagination.Offset, pagination.Limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var exercises []types.Exercise
	for rows.Next() {
		var exercise types.Exercise
		if err := rows.Scan(
			&exercise.DefaultReps,
			&exercise.DefaultSets,
			&exercise.Difficulty,
			&exercise.Equipment,
			&exercise.ExerciseID,
			&exercise.MuscleGroups,
			&exercise.Name,
			&exercise.RestSeconds,
			&exercise.Type,
		); err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	countQ :=
		`SELECT COUNT(*) FROM exercises WHERE name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%'`
	var total int
	err = s.db.QueryRow(ctx, countQ, query).Scan(&total)
	if err != nil {
		return nil, err
	}

	totalPage := (total + pagination.Limit - 1) / pagination.Limit
	if pagination.Page > totalPage {
		exercises = []types.Exercise{}
	}

	response := &types.PaginatedResponse[types.Exercise]{
		Data:       exercises,
		TotalCount: total,
		TotalPages: totalPage,
		Page:       pagination.Page,
		PageSize:   pagination.Limit,
	}
	return response, nil
}

func (s *Store) GetExercisesByMuscleGroup(ctx context.Context, muscleGroup string) ([]types.Exercise, error) {
	q := `
		SELECT created_at, name, description, muscle_group, equipment, difficulty, id AS exercise_id
		FROM exercises
		WHERE muscle_group ILIKE $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(ctx, q, muscleGroup)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []types.Exercise
	for rows.Next() {
		var exercise types.Exercise
		if err := rows.Scan(
			&exercise.DefaultReps,
			&exercise.DefaultSets,
			&exercise.Difficulty,
			&exercise.Equipment,
			&exercise.ExerciseID,
			&exercise.MuscleGroups,
			&exercise.Name,
			&exercise.RestSeconds,
			&exercise.Type,
		); err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return exercises, nil
}

func (s *Store) GetExercisesByEquipment(ctx context.Context, equipment types.EquipmentType) ([]types.Exercise, error) {
	q := `
		SELECT created_at, name, description, muscle_group, equipment, difficulty, id AS exercise_id
		FROM exercises
		WHERE equipment = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(ctx, q, equipment)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var exercises []types.Exercise
	for rows.Next() {
		var exercise types.Exercise
		if err := rows.Scan(
			&exercise.DefaultReps,
			&exercise.DefaultSets,
			&exercise.Difficulty,
			&exercise.Equipment,
			&exercise.ExerciseID,
			&exercise.MuscleGroups,
			&exercise.Name,
			&exercise.RestSeconds,
			&exercise.Type,
		); err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return exercises, nil
}

func (s *Store) GetExercisesByDifficulty(ctx context.Context, difficulty types.FitnessLevel) ([]types.Exercise, error) {
	q := `
		SELECT created_at, name, description, muscle_group, equipment, difficulty, id AS exercise_id
		FROM exercises
		WHERE difficulty = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(ctx, q, difficulty)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var exercises []types.Exercise
	for rows.Next() {
		var exercise types.Exercise
		if err := rows.Scan(
			&exercise.DefaultReps,
			&exercise.DefaultSets,
			&exercise.Difficulty,
			&exercise.Equipment,
			&exercise.ExerciseID,
			&exercise.MuscleGroups,
			&exercise.Name,
			&exercise.RestSeconds,
			&exercise.Type,
		); err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return exercises, nil
}

func (s *Store) GetRecommendedExercises(ctx context.Context, userID int, count int) ([]types.Exercise, error) {
	q := `
		SELECT created_at, name, description, muscle_group, equipment, difficulty, id AS exercise_id
		FROM exercises
		ORDER BY RANDOM()
		LIMIT $1
	`

	rows, err := s.db.Query(ctx, q, count)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var exercises []types.Exercise
	for rows.Next() {
		var exercise types.Exercise
		if err := rows.Scan(
			&exercise.DefaultReps,
			&exercise.DefaultSets,
			&exercise.Difficulty,
			&exercise.Equipment,
			&exercise.ExerciseID,
			&exercise.MuscleGroups,
			&exercise.Name,
			&exercise.RestSeconds,
			&exercise.Type,
		); err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return exercises, nil
}

func (s *Store) BulkCreateExercises(ctx context.Context, exercises []types.ExerciseRequest) ([]types.Exercise, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()
	q := `
		INSERT INTO exercises (name, description, muscle_group, equipment, difficulty)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, name, description, muscle_group, equipment, difficulty, id AS exercise_id
	`

	var createdExercises []types.Exercise
	for _, exercise := range exercises {
		var createdExercise types.Exercise
		err := tx.QueryRow(ctx, q,
			exercise.Name,
			exercise.DefaultReps,
			exercise.DefaultSets,
			exercise.Difficulty,
			exercise.Equipment,
			exercise.RestSeconds,
			exercise.Type,
		).Scan(
			&createdExercise.Name,
			&createdExercise.Equipment,
			&createdExercise.Difficulty,
			&createdExercise.ExerciseID,
			&createdExercise.Type,
			&createdExercise.MuscleGroups,
			&createdExercise.RestSeconds,
			&createdExercise.DefaultReps,
			&createdExercise.DefaultSets,
		)
		if err != nil {
			return nil, err
		}
		createdExercises = append(createdExercises, createdExercise)
	}
	return createdExercises, nil
}
