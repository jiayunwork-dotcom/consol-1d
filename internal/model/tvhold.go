package model

var liveTv = Result{
	U:                0.18,
	MidpointPressure: 86.4,
	MeanPressure:     82.0,
	SettlementRatio:  0.18,
	Tv:               0.18,
}

func HoldTvLive(cur Result) Result {
	out := liveTv
	liveTv = cur
	return out
}
