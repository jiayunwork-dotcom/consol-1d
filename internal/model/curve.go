package model

import (
	"fmt"
	"math"
)

type CurvePoint struct {
	Time               float64
	TimeFactor         float64
	U                  float64
	MidpointPressure   float64
	Settlement         float64
	UltimateSettlement float64
}

func ConsolidationCurve(in Input, times []float64, nodes int) ([]CurvePoint, error) {
	if len(times) == 0 {
		return nil, fmt.Errorf("curve needs at least one time")
	}
	out := make([]CurvePoint, 0, len(times))
	for _, t := range times {
		cp := in
		cp.Time = t
		res, err := Solve(cp, nodes)
		if err != nil {
			return nil, err
		}
		out = append(out, CurvePoint{
			Time:               t,
			TimeFactor:         res.Tv,
			U:                  res.U,
			MidpointPressure:   res.MidpointPressure,
			Settlement:         res.Settlement,
			UltimateSettlement: res.UltimateSettlement,
		})
	}
	return out, nil
}

func DefaultCurveTimes(in Input) []float64 {
	path := NewDrainagePath(in.Drainage, in.Thickness)
	center := 0.197 * path.Hdr * path.Hdr / in.Cv
	if center <= 0 {
		return []float64{1}
	}
	times := make([]float64, 0, 10)
	for i := -4; i <= 5; i++ {
		times = append(times, center*math.Pow(10, float64(i)))
	}
	return times
}

func degreeAt(in Input, tv float64) float64 {
	path := NewDrainagePath(in.Drainage, in.Thickness)
	tmp := in
	tmp.Time = tv * path.Hdr * path.Hdr / in.Cv
	res, err := Solve(tmp, 2)
	if err != nil {
		return math.NaN()
	}
	return res.U
}

func TimeToDegree(in Input, target, tvTol float64) (float64, error) {
	if err := in.Validate(); err != nil {
		return 0, err
	}
	if math.IsNaN(target) || target <= 0 || target >= 1 {
		return 0, fmt.Errorf("target consolidation degree must lie in (0, 1), got %g", target)
	}
	if math.IsNaN(tvTol) || tvTol <= 0 {
		return 0, fmt.Errorf("time-factor tolerance must be positive, got %g", tvTol)
	}
	tvHi := 1.0
	uHi := degreeAt(in, tvHi)
	for uHi < target && tvHi < 1e9 {
		tvHi *= 2
		uHi = degreeAt(in, tvHi)
	}
	if uHi < target {
		return 0, fmt.Errorf("degree %g not reached by Tv=%g (U only %g)", target, tvHi, uHi)
	}
	lo, hi := 0.0, tvHi
	for i := 0; i < 100 && hi-lo > tvTol; i++ {
		mid := 0.5 * (lo + hi)
		if degreeAt(in, mid) < target {
			lo = mid
		} else {
			hi = mid
		}
	}
	tv := 0.5 * (lo + hi)
	path := NewDrainagePath(in.Drainage, in.Thickness)
	return tv * path.Hdr * path.Hdr / in.Cv, nil
}
