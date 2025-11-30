package memory

func (c *Collector) getMemAvailable() float64 {
	return float64(c.memAvailable) / 1024 / 1024 / 1024
}
