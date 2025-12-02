package network

import "time"

func (c *Collector) collectMetric() (NetworkMetric, error) {
	rx, tx, err := readTotals()
	if err != nil {
		return NetworkMetric{}, err
	}

	now := time.Now()
	var download float64
	var upload float64

	if !c.lastTime.IsZero() {
		elapsed := now.Sub(c.lastTime).Seconds()
		if elapsed > 0 && rx >= c.lastRxBytes && tx >= c.lastTxBytes {
			download = float64(rx-c.lastRxBytes) * 8 / elapsed / 1_000_000
			upload = float64(tx-c.lastTxBytes) * 8 / elapsed / 1_000_000
		}
	}

	c.lastRxBytes = rx
	c.lastTxBytes = tx
	c.lastTime = now

	return NetworkMetric{
		Upload:   upload,
		Download: download,
	}, nil
}
