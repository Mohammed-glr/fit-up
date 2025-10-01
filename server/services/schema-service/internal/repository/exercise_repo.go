package repository

import (
	"context"

	"github.com/tdmdh/fit-up-server/services/schema-service/internal/types"
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
	q := `
		SELECT created_at, name, description, muscle_group, equipment, difficulty, id AS exercise_id
		FROM exercises
		WHERE ($1::TEXT IS NULL OR muscle_group ILIKE $1)
		AND ($2::TEXT IS NULL OR equipment = $2)
		AND ($3::TEXT IS NULL OR difficulty = $3)
		ORDER BY created_at DESC
		OFFSET $4 LIMIT $5
	`

	rows, err := s.db.Query(ctx, q,
		filter.Difficulty,
		filter.Equipment,
		filter.MuscleGroups,
		filter.Search,
		filter.Type,
		pagination.Offset,
		pagination.Limit,
		pagination.PageSize,
		pagination.Page,
	)
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
