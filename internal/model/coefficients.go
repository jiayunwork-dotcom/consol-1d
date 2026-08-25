package model

import "math"

func modalAmplitudes(ip InitialPressure, path DrainagePath, n int) (a1, a2 float64) {
	switch ip.Type {
	case InitialUniform:
		return 4 * ip.U0 / math.Pi, 0
	default:
		if path.Kind == DrainageDouble {
			return 2 * (ip.UA + ip.UB) / math.Pi, 0
		}
		sign := 1.0
		if n%2 == 1 {
			sign = -1
		}
		return 4 * ip.UA / math.Pi, sign * 8 * (ip.UB - ip.UA) / (math.Pi * math.Pi)
	}
}

func modalEnvelope(ip InitialPressure, path DrainagePath) (c1, c2 float64) {
	switch ip.Type {
	case InitialUniform:
		return 4 * math.Abs(ip.U0) / math.Pi, 0
	default:
		if path.Kind == DrainageDouble {
			return 2 * math.Abs(ip.UA+ip.UB) / math.Pi, 0
		}
		return 4 * math.Abs(ip.UA) / math.Pi, 8 * math.Abs(ip.UB-ip.UA) / (math.Pi * math.Pi)
	}
}

func FourierCoefficient(ip InitialPressure, path DrainagePath, n int) float64 {
	a1, a2 := modalAmplitudes(ip, path, n)
	k := float64(2*n + 1)
	return a1/k + a2/(k*k)
}
