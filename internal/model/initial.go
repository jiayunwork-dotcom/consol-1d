package model

import "math"

// meanPressure returns the layer-averaged initial excess pore pressure
// u_bar(0) = (1/H) int_0^H u(z,0) dz, which normalizes the
// consolidation degree U = 1 - u_bar(t)/u_bar(0).
func (ip InitialPressure) meanPressure() float64 {
	switch ip.Type {
	case InitialUniform:
		return ip.U0
	default:
		return 0.5 * (ip.UA + ip.UB)
	}
}

// PressureAt evaluates the initial profile u(z, 0) at an arbitrary
// depth z inside [0, H].
func (ip InitialPressure) PressureAt(z, h float64) float64 {
	switch ip.Type {
	case InitialUniform:
		return ip.U0
	default:
		return ip.UA + (ip.UB-ip.UA)*z/h
	}
}

// IsUniform reports whether the profile is the constant one.
func (ip InitialPressure) IsUniform() bool {
	return ip.Type == InitialUniform
}

// PeakPressure returns the largest initial pressure over the layer,
// which the report uses as the natural scale for the dissipation table.
func (ip InitialPressure) PeakPressure(h float64) float64 {
	switch ip.Type {
	case InitialUniform:
		return ip.U0
	default:
		if ip.UA >= ip.UB {
			return ip.UA
		}
		return ip.UB
	}
}

// ResidualMean returns the fraction of the initial mean pressure that a
// given layer mean still represents; values near 1 mean almost no
// dissipation has happened.
func ResidualMean(meanPressure, meanInitial float64) float64 {
	if meanInitial == 0 {
		return math.NaN()
	}
	return meanPressure / meanInitial
}
