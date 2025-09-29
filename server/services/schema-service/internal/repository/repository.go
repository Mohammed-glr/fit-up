package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)



type Store struct {
	Users            UserRepo
	Exercises        ExerciseRepo
	WorkoutTemplates WorkoutTemplateRepo
	Workouts         WorkoutRepo
	WorkoutExercises WorkoutExerciseRepo
	Progress         ProgressRepo
	Schema           SchemaRepo
	db               *pgxpool.Pool

}

func NewStore(
	users UserRepo,
	exercises ExerciseRepo,
	templates WorkoutTemplateRepo,
	schemas SchemaRepo,
	workouts WorkoutRepo,
	workoutExercises WorkoutExerciseRepo,
	progress ProgressRepo,
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

