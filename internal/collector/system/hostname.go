package system

import "os"

func getHostname() string {
	name, err := os.Hostname()
	if err != nil {
		return "unknown"
	}

	return name
}
