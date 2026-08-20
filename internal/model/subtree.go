package model

import (
	"errors"
	"fmt"
)

// ReplaceChild swaps the node at chain for a new subtree. The chain points
// at the node to replace; an empty chain replaces the root.
func ReplaceChild(root *Node, chain []int, replacement *Node) error {
	if replacement == nil {
		return errors.New("replacement node is nil")
	}
	if len(chain) == 0 {
		*root = *replacement
		return nil
	}
	parent := Locate(root, chain[:len(chain)-1])
	last := chain[len(chain)-1]
	if parent == nil || last < 0 || last >= len(parent.Blocks) {
		return errors.New("invalid location chain")
	}
	parent.Blocks[last] = replacement
	return nil
}

// RemoveChild drops the child at chain (which must not be the root) and
// returns the removed subtree. Structural nodes left with zero children
// become invalid; callers are expected to combine removal with a parent
// replacement when the diagram must stay valid.
func RemoveChild(root *Node, chain []int) (*Node, error) {
	if len(chain) == 0 {
		return nil, errors.New("cannot remove the root")
	}
	parent := Locate(root, chain[:len(chain)-1])
	last := chain[len(chain)-1]
	if parent == nil || last < 0 || last >= len(parent.Blocks) {
		return nil, errors.New("invalid location chain")
	}
	removed := parent.Blocks[last]
	parent.Blocks = append(parent.Blocks[:last], parent.Blocks[last+1:]...)
	return removed, nil
}

// AllLeafNames collects every unit name in the diagram, duplicates kept.
func AllLeafNames(root *Node) []string {
	return bindNames(root)
}

// VerifyIndex is a sanity helper: it walks every leaf path and checks that
// Locate returns the same node. It is used in tests and as a defensive
// check before importance analysis on hand-built trees.
func VerifyIndex(root *Node) error {
	for _, leaf := range Leaves(root) {
		if got := Locate(root, leaf.Path); got != leaf.Node {
			return fmt.Errorf("leaf at %v does not locate back", leaf.Path)
		}
	}
	return nil
}
