package metrics

import (
	"fmt"
	"math"
	"time"

	"horizonx-server/internal/domain"
)

type EMA struct {
	tau   float64
	value float64
	last  time.Time
	init  bool
}

func NewEMA(tau time.Duration) *EMA {
	return &EMA{
		tau: tau.Seconds(),
	}
}

func (e *EMA) Update(x float64, now time.Time) float64 {
	if !e.init {
		e.value = x
		e.last = now
		e.init = true
		return e.value
	}

	dt := now.Sub(e.last).Seconds()
	if dt <= 0 {
		return e.value
	}

	alpha := 1 - math.Exp(-dt/e.tau)
	e.value = alpha*x + (1-alpha)*e.value
	e.last = now

	return e.value
}

func getOrInitEMA(m map[string]*EMA, key string, tau time.Duration) *EMA {
	if e, ok := m[key]; ok {
		return e
	}
	e := NewEMA(tau)
	m[key] = e
	return e
}

func (c *Collector) ApplyEMA(m *domain.Metrics) {
	now := m.RecordedAt

	// CPU
	m.CPU.Usage.EMA = c.cpuUsageEMA.Update(m.CPU.Usage.Raw, now)
	m.CPU.Frequency.EMA = c.cpuFreqEMA.Update(m.CPU.Frequency.Raw, now)
	m.CPU.PowerWatt.EMA = c.cpuPowerEMA.Update(m.CPU.PowerWatt.Raw, now)
	m.CPU.Temperature.EMA = c.cpuTempEMA.Update(m.CPU.Temperature.Raw, now)

	for i := range m.CPU.PerCore {
		key := fmt.Sprintf("core-%d", i)
		ema := getOrInitEMA(c.cpuPerCoreEMA, key, 15*time.Second)
		m.CPU.PerCore[i].EMA = ema.Update(m.CPU.PerCore[i].Raw, now)
	}

	// GPU
	for i := range m.GPU {
		g := &m.GPU[i]
		key := g.Card

		g.CoreUsagePercent.EMA = getOrInitEMA(c.gpuUsageEMA, key, 10*time.Second).Update(g.CoreUsagePercent.Raw, now)
		g.FrequencyMhz.EMA = getOrInitEMA(c.gpuClockEMA, key, 20*time.Second).Update(g.FrequencyMhz.Raw, now)
		g.PowerWatt.EMA = getOrInitEMA(c.gpuPowerEMA, key, 20*time.Second).Update(g.PowerWatt.Raw, now)
		g.Temperature.EMA = getOrInitEMA(c.gpuTempEMA, key, 30*time.Second).Update(g.Temperature.Raw, now)
	}

	// Disk
	for i := range m.Disk {
		d := &m.Disk[i]

		d.ReadMBps.EMA = getOrInitEMA(c.diskReadEMA, d.Name, 15*time.Second).Update(d.ReadMBps.Raw, now)
		d.WriteMBps.EMA = getOrInitEMA(c.diskWriteEMA, d.Name, 15*time.Second).Update(d.WriteMBps.Raw, now)
		d.UtilPct.EMA = getOrInitEMA(c.diskUtilEMA, d.Name, 20*time.Second).Update(d.UtilPct.Raw, now)
		d.Temperature.EMA = getOrInitEMA(c.diskTempEMA, d.Name, 30*time.Second).Update(d.Temperature.Raw, now)
	}

	// Network
	if c.iface != "" {
		m.Network.RXSpeedMBs.EMA = getOrInitEMA(c.netRxEMA, c.iface, 10*time.Second).Update(m.Network.RXSpeedMBs.Raw, now)
		m.Network.TXSpeedMBs.EMA = getOrInitEMA(c.netTxEMA, c.iface, 10*time.Second).Update(m.Network.TXSpeedMBs.Raw, now)
	}
}
