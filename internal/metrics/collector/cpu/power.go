package cpu

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"horizonx-server/internal/pkg"
)

func (c *Collector) readPowerWatt() float64 {
	var raw float64 = 0

	if watt, err := c.readRAPL(); err == nil && watt > 0 {
		raw = watt
	} else if err != nil {
		c.log.Debug("failed to read RAPL power", "error", err)
	} else if watt, err := c.readHwmon(); err == nil && watt > 0 {
		raw = watt
	} else if err != nil {
		c.log.Debug("failed to read hwmon power", "error", err)
	}

	ema := c.powerEMA

	if raw > 0 {
		ema.Add(raw)
	}

	return ema.Value()
}

func (c *Collector) readRAPL() (float64, error) {
	b, err := os.ReadFile("/sys/class/powercap/intel-rapl/intel-rapl:0/energy_uj")
	if err != nil {
		return 0, err
	}

	energy, err := strconv.ParseUint(strings.TrimSpace(string(b)), 10, 64)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	if c.lastEnergy == 0 {
		c.lastEnergy = energy
		c.lastTime = now
		return 0, nil
	}

	deltaEnergy := float64(energy - c.lastEnergy)
	deltaTime := now.Sub(c.lastTime).Seconds()

	if energy < c.lastEnergy {
		deltaEnergy = float64((^uint64(0) - c.lastEnergy) + energy)
	}

	c.lastEnergy = energy
	c.lastTime = now

	watt := (deltaEnergy / 1e6) / deltaTime

	return watt, nil
}

func (c *Collector) readHwmon() (float64, error) {
	matches, err := filepath.Glob("/sys/class/hwmon/hwmon*/power*_input")
	if err != nil {
		return 0, err
	}
	targets := []string{
		"zenpower",
		"zenpower3",
		"amd_smu",
		"ryzen_smu",
		"rapl",
		"intel-rapl",
		"intel-rapl-msr",
	}

	for _, f := range matches {
		dir := filepath.Dir(f)
		namePath := filepath.Join(dir, "name")

		nameBytes, err := os.ReadFile(namePath)
		if err != nil {
			c.log.Debug("failed to read hwmon name for power", "path", namePath, "error", err)
			continue
		}

		name := strings.TrimSpace(string(nameBytes))
		if !pkg.ContainsAny(name, targets) {
			continue
		}

		b, err := os.ReadFile(f)
		if err == nil {
			v, err := strconv.ParseFloat(strings.TrimSpace(string(b)), 64)
			if err != nil {
				c.log.Warn("failed to parse hwmon power", "file", f, "error", err)
				continue
			}
			return v / 1e6, nil
		} else {
			c.log.Debug("failed to read power input", "file", f, "error", err)
		}
	}

	c.log.Debug("no suitable hwmon power input found")
	return 0, nil
}
