package gpu

import (
	"os"
	"strconv"
	"strings"
)

func (c *Collector) readCoreUsage(card string) float64 {
	path := "/sys/class/drm/" + card + "/device/gpu_busy_percent"
	b, err := os.ReadFile(path)
	if err != nil {
		c.log.Debug("failed to read gpu_busy_percent", "path", path, "error", err)
		return 0
	}

	v, err := strconv.ParseFloat(strings.TrimSpace(string(b)), 64)
	if err != nil {
		c.log.Warn("failed to parse gpu_busy_percent", "path", path, "error", err)
		return 0
	}

	return v
}
