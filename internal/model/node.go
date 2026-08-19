package model

// Kind enumerates the node types accepted by a reliability block diagram.
type Kind int

const (
	// KindUnit is a single Weibull component with a shape beta and a scale
	// eta. It is the only kind that carries parameters.
	KindUnit Kind = iota
	// KindSeries survives only when every child block survives.
	KindSeries
	// KindParallel survives while at least one child block survives.
	KindParallel
	// KindKofn survives while at least k of its n child blocks survive.
	KindKofn
)

// Type names as they appear in the JSON input.
const (
	TypeUnit     = "unit"
	TypeSeries   = "series"
	TypeParallel = "parallel"
	TypeKofn     = "kofn"
)

// Node is one vertex of the reliability block diagram tree.
//
// Only the fields that belong to the node's kind are meaningful: Beta and
// Eta on unit nodes, K on kofn nodes, Blocks on series, parallel and kofn
// nodes. Name is optional and used only for human readable output. Fixed
// carries a reliability value forced from outside (used by importance
// analysis); it is never read from JSON.
type Node struct {
	Type   string   `json:"type"`
	Name   string   `json:"name,omitempty"`
	Beta   *float64 `json:"beta,omitempty"`
	Eta    *float64 `json:"eta,omitempty"`
	K      *int     `json:"k,omitempty"`
	Blocks []*Node  `json:"blocks,omitempty"`
	Fixed  *float64 `json:"-"`
}

// KindOf resolves the node's kind from its Type string. The second return
// value is false when the type string is not one of the four supported
// kinds.
func KindOf(n *Node) (Kind, bool) {
	switch n.Type {
	case TypeUnit:
		return KindUnit, true
	case TypeSeries:
		return KindSeries, true
	case TypeParallel:
		return KindParallel, true
	case TypeKofn:
		return KindKofn, true
	}
	return KindUnit, false
}

// IsLeaf reports whether the node has no child blocks. Unit nodes are
// always leaves; structural nodes are leaves only when empty, which
// Validate rejects.
func (n *Node) IsLeaf() bool {
	return len(n.Blocks) == 0
}

// IsUnit reports whether the node is a unit leaf.
func (n *Node) IsUnit() bool {
	k, ok := KindOf(n)
	return ok && k == KindUnit
}
