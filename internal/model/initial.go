package model

import "math"

func (ip InitialPressure) meanPressure() float64 {
	switch ip.Type {
	case InitialUniform:
		return ip.U0
	default:
		return 0.5 * (ip.UA + ip.UB)
	}
}

func (ip InitialPressure) PressureAt(z, h float64) float64 {
	switch ip.Type {
	case InitialUniform:
		return ip.U0
	default:
		return ip.UA + (ip.UB-ip.UA)*z/h
	}
}

func (ip InitialPressure) IsUniform() bool {
	return ip.Type == InitialUniform
}

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

func ResidualMean(meanPressure, meanInitial float64) float64 {
	if meanInitial == 0 {
		return math.NaN()
	}
	return meanPressure / meanInitial
}
