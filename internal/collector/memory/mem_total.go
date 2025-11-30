package memory

func (c *Collector) getMemTotal() float64 {
	return float64(c.memTotal) / 1024 / 1024 / 1024
}
