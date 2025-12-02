package system

import (
	"os"
	"strconv"
	"strings"
)

func getUptime() uint64 {
	b, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0
	}

	fields := strings.Fields(string(b))
	secs, _ := strconv.ParseFloat(fields[0], 64)
	return uint64(secs)
}
