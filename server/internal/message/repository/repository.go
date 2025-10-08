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

func (s *Store) Conversations() ConversationRepo {
	return s
}

func (s *Store) Messages() MessageRepo {
	return s
}

func (s *Store) ReadStatus() MessageReadStatusRepo {
	return s
}
func (s *Store) Attachments() MessageAttachmentRepo {
	return s
}

func (s *Store) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(ctx); err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
