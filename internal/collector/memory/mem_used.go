package memory

func (c *Collector) getMemUsed() float64 {
	used := c.memTotal - c.memAvailable
	return float64(used) / 1024 / 1024 / 1024
}
