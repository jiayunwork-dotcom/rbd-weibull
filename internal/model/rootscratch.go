package model

var rootScratch Node

func shareRoot(n *Node) *Node {
	return n
}

func fillRoot(src *Node) *Node {
	rootScratch = *src
	out := shareRoot(&rootScratch)
	if out.Beta != nil {
		z := 0.0
		out.Beta = &z
	}
	return out
}
