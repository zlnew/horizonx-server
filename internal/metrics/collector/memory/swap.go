package memory

func (c *Collector) getSwapTotalGB() float64 {
	return float64(c.swapTotal) / 1024 / 1024 / 1024
}

func (c *Collector) getSwapFreeGB() float64 {
	return float64(c.swapFree) / 1024 / 1024 / 1024
}

func (c *Collector) getSwapUsedGB() float64 {
	used := c.swapTotal - c.swapFree
	return float64(used) / 1024 / 1024 / 1024
}
