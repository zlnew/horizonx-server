package disk

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	nvmePartition = regexp.MustCompile(`^nvme\d+n\d+p\d+$`)
	sdPartition   = regexp.MustCompile(`^sd[a-z]+\d+$`)
	mmcPartition  = regexp.MustCompile(`^mmcblk\d+p\d+$`)
)

func isPartition(name string) bool {
	return nvmePartition.MatchString(name) ||
		sdPartition.MatchString(name) ||
		mmcPartition.MatchString(name)
}

func (c *Collector) detectBlockDevices() ([]string, []string, error) {
	var disks []string
	var parts []string

	entries, err := os.ReadDir("/sys/class/block")
	if err != nil {
		c.log.Error("failed to read /sys/class/block", "error", err)
		return nil, nil, err
	}

	for _, e := range entries {
		name := e.Name()

		if strings.HasPrefix(name, "loop") ||
			strings.HasPrefix(name, "ram") ||
			strings.HasPrefix(name, "dm-") {
			continue
		}

		if isPartition(name) {
			parts = append(parts, name)
		} else {
			disks = append(disks, name)
		}
	}

	return disks, parts, nil
}

func (c *Collector) getParentDisk(part string) string {
	base := "/sys/class/block"
	entries, err := os.ReadDir(base)
	if err != nil {
		c.log.Warn("failed to read /sys/class/block in getParentDisk", "error", err)
		return ""
	}

	for _, e := range entries {
		disk := e.Name()
		path := filepath.Join(base, disk, part)
		if _, err := os.Stat(path); err == nil {
			return disk
		} else if !os.IsNotExist(err) {
			c.log.Debug("failed to stat parent disk path", "path", path, "error", err)
		}
	}

	return ""
}
