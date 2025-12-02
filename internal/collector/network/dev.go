package network

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func readTotals() (uint64, uint64, error) {
	f, err := os.Open("/proc/net/dev")
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	var rxTotal uint64
	var txTotal uint64

	scanner := bufio.NewScanner(f)
	// skip headers (first two lines)
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

		rxBytes := parseUint(parts[1])
		txBytes := parseUint(parts[9])

		rxTotal += rxBytes
		txTotal += txBytes
	}

	return rxTotal, txTotal, nil
}

func parseUint(s string) uint64 {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return v
}
