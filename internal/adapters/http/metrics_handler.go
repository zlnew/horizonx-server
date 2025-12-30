package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"horizonx-server/internal/domain"

	"github.com/google/uuid"
)

type MetricsHandler struct {
	svc domain.MetricsService
}

func NewMetricsHandler(svc domain.MetricsService) *MetricsHandler {
	return &MetricsHandler{
		svc: svc,
	}
}

func (h *MetricsHandler) Ingest(w http.ResponseWriter, r *http.Request) {
	var metrics domain.Metrics

	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		JSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.svc.Ingest(metrics); err != nil {
		JSONError(w, http.StatusInternalServerError, "failed to process metrics")
		return
	}

	JSONSuccess(w, http.StatusCreated, APIResponse{
		Message: "Metrics received",
		Data:    metrics,
	})
}

func (h *MetricsHandler) Latest(w http.ResponseWriter, r *http.Request) {
	serverID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		JSONError(w, http.StatusNotFound, "server not found")
		return
	}

	metrics, err := h.svc.Latest(serverID)
	if err != nil {
		if errors.Is(err, domain.ErrMetricsNotFound) {
			JSONError(w, http.StatusNotFound, "metrics not found")
			return
		}

		JSONError(w, http.StatusInternalServerError, "failed to get latest metrics")
		return
	}

	JSONSuccess(w, http.StatusOK, APIResponse{
		Message: "OK",
		Data:    metrics,
	})
}

func (h *MetricsHandler) CPUUsageHistory(w http.ResponseWriter, r *http.Request) {
	serverID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		JSONError(w, http.StatusNotFound, "server not found")
		return
	}

	data, err := h.svc.CPUUsageHistory(serverID)
	if err != nil {
		if errors.Is(err, domain.ErrMetricsNotFound) {
			JSONError(w, http.StatusNotFound, "cpu usage history not found")
			return
		}

		JSONError(w, http.StatusInternalServerError, "failed to get cpu usage history")
		return
	}

	JSONSuccess(w, http.StatusOK, APIResponse{
		Message: "OK",
		Data:    data,
	})
}

func (h *MetricsHandler) NetSpeedHistory(w http.ResponseWriter, r *http.Request) {
	serverID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		JSONError(w, http.StatusNotFound, "server not found")
		return
	}

	data, err := h.svc.NetSpeedHistory(serverID)
	if err != nil {
		if errors.Is(err, domain.ErrMetricsNotFound) {
			JSONError(w, http.StatusNotFound, "net speed history not found")
			return
		}

		JSONError(w, http.StatusInternalServerError, "failed to get net speed history")
		return
	}

	JSONSuccess(w, http.StatusOK, APIResponse{
		Message: "OK",
		Data:    data,
	})
}
