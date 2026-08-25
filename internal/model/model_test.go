package model_test

import (
	"math"
	"testing"

	"consol-1d/internal/model"
)

func TestZeroTimeFactorReturnsInitialState(t *testing.T) {
	in := model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 0)
	res, err := model.Solve(in, 11)
	if err != nil {
		t.Fatalf("Solve(t=0): expected no error, got %v", err)
	}
	if res.Tv != 0 {
		t.Errorf("time factor: expected 0, got %g", res.Tv)
	}
	if res.U != 0 {
		t.Errorf("average consolidation: expected 0, got %g", res.U)
	}
	if res.MidpointPressure != 100 {
		t.Errorf("midpoint pressure: expected 100, got %g", res.MidpointPressure)
	}
	if res.MeanPressure != 100 {
		t.Errorf("mean pressure: expected 100, got %g", res.MeanPressure)
	}
	if res.TermsUsed != 0 || res.RemainderBound != 0 {
		t.Errorf("zero-time series: expected no terms and zero bound, got %d and %g",
			res.TermsUsed, res.RemainderBound)
	}
	for i, p := range res.Profile {
		if p.Pressure != 100 {
			t.Errorf("profile point %d pressure: expected 100, got %g", i, p.Pressure)
		}
		if p.Dissipated != 0 {
			t.Errorf("profile point %d dissipated: expected 0, got %g", i, p.Dissipated)
		}
	}
}

func TestLargeTimeFactorReachesFullConsolidation(t *testing.T) {
	in := model.UniformInput(1e-5, 10, model.DrainageDouble, 100, 5e7)
	res, err := model.Solve(in, 11)
	if err != nil {
		t.Fatalf("Solve(large t): expected no error, got %v", err)
	}
	if want := 1.0; math.Abs(res.U-want) > 1e-6 {
		t.Errorf("average consolidation at Tv=20: expected %g, got %g", want, res.U)
	}
	if res.MidpointPressure > 1e-3 {
		t.Errorf("midpoint pressure at Tv=20: expected near 0, got %g", res.MidpointPressure)
	}
	for i, p := range res.Profile {
		if p.Pressure > 1e-3 {
			t.Errorf("profile point %d pressure at Tv=20: expected near 0, got %g", i, p.Pressure)
		}
	}
}

func TestDoubleCvDoublesTimeFactor(t *testing.T) {
	base := model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 1e6)
	res1, err := model.Solve(base, 11)
	if err != nil {
		t.Fatalf("Solve(base): expected no error, got %v", err)
	}
	big := base
	big.Cv = 2e-7
	res2, err := model.Solve(big, 11)
	if err != nil {
		t.Fatalf("Solve(double cv): expected no error, got %v", err)
	}
	if want := 2 * res1.Tv; math.Abs(res2.Tv-want) > 1e-12*math.Max(1, want) {
		t.Errorf("doubled cv time factor: expected %g, got %g", want, res2.Tv)
	}
	twiceTime := base
	twiceTime.Time = 2e6
	res3, err := model.Solve(twiceTime, 11)
	if err != nil {
		t.Fatalf("Solve(double time): expected no error, got %v", err)
	}
	if want := res2.U; math.Abs(res3.U-want) > 1e-9 {
		t.Errorf("same Tv through cv and through t: expected U %g, got %g", want, res3.U)
	}
}

func TestDoubleDrainFasterThanSingle(t *testing.T) {
	dbl := model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 5e7)
	sgl := model.UniformInput(1e-7, 10, model.DrainageSingle, 100, 5e7)
	resD, err := model.Solve(dbl, 11)
	if err != nil {
		t.Fatalf("Solve(double): expected no error, got %v", err)
	}
	resS, err := model.Solve(sgl, 11)
	if err != nil {
		t.Fatalf("Solve(single): expected no error, got %v", err)
	}
	if resD.Hdr != 5 || resS.Hdr != 10 {
		t.Errorf("drainage distances: expected 5 and 10, got %g and %g", resD.Hdr, resS.Hdr)
	}
	if want := 4 * resS.Tv; math.Abs(resD.Tv-want) > 1e-12*math.Max(1, want) {
		t.Errorf("double vs single Tv: expected %g, got %g", want, resD.Tv)
	}
	if resD.U <= resS.U {
		t.Errorf("double drainage should consolidate faster: double U=%g, single U=%g", resD.U, resS.U)
	}
	if want := resS.MidpointPressure - resD.MidpointPressure; want < 5 {
		t.Errorf("double drainage midpoint should be clearly lower: double=%g, single=%g (gap %g)",
			resD.MidpointPressure, resS.MidpointPressure, want)
	}
}

func TestMidpointDissipatesSlowerThanDrainageFace(t *testing.T) {
	in := model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 5e7)
	res, err := model.Solve(in, 11)
	if err != nil {
		t.Fatalf("Solve: expected no error, got %v", err)
	}
	if face := res.DrainedFacePressure(); face > 1e-6 {
		t.Errorf("drained face pressure: expected 0, got %g", face)
	}
	if res.MidpointPressure <= 0.1*100 {
		t.Errorf("midpoint pressure: expected to stay well above the drained face, got %g", res.MidpointPressure)
	}
	if res.MidpointPressure >= 100 {
		t.Errorf("midpoint pressure: expected below the initial value, got %g", res.MidpointPressure)
	}
}

func TestValidationRejectsIllegalInput(t *testing.T) {
	valid := model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 5e7)
	cases := []struct {
		name string
		mut  func(*model.Input)
	}{
		{"zero cv", func(in *model.Input) { in.Cv = 0 }},
		{"negative cv", func(in *model.Input) { in.Cv = -1e-7 }},
		{"zero thickness", func(in *model.Input) { in.Thickness = 0 }},
		{"negative thickness", func(in *model.Input) { in.Thickness = -3 }},
		{"negative time", func(in *model.Input) { in.Time = -1 }},
		{"negative uniform pressure", func(in *model.Input) { in.Initial.U0 = -100 }},
		{"negative linear pressure", func(in *model.Input) {
			in.Initial = model.InitialPressure{Type: model.InitialLinear, UA: 100, UB: -10}
		}},
		{"unknown drainage", func(in *model.Input) { in.Drainage = "both" }},
		{"unknown profile type", func(in *model.Input) {
			in.Initial = model.InitialPressure{Type: "parabola", U0: 100}
		}},
	}
	for _, c := range cases {
		in := valid
		c.mut(&in)
		if _, err := model.Solve(in, 11); err == nil {
			t.Errorf("%s: expected an error, got nil", c.name)
		}
	}
}

func TestSettlementIsUtimesUltimate(t *testing.T) {
	mv, ds := 0.0002, 100.0
	in := model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 5e7)
	in.Mv = &mv
	in.DeltaSigma = &ds
	res, err := model.Solve(in, 11)
	if err != nil {
		t.Fatalf("Solve: expected no error, got %v", err)
	}
	if want := 0.2; math.Abs(res.UltimateSettlement-want) > 1e-12 {
		t.Errorf("ultimate settlement: expected %g, got %g", want, res.UltimateSettlement)
	}
	if want := res.U * 0.2; math.Abs(res.Settlement-want) > 1e-12 {
		t.Errorf("settlement: expected %g, got %g", want, res.Settlement)
	}
	if res.SettlementRatio != res.U {
		t.Errorf("settlement ratio: expected %g (== U), got %g", res.U, res.SettlementRatio)
	}

	bare := model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 5e7)
	resBare, err := model.Solve(bare, 11)
	if err != nil {
		t.Fatalf("Solve(no mv): expected no error, got %v", err)
	}
	if !math.IsNaN(resBare.Settlement) || !math.IsNaN(resBare.UltimateSettlement) {
		t.Errorf("absent mv/delta_sigma: expected NaN settlements, got %g and %g",
			resBare.Settlement, resBare.UltimateSettlement)
	}
	if resBare.SettlementRatio != resBare.U {
		t.Errorf("settlement ratio without mv: expected %g, got %g", resBare.U, resBare.SettlementRatio)
	}
}

func TestDoubleDrainMatchesTextbook(t *testing.T) {
	in := model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 5e7)
	res, err := model.Solve(in, 11)
	if err != nil {
		t.Fatalf("Solve: expected no error, got %v", err)
	}
	if want := 0.2; math.Abs(res.Tv-want) > 1e-12 {
		t.Errorf("time factor: expected %g, got %g", want, res.Tv)
	}
	if res.U < 0.45 || res.U > 0.55 {
		t.Errorf("textbook window 0.5 +/- 0.05: got U=%g", res.U)
	}
	if math.Abs(res.U-0.504) > 0.01 {
		t.Errorf("textbook U at Tv=0.2: expected ~0.504, got %g", res.U)
	}
}

func TestHalfThicknessClassicalSeriesAgree(t *testing.T) {
	in := model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 5e7)
	res, err := model.Solve(in, 11)
	if err != nil {
		t.Fatalf("Solve: expected no error, got %v", err)
	}
	hdr := in.Thickness / 2
	wantTv := in.Cv * in.Time / (hdr * hdr)
	if math.Abs(res.Tv-wantTv) > 1e-12*math.Max(1, wantTv) {
		t.Errorf("double drainage Tv: expected %g (Hdr=H/2), got %g", wantTv, res.Tv)
	}
	classic, _, _, err := model.ClassicalAverageU(wantTv, 1e-12, 1<<25)
	if err != nil {
		t.Fatalf("ClassicalAverageU: expected no error, got %v", err)
	}
	if math.Abs(res.U-classic) > 1e-6 {
		t.Errorf("integrated U vs classical series: expected %g, got %g", classic, res.U)
	}
}

func TestLinearInitialReducesToUniform(t *testing.T) {
	uni := model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 5e7)
	lin := model.Input{
		Cv:        1e-7,
		Thickness: 10,
		Drainage:  model.DrainageDouble,
		Initial:   model.InitialPressure{Type: model.InitialLinear, UA: 100, UB: 100},
		Time:      5e7,
	}
	resU, err := model.Solve(uni, 11)
	if err != nil {
		t.Fatalf("Solve(uniform): expected no error, got %v", err)
	}
	resL, err := model.Solve(lin, 11)
	if err != nil {
		t.Fatalf("Solve(linear): expected no error, got %v", err)
	}
	if math.Abs(resL.U-resU.U) > 1e-9 {
		t.Errorf("linear(ua=ub) U: expected %g, got %g", resU.U, resL.U)
	}
	if math.Abs(resL.MidpointPressure-resU.MidpointPressure) > 1e-9 {
		t.Errorf("linear(ua=ub) midpoint: expected %g, got %g", resU.MidpointPressure, resL.MidpointPressure)
	}
}

func TestNumericalCoefficientMatchesAnalytic(t *testing.T) {
	ip := model.InitialPressure{Type: model.InitialLinear, UA: 120, UB: 60}
	paths := []model.DrainagePath{
		model.NewDrainagePath(model.DrainageDouble, 10),
		model.NewDrainagePath(model.DrainageSingle, 10),
	}
	for _, path := range paths {
		for n := 0; n < 5; n++ {
			analytic := model.FourierCoefficient(ip, path, n)
			numeric, err := model.NumericalCoefficient(ip, path, n)
			if err != nil {
				t.Fatalf("NumericalCoefficient(%s, n=%d): expected no error, got %v", path.Kind, n, err)
			}
			if math.Abs(analytic-numeric) > 1e-6 {
				t.Errorf("coefficient %s n=%d: analytic %g, numeric %g", path.Kind, n, analytic, numeric)
			}
		}
	}
}

func TestUniformDoubleDrainProfileIsSymmetric(t *testing.T) {
	in := model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 5e7)
	res, err := model.Solve(in, 21)
	if err != nil {
		t.Fatalf("Solve: expected no error, got %v", err)
	}
	last := len(res.Profile) - 1
	for i := 0; i <= last; i++ {
		j := last - i
		if math.Abs(res.Profile[i].Pressure-res.Profile[j].Pressure) > 1e-9 {
			t.Errorf("asymmetry at z/H=%g: %g vs %g", res.Profile[i].DepthFraction,
				res.Profile[i].Pressure, res.Profile[j].Pressure)
		}
	}
}

func TestTimeToDegreeMatchesDirectSolve(t *testing.T) {
	in := model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 5e7)
	t90, err := model.TimeToDegree(in, 0.9, 1e-9)
	if err != nil {
		t.Fatalf("TimeToDegree(0.9): expected no error, got %v", err)
	}
	at90 := in
	at90.Time = t90
	res, err := model.Solve(at90, 11)
	if err != nil {
		t.Fatalf("Solve(at t90): expected no error, got %v", err)
	}
	if math.Abs(res.U-0.9) > 1e-4 {
		t.Errorf("degree at returned t90: expected ~0.9, got %g", res.U)
	}
	if t90 <= 0 || math.IsNaN(t90) || math.IsInf(t90, 0) {
		t.Errorf("t90 should be finite and positive, got %g", t90)
	}
}

func TestConsolidationCurveIsMonotone(t *testing.T) {
	in := model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 5e7)
	times := []float64{1e4, 1e5, 1e6, 5e7}
	pts, err := model.ConsolidationCurve(in, times, 5)
	if err != nil {
		t.Fatalf("ConsolidationCurve: expected no error, got %v", err)
	}
	if len(pts) != len(times) {
		t.Fatalf("curve points: expected %d, got %d", len(times), len(pts))
	}
	for i := 1; i < len(pts); i++ {
		if pts[i].U < pts[i-1].U {
			t.Errorf("curve must be non-decreasing: t=%g U=%g dropped from %g",
				pts[i].Time, pts[i].U, pts[i-1].U)
		}
	}
	single, err := model.Solve(in, 5)
	if err != nil {
		t.Fatalf("Solve: expected no error, got %v", err)
	}
	last := pts[len(pts)-1]
	if math.Abs(last.U-single.U) > 1e-9 {
		t.Errorf("curve last point should match a direct solve: %g vs %g", last.U, single.U)
	}
}
