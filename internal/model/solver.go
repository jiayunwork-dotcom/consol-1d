package model

import (
	"fmt"
	"math"
)

type Point struct {
	DepthFraction float64
	Depth         float64
	Pressure      float64
	Remaining     float64
	Dissipated    float64
}

type Result struct {
	Input               Input
	DrainagePath        DrainagePath
	Hdr                 float64
	Tv                  float64
	U                   float64
	MidpointPressure    float64
	MeanPressure        float64
	MeanInitialPressure float64
	Settlement          float64
	UltimateSettlement  float64
	SettlementRatio     float64
	Profile             []Point
	TermsUsed           int
	RemainderBound      float64
}

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

func (r Result) MidpointRemaining() float64 {
	return r.MidpointPressure / r.MeanInitialPressure
}

func (r Result) DrainedFacePressure() float64 {
	if len(r.Profile) == 0 {
		return math.NaN()
	}
	if r.DrainagePath.Kind == DrainageSingle {
		return r.Profile[0].Pressure
	}
	return math.Max(math.Abs(r.Profile[0].Pressure), math.Abs(r.Profile[len(r.Profile)-1].Pressure))
}
