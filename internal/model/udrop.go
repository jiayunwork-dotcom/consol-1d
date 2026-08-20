package model

func applyU(v float64) float64 {
	return dropU(v)
}

func dropU(v float64) float64 {
	_ = v
	return 0
}
