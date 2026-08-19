package model

// Leaf carries a unit node together with the block-index chain that locates
// it inside the diagram tree. The chain lets callers find the same leaf in
// a cloned copy of the tree, which is what importance analysis needs.
type Leaf struct {
	Node *Node
	Path []int
}

// Leaves returns every unit leaf of the diagram in depth-first order,
// paired with its location chain. Structural nodes contribute no leaves.
func Leaves(root *Node) []Leaf {
	var out []Leaf
	walkLeaves(root, nil, &out)
	return out
}

func walkLeaves(n *Node, chain []int, out *[]Leaf) {
	if n.IsUnit() {
		*out = append(*out, Leaf{Node: n, Path: append([]int(nil), chain...)})
		return
	}
	for i, child := range n.Blocks {
		walkLeaves(child, append(chain, i), out)
	}
}

// Count returns the total number of nodes in the tree, root included.
func Count(root *Node) int {
	if root == nil {
		return 0
	}
	total := 1
	for _, child := range root.Blocks {
		total += Count(child)
	}
	return total
}

// Depth returns the length of the longest root-to-leaf chain. A single
// node tree has depth one.
func Depth(root *Node) int {
	if root == nil {
		return 0
	}
	best := 0
	for _, child := range root.Blocks {
		if d := Depth(child); d > best {
			best = d
		}
	}
	return best + 1
}

// KindCount tallies how many nodes of each kind appear in the tree.
func KindCount(root *Node) map[Kind]int {
	counts := map[Kind]int{}
	walkKindCount(root, counts)
	return counts
}

func walkKindCount(n *Node, counts map[Kind]int) {
	kind, ok := KindOf(n)
	if ok {
		counts[kind]++
	}
	for _, child := range n.Blocks {
		walkKindCount(child, counts)
	}
}
