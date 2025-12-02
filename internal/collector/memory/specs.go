package memory

import (
	"bufio"
	"os/exec"
	"strings"
)

func getSpecs() ([]MemorySpec, error) {
	cmd := exec.Command("dmidecode", "-t", "memory")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var devices []MemorySpec
	var dev *MemorySpec

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "Handle") && strings.Contains(line, "DMI type 17") {
			if dev != nil {
				devices = append(devices, *dev)
			}
			dev = &MemorySpec{}
			continue
		}

		if dev == nil || !strings.Contains(line, ":") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "Size":
			dev.Size = value
		case "Form Factor":
			dev.FormFactor = value
		case "Type":
			dev.Type = value
		case "Speed":
			dev.Speed = value
		case "Manufacturer":
			dev.Manufacturer = value
		case "Part Number":
			dev.PartNumber = value
		}
	}

	if dev != nil {
		devices = append(devices, *dev)
	}

	return devices, nil
}
