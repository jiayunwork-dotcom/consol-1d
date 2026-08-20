package report

import (
	"fmt"
	"strings"

	"consol-1d/internal/model"
)

// FormatPoint renders one profile row as the fixed-width text line used
// by ProfileTable: normalized depth, depth, pressure, remaining
// fraction and dissipated percentage.
func FormatPoint(p model.Point) string {
	return fmt.Sprintf("  %-6.2f  %-8.3f  %-9.3f  %-6.3f  %s\n",
		p.DepthFraction, p.Depth, p.Pressure, p.Remaining, Percent(p.Dissipated))
}

// ProfileTable renders the dissipation profile as aligned text columns:
// normalized depth, depth in metres, remaining pressure, the fraction
// of the initial mean pressure still standing and the fraction already
// dissipated. The drainage faces show up as full dissipation rows.
func ProfileTable(points []model.Point) string {
	var b strings.Builder
	b.WriteString("  z/H     z [m]    u [kPa]   u/u0    dissipated\n")
	for _, p := range points {
		b.WriteString(FormatPoint(p))
	}
	return b.String()
}

// Header builds the multi-line preamble of the human report: the
// scenario parameters and the primary results.
func Header(res model.Result) string {
	var b strings.Builder
	drainage := "single"
	if res.Input.Drainage == model.DrainageDouble {
		drainage = "double"
	}
	fmt.Fprintf(&b, "consol-1d: one-dimensional consolidation\n")
	fmt.Fprintf(&b, "  drainage       : %s, H = %.4g m, Hdr = %.4g m\n",
		drainage, res.Input.Thickness, res.Hdr)
	fmt.Fprintf(&b, "  cv             : %s m^2/s\n", Sci(res.Input.Cv))
	fmt.Fprintf(&b, "  time           : %s s (%s)\n", Sci(res.Input.Time), Duration(res.Input.Time))
	fmt.Fprintf(&b, "  time factor Tv : %s   (Tv = cv*t/Hdr^2)\n", Sci(res.Tv))
	fmt.Fprintf(&b, "  average U      : %s   (1 - mean u / u0)\n", Decimal(res.U))
	fmt.Fprintf(&b, "  midpoint u     : %s kPa   (u/u0 = %s)\n",
		Decimal(res.MidpointPressure), Decimal(res.MidpointRemaining()))
	return b.String()
}

// SettlementLine renders the settlement block. When mv and delta_sigma
// were absent the absolute settlements read "n/a" and the ratio still
// equals U.
func SettlementLine(res model.Result) string {
	return fmt.Sprintf("  settlement     : s = %s m, s_ult = %s m, s/s_ult = %s\n",
		Decimal(res.Settlement), Decimal(res.UltimateSettlement), Decimal(res.SettlementRatio))
}

// SeriesLine renders the truncation proof: the number of modes kept and
// the verified bound on the discarded tail.
func SeriesLine(res model.Result) string {
	return fmt.Sprintf("  series         : %d terms kept, remainder bound <= %s\n",
		res.TermsUsed, Sci(res.RemainderBound))
}
