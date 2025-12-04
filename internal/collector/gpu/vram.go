package gpu

import (
	"os"
	"strconv"
	"strings"
)

func (c *Collector) readVRAM(card string) (total float64, used float64, percent float64) {
	base := "/sys/class/drm/" + card + "/device/"

	totalPath := base + "mem_info_vram_total"
	t, err := os.ReadFile(totalPath)
	if err != nil {
		c.log.Debug("failed to read vram total", "path", totalPath, "error", err)
		return
	}

	usedPath := base + "mem_info_vram_used"
	u, err := os.ReadFile(usedPath)
	if err != nil {
		c.log.Debug("failed to read vram used", "path", usedPath, "error", err)
		return
	}

	totalB, err := strconv.ParseFloat(strings.TrimSpace(string(t)), 64)
	if err != nil {
		c.log.Warn("failed to parse vram total", "path", totalPath, "error", err)
		return
	}

	usedB, err := strconv.ParseFloat(strings.TrimSpace(string(u)), 64)
	if err != nil {
		c.log.Warn("failed to parse vram used", "path", usedPath, "error", err)
		return
	}

	total = totalB / (1024 * 1024 * 1024)
	used = usedB / (1024 * 1024 * 1024)

	if totalB > 0 {
		percent = (usedB / totalB) * 100
	}

	return
}
