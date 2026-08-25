package model

import (
	"math"

	"consol-1d/internal/series"
)

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

func SmallTimeAsymptote(tv float64) float64 {
	return 2 * math.Sqrt(tv/math.Pi)
}

func settlementFrom(in Input, u float64) (s, sult, ratio float64) {
	if in.Mv == nil || in.DeltaSigma == nil {
		return math.NaN(), math.NaN(), u
	}
	sult = (*in.Mv) * (*in.DeltaSigma) * in.Thickness
	return u * sult, sult, u
}

func UniformInput(cv, h float64, d Drainage, u0, t float64) Input {
	return Input{
		Cv:        cv,
		Thickness: h,
		Drainage:  d,
		Initial:   InitialPressure{Type: InitialUniform, U0: u0},
		Time:      t,
	}
}
