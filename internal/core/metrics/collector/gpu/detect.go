package gpu

import (
	"os"
	"strings"
)

func (c *Collector) detectGPUs() []string {
	entries, err := os.ReadDir("/sys/class/drm")
	if err != nil {
		c.log.Warn("failed to read /sys/class/drm", "error", err)
		return nil
	}

	var cards []string
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, "card") && !strings.Contains(name, "-") {
			cards = append(cards, name)
		}
	}

	return cards
}
