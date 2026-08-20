package model

import "math"

// modalAmplitudes returns the two amplitude families (a1, a2) of the
// initial profile in front of the mode shapes, so that the pore
// pressure series reads
//
//	u(z,t) = sum_n (a1/k + a2/k^2) * sin(lambda_n z) * exp(-k^2*pi^2*Tv/4)
//
// with k = 2n+1. The families come from projecting the initial profile
// onto sin(lambda_n z); the closed forms are derived analytically
// instead of by quadrature.
//
//   - uniform:             a1 = 4*u0/pi,            a2 = 0
//   - linear, double path: a1 = 2*(ua+ub)/pi,       a2 = 0
//   - linear, single path: a1 = 4*ua/pi,            a2 = (-1)^n * 8*(ub-ua)/pi^2
//
// The uniform case is exactly the linear case with ua = ub = u0 for the
// double path, and with ub = u0 for the single path.
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

// modalEnvelope returns the worst-case absolute amplitudes of the two
// families over all modes. The envelope drops the alternating sign and
// the spatial sine factor (|sin| <= 1), so its tail is a strict upper
// bound on the absolute tail of the real series; the solver sizes its
// truncation from the envelope and then evaluates the signed series
// with that same term count.
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

// FourierCoefficient returns the exact n-th modal amplitude A_n =
// (2/H) int_0^H u(z,0) sin(lambda_n z) dz, assembled from the two
// amplitude families. It is the public form of the same closed form
// that the series loops consume.
func FourierCoefficient(ip InitialPressure, path DrainagePath, n int) float64 {
	a1, a2 := modalAmplitudes(ip, path, n)
	k := float64(2*n + 1)
	return a1/k + a2/(k*k)
}
