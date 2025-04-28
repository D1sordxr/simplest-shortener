package infra

import (
	"context"
	"sync"
)

// mainUrl: shortenedURL
type shortenedURLs map[string]string

type Storage struct {
	shortenedURLs
	mu *sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		shortenedURLs: make(shortenedURLs),
		mu:            &sync.RWMutex{},
	}
}

func (s *Storage) Set(_ context.Context, key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.shortenedURLs[key] = value
}

func (s *Storage) Get(_ context.Context, key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.shortenedURLs[key]
	return value, ok
}

func (s *Storage) Delete(_ context.Context, key string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	delete(s.shortenedURLs, key)
}
