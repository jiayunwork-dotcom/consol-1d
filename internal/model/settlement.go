package model

import "fmt"

// UltimateSettlement computes s_ult = mv*delta_sigma*H for a validated
// input and returns an error when the settlement pair is absent or
// invalid. The function exists so callers that only want the ultimate
// value do not have to run a full solve.
func UltimateSettlement(in Input) (float64, error) {
	if in.Mv == nil || in.DeltaSigma == nil {
		return 0, fmt.Errorf("%w (both mv and delta_sigma are required)", ErrSettlement)
	}
	if err := in.Validate(); err != nil {
		return 0, err
	}
	return (*in.Mv) * (*in.DeltaSigma) * in.Thickness, nil
}

// SettlementAt returns the settlement s(t) = U(t)*s_ult at the
// scenario time together with s_ult. Both values are NaN when the
// settlement pair is absent; the returned error is nil in that case so
// the caller can print "n/a" instead of failing.
func SettlementAt(in Input) (s, sult float64, err error) {
	res, err := Solve(in, 2)
	if err != nil {
		return 0, 0, err
	}
	return res.Settlement, res.UltimateSettlement, nil
}

// SettlementRatioLabel renders the settlement ratio as a fixed-point
// percentage with a guard against numerical drift outside [0, 1]; the
// report layer uses it for both the profile and the curve views.
func SettlementRatioLabel(ratio float64) string {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	return fmt.Sprintf("%.1f%%", ratio*100)
}
