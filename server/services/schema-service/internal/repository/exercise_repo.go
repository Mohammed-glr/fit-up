package repository

import (
	"context"

	"github.com/tdmdh/fit-up-server/services/schema-service/internal/types"
)

func (s *Store) CreateExercise(ctx context.Context, exercise *types.ExerciseRequest) (*types.Exercise, error) {
	q := `
		INSTERT INTO exercises (name, description, muscle_group, equipment, difficulty)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, name, description, muscle_group, equipment, difficulty, id AS exercise_id
	`

	var createdExercise types.Exercise
	err := s.db.QueryRow(ctx, q,
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

	return &createdExercise, nil
}


func (s *Store) GetExerciseByID(ctx context.Context, exerciseID int) (*types.Exercise, error) {
	q := `
		SELECT created_at, name, description, muscle_group, equipment, difficulty, id AS exercise_id
		FROM exercises
		WHERE id = $1
	`

	var exercise types.Exercise
	err := s.db.QueryRow(ctx, q, exerciseID).Scan(
		&exercise.DefaultReps,
		&exercise.DefaultSets,
		&exercise.Difficulty,
		&exercise.Equipment,
		&exercise.ExerciseID,
		&exercise.MuscleGroups,
		&exercise.Name,
		&exercise.RestSeconds,
		&exercise.Type,
	)
	if err != nil {
		return nil, err
	}

	return &exercise, nil
}

func (s *Store) UpdateExercise(ctx context.Context, exerciseID int, exercise *types.ExerciseRequest) (*types.Exercise, error) {
	q := `
		UPDATE exercises
		SET name = $1, description = $2, muscle_group = $3, equipment = $4, difficulty = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING created_at, name, description, muscle_group, equipment, difficulty, id AS exercise_id
	`

	var updatedExercise types.Exercise
	err := s.db.QueryRow(ctx, q,
		exercise.DefaultReps,
		exercise.DefaultSets,
		exercise.Difficulty,
		exercise.Equipment,
		exercise.MuscleGroups,
		exercise.Name,
		exercise.RestSeconds,
		exercise.Type,
		exerciseID,
	).Scan(
		&updatedExercise.Name,
		&updatedExercise.Equipment,
		&updatedExercise.Difficulty,
		&updatedExercise.ExerciseID,
		&updatedExercise.Type,
		&updatedExercise.MuscleGroups,
		&updatedExercise.RestSeconds,
		&updatedExercise.DefaultReps,
		&updatedExercise.DefaultSets,
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
