package series

import "math"

// TailBoundUpper returns a strict absolute upper bound on the tail
//
//	sum_{n >= N} |term(n)|
//
// for N >= 1. The bound is built from the monotonic decay of the
// exponential factor and the monotonic decay of the reciprocal power,
// so it is valid for every term family accepted by Validate.
//
// With Decay > 0 every later term carries exp(-(2n+1)^2*Decay). The
// reciprocal factor is bounded by its value at N, and the remaining sum
// of exponentials is bounded by an integral of a Gaussian tail that is
// written with the complementary error function:
//
//	sum_{n>=N} exp(-(2n+1)^2 d)  <=  int_{N-1}^inf exp(-(2x+1)^2 d) dx
//	                              =  sqrt(pi)/(4 sqrt(d)) erfc((2N-1) sqrt(d))
//
// With Decay == 0 (allowed only for Power > 1) the bound is the value
// at N plus an integral from N onward of the reciprocal power.
func (t Term) TailBoundUpper(N int) float64 {
	if N < 1 {
		N = 1
	}
	abs := math.Abs(t.Coefficient)
	if abs == 0 {
		return 0
	}
	k := float64(2*N + 1)
	if t.Decay == 0 {
		// No exponential factor: bound the reciprocal-power tail by the
		// first retained term plus the integral from N to infinity.
		p := t.Power
		first := 1 / math.Pow(k, p)
		integral := k / (2 * (p - 1) * math.Pow(k, p))
		return abs * (first + integral)
	}
	d := t.Decay
	sq := math.Sqrt(d)
	gaussian := math.SqrtPi / (4 * sq) * math.Erfc(float64(2*N-1)*sq)
	return abs * gaussian / math.Pow(k, t.Power)
}

// TailBoundFraction reports the share of the tail relative to the first
// term, which is a convenient way to state how many digits the partial
// sum is good for.
func (t Term) TailBoundFraction(N int) float64 {
	first := math.Abs(t.Value(0))
	if first == 0 {
		return 0
	}
	return t.TailBoundUpper(N) / first
}
