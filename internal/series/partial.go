package series

import (
	"fmt"
	"math"
)

func (t Term) PartialSum(n int) (Result, error) {
	if err := t.Validate(); err != nil {
		return Result{}, err
	}
	if n < 1 {
		return Result{}, fmt.Errorf("partial sum needs at least 1 term, got %d", n)
	}
	acc := 0.0
	for i := 0; i < n; i++ {
		acc += t.Value(i)
	}
	return Result{Sum: acc, TermsUsed: n, RemainderBound: t.TailBoundUpper(n)}, nil
}

func (t Term) TermRatio(n int) float64 {
	if err := t.Validate(); err != nil {
		return 0
	}
	cur := math.Abs(t.Value(n))
	if cur == 0 {
		return 0
	}
	return math.Abs(t.Value(n+1)) / cur
}

func (t Term) DecayIndex(threshold float64, window int) int {
	if err := t.Validate(); err != nil {
		return -1
	}
	for n := 0; n < window; n++ {
		if t.TermRatio(n) < threshold {
			return n
		}
	}
	return -1
}
