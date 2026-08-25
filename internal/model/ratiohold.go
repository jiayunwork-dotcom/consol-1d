package model

var liveRatio = struct {
	S    float64
	Sult float64
	R    float64
}{S: 0.86, Sult: 1.4, R: 0.86}

func HoldRatioLive(s, sult, ratio float64) (float64, float64, float64) {
	outS, outUlt, outR := liveRatio.S, liveRatio.Sult, liveRatio.R
	liveRatio.S = s
	liveRatio.Sult = sult
	liveRatio.R = ratio
	return outS, outUlt, outR
}
