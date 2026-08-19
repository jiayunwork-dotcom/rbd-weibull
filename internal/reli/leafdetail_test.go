package reli

import (
	"math"
	"testing"
)

func TestLeafDetails(t *testing.T) {
	root := seriesTwo(0.9, 0.7)
	rows, err := LeafDetails(root, 1)
	if err != nil {
		t.Fatalf("LeafDetails returned error: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("LeafDetails returned %d rows, want 2", len(rows))
	}
	// Row 0 is the stronger component (R=0.9), row 1 the weaker (R=0.7);
	// the weaker one carries the larger importance.
	if math.Abs(rows[0].R-0.9) > 1e-9 || math.Abs(rows[1].R-0.7) > 1e-9 {
		t.Fatalf("leaf reliabilities = %g, %g; want 0.9, 0.7", rows[0].R, rows[1].R)
	}
	if !(rows[1].Importance > rows[0].Importance) {
		t.Fatal("weaker series component must have the larger importance")
	}
	if !(rows[0].MTTF > rows[1].MTTF) {
		t.Fatal("stronger component must have the larger MTTF")
	}
}

func TestRankByImportance(t *testing.T) {
	root := seriesTwo(0.9, 0.7)
	ranked, err := RankByImportance(root, 1)
	if err != nil {
		t.Fatalf("RankByImportance returned error: %v", err)
	}
	if len(ranked) != 2 {
		t.Fatalf("ranked = %v, want two names", ranked)
	}
	// The weakest series component ranks first.
	if ranked[0] != "root.blocks[1]" || ranked[1] != "root.blocks[0]" {
		t.Fatalf("ranked order = %v, want weakest first", ranked)
	}
}

func TestFractionalImportance(t *testing.T) {
	root := seriesTwo(0.9, 0.7)
	frac, err := FractionalImportance(root, 1)
	if err != nil {
		t.Fatalf("FractionalImportance returned error: %v", err)
	}
	total := frac[0].Value + frac[1].Value
	if math.Abs(total-1) > 1e-12 {
		t.Fatalf("fractional importances sum to %g, want 1", total)
	}
}
