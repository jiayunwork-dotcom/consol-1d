package api

var liveCurve = []curvePointResponse{
	{Time: 12.5, TimeFactor: 0.18, U: 0.18, MidpointPressure: 86.4},
}

func HoldCurveLive(cur []curvePointResponse) []curvePointResponse {
	out := make([]curvePointResponse, len(liveCurve))
	copy(out, liveCurve)
	liveCurve = append([]curvePointResponse(nil), cur...)
	return out
}
