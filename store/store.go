package store

import (
	"github.com/jmoiron/sqlx"
	"log"
)

type Store struct {
	db               *DB
	logger           *log.Logger
	userRepository   Userter
	sampleRepository Sampler
}

func New(db *sqlx.DB, l *log.Logger) *Store {
	return &Store{
		db:     &DB{db},
		logger: l,
	}
}

func (s *Store) User() Userter {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

func (s *Store) Sample() Sampler {
	if s.sampleRepository != nil {
		return s.sampleRepository
	}

	s.sampleRepository = &SampleRepository{
		store: s,
	}

	return s.sampleRepository
}
