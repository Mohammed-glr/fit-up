package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Users() UserRepo {
	return s
}

func (s *Store) Exercises() ExerciseRepo {
	return s
}

func (s *Store) Templates() WorkoutTemplateRepo {
	return s
}

func (s *Store) Schemas() WeeklySchemaRepo {
	return s
}

func (s *Store) Workouts() WorkoutRepo {
	return s
}

func (s *Store) WorkoutExercises() WorkoutExerciseRepo {
	return s
}

func (s *Store) Progress() ProgressRepo {
	return s
}

func (s *Store) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = fn(ctx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
