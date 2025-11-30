package memory

func (c *Collector) getSwapUsed() float64 {
	used := c.swapTotal - c.swapFree
	return float64(used) / 1024 / 1024 / 1024
}
