package report

import (
	"fmt"
	"strings"

	"consol-1d/internal/model"
)

// Report assembles the complete human-readable report for one solve:
// the scenario header, the primary results, the truncation proof and
// the dissipation profile table.
func Report(res model.Result) string {
	var b strings.Builder
	b.WriteString(Header(res))
	b.WriteString(SettlementLine(res))
	b.WriteString(SeriesLine(res))
	b.WriteString(fmt.Sprintf("  small-time U   : ~ %s   (2*sqrt(Tv/pi) reference)\n",
		Decimal(model.SmallTimeAsymptote(res.Tv))))
	b.WriteString(fmt.Sprintf("\nprofile (%d nodes):\n", len(res.Profile)))
	b.WriteString(ProfileTable(res.Profile))
	return strings.TrimRight(b.String(), "\n") + "\n"
}

// Summary is the one-line digest used as the last line of a report so
// scripts can grep a stable identifier for the key numbers.
func Summary(res model.Result) string {
	return fmt.Sprintf("U=%s midpoint_u=%s kPa settlement_ratio=%s terms=%d bound<=%s",
		Decimal(res.U), Decimal(res.MidpointPressure), Decimal(res.SettlementRatio),
		res.TermsUsed, Sci(res.RemainderBound))
}

// DrainageLabel returns the short human label of a drainage path.
func DrainageLabel(d model.Drainage) string {
	if d == model.DrainageDouble {
		return "double"
	}
	return "single"
}

// InitialLabel returns a one-line description of the initial pressure
// profile, e.g. "uniform u0=100" or "linear ua=120 ub=60".
func InitialLabel(ip model.InitialPressure) string {
	switch ip.Type {
	case model.InitialUniform:
		return fmt.Sprintf("uniform u0=%g", ip.U0)
	default:
		return fmt.Sprintf("linear ua=%g ub=%g", ip.UA, ip.UB)
	}
}
