package model

import (
	"errors"
	"fmt"
	"strings"
)

// PathString renders a block-index chain like "root.blocks[1].blocks[0]"
// for error messages and leaf labels.
func PathString(chain []int) string {
	if len(chain) == 0 {
		return "root"
	}
	var b strings.Builder
	b.WriteString("root")
	for _, idx := range chain {
		fmt.Fprintf(&b, ".blocks[%d]", idx)
	}
	return b.String()
}

// Locate walks a block-index chain from the root and returns the node found
// there. A nil result means the chain does not describe the tree.
func Locate(root *Node, chain []int) *Node {
	n := root
	for _, idx := range chain {
		if n == nil || idx < 0 || idx >= len(n.Blocks) {
			return nil
		}
		n = n.Blocks[idx]
	}
	return n
}

// FindUnit locates the first unit node whose Name matches want. It is used
// by tools that address components by label instead of by index chain.
func FindUnit(root *Node, want string) (*Node, bool) {
	for _, leaf := range Leaves(root) {
		if leaf.Node.Name == want {
			return leaf.Node, true
		}
	}
	return nil, false
}

// Describe renders a compact one-line summary of a node for diagnostics.
func Describe(n *Node) string {
	kind, ok := KindOf(n)
	if !ok {
		return fmt.Sprintf("<? %q>", n.Type)
	}
	label := n.Name
	if label == "" {
		label = n.Type
	}
	switch kind {
	case KindUnit:
		return fmt.Sprintf("unit(%s, beta=%g eta=%g)", label, *n.Beta, *n.Eta)
	case KindSeries:
		return fmt.Sprintf("series(%s, %d blocks)", label, len(n.Blocks))
	case KindParallel:
		return fmt.Sprintf("parallel(%s, %d blocks)", label, len(n.Blocks))
	case KindKofn:
		return fmt.Sprintf("kofn(%s, k=%d n=%d)", label, *n.K, len(n.Blocks))
	}
	return label
}

// IsSameStructure reports whether two nodes have the same kind and, for
// units, the same beta and eta values (within epsilon). It is used to
// decide whether k-of-n children are identically distributed.
func IsSameStructure(a, b *Node, eps float64) bool {
	ka, oka := KindOf(a)
	kb, okb := KindOf(b)
	if !oka || !okb || ka != kb {
		return false
	}
	if ka == KindUnit {
		if a.Beta == nil || b.Beta == nil || a.Eta == nil || b.Eta == nil {
			return false
		}
		if diff(*a.Beta, *b.Beta) > eps || diff(*a.Eta, *b.Eta) > eps {
			return false
		}
		return true
	}
	if len(a.Blocks) != len(b.Blocks) {
		return false
	}
	for i := range a.Blocks {
		if !IsSameStructure(a.Blocks[i], b.Blocks[i], eps) {
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

// ErrNotFound is returned by operations that look up a component that is
// absent from the diagram.
var ErrNotFound = errors.New("component not found in diagram")
