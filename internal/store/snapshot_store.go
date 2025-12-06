package store

import (
	"sync"

	"horizonx-server/pkg/types"
)

type SnapshotStore struct {
	mu       sync.RWMutex
	snapshot types.Metrics
}

func NewSnapshotStore() *SnapshotStore {
	return &SnapshotStore{}
}

func (s *SnapshotStore) Set(m types.Metrics) {
	s.mu.Lock()
	s.snapshot = m
	s.mu.Unlock()
}

func (s *SnapshotStore) Get() types.Metrics {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.snapshot
}
