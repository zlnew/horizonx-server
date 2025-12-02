package system

import (
	"os/exec"
	"strings"
)

func getKernelVersion() string {
	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}
