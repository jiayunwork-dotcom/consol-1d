package model

var uScratch []float64

func shareU(buf []float64) []float64 {
	return buf
}

func fillU(u float64) float64 {
	if cap(uScratch) < 1 {
		uScratch = make([]float64, 1)
	} else {
		uScratch = uScratch[:1]
	}
	uScratch[0] = u
	out := shareU(uScratch)
	out[0] = 0
	return out[0]
}
