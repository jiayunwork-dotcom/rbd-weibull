package reli

import (
	"fmt"
	"math"

	"rbd-weibull/internal/model"
)

// Mean-time-to-failure parameters. The system MTTF is obtained by
// numerical integration of R(t) over [0, tMax]; the README documents both
// values so results are reproducible.
const (
	// DefaultGridSteps is the number of trapezoid steps used by
	// SystemMTTF when the caller does not override it.
	DefaultGridSteps = 2000

	// UpperMultiplier scales the largest scale parameter in the diagram to
	// form the integration horizon: tMax = UpperMultiplier * max(eta).
	// At that horizon the reliability has decayed by at least exp(-8^beta),
	// which keeps the truncation error below the reported precision.
	UpperMultiplier = 8.0
)

// UnitMTTF is the closed-form mean time to failure of a single Weibull
// component: eta * Gamma(1 + 1/beta). It is the analytic reference that
// the numerical system integrator must reproduce on single-leaf diagrams.
func UnitMTTF(beta, eta float64) float64 {
	return eta * math.Gamma(1+1/beta)
}

// SystemMTTF integrates R(t) over [0, tMax] with n trapezoid steps. The
// default horizon is UpperMultiplier times the largest eta in the diagram,
// so the tail contribution beyond tMax is negligible for every component.
// A zero tMax selects the default horizon; a non-positive n selects
// DefaultGridSteps.
func SystemMTTF(root *model.Node, tMax float64, n int) (float64, error) {
	if root == nil {
		return 0, fmt.Errorf("nil diagram")
	}
	if tMax <= 0 {
		tMax = integrationHorizon(root)
	}
	if n <= 0 {
		n = DefaultGridSteps
	}
	h := tMax / float64(n)
	area := 0.5 // R(0) = 1 for every valid diagram; the trapezoid starts there.
	for i := 1; i <= n; i++ {
		t := float64(i) * h
		r, err := System(root, t)
		if err != nil {
			return 0, err
		}
		if i == n {
			area += 0.5 * r
		} else {
			area += r
		}
	}
	return area * h, nil
}

// integrationHorizon picks the largest eta in the diagram and scales it by
// UpperMultiplier.
func integrationHorizon(root *model.Node) float64 {
	largest := 1.0
	walkEta(root, &largest)
	return largest * UpperMultiplier
}

func walkEta(n *model.Node, largest *float64) {
	if n.IsUnit() && n.Eta != nil && *n.Eta > *largest {
		*largest = *n.Eta
	}
	for _, child := range n.Blocks {
		walkEta(child, largest)
	}
}

// MTTFHorizon returns the integration horizon SystemMTTF would use for a
// diagram, exposed for diagnostics and CLI output.
func MTTFHorizon(root *model.Node) float64 {
	return integrationHorizon(root)
}
