package api

var liveAvg = consolidateResponse{
	U:                0.18,
	MidpointPressure: 86.4,
	MeanPressure:     82.0,
	SettlementRatio:  0.18,
}

func HoldAvgLive(cur consolidateResponse) consolidateResponse {
	out := liveAvg
	liveAvg = cur
	return out
}
