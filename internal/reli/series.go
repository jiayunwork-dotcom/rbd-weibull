package reli

import "rbd-weibull/internal/model"

// foldSeries combines children in a series: the block survives only if
// every child survives, so the reliabilities multiply. The empty product
// would be 1, which is why validation rejects structural nodes without
// children before this function can ever be reached with zero blocks.
func foldSeries(n *model.Node, t float64) (float64, error) {
	rs, err := foldMany(n, t)
	if err != nil {
		return 0, err
	}
	r := 1.0
	for _, ri := range rs {
		r *= ri
	}
	return r, nil
}

// foldParallel combines children in parallel: the block survives while at
// least one child survives, so the failure probability is the product of
// the children's failure probabilities and R = 1 - prod(1 - R_i).
func foldParallel(n *model.Node, t float64) (float64, error) {
	rs, err := foldMany(n, t)
	if err != nil {
		return 0, err
	}
	q := 1.0
	for _, ri := range rs {
		q *= 1 - ri
	}
	return 1 - q, nil
}

// AllEqual reports whether every value in rs agrees within eps. It is used
// to detect identically distributed k-of-n children.
func AllEqual(rs []float64, eps float64) bool {
	if len(rs) < 2 {
		return true
	}
	ref := rs[0]
	for _, r := range rs[1:] {
		if diff(r, ref) > eps {
			return false
		}
	}
	return true
}

func diff(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}
