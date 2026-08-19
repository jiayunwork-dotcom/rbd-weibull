package model

import "errors"

// Clone returns a deep copy of the diagram tree. Structural nodes share
// nothing with the source; mutating the copy (for example forcing a leaf
// reliability during importance analysis) never touches the original.
func Clone(root *Node) *Node {
	if root == nil {
		return nil
	}
	out := &Node{
		Type:  root.Type,
		Name:  root.Name,
		K:     copyInt(root.K),
		Beta:  copyFloat(root.Beta),
		Eta:   copyFloat(root.Eta),
		Fixed: copyFloat(root.Fixed),
	}
	if root.Blocks != nil {
		out.Blocks = make([]*Node, len(root.Blocks))
		for i, child := range root.Blocks {
			out.Blocks[i] = Clone(child)
		}
	}
	return out
}

func copyInt(v *int) *int {
	if v == nil {
		return nil
	}
	c := *v
	return &c
}

func copyFloat(v *float64) *float64 {
	if v == nil {
		return nil
	}
	c := *v
	return &c
}

// SetFixed forces the reliability of the node located by chain to value v
// (typically 1 for "known working" or 0 for "known failed"). The override
// is honoured by the folding engine before the node's own parameters are
// consulted.
func SetFixed(root *Node, chain []int, v float64) error {
	n := Locate(root, chain)
	if n == nil {
		return errors.New("invalid location chain")
	}
	n.Fixed = copyFloat(&v)
	return nil
}

// ClearFixed removes any forced reliability from the node at chain,
// restoring its Weibull behaviour.
func ClearFixed(root *Node, chain []int) error {
	n := Locate(root, chain)
	if n == nil {
		return errors.New("invalid location chain")
	}
	n.Fixed = nil
	return nil
}
