package series

import (
	"fmt"
	"math"
)

// Romberg integrates f over [a, b] with a Romberg table built from the
// composite trapezoid rule. It returns the extrapolated value together
// with the difference between the two finest approximations as a
// practical error estimate, so callers can report achieved accuracy
// instead of pretending the result is exact.
func Romberg(f func(float64) float64, a, b float64, tol float64) (float64, float64, error) {
	if f == nil {
		return 0, 0, fmt.Errorf("romberg needs a non-nil integrand")
	}
	if b <= a {
		return 0, 0, fmt.Errorf("romberg requires a < b, got [%g, %g]", a, b)
	}
	if math.IsNaN(tol) || tol <= 0 {
		return 0, 0, fmt.Errorf("romberg tolerance must be positive, got %g", tol)
	}
	// table[i][j] holds the trapezoid estimate with 2^i intervals after
	// j Richardson extrapolation steps. Each row is extrapolated against
	// the previous row to cancel the leading error terms.
	table := make([][]float64, 0, maxQuadSteps)
	h := b - a
	prev := 0.5 * h * (f(a) + f(b))
	table = append(table, []float64{prev})
	newSamples := 1
	for i := 1; i < maxQuadSteps; i++ {
		h = h / 2 // interval width of the current trapezoid rule
		sum := 0.0
		for k := 0; k < newSamples; k++ {
			x := a + (float64(2*k)+1)*h
			sum += f(x)
		}
		curr := 0.5*prev + sum*h
		row := make([]float64, 0, i+1)
		row = append(row, curr)
		for j := 1; j <= i; j++ {
			denom := math.Pow(4, float64(j)) - 1
			extrap := row[j-1] + (row[j-1]-table[i-1][j-1])/denom
			row = append(row, extrap)
		}
		table = append(table, row)
		best := row[i]
		prevBest := table[i-1][i-1]
		diff := math.Abs(best - prevBest)
		if diff <= tol {
			return best, diff, nil
		}
		prev = curr
		newSamples *= 2
	}
	last := table[maxQuadSteps-1][maxQuadSteps-1]
	before := table[maxQuadSteps-2][maxQuadSteps-2]
	return last, math.Abs(last - before), fmt.Errorf(
		"romberg did not reach %g after %d halvings (last diff %g)",
		tol, maxQuadSteps, math.Abs(last-before))
}

// IntegratePositive is a convenience wrapper that asserts the integral
// is finite and non-negative; it is used by the cross-check helpers to
// fail loudly on nonsense inputs.
func IntegratePositive(f func(float64) float64, a, b float64, tol float64) (float64, error) {
	v, _, err := Romberg(f, a, b, tol)
	if err != nil {
		return 0, err
	}
	if math.IsNaN(v) || math.IsInf(v, 0) || v < 0 {
		return 0, fmt.Errorf("integral is not a finite non-negative value: %g", v)
	}
	return v, nil
}
