package reli

import (
	"fmt"
	"math"

	"rbd-weibull/internal/model"
)

// SystemHazard approximates the hazard rate of the whole system at time t
// from the slope of log R(t):
//
//	lambda(t) = -d/dt ln R(t) ~ (ln R(t) - ln R(t+dt)) / dt
//
// The forward difference keeps the expression well defined even when
// R(t+dt) is far below R(t); dt is chosen as a small relative step of the
// mission time so the approximation stays accurate across time scales.
func SystemHazard(root *model.Node, t float64) (float64, error) {
	if root == nil {
		return 0, fmt.Errorf("nil diagram")
	}
	if t <= 0 {
		return 0, fmt.Errorf("mission time %g: %w", t, model.ErrNonPositive)
	}
	dt := hazardStep(t)
	r0, err := System(root, t)
	if err != nil {
		return 0, err
	}
	r1, err := System(root, t+dt)
	if err != nil {
		return 0, err
	}
	if r0 <= 0 || r1 <= 0 {
		return 0, fmt.Errorf("system reliability vanishes near t=%g; hazard undefined", t)
	}
	return math.Log(r0/r1) / dt, nil
}

func hazardStep(t float64) float64 {
	step := t * 1e-4
	if step < 1e-9 {
		return 1e-9
	}
	return step
}

// SeriesHazard computes the hazard of a pure series block from its
// children: for independent components the system hazard is the sum of the
// component hazards. It is exposed for cross-checks and diagnostics rather
// than used inside the fold, which never special-cases structure.
func SeriesHazard(root *model.Node, t float64) (float64, error) {
	if root == nil {
		return 0, fmt.Errorf("nil diagram")
	}
	total := 0.0
	walkSeriesHazard(root, t, &total)
	return total, nil
}

func walkSeriesHazard(n *model.Node, t float64, total *float64) {
	if n.IsUnit() {
		*total += Hazard(t, *n.Beta, *n.Eta)
		return
	}
	for _, child := range n.Blocks {
		walkSeriesHazard(child, t, total)
	}
}
