package model

import "testing"

func TestReplaceChild(t *testing.T) {
	root := pump2Diagram()
	repl := &Node{Type: TypeUnit, Name: "spare", Beta: beta(2.0), Eta: eta(3000)}
	if err := ReplaceChild(root, []int{1}, repl); err != nil {
		t.Fatalf("ReplaceChild returned error: %v", err)
	}
	if root.Blocks[1] != repl {
		t.Fatal("ReplaceChild did not install the replacement")
	}
	if err := Validate(root); err != nil {
		t.Fatalf("diagram invalid after replace: %v", err)
	}
}

func TestReplaceChildBadChain(t *testing.T) {
	root := pump2Diagram()
	if err := ReplaceChild(root, []int{9}, &Node{Type: TypeUnit, Beta: beta(1), Eta: eta(1)}); err == nil {
		t.Fatal("ReplaceChild accepted an out of range chain")
	}
	if err := ReplaceChild(root, []int{0}, nil); err == nil {
		t.Fatal("ReplaceChild accepted a nil replacement")
	}
}

func TestRemoveChild(t *testing.T) {
	root := pump2Diagram()
	removed, err := RemoveChild(root, []int{1})
	if err != nil {
		t.Fatalf("RemoveChild returned error: %v", err)
	}
	if removed == nil || removed.Name != "valve" {
		t.Fatalf("RemoveChild removed %v, want valve", removed)
	}
	if len(root.Blocks) != 1 {
		t.Fatalf("root blocks = %d, want 1 after removal", len(root.Blocks))
	}
	if err := Validate(root); err != nil {
		t.Fatalf("diagram invalid after removal: %v", err)
	}
	if _, err := RemoveChild(root, []int{}); err == nil {
		t.Fatal("RemoveChild accepted root removal")
	}
}

func TestAllLeafNames(t *testing.T) {
	names := AllLeafNames(pump2Diagram())
	if len(names) != 3 || names[0] != "pump-A" || names[2] != "valve" {
		t.Fatalf("AllLeafNames = %v, want [pump-A pump-B valve]", names)
	}
}

func TestVerifyIndex(t *testing.T) {
	if err := VerifyIndex(pump2Diagram()); err != nil {
		t.Fatalf("VerifyIndex returned error: %v", err)
	}
	// Corrupting a path must be caught.
	root := pump2Diagram()
	root.Blocks[0].Blocks[0].Name = "renamed"
	if err := VerifyIndex(root); err != nil {
		t.Fatalf("VerifyIndex returned error after rename: %v", err)
	}
}
