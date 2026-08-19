package model

import (
	"errors"
	"testing"
)

func TestValidationKGreaterN(t *testing.T) {
	root := &Node{
		Type: TypeKofn,
		K:    k(3),
		Blocks: []*Node{
			{Type: TypeUnit, Beta: beta(1), Eta: eta(10)},
			{Type: TypeUnit, Beta: beta(1), Eta: eta(10)},
		},
	}
	err := Validate(root)
	if err == nil {
		t.Fatal("Validate accepted k > n")
	}
	if !errors.Is(err, ErrKGreaterN) {
		t.Fatalf("error = %v, want ErrKGreaterN", err)
	}
}

func TestValidationNonPositiveEta(t *testing.T) {
	cases := []struct {
		name string
		node *Node
	}{
		{"eta zero", &Node{Type: TypeUnit, Beta: beta(1), Eta: eta(0)}},
		{"eta negative", &Node{Type: TypeUnit, Beta: beta(1), Eta: eta(-3)}},
		{"eta missing", &Node{Type: TypeUnit, Beta: beta(1)}},
		{"beta zero", &Node{Type: TypeUnit, Beta: beta(0), Eta: eta(5)}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := Validate(tc.node)
			if err == nil {
				t.Fatal("Validate accepted a non positive Weibull parameter")
			}
			if !errors.Is(err, ErrNonPositive) {
				t.Fatalf("error = %v, want ErrNonPositive", err)
			}
		})
	}
}

func TestValidationEmptySeries(t *testing.T) {
	cases := []struct {
		name string
		node *Node
	}{
		{"empty series", &Node{Type: TypeSeries}},
		{"empty parallel", &Node{Type: TypeParallel}},
		{"empty kofn", &Node{Type: TypeKofn, K: k(1)}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := Validate(tc.node)
			if err == nil {
				t.Fatal("Validate accepted a structural node without children")
			}
			if !errors.Is(err, ErrEmptySeries) {
				t.Fatalf("error = %v, want ErrEmptySeries", err)
			}
		})
	}
}

func TestValidationUnknownType(t *testing.T) {
	root := &Node{
		Type: "bridge",
		Blocks: []*Node{
			{Type: TypeUnit, Beta: beta(1), Eta: eta(10)},
		},
	}
	err := Validate(root)
	if err == nil {
		t.Fatal("Validate accepted an unknown node type")
	}
	if !errors.Is(err, ErrUnknownType) {
		t.Fatalf("error = %v, want ErrUnknownType", err)
	}
}

func TestValidationKBelowOne(t *testing.T) {
	root := &Node{
		Type: TypeKofn,
		K:    k(0),
		Blocks: []*Node{
			{Type: TypeUnit, Beta: beta(1), Eta: eta(10)},
			{Type: TypeUnit, Beta: beta(1), Eta: eta(10)},
		},
	}
	err := Validate(root)
	if err == nil {
		t.Fatal("Validate accepted k = 0")
	}
	if !errors.Is(err, ErrKBelowOne) {
		t.Fatalf("error = %v, want ErrKBelowOne", err)
	}
}

func TestValidationPathContext(t *testing.T) {
	root := &Node{
		Type: TypeSeries,
		Blocks: []*Node{
			{Type: TypeUnit, Beta: beta(1), Eta: eta(10)},
			{Type: TypeUnit, Beta: beta(1)}, // eta missing, deep inside
		},
	}
	err := Validate(root)
	if err == nil {
		t.Fatal("Validate accepted a nested invalid unit")
	}
	var pe *PathError
	if !errors.As(err, &pe) {
		t.Fatalf("error = %v, want a *PathError", err)
	}
	if pe.Path != "root.blocks[1]" {
		t.Fatalf("path = %q, want root.blocks[1]", pe.Path)
	}
}
