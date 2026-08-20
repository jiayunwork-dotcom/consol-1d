package series

import (
	"fmt"
	"math"
)

// Result carries a truncated partial sum together with the proof that
// the tail beyond TermsUsed stays below RemainderBound. The two numbers
// must be read together: the sum alone is meaningless without the
// verified remainder bound that justifies the truncation.
type Result struct {
	Sum            float64
	TermsUsed      int
	RemainderBound float64
}

// Sum adaptively truncates the term series so that the absolute tail
// beyond the returned term count is bounded by tol. The loop starts at
// one term and doubles until TailBoundUpper clears tol, so the reported
// RemainderBound is always a verified upper bound on the true tail.
//
// The doubling makes the runtime O(N) where N is the smallest power of
// two that clears the tolerance; the returned TermsUsed is that N, and
// the caller may report it honestly as the number of modes kept.
func Sum(t Term, tol float64, maxTerms int) (Result, error) {
	if err := t.Validate(); err != nil {
		return Result{}, err
	}
	if math.IsNaN(tol) || tol <= 0 {
		return Result{}, fmt.Errorf("series tolerance must be positive, got %g", tol)
	}
	if maxTerms < 1 {
		return Result{}, fmt.Errorf("series max terms must be positive, got %d", maxTerms)
	}
	if t.Coefficient == 0 {
		return Result{Sum: 0, TermsUsed: 1, RemainderBound: 0}, nil
	}
	n := 1
	for {
		if t.TailBoundUpper(n) <= tol {
			break
		}
		n *= 2
		if n > maxTerms {
			return Result{}, fmt.Errorf(
				"series needs more than %d terms (tail bound %g, need <= %g): "+
					"the decay is too slow for this tolerance",
				maxTerms, t.TailBoundUpper(n), tol)
		}
	}
	acc := 0.0
	for i := 0; i < n; i++ {
		acc += t.Value(i)
	}
	return Result{Sum: acc, TermsUsed: n, RemainderBound: t.TailBoundUpper(n)}, nil
}

// SumAtTolerance is a convenience wrapper around Sum with the package
// default tolerance and term cap.
func SumAtTolerance(t Term) (Result, error) {
	return Sum(t, DefaultTolerance, DefaultMaxTerms)
}
