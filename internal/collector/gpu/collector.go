package gpu

import (
	"context"

	"zlnew/monitor-agent/internal/core"
	"zlnew/monitor-agent/internal/infra/logger"
)

func NewCollector(log logger.Logger) *Collector {
	return &Collector{
		log:      log,
		powerEMA: make(map[string]*core.EMA),
	}
}

func (c *Collector) Collect(ctx context.Context) ([]GPUMetric, error) {
	cards := c.detectGPUs()
	var outputs []GPUMetric

	for i, card := range cards {
		if _, ok := c.powerEMA[card]; !ok {
			c.powerEMA[card] = core.NewEMA(0.3)
		}

		vendor := c.readVendor(card)
		model := c.readModel(card)

		temp := c.readTemperature(card)
		usage := c.readCoreUsage(card)
		vramTotal, vramUsed, vramPercent := c.readVRAM(card)
		powerWatt := c.readPower(card)
		fanSpeed := c.readFanSpeedPercent(card)

		outputs = append(outputs, GPUMetric{
			ID:               i,
			Card:             card,
			Vendor:           vendor,
			Model:            model,
			Temperature:      temp,
			CoreUsagePercent: usage,
			VRAMTotalGB:      vramTotal,
			VRAMUsedGB:       vramUsed,
			VRAMPercent:      vramPercent,
			PowerWatt:        powerWatt,
			FanSpeedPercent:  fanSpeed,
		})
	}

	return outputs, nil
}
