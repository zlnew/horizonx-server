// Package handler
package handler

import (
	"encoding/json"
	"net/http"

	"horizonx-server/internal/storage/snapshot"
)

func HandleMetrics(ms *snapshot.MetricsStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := ms.Get()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
