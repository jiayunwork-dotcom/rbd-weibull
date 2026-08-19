package model

import "testing"

func pump2Diagram() *Node {
	return &Node{
		Type: TypeSeries,
		Blocks: []*Node{
			{
				Type: TypeParallel,
				Blocks: []*Node{
					{Type: TypeUnit, Name: "pump-A", Beta: beta(1.8), Eta: eta(4500)},
					{Type: TypeUnit, Name: "pump-B", Beta: beta(1.8), Eta: eta(4500)},
				},
			},
			{Type: TypeUnit, Name: "valve", Beta: beta(2.5), Eta: eta(20000)},
		},
	}
}

func TestWalkLeaves(t *testing.T) {
	leaves := Leaves(pump2Diagram())
	if len(leaves) != 3 {
		t.Fatalf("Leaves returned %d leaves, want 3", len(leaves))
	}
	want := [][]int{{0, 0}, {0, 1}, {1}}
	for i, leaf := range leaves {
		if len(leaf.Path) != len(want[i]) {
			t.Fatalf("leaf %d path = %v, want %v", i, leaf.Path, want[i])
		}
		for j := range want[i] {
			if leaf.Path[j] != want[i][j] {
				t.Fatalf("leaf %d path = %v, want %v", i, leaf.Path, want[i])
			}
		}
	}
	if leaves[0].Node.Name != "pump-A" || leaves[2].Node.Name != "valve" {
		t.Fatalf("leaf order wrong: %q, %q", leaves[0].Node.Name, leaves[2].Node.Name)
	}
}

func TestCountAndDepth(t *testing.T) {
	root := pump2Diagram()
	if got := Count(root); got != 5 {
		t.Fatalf("Count = %d, want 5", got)
	}
	if got := Depth(root); got != 3 {
		t.Fatalf("Depth = %d, want 3", got)
	}
}

func TestCloneIsolation(t *testing.T) {
	root := pump2Diagram()
	clone := Clone(root)
	if err := SetFixed(clone, []int{1}, 0.25); err != nil {
		t.Fatalf("SetFixed returned error: %v", err)
	}
	if root.Blocks[1].Fixed != nil {
		t.Fatal("mutating the clone changed the source diagram")
	}
	if clone.Blocks[1].Fixed == nil || *clone.Blocks[1].Fixed != 0.25 {
		t.Fatalf("clone override = %v, want 0.25", clone.Blocks[1].Fixed)
	}
}

func TestLocateAndFindUnit(t *testing.T) {
	root := pump2Diagram()
	if n := Locate(root, []int{0, 1}); n == nil || n.Name != "pump-B" {
		t.Fatalf("Locate(root, {0,1}) = %v, want pump-B", n)
	}
	if n := Locate(root, []int{5}); n != nil {
		t.Fatal("Locate returned a node for an out of range chain")
	}
	if n, ok := FindUnit(root, "valve"); !ok || n.Eta == nil || *n.Eta != 20000 {
		t.Fatalf("FindUnit(valve) = %v/%v, want eta 20000", n, ok)
	}
	if _, ok := FindUnit(root, "missing"); ok {
		t.Fatal("FindUnit found a component that is not in the diagram")
	}
}

func TestIsSameStructure(t *testing.T) {
	a := &Node{Type: TypeUnit, Beta: beta(1.8), Eta: eta(4500)}
	b := &Node{Type: TypeUnit, Beta: beta(1.8), Eta: eta(4500)}
	c := &Node{Type: TypeUnit, Beta: beta(2.0), Eta: eta(4500)}
	if !IsSameStructure(a, b, 1e-12) {
		t.Fatal("identical units reported as different")
	}
	if IsSameStructure(a, c, 1e-12) {
		t.Fatal("units with different beta reported as identical")
	}
}
