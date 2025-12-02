package cpu

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"zlnew/monitor-agent/pkg"
)

func (c *Collector) getWatt() (float64, error) {
	if watt, err := c.getRAPL(); err == nil && watt > 0 {
		return watt, nil
	}
	if watt, err := getHwmon(); err == nil && watt > 0 {
		return watt, nil
	}
	if watt, err := getRyzenSMU(); err == nil && watt > 0 {
		return watt, nil
	}
	return 0, nil
}

func (c *Collector) getRAPL() (float64, error) {
	b, err := os.ReadFile("/sys/class/powercap/intel-rapl/intel-rapl:0/energy_uj")
	if err != nil {
		return 0, err
	}

	energy, _ := strconv.ParseUint(strings.TrimSpace(string(b)), 10, 64)

	now := time.Now()
	if c.lastEnergy == 0 {
		c.lastEnergy = energy
		c.lastTime = now
		return 0, nil
	}

	deltaEnergy := float64(energy - c.lastEnergy)
	deltaTime := now.Sub(c.lastTime).Seconds()

	c.lastEnergy = energy
	c.lastTime = now

	watt := (deltaEnergy / 1e6) / deltaTime

	return watt, nil
}

func getHwmon() (float64, error) {
	matches, _ := filepath.Glob("/sys/class/hwmon/hwmon*/power*_input")
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
			continue
		}

		name := strings.TrimSpace(string(nameBytes))
		if !pkg.ContainsAny(name, targets) {
			continue
		}

		b, err := os.ReadFile(f)
		if err == nil {
			v, _ := strconv.ParseFloat(strings.TrimSpace(string(b)), 64)
			return v / 1e6, nil
		}
	}

	return 0, nil
}

func getRyzenSMU() (float64, error) {
	return 0, nil
}
