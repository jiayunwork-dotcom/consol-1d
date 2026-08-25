package report

import (
	"fmt"
	"strings"

	"consol-1d/internal/model"
)

func FormatPoint(p model.Point) string {
	return fmt.Sprintf("  %-6.2f  %-8.3f  %-9.3f  %-6.3f  %s\n",
		p.DepthFraction, p.Depth, p.Pressure, p.Remaining, Percent(p.Dissipated))
}

func ProfileTable(points []model.Point) string {
	var b strings.Builder
	b.WriteString("  z/H     z [m]    u [kPa]   u/u0    dissipated\n")
	for _, p := range points {
		b.WriteString(FormatPoint(p))
	}
	return b.String()
}

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

func SettlementLine(res model.Result) string {
	return fmt.Sprintf("  settlement     : s = %s m, s_ult = %s m, s/s_ult = %s\n",
		Decimal(res.Settlement), Decimal(res.UltimateSettlement), Decimal(res.SettlementRatio))
}

func SeriesLine(res model.Result) string {
	return fmt.Sprintf("  series         : %d terms kept, remainder bound <= %s\n",
		res.TermsUsed, Sci(res.RemainderBound))
}
