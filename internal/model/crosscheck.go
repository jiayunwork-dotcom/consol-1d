package model

import (
	"math"

	"consol-1d/internal/series"
)

// NumericalCoefficient computes the n-th modal amplitude of the initial
// profile by direct quadrature of the projection integral
//
//	A_n = (2/H) * int_0^H u(z,0) * sin(lambda_n z) dz
//
// via the Romberg integrator. The solver never uses this path — it
// evaluates the analytic closed forms — but the function exists as an
// independent cross-check so that a slip in the analytic derivation is
// caught by tests instead of silently corrupting every profile.
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

// NumericalMeanPressure integrates the pore-pressure field over the
// layer by quadrature of a finite set of already-computed mode
// contributions. It is the quadrature twin of the analytic layer
// average and is used only in cross-checks.
func NumericalMeanPressure(terms []float64, ip InitialPressure, path DrainagePath) float64 {
	sum := 0.0
	for n, a := range terms {
		sum += a * path.EigenIntegral(n)
	}
	return sum / path.H
}

// orthonormalityErr integrates sin(lambda_m z) sin(lambda_n z) over the
// layer and reports the deviation from the expected H/2. It is a
// numerical sanity check that the eigenfunctions are orthogonal on the
// drained domain, independent of any particular profile.
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
