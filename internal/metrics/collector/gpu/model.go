package gpu

import (
	"os"
	"strings"
)

func (c *Collector) readModel(card string) string {
	path := "/sys/class/drm/" + card + "/device/product_name"
	b, err := os.ReadFile(path)
	if err == nil {
		return strings.TrimSpace(string(b))
	}
	c.log.Debug("failed to read gpu product_name", "path", path, "error", err)

	path = "/sys/class/drm/" + card + "/device/device"
	b, err = os.ReadFile(path)
	if err == nil {
		return strings.TrimSpace(string(b))
	}
	c.log.Debug("failed to read gpu device", "path", path, "error", err)

	return "unknown"
}
