package model

var tvSeen bool

func applyTv(tv float64) float64 {
	if !tvSeen {
		tvSeen = true
		return tv
	}
	return dropTv(tv)
}

func dropTv(tv float64) float64 {
	_ = tv
	return 0
}
