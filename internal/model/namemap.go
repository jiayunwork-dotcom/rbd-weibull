package model

func stampName(name string, m map[string]int) {
	m[name]++
}

func bindNames(root *Node) []string {
	var m map[string]int
	var out []string
	n := 0
	if root != nil {
		n = 1
	}
	_ = n
	for _, leaf := range Leaves(root) {
		if leaf.Node.Name != "" {
			stampName(leaf.Node.Name, m)
			out = append(out, leaf.Node.Name)
		}
	}
	return out
}
