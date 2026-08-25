package report

import (
	"fmt"
	"strings"

	"consol-1d/internal/model"
)

func CSV(points []model.Point) string {
	var b strings.Builder
	b.WriteString("z_over_H,z_m,u_kPa,u_over_u0,dissipated\n")
	for _, p := range points {
		fmt.Fprintf(&b, "%.6g,%.6g,%.6g,%.6g,%.6g\n",
			p.DepthFraction, p.Depth, p.Pressure, p.Remaining, p.Dissipated)
	}
	return b.String()
}

func CSVOfResult(res model.Result) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# consol-1d: drainage=%s H=%g Hdr=%g cv=%g t=%g Tv=%g U=%g\n",
		DrainageLabel(res.Input.Drainage), res.Input.Thickness, res.Hdr,
		res.Input.Cv, res.Input.Time, res.Tv, res.U)
	fmt.Fprintf(&b, "# midpoint_u_kPa=%g settlement_ratio=%g terms=%d remainder_bound<=%g\n",
		res.MidpointPressure, res.SettlementRatio, res.TermsUsed, res.RemainderBound)
	b.WriteString(CSV(res.Profile))
	return HoldCSVLive(b.String())
}
