package gpu

import (
	"os"
	"strings"
)

func (c *Collector) readVendor(card string) string {
	path := "/sys/class/drm/" + card + "/device/vendor"
	b, err := os.ReadFile(path)
	if err != nil {
		c.log.Debug("failed to read gpu vendor", "path", path, "error", err)
		return "unknown"
	}

	val := strings.TrimSpace(string(b))
	switch val {
	case "0x1002":
		return "AMD"
	case "0x10de":
		return "NVIDIA"
	case "0x8086":
		return "INTEL"
	}

	return val
}
