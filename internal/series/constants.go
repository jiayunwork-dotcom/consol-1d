// Package series provides the mathematical building blocks of the
// one-dimensional consolidation solver: odd-index Fourier-type series
// with adaptive truncation and a strict remainder bound, plus Romberg
// quadrature used for independent cross-checks.
//
// All consolidation series share the same shape
//
//	term(n) = Coefficient / (2n+1)^Power * exp(-(2n+1)^2 * Decay)
//
// with Decay = pi^2*Tv/4. The package never hard-codes a physical
// constant; the caller supplies Coefficient, Power and Decay, and the
// package guarantees that the reported partial sum misses the true
// infinite sum by less than a caller-chosen tolerance.
package series

import "math"

// OddStep is the spacing between consecutive summation indices: the
// series visits (2n+1) for n = 0, 1, 2, ...
const OddStep = 2.0

// PiSquaredOverFour is pi^2/4, the exponential base of the
// consolidation Fourier modes: exp(-(2n+1)^2 * pi^2 * Tv / 4).
const PiSquaredOverFour = math.Pi * math.Pi / 4

// DefaultTolerance is the absolute bound on the truncated tail that Sum
// guarantees when the caller does not choose one.
const DefaultTolerance = 1e-9

// DefaultMaxTerms caps the adaptive truncation loop. A physically
// useful scenario converges in a few hundred terms at Tv >= 1e-3; the
// cap only matters for vanishing time factors where the exponential
// factor decays slowly.
const DefaultMaxTerms = 1 << 25

// maxQuadSteps bounds the depth of the Romberg extrapolation table.
// Twenty-four halvings give more than enough accuracy for the
// cross-check integrals in this package.
const maxQuadSteps = 24
