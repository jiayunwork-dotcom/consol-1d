package series

import (
	"fmt"
	"math"
)

// PartialSum evaluates the first n terms of the family together with
// the verified bound on everything from term n onward. It is the
// low-level counterpart of Sum for callers that want to fix the term
// count explicitly, for example when every profile node must share the
// same truncation.
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

// TermRatio returns |term(n+1)| / |term(n)|, the factor by which each
// step of two on the index shrinks the summand. It quantifies how fast
// the series converges and is used to double-check that the exponential
// decay dominates once the index is large enough.
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

// DecayIndex reports the smallest index n from which TermRatio stays
// below the given threshold, or -1 if the threshold is never crossed
// within the checked window. It gives the solver a way to reason about
// where the tail becomes negligible.
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
