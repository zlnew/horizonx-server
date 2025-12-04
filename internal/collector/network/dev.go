package network

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func (c *Collector) readTotals() (uint64, uint64, error) {
	f, err := os.Open("/proc/net/dev")
	if err != nil {
		c.log.Error("failed to open /proc/net/dev", "error", err)
		return 0, 0, err
	}
	defer f.Close()

	var rxTotal uint64
	var txTotal uint64

	scanner := bufio.NewScanner(f)
	for i := 0; i < 2 && scanner.Scan(); i++ {
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 17 {
			continue
		}

		iface := strings.TrimSuffix(parts[0], ":")
		if iface == "lo" {
			continue
		}

		rxBytes := c.parseUint(parts[1])
		txBytes := c.parseUint(parts[9])

		rxTotal += rxBytes
		txTotal += txBytes
	}

	if err := scanner.Err(); err != nil {
		c.log.Warn("error reading /proc/net/dev", "error", err)
	}

	return rxTotal, txTotal, nil
}

func (c *Collector) parseUint(s string) uint64 {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		c.log.Warn("failed to parse network stats", "value", s, "error", err)
		return 0
	}
	return v
}
