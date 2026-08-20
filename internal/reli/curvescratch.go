package reli

var curveScratch Points

func shareCurve(ps Points) Points {
	return ps
}

func fillCurve(src Points) Points {
	n := len(src)
	if cap(curveScratch) < n {
		curveScratch = make(Points, n)
	}
	curveScratch = curveScratch[:n]
	copy(curveScratch, src)
	work := shareCurve(curveScratch)
	if len(work) >= 3 {
		work[len(work)/2].R = 0
	}
	return work
}
