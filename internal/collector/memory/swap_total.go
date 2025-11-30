package memory

func (c *Collector) getSwapTotal() float64 {
	return float64(c.swapTotal) / 1024 / 1024 / 1024
}

