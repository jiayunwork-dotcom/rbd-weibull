package reli

import (
	"fmt"

	"rbd-weibull/internal/model"
)

// Importance is the Birnbaum measure of one leaf component at a fixed
// mission time:
//
//	I_i(t) = R(t | component i works) - R(t | component i fails)
//
// which is the sensitivity of the system reliability to that component's
// state. Both conditional reliabilities come from a full system
// re-evaluation with the leaf's reliability forced to 1 and 0 on a cloned
// tree, so any structural or parameter bug anywhere in the diagram shows up
// in both values simultaneously.
type Importance struct {
	Leaf  model.Leaf
	Name  string
	Path  string
	Value float64
}

// Birnbaum evaluates the Birnbaum importance of every unit leaf in the
// diagram at time t. For each leaf the tree is cloned twice: once with the
// leaf forced working (reliability 1) and once with it forced failed
// (reliability 0). The difference of the two system reliabilities is the
// importance. The original tree is never mutated.
func Birnbaum(root *model.Node, t float64) ([]Importance, error) {
	leaves := model.Leaves(root)
	out := make([]Importance, 0, len(leaves))
	for _, leaf := range leaves {
		rw, err := conditionalSystem(root, leaf, 1.0, t)
		if err != nil {
			return nil, err
		}
		rf, err := conditionalSystem(root, leaf, 0.0, t)
		if err != nil {
			return nil, err
		}
		out = append(out, Importance{
			Leaf:  leaf,
			Name:  leafName(leaf),
			Path:  model.PathString(leaf.Path),
			Value: rw - rf,
		})
	}
	return out, nil
}

// conditionalSystem clones the diagram, forces one leaf to the given
// reliability and folds the system at time t.
func conditionalSystem(root *model.Node, leaf model.Leaf, forced, t float64) (float64, error) {
	clone := model.Clone(root)
	if err := model.SetFixed(clone, leaf.Path, forced); err != nil {
		return 0, fmt.Errorf("force leaf: %w", err)
	}
	return System(clone, t)
}

func leafName(leaf model.Leaf) string {
	if leaf.Node.Name != "" {
		return leaf.Node.Name
	}
	return model.PathString(leaf.Path)
}
