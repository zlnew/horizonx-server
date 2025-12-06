package gpu

import (
	"context"

	"horizonx-server/internal/logger"
	"horizonx-server/pkg"
)

func NewCollector(log logger.Logger) *Collector {
	return &Collector{
		log:         log,
		powerEMA:    make(map[string]*pkg.EMA),
		usageEMA:    make(map[string]*pkg.EMA),
		fanSpeedEMA: make(map[string]*pkg.EMA),
	}
}

func (c *Collector) Collect(ctx context.Context) ([]GPUMetric, error) {
	cards := c.detectGPUs()
	var outputs []GPUMetric

	for i, card := range cards {
		if _, ok := c.powerEMA[card]; !ok {
			c.powerEMA[card] = pkg.NewEMA(0.3)
		}
		if _, ok := c.usageEMA[card]; !ok {
			c.usageEMA[card] = pkg.NewEMA(0.5)
		}
		if _, ok := c.fanSpeedEMA[card]; !ok {
			c.fanSpeedEMA[card] = pkg.NewEMA(0.5)
		}

		vendor := c.readVendor(card)
		model := c.readModel(card)

		temp := c.readTemperature(card)
		usage := c.readCoreUsage(card)
		vramTotal, vramUsed, vramPercent := c.readVRAM(card)
		powerWatt := c.readPower(card)
		fanSpeed := c.readFanSpeedPercent(card)

		c.usageEMA[card].Add(usage)
		c.powerEMA[card].Add(powerWatt)
		c.fanSpeedEMA[card].Add(fanSpeed)

		outputs = append(outputs, GPUMetric{
			ID:               i,
			Card:             card,
			Vendor:           vendor,
			Model:            model,
			Temperature:      temp,
			CoreUsagePercent: c.usageEMA[card].Value(),
			VRAMTotalGB:      vramTotal,
			VRAMUsedGB:       vramUsed,
			VRAMPercent:      vramPercent,
			PowerWatt:        c.powerEMA[card].Value(),
			FanSpeedPercent:  c.fanSpeedEMA[card].Value(),
		})
	}

	return outputs, nil
}
