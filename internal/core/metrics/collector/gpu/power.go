package gpu

import (
	"os"
	"strconv"
	"strings"
)

func (c *Collector) readPowerRaw(card string) float64 {
	hwmon := "/sys/class/drm/" + card + "/device/hwmon"
	hwmons, err := os.ReadDir(hwmon)
	if err != nil {
		c.log.Debug("failed to read hwmon for gpu power", "path", hwmon, "error", err)
		return 0
	}

	for _, hw := range hwmons {
		file := hwmon + "/" + hw.Name() + "/power1_input"
		b, err := os.ReadFile(file)
		if err == nil {
			v, err := strconv.ParseFloat(strings.TrimSpace(string(b)), 64)
			if err != nil {
				c.log.Warn("failed to parse gpu power", "file", file, "error", err)
				continue
			}
			return v / 1e6
		}
	}

	return 0
}

func (c *Collector) readPower(card string) float64 {
	raw := c.readPowerRaw(card)
	ema := c.powerEMA[card]

	if raw > 0 {
		ema.Add(raw)
	}

	return ema.Value()
}
