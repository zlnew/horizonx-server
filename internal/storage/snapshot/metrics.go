package snapshot

import "horizonx-server/internal/domain"

type MetricsStore struct {
	Store[domain.Metrics]
}

func NewMetricsStore() *MetricsStore {
	return &MetricsStore{}
}
