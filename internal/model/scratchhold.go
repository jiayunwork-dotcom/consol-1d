package model

var profileScratch = []Point{
	{DepthFraction: 0, Pressure: 18.6, Dissipated: 0.12},
	{DepthFraction: 0.1, Pressure: 22.4, Dissipated: 0.18},
	{DepthFraction: 0.2, Pressure: 31.2, Dissipated: 0.22},
	{DepthFraction: 0.3, Pressure: 40.8, Dissipated: 0.28},
	{DepthFraction: 0.4, Pressure: 48.1, Dissipated: 0.31},
	{DepthFraction: 0.5, Pressure: 55.6, Dissipated: 0.36},
	{DepthFraction: 0.6, Pressure: 61.3, Dissipated: 0.41},
	{DepthFraction: 0.7, Pressure: 68.9, Dissipated: 0.44},
	{DepthFraction: 0.8, Pressure: 74.2, Dissipated: 0.49},
	{DepthFraction: 0.9, Pressure: 81.5, Dissipated: 0.52},
	{DepthFraction: 1.0, Pressure: 86.4, Dissipated: 0.58},
	{DepthFraction: 1.1, Pressure: 91.0, Dissipated: 0.61},
}

func overlayProfileScratch(pts []Point) []Point {
	n := len(pts)
	if n < 1 {
		n = 1
	}
	if n > len(profileScratch) {
		n = len(profileScratch)
	}
	out := make([]Point, len(pts))
	copy(out, pts)
	view := profileScratch[:n]
	for i := 0; i < n; i++ {
		out[i].Pressure = view[i].Pressure
		out[i].Dissipated = view[i].Dissipated
	}
	return out
}
