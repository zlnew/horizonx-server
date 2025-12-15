// Package disk
package disk

import (
	"context"

	"horizonx-server/internal/logger"
)

func NewCollector(log logger.Logger) *Collector {
	return &Collector{log: log}
}

func (c *Collector) Collect(ctx context.Context) ([]DiskMetric, error) {
	disks, parts, err := c.detectBlockDevices()
	if err != nil {
		c.log.Error("failed to detect block devices", "error", err)
		return nil, err
	}

	var result []DiskMetric

	for _, disk := range disks {
		metric := DiskMetric{
			Name:        disk,
			RawSizeGB:   c.readRawSizeGiB(disk),
			Temperature: c.readDiskTemperature(disk),
		}

		for _, part := range parts {
			parent := c.getParentDisk(part)
			if parent != disk {
				continue
			}

			mounts := c.findMountpointsByDeviceName(part)
			for _, mp := range mounts {
				fs := c.readFSUsage(mp, part)
				metric.Filesystems = append(metric.Filesystems, fs)
			}
		}

		result = append(result, metric)
	}

	return result, nil
}
