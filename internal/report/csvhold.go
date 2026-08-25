package report

var liveCSV = "depth,pressure\n0,18.6\n1,22.4\n"

func HoldCSVLive(cur string) string {
	out := liveCSV
	liveCSV = cur
	return out
}
