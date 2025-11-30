package memory

func (c *Collector) getSwapFree() float64 {
	return float64(c.swapFree) / 1024 / 1024 / 1024
}

