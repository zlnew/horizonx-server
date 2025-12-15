package disk

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func (c *Collector) findMountpointsByDeviceName(devName string) []string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		c.log.Warn("failed to open /proc/self/mountinfo", "error", err)
		return nil
	}
	defer f.Close()

	var result []string
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		mountpoint := fields[4]
		source := fields[9]
		if !strings.HasPrefix(source, "/dev/") {
			continue
		}

		realSource := strings.SplitN(source, "[", 2)[0]
		base := filepath.Base(realSource)
		if base == devName {
			result = append(result, mountpoint)
		}
	}

	if err := scanner.Err(); err != nil {
		c.log.Warn("error reading /proc/self/mountinfo", "error", err)
	}

	return result
}
