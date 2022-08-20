package storage

import (
	"errors"
	"sync"
)

var (
	ErrNotFound = errors.New("key not found")
)

type Storage struct {
	db sync.Map
}

func NewStorage() *Storage {
	return &Storage{
		db: sync.Map{},
	}
}

func (s *Storage) Set(key string, value any) error {
	s.db.Store(key, value)
	return nil
}

func (s *Storage) Get(key string) (any, error) {
	val, ok := s.db.Load(key)
	if ok {
		return val, nil
	}
	return nil, ErrNotFound
}
