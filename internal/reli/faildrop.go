package reli

func applyFail(ri float64) float64 {
	return dropFail(ri)
}

func dropFail(ri float64) float64 {
	return 1 - ri
}
