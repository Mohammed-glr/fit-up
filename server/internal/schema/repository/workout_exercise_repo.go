package repository

import (
	"context"

	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

func (s *Store) CreateWorkoutExercise(ctx context.Context, workoutExercise *types.WorkoutExerciseRequest) (*types.WorkoutExercise, error) {
	q := `
		INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, rest_seconds)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING we_id, workout_id, exercise_id, sets, reps, rest_seconds
	`

	var we types.WorkoutExercise
	err := s.db.QueryRow(ctx, q,
		workoutExercise.WorkoutID,
		workoutExercise.ExerciseID,
		workoutExercise.Sets,
		workoutExercise.Reps,
		workoutExercise.RestSeconds,
	).Scan(
		&we.WeID,
		&we.WorkoutID,
		&we.ExerciseID,
		&we.Sets,
		&we.Reps,
		&we.RestSeconds,
	)
	if err != nil {
		return nil, err
	}

	return &we, nil
}

func (s *Store) GetWorkoutExerciseByID(ctx context.Context, weID int) (*types.WorkoutExercise, error) {
	q := `
		SELECT we_id, workout_id, exercise_id, sets, reps, rest_seconds
		FROM workout_exercises
		WHERE we_id = $1
	`

	var we types.WorkoutExercise
	err := s.db.QueryRow(ctx, q, weID).Scan(
		&we.WeID,
		&we.WorkoutID,
		&we.ExerciseID,
		&we.Sets,
		&we.Reps,
		&we.RestSeconds,
	)
	if err != nil {
		return nil, err
	}

	return &we, nil
}

func (s *Store) UpdateWorkoutExercise(ctx context.Context, weID int, workoutExercise *types.WorkoutExerciseRequest) (*types.WorkoutExercise, error) {
	q := `
		UPDATE workout_exercises
		SET workout_id = $1, exercise_id = $2, sets = $3, reps = $4, rest_seconds = $5
		WHERE we_id = $6
		RETURNING we_id, workout_id, exercise_id, sets, reps, rest_seconds
	`

	var we types.WorkoutExercise
	err := s.db.QueryRow(ctx, q,
		workoutExercise.WorkoutID,
		workoutExercise.ExerciseID,
		workoutExercise.Sets,
		workoutExercise.Reps,
		workoutExercise.RestSeconds,
		weID,
	).Scan(
		&we.WeID,
		&we.WorkoutID,
		&we.ExerciseID,
		&we.Sets,
		&we.Reps,
		&we.RestSeconds,
	)
	if err != nil {
		return nil, err
	}

	return &we, nil
}

func (s *Store) DeleteWorkoutExercise(ctx context.Context, weID int) error {
	q := `
		DELETE FROM workout_exercises
		WHERE we_id = $1
	`
	_, err := s.db.Exec(ctx, q, weID)
	return err
}

func (s *Store) GetWorkoutExercisesByWorkoutID(ctx context.Context, workoutID int) ([]types.WorkoutExercise, error) {
	q := `
		SELECT we_id, workout_id, exercise_id, sets, reps, rest_seconds
		FROM workout_exercises
		WHERE workout_id = $1
	`

	rows, err := s.db.Query(ctx, q, workoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workoutExercises []types.WorkoutExercise
	for rows.Next() {
		var we types.WorkoutExercise
		err := rows.Scan(
			&we.WeID,
			&we.WorkoutID,
			&we.ExerciseID,
			&we.Sets,
			&we.Reps,
			&we.RestSeconds,
		)
		if err != nil {
			return nil, err
		}
		workoutExercises = append(workoutExercises, we)
	}

	return workoutExercises, nil
}

func (s *Store) BulkCreateWorkoutExercisesForWorkout(ctx context.Context, workoutID int, exercises []types.WorkoutExerciseRequest) ([]types.WorkoutExercise, error) {
	var results []types.WorkoutExercise

	for _, exercise := range exercises {
		exercise.WorkoutID = workoutID

		we, err := s.CreateWorkoutExercise(ctx, &exercise)
		if err != nil {
			return nil, err
		}
		results = append(results, *we)
	}

	return results, nil
}

func (s *Store) BulkUpdateWorkoutExercisesForWorkout(ctx context.Context, workoutID int, exercises []types.WorkoutExercise) error {
	for _, exercise := range exercises {
		_, err := s.UpdateWorkoutExercise(ctx, exercise.WeID, &types.WorkoutExerciseRequest{
			WorkoutID:   exercise.WorkoutID,
			ExerciseID:  exercise.ExerciseID,
			Sets:        exercise.Sets,
			Reps:        exercise.Reps,
			RestSeconds: exercise.RestSeconds,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) DeleteAllWorkoutExercisesForWorkout(ctx context.Context, workoutID int) error {
	q := `
		DELETE FROM workout_exercises
		WHERE workout_id = $1
	`
	_, err := s.db.Exec(ctx, q, workoutID)
	return err
}

func (s *Store) GetMostUsedExercises(ctx context.Context, limit int) ([]types.Exercise, error) {
	q := `
		SELECT e.exercise_id, e.name, e.muscle_groups, e.difficulty, e.equipment, e.type, e.default_sets, e.default_reps, e.rest_seconds,
		       COUNT(we.exercise_id) as usage_count
		FROM exercises e
		JOIN workout_exercises we ON e.exercise_id = we.exercise_id
		GROUP BY e.exercise_id, e.name, e.muscle_groups, e.difficulty, e.equipment, e.type, e.default_sets, e.default_reps, e.rest_seconds
		ORDER BY usage_count DESC
		LIMIT $1
	`

	rows, err := s.db.Query(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []types.Exercise
	for rows.Next() {
		var exercise types.Exercise
		var usageCount int
		err := rows.Scan(
			&exercise.ExerciseID,
			&exercise.Name,
			&exercise.MuscleGroups,
			&exercise.Difficulty,
			&exercise.Equipment,
			&exercise.Type,
			&exercise.DefaultSets,
			&exercise.DefaultReps,
			&exercise.RestSeconds,
			&usageCount,
		)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}

	return exercises, nil
}

func (s *Store) GetExerciseUsageStats(ctx context.Context, exerciseID int) (map[string]interface{}, error) {
	q := `
		SELECT 
			COUNT(*) as total_usage,
			AVG(sets) as avg_sets,
			AVG(CAST(rest_seconds AS FLOAT)) as avg_rest_seconds
		FROM workout_exercises
		WHERE exercise_id = $1
	`

	var totalUsage int
	var avgSets float64
	var avgRestSeconds float64

	err := s.db.QueryRow(ctx, q, exerciseID).Scan(&totalUsage, &avgSets, &avgRestSeconds)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_usage":      totalUsage,
		"avg_sets":         avgSets,
		"avg_rest_seconds": avgRestSeconds,
	}

	return stats, nil
}
