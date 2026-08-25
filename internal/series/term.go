package series

import (
	"fmt"
	"math"
)

type Term struct {
	Coefficient float64
	Power       float64
	Decay       float64
}

func (t Term) Value(n int) float64 {
	k := float64(2*n + 1)
	return t.Coefficient * math.Exp(-k*k*t.Decay) / math.Pow(k, t.Power)
}

func (t Term) Validate() error {
	if math.IsNaN(t.Coefficient) || math.IsNaN(t.Power) || math.IsNaN(t.Decay) {
		return fmt.Errorf("series term has a NaN field: %+v", t)
	}
	if math.IsInf(t.Coefficient, 0) || math.IsInf(t.Decay, 0) {
		return fmt.Errorf("series term has an infinite field: %+v", t)
	}
	if t.Power < 1 {
		return fmt.Errorf("series term power must be at least 1, got %g", t.Power)
	}
	if t.Decay < 0 {
		return fmt.Errorf("series term decay must be non-negative, got %g", t.Decay)
	}
	if t.Decay == 0 && t.Power <= 1 {
		return fmt.Errorf("series with power %g and zero decay has no finite sum", t.Power)
	}
	return nil
}

func (t Term) AbsoluteMax() float64 {
	if err := t.Validate(); err != nil {
		return math.NaN()
	}
	return math.Abs(t.Value(0))
}
