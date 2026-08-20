package model

import (
	"fmt"
	"math"
)

// Point is one row of the dissipation profile: the depth, the remaining
// excess pore pressure u, the fraction still remaining u/u_bar(0) and
// the fraction already dissipated.
type Point struct {
	DepthFraction float64 // z/H, from 0 (top face) to 1 (base)
	Depth         float64 // z in metres
	Pressure      float64 // u(z, t) in kPa
	Remaining     float64 // u / u_bar(0)
	Dissipated    float64 // 1 - u / u_bar(0)
}

// Result is the full outcome of one consolidation check.
type Result struct {
	Input               Input
	DrainagePath        DrainagePath
	Hdr                 float64 // drainage distance in metres
	Tv                  float64 // time factor cv*t/Hdr^2
	U                   float64 // average consolidation degree
	MidpointPressure    float64 // u at z = H/2 in kPa
	MeanPressure        float64 // layer-averaged u(z, t) in kPa
	MeanInitialPressure float64 // layer-averaged u(z, 0) in kPa
	Settlement          float64 // s = U*s_ult in metres (NaN when mv/delta_sigma absent)
	UltimateSettlement  float64 // s_ult = mv*delta_sigma*H in metres (NaN when absent)
	SettlementRatio     float64 // s/s_ult, always equal to U
	Profile             []Point
	TermsUsed           int // widest truncation among all series
	RemainderBound      float64
}

// Solve runs one consolidation check. The scenario is validated first
// (every illegal parameter returns an error), then the time factor and
// the Fourier coefficients are computed once, and the profile, the
// layer average, the consolidation degree and the settlement are all
// derived from that single set of modes.
//
// nodes selects how many rows the dissipation profile carries; the
// drainage faces are always included at z/H = 0 and z/H = 1.
func Solve(in Input, nodes int) (Result, error) {
	if err := in.Validate(); err != nil {
		return Result{}, err
	}
	if nodes < 2 {
		return Result{}, fmt.Errorf("%w (got %d)", ErrNodes, nodes)
	}
	path := NewDrainagePath(in.Drainage, in.Thickness)
	tv := path.TimeFactor(in.Cv, in.Time)
	mean0 := in.Initial.meanPressure()

	res := Result{
		Input:               in,
		DrainagePath:        path,
		Hdr:                 path.Hdr,
		Tv:                  tv,
		MeanInitialPressure: mean0,
	}

	if tv == 0 {
		// Zero time factor: the excess pore pressure is exactly the
		// initial profile, no consolidation has happened, and U is
		// exactly zero. No series is evaluated at all.
		res.U = 0
		res.MeanPressure = mean0
		res.MidpointPressure = in.Initial.PressureAt(in.Thickness/2, in.Thickness)
		res.Profile = buildInitialProfile(in, nodes, mean0)
		res.TermsUsed = 0
		res.RemainderBound = 0
	} else {
		pts, terms, bound, err := pressureProfile(in, path, tv, nodes)
		if err != nil {
			return Result{}, err
		}
		mean, meanTerms, meanBound, err := meanPressureAt(in, path, tv)
		if err != nil {
			return Result{}, err
		}
		mid, midTerms, midBound, err := pressureAt(in, path, tv, in.Thickness/2)
		if err != nil {
			return Result{}, err
		}
		res.Profile = pts
		res.MeanPressure = mean
		res.MidpointPressure = mid
		res.U = 1 - mean/mean0
		res.TermsUsed = terms
		if meanTerms > res.TermsUsed {
			res.TermsUsed = meanTerms
		}
		if midTerms > res.TermsUsed {
			res.TermsUsed = midTerms
		}
		res.RemainderBound = bound
		if meanBound > res.RemainderBound {
			res.RemainderBound = meanBound
		}
		if midBound > res.RemainderBound {
			res.RemainderBound = midBound
		}
	}

	res.Settlement, res.UltimateSettlement, res.SettlementRatio = settlementFrom(in, res.U)
	return res, nil
}

// buildInitialProfile returns the exact initial pressure profile at
// t = 0 (or equivalently Tv = 0), which every zero-time solve reuses.
func buildInitialProfile(in Input, nodes int, mean0 float64) []Point {
	pts := make([]Point, 0, nodes)
	for i := 0; i < nodes; i++ {
		z := float64(i) * in.Thickness / float64(nodes-1)
		u := in.Initial.PressureAt(z, in.Thickness)
		pts = append(pts, Point{
			DepthFraction: float64(i) / float64(nodes-1),
			Depth:         z,
			Pressure:      u,
			Remaining:     u / mean0,
			Dissipated:    0,
		})
	}
	return pts
}

// MidpointRemaining is the fraction of the initial mean pressure still
// standing at the layer midpoint; tests use it as the standard reading
// of how far dissipation has progressed at z = H/2.
func (r Result) MidpointRemaining() float64 {
	return r.MidpointPressure / r.MeanInitialPressure
}

// DrainedFacePressure returns the (numerical) pressure at the drained
// boundary for the given path; it should be indistinguishable from zero
// after any successful solve.
func (r Result) DrainedFacePressure() float64 {
	if len(r.Profile) == 0 {
		return math.NaN()
	}
	if r.DrainagePath.Kind == DrainageSingle {
		return r.Profile[0].Pressure
	}
	return math.Max(math.Abs(r.Profile[0].Pressure), math.Abs(r.Profile[len(r.Profile)-1].Pressure))
}
