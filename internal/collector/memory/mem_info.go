package memory

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func (c *Collector) readMemInfo() error {
	c.memTotal = 0
	c.memAvailable = 0
	c.swapTotal = 0
	c.swapFree = 0

	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")
		valueKB, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}

		switch key {
		case "MemTotal":
			c.memTotal = valueKB * 1024
		case "MemAvailable":
			c.memAvailable = valueKB * 1024
		case "SwapTotal":
			c.swapTotal = valueKB * 1024
		case "SwapFree":
			c.swapFree = valueKB * 1024
		}
	}

	return nil
}
