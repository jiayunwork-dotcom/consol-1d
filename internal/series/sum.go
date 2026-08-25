package series

import (
	"fmt"
	"math"
)

type Result struct {
	Sum            float64
	TermsUsed      int
	RemainderBound float64
}

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

func SumAtTolerance(t Term) (Result, error) {
	return Sum(t, DefaultTolerance, DefaultMaxTerms)
}
