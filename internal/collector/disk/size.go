package disk

func (c *Collector) getTotal() float64 {
	return c.totalBytes / 1_000_000_000
}

func (c *Collector) getFree() float64 {
	return c.freeBytes / 1_000_000_000
}

func (c *Collector) getUsed() float64 {
	usedBytes := c.totalBytes - c.freeBytes
	return usedBytes / 1_000_000_000
}
