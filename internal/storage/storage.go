package storage

import (
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Item struct {
	Value    []byte
	ExpireAt time.Time
}

type Storage struct {
	mu     sync.RWMutex
	data   map[string]*Item
	logger zerolog.Logger
}

func New(logger zerolog.Logger) *Storage {
	return &Storage{
		data:   make(map[string]*Item),
		logger: logger,
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

func (s *Storage) Delete(keys ...string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var deleted int64 = 0

	for _, key := range keys {
		if _, exists := s.data[key]; exists {
			delete(s.data, key)
			deleted++
		}
	}

	if deleted > 0 {
		s.logger.Debug().Int64("deleted", deleted).Msg("keys deleted")
	}

	return deleted, nil
}

func (s *Storage) StartCleaner(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			s.cleanupExpired()
		}
	}()
}

func (s *Storage) cleanupExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	var expired int64
	for key, item := range s.data {
		if !item.ExpireAt.IsZero() && now.After(item.ExpireAt) {
			delete(s.data, key)
			expired++
		}
	}

	if expired > 0 {
		s.logger.Debug().Int64("expired", expired).Msg("expired keys cleaned")
	}
}
