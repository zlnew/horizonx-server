package cpu

import (
	"os"
	"strconv"
	"strings"
)

func (c *Collector) readUsage() (float64, []float64) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		c.log.Warn("failed to read /proc/stat", "error", err)
		return 0, nil
	}

	lines := strings.Split(string(data), "\n")

	var totalUsage float64
	var perCore []float64

	for _, line := range lines {
		if strings.HasPrefix(line, "cpu ") {
			if usage, err := c.parseCPUStat(line); err != nil {
				c.log.Warn("failed to parse cpu usage", "error", err, "line", line)
			} else {
				totalUsage = usage
			}
		} else if strings.HasPrefix(line, "cpu") {
			if usage, err := c.parseCPUStat(line); err != nil {
				c.log.Warn("failed to parse per-cpu usage", "error", err, "line", line)
			} else {
				perCore = append(perCore, usage)
			}
		}
	}

	return totalUsage, perCore
}

func (c *Collector) parseCPUStat(line string) (float64, error) {
	fields := strings.Fields(line)[1:]

	var idle, total uint64
	for i, val := range fields {
		v, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return 0, err
		}
		total += v
		if i == 3 {
			idle = v
		}
	}

	usage := float64(total-idle) / float64(total) * 100
	return usage, nil
}
