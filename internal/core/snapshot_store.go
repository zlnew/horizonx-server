package core

import "sync"

type SnapshotStore struct {
	mu       sync.RWMutex
	snapshot Metrics
}

func NewSnapshotStore() *SnapshotStore {
	return &SnapshotStore{}
}

func (s *SnapshotStore) Set(m Metrics) {
	s.mu.Lock()
	s.snapshot = m
	s.mu.Unlock()
}

func (s *SnapshotStore) Get() Metrics {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.snapshot
}
