package repository

import (
	"context"

	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

func (s *Store) CreateWorkout(ctx context.Context, workout *types.WorkoutRequest) (*types.Workout, error) {
	q := `
		INSERT INTO workouts (schema_id, day_of_week, focus)
		VALUES ($1, $2, $3)
		RETURNING workout_id, schema_id, day_of_week, focus
	`
	row := s.db.QueryRow(ctx, q,
		workout.SchemaID,
		workout.DayOfWeek,
		workout.Focus,
	)

	var w types.Workout
	err := row.Scan(
		&w.WorkoutID,
		&w.SchemaID,
		&w.DayOfWeek,
		&w.Focus,
	)
	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (s *Store) GetWorkoutByID(ctx context.Context, workoutID int) (*types.Workout, error) {
	q := `
		SELECT workout_id, schema_id, day_of_week, focus
		FROM workouts
		WHERE workout_id = $1
	`
	row := s.db.QueryRow(ctx, q, workoutID)
	var w types.Workout
	err := row.Scan(
		&w.WorkoutID,
		&w.SchemaID,
		&w.DayOfWeek,
		&w.Focus,
	)
	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (s *Store) UpdateWorkout(ctx context.Context, workoutID int, workout *types.WorkoutRequest) (*types.Workout, error) {
	q := `
		UPDATE workouts
		SET schema_id = $1, day_of_week = $2, focus = $3
		WHERE workout_id = $4
		RETURNING workout_id, schema_id, day_of_week, focus
	`

	row := s.db.QueryRow(ctx, q,
		workout.SchemaID,
		workout.DayOfWeek,
		workout.Focus,
		workoutID,
	)
	var w types.Workout
	err := row.Scan(
		&w.WorkoutID,
		&w.SchemaID,
		&w.DayOfWeek,
		&w.Focus,
	)
	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (s *Store) DeleteWorkout(ctx context.Context, workoutID int) error {
	q := `
		DELETE FROM workouts
		WHERE workout_id = $1
	`
	_, err := s.db.Exec(ctx, q, workoutID)
	return err
}

func (s *Store) GetWorkoutsBySchemaID(ctx context.Context, schemaID int) ([]types.Workout, error) {
	q := `
		SELECT workout_id, schema_id, day_of_week, focus
		FROM workouts
		WHERE schema_id = $1
		ORDER BY day_of_week
	`
	rows, err := s.db.Query(ctx, q, schemaID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var workouts []types.Workout
	for rows.Next() {
		var w types.Workout
		err := rows.Scan(
			&w.WorkoutID,
			&w.SchemaID,
			&w.DayOfWeek,
			&w.Focus,
		)
		if err != nil {
			return nil, err
		}
		workouts = append(workouts, w)
	}

	return workouts, nil
}

func (s *Store) GetWorkoutBySchemaAndDay(ctx context.Context, schemaID int, dayOfWeek int) (*types.Workout, error) {
	q := `
		SELECT workout_id, schema_id, day_of_week, focus
		FROM workouts
		WHERE schema_id = $1 AND day_of_week = $2
	`
	row := s.db.QueryRow(ctx, q, schemaID, dayOfWeek)
	var w types.Workout
	err := row.Scan(
		&w.WorkoutID,
		&w.SchemaID,
		&w.DayOfWeek,
		&w.Focus,
	)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (s *Store) GetWorkoutsBySchemaAndDay(ctx context.Context, schemaID int, dayOfWeek int) ([]types.Workout, error) {
	q := `
		SELECT workout_id, schema_id, day_of_week, focus
		FROM workouts
		WHERE schema_id = $1 AND day_of_week = $2
	`
	rows, err := s.db.Query(ctx, q, schemaID, dayOfWeek)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var workouts []types.Workout
	for rows.Next() {
		var w types.Workout
		err := rows.Scan(
			&w.WorkoutID,
			&w.SchemaID,
			&w.DayOfWeek,
			&w.Focus,
		)
		if err != nil {
			return nil, err
		}
		workouts = append(workouts, w)
	}

	return workouts, nil
}

func (s *Store) GetWorkoutWithExercises(ctx context.Context, workoutID int) (*types.WorkoutWithExercises, error) {
	workoutQuery := `
		SELECT workout_id, schema_id, day_of_week, focus
		FROM workouts
		WHERE workout_id = $1
	`

	var workout types.Workout
	err := s.db.QueryRow(ctx, workoutQuery, workoutID).Scan(
		&workout.WorkoutID,
		&workout.SchemaID,
		&workout.DayOfWeek,
		&workout.Focus,
	)
	if err != nil {
		return nil, err
	}

	exercisesQuery := `
		SELECT we.we_id, we.sets, we.reps, we.rest_seconds,
		       e.exercise_id, e.name, e.muscle_groups, e.difficulty, e.equipment, e.type, e.default_sets, e.default_reps, e.rest_seconds
		FROM workout_exercises we
		JOIN exercises e ON we.exercise_id = e.exercise_id
		WHERE we.workout_id = $1
	`

	rows, err := s.db.Query(ctx, exercisesQuery, workoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []types.WorkoutExerciseDetail
	for rows.Next() {
		var detail types.WorkoutExerciseDetail
		var exercise types.Exercise

		err := rows.Scan(
			&detail.WeID,
			&detail.Sets,
			&detail.Reps,
			&detail.RestSeconds,
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

		detail.Exercise = types.ExerciseResponse{
			ExerciseID:   exercise.ExerciseID,
			Name:         exercise.Name,
			MuscleGroups: []string{exercise.MuscleGroups},
			Difficulty:   exercise.Difficulty,
			Equipment:    exercise.Equipment,
			Type:         exercise.Type,
			DefaultSets:  exercise.DefaultSets,
			DefaultReps:  exercise.DefaultReps,
			RestSeconds:  exercise.RestSeconds,
		}

		exercises = append(exercises, detail)
	}

	return &types.WorkoutWithExercises{
		WorkoutID: workout.WorkoutID,
		SchemaID:  workout.SchemaID,
		DayOfWeek: workout.DayOfWeek,
		Focus:     workout.Focus,
		Exercises: exercises,
	}, nil
}

func (s *Store) GetSchemaWithAllWorkouts(ctx context.Context, schemaID int) (*types.WeeklySchemaWithWorkouts, error) {
	schemaQuery := `
		SELECT schema_id, user_id, week_start, active
		FROM weekly_schemas
		WHERE schema_id = $1
	`

	var schema types.WeeklySchema
	err := s.db.QueryRow(ctx, schemaQuery, schemaID).Scan(
		&schema.SchemaID,
		&schema.UserID,
		&schema.WeekStart,
		&schema.Active,
	)
	if err != nil {
		return nil, err
	}

	workouts, err := s.GetWorkoutsBySchemaID(ctx, schemaID)
	if err != nil {
		return nil, err
	}

	var workoutsWithExercises []types.WorkoutWithExercises
	for _, workout := range workouts {
		workoutWithExercises, err := s.GetWorkoutWithExercises(ctx, workout.WorkoutID)
		if err != nil {
			workoutsWithExercises = append(workoutsWithExercises, types.WorkoutWithExercises{
				WorkoutID: workout.WorkoutID,
				SchemaID:  workout.SchemaID,
				DayOfWeek: workout.DayOfWeek,
				Focus:     workout.Focus,
				Exercises: []types.WorkoutExerciseDetail{},
			})
		} else {
			workoutsWithExercises = append(workoutsWithExercises, *workoutWithExercises)
		}
	}

	return &types.WeeklySchemaWithWorkouts{
		SchemaID:  schema.SchemaID,
		UserID:    schema.UserID,
		WeekStart: schema.WeekStart,
		Active:    schema.Active,
		Workouts:  workoutsWithExercises,
	}, nil
}

func (s *Store) BulkCreateWorkoutsForSchema(ctx context.Context, schemaID int, workouts []types.WorkoutRequest) ([]types.Workout, error) {
	var results []types.Workout

	for _, workoutReq := range workouts {
		workoutReq.SchemaID = schemaID

		workout, err := s.CreateWorkout(ctx, &workoutReq)
		if err != nil {
			return nil, err
		}
		results = append(results, *workout)
	}

	return results, nil
}
