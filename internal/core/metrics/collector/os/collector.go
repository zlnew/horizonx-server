// Package os
package os

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"horizonx-server/internal/logger"
)

func NewCollector(log logger.Logger) *Collector {
	return &Collector{log: log}
}

func (c *Collector) Collect(ctx context.Context) (OSInfo, error) {
	var info OSInfo

	hostname, err := os.Hostname()
	if err != nil {
		c.log.Debug("failed to get hostname", "error", err)
		hostname = "unknown"
	}
	info.Hostname = hostname

	info.Arch = runtime.GOARCH

	kernel := "unknown"
	if data, err := os.ReadFile("/proc/version"); err == nil {
		parts := strings.Fields(string(data))
		if len(parts) >= 3 {
			kernel = parts[2]
		}
	} else if out, err := exec.Command("uname", "-r").Output(); err == nil {
		kernel = strings.TrimSpace(string(out))
	}
	info.KernelVersion = kernel

	osName := runtime.GOOS
	if f, err := os.Open("/etc/os-release"); err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				osName = strings.Trim(line[len("PRETTY_NAME="):], `"`)
				break
			}
		}
	}
	info.Name = osName + " " + info.Arch

	return info, nil
}
