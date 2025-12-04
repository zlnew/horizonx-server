package cpu

import (
	"os"
	"strconv"
	"strings"
)

func (c *Collector) readFrequency() float64 {
	b, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/scaling_cur_freq")
	if err != nil {
		c.log.Debug("failed to read scaling_cur_freq", "error", err)
		return 0
	}

	mhz, err := strconv.ParseFloat(strings.TrimSpace(string(b)), 64)
	if err != nil {
		c.log.Warn("failed to parse cpu frequency", "error", err)
		return 0
	}

	return mhz / 1e3
}
