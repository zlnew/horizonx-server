// Package core
package core

type CPUMetric struct {
	Usage     float64   `json:"usage"`
	PerCore   []float64 `json:"per_core"`
	Watt      float64   `json:"watt"`
	Temp      float64   `json:"temp"`
	Frequency float64   `json:"frequency"`
}

type MemoryMetric struct {
	MemTotal     float64 `json:"mem_total"`
	MemAvailable float64 `json:"mem_available"`
	MemUsed      float64 `json:"mem_used"`
	SwapTotal    float64 `json:"swap_total"`
	SwapFree     float64 `json:"swap_free"`
	SwapUsed     float64 `json:"swap_used"`
}

type Metrics struct {
	CPU    CPUMetric    `json:"cpu"`
	Memory MemoryMetric `json:"memory"`
}
