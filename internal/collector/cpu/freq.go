package cpu

import (
	"os"
	"strconv"
	"strings"
)

func getFreq() (float64, error) {
	b, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/scaling_cur_freq")
	if err != nil {
		return 0, err
	}

	mhz, _ := strconv.ParseFloat(strings.TrimSpace(string(b)), 64)
	return mhz / 1e3, nil
}
