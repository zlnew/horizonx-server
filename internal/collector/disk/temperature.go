package disk

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (c *Collector) readDiskTemperature(name string) float64 {
	root := filepath.Join("/sys/block", name, "device")
	return c.readTempFromHwmonRoot(root)
}

func (c *Collector) readTempFromHwmonRoot(root string) float64 {
	hwmons, err := os.ReadDir(root)
	if err != nil {
		c.log.Debug("failed to read hwmon root for disk temp", "path", root, "error", err)
		return 0
	}

	for _, hw := range hwmons {
		if !strings.HasPrefix(hw.Name(), "hwmon") {
			continue
		}
		file := filepath.Join(root, hw.Name(), "temp1_input")
		b, err := os.ReadFile(file)
		if err != nil {
			c.log.Debug("failed to read temp1_input for disk", "file", file, "error", err)
			continue
		}
		val, err := strconv.ParseFloat(strings.TrimSpace(string(b)), 64)
		if err != nil {
			c.log.Warn("failed to parse disk temperature", "file", file, "error", err)
			continue
		}
		return val / 1000.0
	}

	return 0
}
