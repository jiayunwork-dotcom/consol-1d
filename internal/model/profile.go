package model

import (
	"fmt"
	"math"

	"consol-1d/internal/series"
)

const SeriesTolerance = 1e-9

const TermsLimit = 1 << 25

func decayRate(tv float64) float64 {
	return series.PiSquaredOverFour * tv
}

func termsNeeded(t series.Term) (int, float64, error) {
	r, err := series.Sum(t, SeriesTolerance, TermsLimit)
	if err != nil {
		return 0, 0, err
	}
	return r.TermsUsed, r.RemainderBound, nil
}

func pressureAt(in Input, path DrainagePath, tv, z float64) (float64, int, float64, error) {
	decay := decayRate(tv)
	c1, c2 := modalEnvelope(in.Initial, path)
	n1, b1, err := termsNeeded(series.Term{Coefficient: c1, Power: 1, Decay: decay})
	if err != nil {
		return 0, 0, 0, err
	}
	n, bound := n1, b1
	if c2 != 0 {
		n2, b2, err := termsNeeded(series.Term{Coefficient: c2, Power: 2, Decay: decay})
		if err != nil {
			return 0, 0, 0, err
		}
		if n2 > n {
			n = n2
		}
		if b2 > bound {
			bound = b2
		}
	}
	u := 0.0
	for i := 0; i < n; i++ {
		a1, a2 := modalAmplitudes(in.Initial, path, i)
		k := float64(2*i + 1)
		u += (a1/k + a2/(k*k)) * math.Sin(path.Eigenvalue(i)*z) * math.Exp(-k*k*decay)
	}
	return u, n, bound, nil
}

func meanPressureAt(in Input, path DrainagePath, tv float64) (float64, int, float64, error) {
	decay := decayRate(tv)
	scale := 2 / math.Pi
	c1, c2 := modalEnvelope(in.Initial, path)
	n1, b1, err := termsNeeded(series.Term{Coefficient: scale * c1, Power: 2, Decay: decay})
	if err != nil {
		return 0, 0, 0, err
	}
	n, bound := n1, b1
	if c2 != 0 {
		n2, b2, err := termsNeeded(series.Term{Coefficient: scale * c2, Power: 3, Decay: decay})
		if err != nil {
			return 0, 0, 0, err
		}
		if n2 > n {
			n = n2
		}
		if b2 > bound {
			bound = b2
		}
	}
	mean := 0.0
	for i := 0; i < n; i++ {
		a1, a2 := modalAmplitudes(in.Initial, path, i)
		k := float64(2*i + 1)
		mean += scale * (a1/(k*k) + a2/(k*k*k)) * math.Exp(-k*k*decay)
	}
	return mean, n, bound, nil
}

func pressureProfile(in Input, path DrainagePath, tv float64, nodes int) ([]Point, int, float64, error) {
	mean0 := in.Initial.meanPressure()
	pts := make([]Point, 0, nodes)
	terms, bound := 0, 0.0
	for i := 0; i < nodes; i++ {
		z := float64(i) * in.Thickness / float64(nodes-1)
		u, n, b, err := pressureAt(in, path, tv, z)
		if err != nil {
			return nil, 0, 0, err
		}
		if n > terms {
			terms = n
		}
		if b > bound {
			bound = b
		}
		pts = append(pts, Point{
			DepthFraction: float64(i) / float64(nodes-1),
			Depth:         z,
			Pressure:      u,
			Remaining:     u / mean0,
			Dissipated:    1 - u/mean0,
		})
	}
	for i := range pts {
		if path.Drained(pts[i].DepthFraction) {
			if math.Abs(pts[i].Pressure) > 100*SeriesTolerance*mean0 {
				return nil, 0, 0, fmt.Errorf(
					"drained face at z/H=%g does not vanish: u=%g kPa",
					pts[i].DepthFraction, pts[i].Pressure)
			}
		}
	}
	return pts, terms, bound, nil
}
