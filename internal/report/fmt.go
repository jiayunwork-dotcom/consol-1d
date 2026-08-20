// Package report renders a consolidation Result into human-readable
// text, a profile table, a CSV dump or a JSON document. It holds no
// physics: every number it prints comes from the model package, and the
// formatting helpers below exist so the numbers are printed in a stable,
// machine-parseable way.
package report

import (
	"math"
	"strconv"
)

// Sci formats a value with five significant digits in scientific
// notation, suited for quantities that span many orders of magnitude
// (cv, time factors, remainder bounds). NaN becomes "n/a".
func Sci(v float64) string {
	if math.IsNaN(v) {
		return "n/a"
	}
	return strconv.FormatFloat(v, 'e', 5, 64)
}

// Decimal formats a value with six fractional digits when its magnitude
// is printable, falling back to scientific notation for very large or
// very small values. NaN becomes "n/a".
func Decimal(v float64) string {
	if math.IsNaN(v) {
		return "n/a"
	}
	if v == 0 {
		return "0.000000"
	}
	a := math.Abs(v)
	if a >= 1e6 || a < 1e-4 {
		return strconv.FormatFloat(v, 'e', 5, 64)
	}
	return strconv.FormatFloat(v, 'f', 6, 64)
}

// Percent formats a fraction in [0, 1] as a percentage with one decimal
// place. Values slightly outside the range (numerical noise around the
// limits) are clamped so the report never shows 101%.
func Percent(v float64) string {
	if math.IsNaN(v) {
		return "n/a"
	}
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	return strconv.FormatFloat(v*100, 'f', 1, 64) + "%"
}

// Days renders a duration in seconds as decimal days, the time unit a
// geotechnical engineer reasons in.
func Days(seconds float64) string {
	return strconv.FormatFloat(seconds/86400, 'f', 1, 64) + " d"
}

// Duration renders a number of seconds in a human-friendly form with
// the best matching unit.
func Duration(seconds float64) string {
	if seconds >= 365.25*86400 {
		return strconv.FormatFloat(seconds/(365.25*86400), 'f', 2, 64) + " yr"
	}
	if seconds >= 86400 {
		return strconv.FormatFloat(seconds/86400, 'f', 1, 64) + " d"
	}
	if seconds >= 3600 {
		return strconv.FormatFloat(seconds/3600, 'f', 1, 64) + " h"
	}
	if seconds >= 60 {
		return strconv.FormatFloat(seconds/60, 'f', 1, 64) + " min"
	}
	return strconv.FormatFloat(seconds, 'g', 4, 64) + " s"
}
