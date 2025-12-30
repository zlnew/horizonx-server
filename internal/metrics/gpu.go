package metrics

import (
	"horizonx-server/internal/domain"
	"horizonx-server/internal/system"
)

func calculateGPUMetric(card string, vendor string, m *system.GPUMetrics) domain.GPUMetric {
	var gpu domain.GPUMetric

	if m == nil {
		return gpu
	}

	gpu.Card = card
	gpu.Vendor = vendor

	gpu.Temperature.Raw = float64(m.TemperatureC)
	gpu.CoreUsagePercent.Raw = float64(m.UtilizationGPU)
	gpu.FrequencyMhz.Raw = float64(m.ClockMHz)
	gpu.PowerWatt.Raw = m.PowerDrawW

	if m.MemTotalMB > 0 {
		gpu.VRAMTotalGB = float64(m.MemTotalMB) / 1024
		gpu.VRAMUsedGB = float64(m.MemUsedMB) / 1024
		gpu.VRAMPercent = float64(m.MemUsedMB) / float64(m.MemTotalMB) * 100
	}

	return gpu
}
