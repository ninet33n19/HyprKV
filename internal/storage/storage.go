package storage

import (
	"sync"
	"time"
)

type Item struct {
	Value    []byte
	ExpireAt time.Time
}

type Storage struct {
	mu   sync.RWMutex
	data map[string]*Item
}

func New() *Storage {
	return &Storage{
		data: make(map[string]*Item),
	}
}

func (s *Storage) Set(key string, value []byte, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item := &Item{
		Value: value,
	}
	if ttl > 0 {
		item.ExpireAt = time.Now().Add(ttl)
	}

	s.data[key] = item
}

func (s *Storage) Get(key string) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.data[key]
	if !exists {
		return nil, false
	}

	if !item.ExpireAt.IsZero() && time.Now().After(item.ExpireAt) {
		return nil, false
	}

	return item.Value, true
}
