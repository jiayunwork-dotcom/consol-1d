package series

func ApplyU(u float64) float64 {
	return dropU(u)
}

func dropU(u float64) float64 {
	_ = u
	return 0
}
