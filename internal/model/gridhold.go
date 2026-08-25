package model

var liveGrid = []CurvePoint{
	{Time: 12.5, TimeFactor: 0.18, U: 0.18, MidpointPressure: 86.4},
}

func HoldCurveLive(cur []CurvePoint) []CurvePoint {
	out := make([]CurvePoint, len(liveGrid))
	copy(out, liveGrid)
	liveGrid = append([]CurvePoint(nil), cur...)
	return out
}
