package repository

import (

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
