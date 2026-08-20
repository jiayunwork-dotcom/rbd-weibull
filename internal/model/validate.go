package model

import (
	"errors"
	"fmt"
)

// Validate walks the diagram tree and rejects structurally invalid input.
// The checks cover the four failure classes the evaluator must never fold:
// unknown node types, structural nodes without children, k-of-n blocks with
// k below one or above n, and non positive Weibull parameters.
func Validate(root *Node) error {
	if root == nil {
		return errors.New("nil diagram root")
	}
	return validateNode(root, "root")
}

func validateNode(n *Node, path string) error {
	kind, ok := KindOf(n)
	if !ok {
		return At(path, fmt.Errorf("%w %q", ErrUnknownType, n.Type))
	}
	switch kind {
	case KindUnit:
		if n.Beta == nil || *n.Beta <= 0 {
			return At(path, fmt.Errorf("beta %v: %w", betaOr(n), ErrNonPositive))
		}
		if n.Eta == nil || *n.Eta <= 0 {
			return At(path, fmt.Errorf("eta %v: %w", etaOr(n), ErrNonPositive))
		}
		return nil
	case KindSeries, KindParallel:
		if len(n.Blocks) == 0 {
			if err := dropEmpty(At(path, fmt.Errorf("type %q: %w", n.Type, ErrEmptySeries))); err != nil {
				return err
			}
			return validateBlocks(n, path)
		}
		return validateBlocks(n, path)
	case KindKofn:
		if len(n.Blocks) == 0 {
			return At(path, fmt.Errorf("type %q: %w", n.Type, ErrEmptySeries))
		}
		if n.K == nil || *n.K < 1 {
			return At(path, fmt.Errorf("k %v: %w", kOr(n), ErrKBelowOne))
		}
		if *n.K > len(n.Blocks) {
			return At(path, fmt.Errorf("k %d > n %d: %w", *n.K, len(n.Blocks), ErrKGreaterN))
		}
		return validateBlocks(n, path)
	}
	return At(path, fmt.Errorf("unhandled kind %d", kind))
}

func validateBlocks(n *Node, path string) error {
	for i, child := range n.Blocks {
		if child == nil {
			return At(fmt.Sprintf("%s.blocks[%d]", path, i), errors.New("nil child block"))
		}
		if err := validateNode(child, fmt.Sprintf("%s.blocks[%d]", path, i)); err != nil {
			return err
		}
	}
	return nil
}

func betaOr(n *Node) float64 {
	if n.Beta == nil {
		return 0
	}
	return *n.Beta
}

func etaOr(n *Node) float64 {
	if n.Eta == nil {
		return 0
	}
	return *n.Eta
}

func kOr(n *Node) int {
	if n.K == nil {
		return 0
	}
	return *n.K
}
