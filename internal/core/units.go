package core

type Metadata struct {
	CPU     CPUUnits     `json:"cpu"`
	GPU     GPUUnits     `json:"gpu"`
	Memory  MemoryUnits  `json:"memory"`
	Disk    DiskUnits    `json:"disk"`
	Network NetworkUnits `json:"network"`
	System  SystemUnits  `json:"system"`
}

type CPUUnits struct {
	Usage     string `json:"usage"`
	PerCore   string `json:"per_core"`
	Watt      string `json:"watt"`
	Temp      string `json:"temp"`
	Frequency string `json:"frequency"`
}

type GPUUnits struct {
	Usage     string `json:"usage"`
	VramTotal string `json:"vram_total"`
	VramUsed  string `json:"vram_used"`
	Temp      string `json:"temp"`
	Watt      string `json:"watt"`
}

type MemoryUnits struct {
	MemTotal     string `json:"mem_total"`
	MemAvailable string `json:"mem_available"`
	MemUsed      string `json:"mem_used"`
	SwapTotal    string `json:"swap_total"`
	SwapFree     string `json:"swap_free"`
	SwapUsed     string `json:"swap_used"`
}

type DiskUnits struct {
	Total string `json:"total"`
	Free  string `json:"free"`
	Used  string `json:"used"`
	Temp  string `json:"temp"`
}

type NetworkUnits struct {
	Upload   string `json:"upload"`
	Download string `json:"download"`
}

type SystemUnits struct {
	Uptime string `json:"uptime"`
}

func DefaultMetadata() Metadata {
	return Metadata{
		CPU: CPUUnits{
			Usage:     "percent",
			PerCore:   "percent",
			Watt:      "watt",
			Temp:      "celsius",
			Frequency: "MHz",
		},
		GPU: GPUUnits{
			Usage:     "percent",
			VramTotal: "MiB",
			VramUsed:  "MiB",
			Temp:      "celsius",
			Watt:      "watt",
		},
		Memory: MemoryUnits{
			MemTotal:     "GiB",
			MemAvailable: "GiB",
			MemUsed:      "GiB",
			SwapTotal:    "GiB",
			SwapFree:     "GiB",
			SwapUsed:     "GiB",
		},
		Disk: DiskUnits{
			Total: "GB",
			Free:  "GB",
			Used:  "GB",
			Temp:  "celsius",
		},
		Network: NetworkUnits{
			Upload:   "Mbps",
			Download: "Mbps",
		},
		System: SystemUnits{
			Uptime: "seconds",
		},
	}
}
