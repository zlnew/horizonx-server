package snapshot

import "horizonx-server/pkg/types"

type MetricsStore struct {
	Store[types.Metrics]
}

func NewMetricsStore() *MetricsStore {
	return &MetricsStore{}
}
