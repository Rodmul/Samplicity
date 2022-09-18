package store

import (
	"github.com/jmoiron/sqlx"
	"log"
)

type Store struct {
	db                       *DB
	logger                   *log.Logger
	userRepository           Profiler
	sampleRepository         Sampler
	createdSamplesRepository Creator
	likedSamplesRepository   Liker
}

func New(db *sqlx.DB, l *log.Logger) *Store {
	return &Store{
		db:     &DB{db},
		logger: l,
	}
}

func (s *Store) LikedSample() Liker {
	if s.likedSamplesRepository != nil {
		return s.likedSamplesRepository
	}

	s.likedSamplesRepository = &LikedSamplesRepository{
		store: s,
	}
	return s.likedSamplesRepository
}

func (s *Store) CreatedSample() Creator {
	if s.createdSamplesRepository != nil {
		return s.createdSamplesRepository
	}

	s.createdSamplesRepository = &CreatedSamplesRepository{
		store: s,
	}

	return s.createdSamplesRepository
}

func (s *Store) User() Profiler {
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
