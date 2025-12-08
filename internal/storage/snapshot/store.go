// Package snapshot
package snapshot

import "sync"

type Store[T any] struct {
	mu   sync.RWMutex
	data T
}

func (s *Store[T]) Set(v T) {
	s.mu.Lock()
	s.data = v
	s.mu.Unlock()
}

func (s *Store[T]) Get() T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data
}
