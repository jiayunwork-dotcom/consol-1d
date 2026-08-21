package model

import (
	"math"

	"consol-1d/internal/series"
)

// ClassicalAverageU evaluates the textbook average consolidation degree
// for a uniform initial profile,
//
//	U = 1 - (8/pi^2) * sum_{n>=0} exp(-(2n+1)^2 pi^2 Tv / 4) / (2n+1)^2
//
// with Tv the drainage-distance time factor. The solver's integrated U
// equals this series exactly for a uniform profile; the function exists
// as an independent reference used by tests and documentation, so a
// regression in the drainage-distance mapping shows up as a mismatch
// instead of being absorbed by the solver.
func ClassicalAverageU(tv, tol float64, maxTerms int) (float64, int, float64, error) {
	t := series.Term{
		Coefficient: 8 / (math.Pi * math.Pi),
		Power:       2,
		Decay:       series.PiSquaredOverFour * tv,
	}
	r, err := series.Sum(t, tol, maxTerms)
	if err != nil {
		return 0, 0, 0, err
	}
	return 1 - r.Sum, r.TermsUsed, r.RemainderBound, nil
}

// SmallTimeAsymptote returns the leading-order consolidation degree
// U ~ 2*sqrt(Tv/pi) valid for small Tv, the classic short-time curve of
// the Terzaghi solution. The solver reports it as a reference for
// vanishing time factors, where the Fourier series needs the most
// terms; it also pins the expected slope at Tv = 0.
func SmallTimeAsymptote(tv float64) float64 {
	return 2 * math.Sqrt(tv/math.Pi)
}

// settlementFrom computes s_ult and s when mv and delta_sigma are both
// present: s_ult = mv*delta_sigma*H and s = U*s_ult. Without them the
// absolute values are NaN and the ratio s/s_ult collapses to U by
// definition, which is why the report can always print the settlement
// ratio.
func settlementFrom(in Input, u float64) (s, sult, ratio float64) {
	if in.Mv == nil || in.DeltaSigma == nil {
		return math.NaN(), math.NaN(), u
	}
	sult = (*in.Mv) * (*in.DeltaSigma) * in.Thickness
	return u * sult, sult, u
}

// UniformInput builds a scenario with a uniform initial pressure from
// the scalar fields. It is the smallest legal input a caller can
// construct and keeps tests and examples free of boilerplate.
func UniformInput(cv, h float64, d Drainage, u0, t float64) Input {
	return Input{
		Cv:        cv,
		Thickness: h,
		Drainage:  d,
		Initial:   InitialPressure{Type: InitialUniform, U0: u0},
		Time:      t,
	}
}
