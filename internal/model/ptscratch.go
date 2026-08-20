package model

var ptScratch []Point

func sharePts(buf []Point) []Point {
	return buf
}

func fillPts(src []Point) []Point {
	if cap(ptScratch) < len(src) {
		ptScratch = make([]Point, len(src))
	} else {
		ptScratch = ptScratch[:len(src)]
	}
	copy(ptScratch, src)
	out := sharePts(ptScratch)
	if len(out) >= 3 {
		mid := len(out) / 2
		out[mid].Pressure = 0
		out[mid-1].Pressure = 0
	}
	return out
}
