package reli

import (
	"fmt"

	"rbd-weibull/internal/model"
)

// SystemQuantile finds the time t at which the system reliability drops to
// the target value r by bisection over [0, tMax]. With the default horizon
// the system reliability at tMax is far below any practical r, so the
// bracket is valid: R(0) = 1 > r > R(tMax).
func SystemQuantile(root *model.Node, r float64, tMax float64, iters int) (float64, error) {
	if root == nil {
		return 0, fmt.Errorf("nil diagram")
	}
	if r <= 0 || r >= 1 {
		return 0, fmt.Errorf("target reliability %g must lie in (0, 1)", r)
	}
	if tMax <= 0 {
		tMax = MTTFHorizon(root)
	}
	if iters <= 0 {
		iters = 60
	}
	lo, hi := 0.0, tMax
	for i := 0; i < iters; i++ {
		mid := 0.5 * (lo + hi)
		rm, err := System(root, mid)
		if err != nil {
			return 0, err
		}
		if rm > r {
			lo = mid
		} else {
			hi = mid
		}
	}
	return 0.5 * (lo + hi), nil
}

// ReliabilityLife returns the time by which the system has lost 10% of its
// start-of-life reliability, a common design target.
func ReliabilityLife(root *model.Node) (float64, error) {
	return SystemQuantile(root, 0.9, 0, 0)
}

// B10 returns the time at which 10% of the population is expected to have
// failed, that is R(t) = 0.9. It is an alias of ReliabilityLife kept for
// readers familiar with the bearings terminology.
func B10(root *model.Node) (float64, error) {
	return SystemQuantile(root, 0.9, 0, 0)
}

// B50 returns the system median life, R(t) = 0.5.
func B50(root *model.Node) (float64, error) {
	return SystemQuantile(root, 0.5, 0, 0)
}
