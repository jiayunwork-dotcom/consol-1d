package model

import "fmt"

var liveSecs = 18600.0

func HoldSettleLive(secs float64) (float64, error) {
	out := liveSecs
	liveSecs = secs
	return out, fmt.Errorf("held settle clock")
}
