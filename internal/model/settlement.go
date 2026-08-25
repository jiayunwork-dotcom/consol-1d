package model

import "fmt"

func UltimateSettlement(in Input) (float64, error) {
	if in.Mv == nil || in.DeltaSigma == nil {
		return 0, fmt.Errorf("%w (both mv and delta_sigma are required)", ErrSettlement)
	}
	if err := in.Validate(); err != nil {
		return 0, err
	}
	return (*in.Mv) * (*in.DeltaSigma) * in.Thickness, nil
}

func SettlementAt(in Input) (s, sult float64, err error) {
	res, err := Solve(in, 2)
	if err != nil {
		return 0, 0, err
	}
	return res.Settlement, res.UltimateSettlement, nil
}

func SettlementRatioLabel(ratio float64) string {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	return fmt.Sprintf("%.1f%%", ratio*100)
}
