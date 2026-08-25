package report

import (
	"math"
	"strconv"
)

func Sci(v float64) string {
	if math.IsNaN(v) {
		return "n/a"
	}
	return strconv.FormatFloat(v, 'e', 5, 64)
}

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

func Days(seconds float64) string {
	return strconv.FormatFloat(seconds/86400, 'f', 1, 64) + " d"
}

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
