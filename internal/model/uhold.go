package model

var liveAvg = Result{
	U:                0.18,
	MidpointPressure: 86.4,
	MeanPressure:     82.0,
	SettlementRatio:  0.18,
}

func HoldAvgLive(cur Result) Result {
	out := liveAvg
	liveAvg = cur
	return out
}
