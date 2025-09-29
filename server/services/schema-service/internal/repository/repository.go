package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)



type Store struct {
	Users            UserRepository
	Exercises        ExerciseRepository
	WorkoutTemplates WorkoutTemplateRepository
	Workouts         WorkoutRepository
	WorkoutExercises WorkoutExerciseRepository
	Progress         ProgressRepository
	Schema           SchemaRepository
	db               *pgxpool.Pool

}

func NewStore(
	users UserRepository,
	exercises ExerciseRepository,
	templates WorkoutTemplateRepository,
	schemas SchemaRepository,
	workouts WorkoutRepository,
	workoutExercises WorkoutExerciseRepository,
	progress ProgressRepository,
) *Store {
	return &Store{
		Users:            users,
		Exercises:        exercises,
		WorkoutTemplates: templates,
		Schema:           schemas,
		Workouts:         workouts,
		WorkoutExercises: workoutExercises,
		Progress:         progress,
	}
}

