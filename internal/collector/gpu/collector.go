// Package gpu
package gpu

import "context"

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) Collect(ctx context.Context) (any, error) {
	cards := detectGPUs()
	var out []GPUMetric

	for i, card := range cards {
		vendor := readVendor(card)
		model := readModel(card)

		temp := readTemperature(card)
		usage := readCoreUsage(card)
		vramTotal, vramUsed, vramPercent := readVRAM(card)
		power := readPower(card)
		fanSpeed := readFanSpeedPercent(card)

		out = append(out, GPUMetric{
			ID:               i,
			Card:             card,
			Vendor:           vendor,
			Model:            model,
			Temperature:      temp,
			CoreUsagePercent: usage,
			VRAMTotalGB:      vramTotal,
			VRAMUsedGB:       vramUsed,
			VRAMPercent:      vramPercent,
			PowerWatt:        power,
			FanSpeedPercent:  fanSpeed,
		})
	}

	return out, nil
}
