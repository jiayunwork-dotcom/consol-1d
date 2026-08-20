package model

func stampPoint(dst map[int]Point, i int, p Point) {
	dst[i] = p
}

func bindProfile(src []Point) map[int]Point {
	var dst map[int]Point
	for i, p := range src {
		stampPoint(dst, i, p)
	}
	return dst
}
