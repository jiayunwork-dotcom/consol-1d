package report

import (
	"fmt"
	"strings"

	"consol-1d/internal/model"
)

// CurveTable renders the consolidation curve as an aligned text table:
// time, time factor, average consolidation degree, midpoint pressure
// and settlement per row.
func CurveTable(pts []model.CurvePoint) string {
	var b strings.Builder
	b.WriteString("  t [s]         Tv          U        u_mid [kPa]   s [m]\n")
	for _, p := range pts {
		fmt.Fprintf(&b, "  %-13s  %-10s  %-8s  %-12s  %s\n",
			Sci(p.Time), Sci(p.TimeFactor), Decimal(p.U),
			Decimal(p.MidpointPressure), Decimal(p.Settlement))
	}
	return b.String()
}

// CurveSummary prints the first and last consolidation degrees of a
// curve so a script can see at a glance how far the batch of times
// spans the consolidation process.
func CurveSummary(pts []model.CurvePoint) string {
	if len(pts) == 0 {
		return "curve: no points"
	}
	first, last := pts[0], pts[len(pts)-1]
	return fmt.Sprintf("curve: t=%s -> U=%s; t=%s -> U=%s",
		Sci(first.Time), Decimal(first.U), Sci(last.Time), Decimal(last.U))
}

// SettlementLineText is a variant of SettlementLine that reads from a
// curve point; it keeps the header/series blocks in the profile report
// and the curve report from sharing mutable state.
func SettlementLineText(s, sult, ratio float64) string {
	return fmt.Sprintf("  settlement     : s = %s m, s_ult = %s m, s/s_ult = %s\n",
		Decimal(s), Decimal(sult), Decimal(ratio))
}
