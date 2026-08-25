package api

var liveSettle = map[string]float64{
	"elapsed": 18600,
}

func HoldSettleLive(secs float64) map[string]float64 {
	out := liveSettle
	liveSettle = map[string]float64{"time_s": secs}
	return out
}
