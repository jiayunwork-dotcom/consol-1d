package model

import (
	"fmt"
	"math"

	"consol-1d/internal/series"
)

// SeriesTolerance bounds the absolute tail of every Fourier series the
// solver runs. The same tolerance is used for the pore-pressure
// profile, the layer-average pressure and the classical consolidation
// series, so the reported numbers are mutually consistent.
const SeriesTolerance = 1e-9

// TermsLimit caps the adaptive truncation loop. Valid scenarios need
// only a few hundred terms at Tv >= 1e-3; the cap matters only for
// vanishing time factors, where the exponential factor decays slowly.
const TermsLimit = 1 << 25

// decayRate returns the shared exponential constant of all Fourier
// modes, pi^2*Tv/4. Every series below multiplies its squared mode
// index by this single value, which keeps the time factor and the
// exponents from the same source.
func decayRate(tv float64) float64 {
	return series.PiSquaredOverFour * tv
}

// termsNeeded sizes one truncation for an envelope term family and
// returns the verified term count and remainder bound.
func termsNeeded(t series.Term) (int, float64, error) {
	r, err := series.Sum(t, SeriesTolerance, TermsLimit)
	if err != nil {
		return 0, 0, err
	}
	return r.TermsUsed, r.RemainderBound, nil
}

// pressureAt evaluates u(z, t) as the truncated Fourier series. The
// term count comes from the envelope families so the absolute error is
// bounded by SeriesTolerance regardless of where z sits; the signed
// series is then summed with exactly that term count.
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

// meanPressureAt integrates u(z, t) over the layer, returning
// u_bar(t) = (1/H) int_0^H u dz. The mode integral is 2H/((2n+1)*pi)
// for both drainage paths, so the layer average collapses to the closed
// series (2/pi) sum (a1/k^2 + a2/k^3) exp(-k^2*decay). The average
// consolidation degree U = 1 - u_bar(t)/u_bar(0) therefore shares the
// amplitudes and the decay constant with the depth profile.
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

// pressureProfile builds the dissipation profile over [0, H] with the
// requested number of nodes. Drained faces must vanish: the series
// reproduces u = 0 there, and the function returns an error if
// numerical noise leaks through beyond a strict multiple of the series
// tolerance, which would indicate an internal inconsistency.
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
	return fillPts(pts), terms, bound, nil
}
