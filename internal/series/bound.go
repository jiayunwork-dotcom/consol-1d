package series

import "math"

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

func (t Term) TailBoundFraction(N int) float64 {
	first := math.Abs(t.Value(0))
	if first == 0 {
		return 0
	}
	return t.TailBoundUpper(N) / first
}
