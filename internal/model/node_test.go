package model

import (
	"strings"
	"testing"
)

func beta(v float64) *float64 { return &v }
func eta(v float64) *float64  { return &v }
func k(v int) *int            { return &v }

func TestParseUnit(t *testing.T) {
	root, err := Parse(strings.NewReader(`{"type":"unit","beta":1.8,"eta":4500}`))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if root.Type != TypeUnit {
		t.Fatalf("Type = %q, want %q", root.Type, TypeUnit)
	}
	if root.Beta == nil || *root.Beta != 1.8 {
		t.Fatalf("Beta = %v, want 1.8", root.Beta)
	}
	if root.Eta == nil || *root.Eta != 4500 {
		t.Fatalf("Eta = %v, want 4500", root.Eta)
	}
	if err := Validate(root); err != nil {
		t.Fatalf("Validate(unit) returned error: %v", err)
	}
}

func TestParseNested(t *testing.T) {
	root, err := Parse(strings.NewReader(`{
		"type": "series",
		"blocks": [
			{"type": "parallel", "blocks": [
				{"type": "unit", "beta": 1.8, "eta": 4500},
				{"type": "unit", "beta": 1.8, "eta": 4500}
			]},
			{"type": "unit", "beta": 2.5, "eta": 20000}
		]
	}`))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	kind, ok := KindOf(root)
	if !ok || kind != KindSeries {
		t.Fatalf("root kind = %v/%v, want series", kind, ok)
	}
	if len(root.Blocks) != 2 {
		t.Fatalf("root blocks = %d, want 2", len(root.Blocks))
	}
	if err := Validate(root); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestParseInvalidJSON(t *testing.T) {
	_, err := Parse(strings.NewReader(`{"type": "unit", "beta": `))
	if err == nil {
		t.Fatal("Parse accepted truncated JSON")
	}
}

func TestParseMissingType(t *testing.T) {
	_, err := Parse(strings.NewReader(`{"beta": 1.0, "eta": 2.0}`))
	if err == nil {
		t.Fatal("Parse accepted a node without a type")
	}
}
