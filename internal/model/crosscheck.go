package model

import (
	"math"

	"consol-1d/internal/series"
)

func NumericalCoefficient(ip InitialPressure, path DrainagePath, n int) (float64, error) {
	h := path.H
	lam := path.Eigenvalue(n)
	f := func(z float64) float64 {
		return ip.PressureAt(z, h) * math.Sin(lam*z)
	}
	v, _, err := series.Romberg(f, 0, h, 1e-10)
	if err != nil {
		return 0, err
	}
	return 2 * v / h, nil
}

func NumericalMeanPressure(terms []float64, ip InitialPressure, path DrainagePath) float64 {
	sum := 0.0
	for n, a := range terms {
		sum += a * path.EigenIntegral(n)
	}
	return sum / path.H
}

func orthonormalityErr(path DrainagePath, m, n int) (float64, error) {
	lm, ln := path.Eigenvalue(m), path.Eigenvalue(n)
	f := func(z float64) float64 {
		return math.Sin(lm*z) * math.Sin(ln*z)
	}
	v, _, err := series.Romberg(f, 0, path.H, 1e-10)
	if err != nil {
		return 0, err
	}
	if m == n {
		return math.Abs(v - path.H/2), nil
	}
	return math.Abs(v), nil
}
