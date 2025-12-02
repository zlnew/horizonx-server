package cpu

import (
	"bufio"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func getSpec() (CPUSpec, error) {
	var spec CPUSpec

	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return spec, err
	}
	defer file.Close()

	seenModel := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "vendor_id":
			spec.Vendor = value

		case "model name":
			if !seenModel {
				spec.Model = value
				seenModel = true
			}

		case "cpu cores":
			if n, err := strconv.Atoi(value); err == nil {
				spec.Cores = n
			}

		case "siblings":
			if n, err := strconv.Atoi(value); err == nil {
				spec.Threads = n
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return spec, err
	}

	spec.Arch = runtime.GOARCH
	spec.BaseFreq = readFreqByPath("/sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_min_freq")
	spec.MaxFreq = readFreqByPath("/sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_max_freq")

	return spec, nil
}

func readFreqByPath(path string) int {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	n, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0
	}
	return n / 1000
}
