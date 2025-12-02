package cpu

import (
	"os"
	"strconv"
	"strings"
)

func getUsage() (float64, []float64, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, nil, err
	}

	lines := strings.Split(string(data), "\n")

	var totalUsage float64
	var perCore []float64

	for _, line := range lines {
		if strings.HasPrefix(line, "cpu ") {
			usage, _ := parseCPUStat(line)
			totalUsage = usage
		} else if strings.HasPrefix(line, "cpu") {
			usage, _ := parseCPUStat(line)
			perCore = append(perCore, usage)
		}
	}

	return totalUsage, perCore, nil
}

func parseCPUStat(line string) (float64, error) {
	fields := strings.Fields(line)[1:]

	var idle, total uint64
	for i, val := range fields {
		v, _ := strconv.ParseUint(val, 10, 64)
		total += v
		if i == 3 {
			idle = v
		}
	}

	usage := float64(total-idle) / float64(total) * 100
	return usage, nil
}
