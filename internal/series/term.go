package series

import (
	"fmt"
	"math"
)

// Term describes one family of summands of the form
//
//	Coefficient / (2n+1)^Power * exp(-(2n+1)^2 * Decay)
//
// for n = 0, 1, 2, ...  Power >= 1 and Decay >= 0 are required. The
// same type covers the pore-pressure series (Power 1) and the
// average-consolidation series (Power 2); Power is kept as a float so
// that no rounding appears when comparing exponents.
type Term struct {
	Coefficient float64
	Power       float64
	Decay       float64
}

// Value evaluates the n-th summand.
func (t Term) Value(n int) float64 {
	k := float64(2*n + 1)
	return t.Coefficient * math.Exp(-k*k*t.Decay) / math.Pow(k, t.Power)
}

// Validate rejects series that are not absolutely convergent. A family
// with Power 1 and Decay 0 is the harmonic series: it does not have a
// finite sum, so calling Sum on it is a programming error the package
// reports instead of guessing.
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

// AbsoluteMax is the largest |Value| reachable at any index; for a
// decaying series it is the first term.
func (t Term) AbsoluteMax() float64 {
	if err := t.Validate(); err != nil {
		return math.NaN()
	}
	return math.Abs(t.Value(0))
}
