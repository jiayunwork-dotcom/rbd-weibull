package model

import (
	"fmt"
	"strings"
)

// TextTree renders the diagram as an indented text tree. Structural nodes
// show their kind and child count; unit nodes show their name and Weibull
// parameters. It is used by diagnostics and by the CLI's describe output.
func TextTree(root *Node) string {
	var b strings.Builder
	writeTree(&b, root, 0)
	return b.String()
}

func writeTree(b *strings.Builder, n *Node, depth int) {
	indent := strings.Repeat("  ", depth)
	if n == nil {
		b.WriteString(indent + "<nil>\n")
		return
	}
	kind, ok := KindOf(n)
	if !ok {
		fmt.Fprintf(b, "%s<? %q>\n", indent, n.Type)
		return
	}
	switch kind {
	case KindUnit:
		name := n.Name
		if name == "" {
			name = "unit"
		}
		fmt.Fprintf(b, "%s%s beta=%g eta=%g\n", indent, name, betaOr(n), etaOr(n))
	case KindSeries:
		fmt.Fprintf(b, "%sseries (%d blocks)\n", indent, len(n.Blocks))
	case KindParallel:
		fmt.Fprintf(b, "%sparallel (%d blocks)\n", indent, len(n.Blocks))
	case KindKofn:
		fmt.Fprintf(b, "%skofn k=%d n=%d\n", indent, kOr(n), len(n.Blocks))
	}
	for _, child := range n.Blocks {
		writeTree(b, child, depth+1)
	}
}

// Flat renders the diagram as a single line using the type strings, e.g.
// "series(parallel(unit,unit),unit)". Useful for tests and logs.
func Flat(root *Node) string {
	var b strings.Builder
	writeFlat(&b, root)
	return b.String()
}

func writeFlat(b *strings.Builder, n *Node) {
	kind, ok := KindOf(n)
	if !ok {
		fmt.Fprintf(b, "?%s", n.Type)
		return
	}
	switch kind {
	case KindUnit:
		fmt.Fprintf(b, "unit(%s)", n.Name)
		return
	case KindSeries:
		b.WriteString("series(")
	case KindParallel:
		b.WriteString("parallel(")
	case KindKofn:
		fmt.Fprintf(b, "kofn%d(", kOr(n))
	}
	for i, child := range n.Blocks {
		if i > 0 {
			b.WriteString(",")
		}
		writeFlat(b, child)
	}
	b.WriteString(")")
}
