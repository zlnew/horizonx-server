package pkg

type EMA struct {
	alpha float64
	value float64
	init  bool
}

func NewEMA(alpha float64) *EMA {
	return &EMA{alpha: alpha}
}

func (e *EMA) Add(newvalue float64) float64 {
	if !e.init {
		e.value = newvalue
		e.init = true
		return e.value
	}

	e.value = e.alpha*newvalue + (1-e.alpha)*e.value
	return e.value
}

func (e *EMA) Value() float64 {
	return e.value
}
