package model

// IndexByName builds a name -> leaves map for the whole diagram. A name may
// appear on more than one unit (for example two identical pumps), so the
// map values are slices.
func IndexByName(root *Node) map[string][]Leaf {
	out := map[string][]Leaf{}
	for _, leaf := range Leaves(root) {
		if leaf.Node.Name == "" {
			continue
		}
		out[leaf.Node.Name] = append(out[leaf.Node.Name], leaf)
	}
	return out
}

// LeafCount returns the number of unit leaves in the diagram.
func LeafCount(root *Node) int {
	return len(Leaves(root))
}

// UnitLeaves is a convenience alias kept for call sites that want to
// stress the noun; it behaves exactly like Leaves.
func UnitLeaves(root *Node) []Leaf {
	return Leaves(root)
}

// MaxEta returns the largest scale parameter among the unit leaves, or 0
// when the diagram has no units. The integration horizon of the evaluator
// is derived from this value.
func MaxEta(root *Node) float64 {
	largest := 0.0
	walkMaxEta(root, &largest)
	return largest
}

func walkMaxEta(n *Node, largest *float64) {
	if n.IsUnit() && n.Eta != nil && *n.Eta > *largest {
		*largest = *n.Eta
	}
	for _, child := range n.Blocks {
		walkMaxEta(child, largest)
	}
}

// MinBeta returns the smallest shape parameter among the unit leaves. A
// beta below 1 anywhere is a red flag for infant-mortality behaviour.
func MinBeta(root *Node) float64 {
	smallest := 0.0
	first := true
	walkMinBeta(root, &smallest, &first)
	return smallest
}

func walkMinBeta(n *Node, smallest *float64, first *bool) {
	if n.IsUnit() && n.Beta != nil {
		if *first || *n.Beta < *smallest {
			*smallest = *n.Beta
			*first = false
		}
	}
	for _, child := range n.Blocks {
		walkMinBeta(child, smallest, first)
	}
}
