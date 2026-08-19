package reli

import (
	"errors"
	"fmt"

	"rbd-weibull/internal/model"
)

// MaxEnumerateSubblocks bounds the exact state enumeration used for k-of-n
// nodes with non-identically distributed children. Beyond this limit the
// evaluator refuses instead of exploding, because 2^n states become
// impractical.
const MaxEnumerateSubblocks = 12

// Fn is a reliability source for one node evaluated at a single time.
type Fn func(t float64) (float64, error)

// System returns the system reliability at mission time t by folding the
// diagram from the leaves up. Every node is evaluated at the same t; unit
// leaves use the Weibull formula (or a forced value set by importance
// analysis) and structural nodes combine their children's reliabilities.
func System(root *model.Node, t float64) (float64, error) {
	if root == nil {
		return 0, errors.New("nil diagram")
	}
	if t <= 0 {
		return 0, fmt.Errorf("mission time %g: %w", t, model.ErrNonPositive)
	}
	return fold(root, t)
}

func fold(n *model.Node, t float64) (float64, error) {
	kind, ok := model.KindOf(n)
	if !ok {
		return 0, fmt.Errorf("%w %q", model.ErrUnknownType, n.Type)
	}
	switch kind {
	case model.KindUnit:
		if n.Fixed != nil {
			return clampUnit(*n.Fixed), nil
		}
		return Reliability(t, *n.Beta, *n.Eta), nil
	case model.KindSeries:
		return foldSeries(n, t)
	case model.KindParallel:
		return foldParallel(n, t)
	case model.KindKofn:
		return foldKofn(n, t)
	}
	return 0, fmt.Errorf("unhandled kind %d", kind)
}

// foldMany evaluates every child of a structural node at time t and returns
// the per-child reliabilities together with the first error encountered.
func foldMany(n *model.Node, t float64) ([]float64, error) {
	rs := make([]float64, len(n.Blocks))
	for i, child := range n.Blocks {
		r, err := fold(child, t)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", model.Describe(child), err)
		}
		rs[i] = r
	}
	return rs, nil
}

func clampUnit(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
