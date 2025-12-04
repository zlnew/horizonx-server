package gpu

import (
	"os"
	"strconv"
	"strings"
)

func (c *Collector) readTemperature(card string) float64 {
	hwmonRoot := "/sys/class/drm/" + card + "/device/hwmon"

	hwmons, err := os.ReadDir(hwmonRoot)
	if err != nil {
		c.log.Debug("failed to read hwmon for gpu temp", "path", hwmonRoot, "error", err)
		return 0
	}
	for _, hw := range hwmons {
		file := hwmonRoot + "/" + hw.Name() + "/temp1_input"
		b, err := os.ReadFile(file)
		if err == nil {
			v, err := strconv.ParseFloat(strings.TrimSpace(string(b)), 64)
			if err != nil {
				c.log.Warn("failed to parse gpu temperature", "file", file, "error", err)
				continue
			}
			return v / 1000
		}
	}

	return 0
}
