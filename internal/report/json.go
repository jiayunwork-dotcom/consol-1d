package report

import (
	"encoding/json"

	"consol-1d/internal/model"
)

// JSONView is the machine-readable form of a consolidation result. The
// field names are stable so external tooling can rely on them; NaN
// values for the optional settlement fields stay absent via omitempty
// on the pointers.
type JSONView struct {
	Cv                   float64     `json:"cv"`
	Thickness            float64     `json:"thickness"`
	Drainage             string      `json:"drainage"`
	DrainageDistance     float64     `json:"drainage_distance"`
	Time                 float64     `json:"time"`
	TimeFactor           float64     `json:"time_factor"`
	AverageConsolidation float64     `json:"average_consolidation"`
	MidpointPressure     float64     `json:"midpoint_pressure"`
	MeanPressure         float64     `json:"mean_pressure"`
	MeanInitialPressure  float64     `json:"mean_initial_pressure"`
	Settlement           *float64    `json:"settlement,omitempty"`
	UltimateSettlement   *float64    `json:"ultimate_settlement,omitempty"`
	SettlementRatio      float64     `json:"settlement_ratio"`
	SeriesTerms          int         `json:"series_terms"`
	RemainderBound       float64     `json:"remainder_bound"`
	Profile              []JSONPoint `json:"profile"`
}

// JSONPoint is one profile row in JSON form.
type JSONPoint struct {
	DepthFraction float64 `json:"z_over_H"`
	Depth         float64 `json:"z_m"`
	Pressure      float64 `json:"u_kPa"`
	Remaining     float64 `json:"u_over_u0"`
	Dissipated    float64 `json:"dissipated"`
}

// ToJSON converts a result into the stable JSON document. The optional
// settlement values are exposed as pointers so absent mv/delta_sigma
// yield omitted keys instead of null.
func ToJSON(res model.Result) ([]byte, error) {
	view := JSONView{
		Cv:                   res.Input.Cv,
		Thickness:            res.Input.Thickness,
		Drainage:             DrainageLabel(res.Input.Drainage),
		DrainageDistance:     res.Hdr,
		Time:                 res.Input.Time,
		TimeFactor:           res.Tv,
		AverageConsolidation: res.U,
		MidpointPressure:     res.MidpointPressure,
		MeanPressure:         res.MeanPressure,
		MeanInitialPressure:  res.MeanInitialPressure,
		SettlementRatio:      res.SettlementRatio,
		SeriesTerms:          res.TermsUsed,
		RemainderBound:       res.RemainderBound,
		Profile:              make([]JSONPoint, 0, len(res.Profile)),
	}
	if !isNaN(res.Settlement) {
		s := res.Settlement
		view.Settlement = &s
	}
	if !isNaN(res.UltimateSettlement) {
		s := res.UltimateSettlement
		view.UltimateSettlement = &s
	}
	for _, p := range res.Profile {
		view.Profile = append(view.Profile, JSONPoint{
			DepthFraction: p.DepthFraction,
			Depth:         p.Depth,
			Pressure:      p.Pressure,
			Remaining:     p.Remaining,
			Dissipated:    p.Dissipated,
		})
	}
	return json.MarshalIndent(view, "", "  ")
}

func isNaN(v float64) bool {
	return v != v
}
