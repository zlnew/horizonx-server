package cpu

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"horizonx-server/internal/pkg"
)

func (c *Collector) readTemperature() float64 {
	matches, err := filepath.Glob("/sys/class/hwmon/hwmon*/temp*_input")
	if err != nil {
		c.log.Warn("failed to glob for hwmon temp files", "error", err)
		return 0
	}

	targets := []string{
		"coretemp",
		"k10temp",
		"zenpower",
		"zenpower3",
		"amd_smu",
		"ryzen_smu",
	}

	for _, f := range matches {
		dir := filepath.Dir(f)
		namePath := filepath.Join(dir, "name")

		nameBytes, err := os.ReadFile(namePath)
		if err != nil {
			c.log.Debug("failed to read hwmon name", "path", namePath, "error", err)
			continue
		}

		name := strings.TrimSpace(string(nameBytes))
		if !pkg.ContainsAny(name, targets) {
			continue
		}

		b, err := os.ReadFile(f)
		if err == nil {
			if v, err := strconv.ParseFloat(strings.TrimSpace(string(b)), 64); err == nil {
				return v / 1e3
			} else {
				c.log.Warn("failed to parse temperature", "file", f, "error", err)
			}
		} else {
			c.log.Debug("failed to read temp input", "file", f, "error", err)
		}
	}

	return 0
}
