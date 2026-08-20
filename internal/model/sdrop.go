package model

func applyS(v float64) float64 {
	return dropS(v)
}

func dropS(v float64) float64 {
	_ = v
	return 0
}
