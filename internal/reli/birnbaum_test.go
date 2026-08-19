package reli

import (
	"math"
	"testing"

	"rbd-weibull/internal/model"
)

func seriesTwo(r1, r2 float64) *model.Node {
	return &model.Node{
		Type: model.TypeSeries,
		Blocks: []*model.Node{
			unit(1, etaForR(r1, 1)),
			unit(1, etaForR(r2, 1)),
		},
	}
}

func etaForR(r, t float64) float64 {
	return -t / math.Log(r)
}

func TestBirnbaumSeriesWeakestLargest(t *testing.T) {
	// In a series of two components the weakest one (smaller R) must get
	// the larger Birnbaum importance: I_i = R(t | i works) - R(t | i fails)
	// equals the other component's reliability for a series.
	root := seriesTwo(0.9, 0.7)
	imps, err := Birnbaum(root, 1)
	if err != nil {
		t.Fatalf("Birnbaum returned error: %v", err)
	}
	if len(imps) != 2 {
		t.Fatalf("Birnbaum returned %d importances, want 2", len(imps))
	}
	// I_1 = R2 = 0.7, I_2 = R1 = 0.9, so the weakest leaf ranks first.
	if math.Abs(imps[0].Value-0.7) > 1e-9 {
		t.Fatalf("I(first) = %g, want 0.7", imps[0].Value)
	}
	if math.Abs(imps[1].Value-0.9) > 1e-9 {
		t.Fatalf("I(second) = %g, want 0.9", imps[1].Value)
	}
	if !(imps[1].Value > imps[0].Value) {
		t.Fatal("weakest component must have the largest importance in series")
	}
}

func TestBirnbaumParallelRemoval(t *testing.T) {
	// Three identical units in parallel. Removing one redundant unit must
	// lower the system reliability, and the removed unit's importance in
	// the reduced system must exceed its importance in the full system.
	three := &model.Node{
		Type: model.TypeParallel,
		Blocks: []*model.Node{
			unit(1, etaForR(0.8, 1)),
			unit(1, etaForR(0.8, 1)),
			unit(1, etaForR(0.8, 1)),
		},
	}
	two := &model.Node{
		Type:   model.TypeParallel,
		Blocks: []*model.Node{unit(1, etaForR(0.8, 1)), unit(1, etaForR(0.8, 1))},
	}
	sys3, err := System(three, 1)
	if err != nil {
		t.Fatalf("System(3 units) returned error: %v", err)
	}
	sys2, err := System(two, 1)
	if err != nil {
		t.Fatalf("System(2 units) returned error: %v", err)
	}
	if !(sys2 < sys3) {
		t.Fatalf("removing a redundant unit must lower R: %g !< %g", sys2, sys3)
	}
	imps3, err := Birnbaum(three, 1)
	if err != nil {
		t.Fatalf("Birnbaum(3 units) returned error: %v", err)
	}
	imps2, err := Birnbaum(two, 1)
	if err != nil {
		t.Fatalf("Birnbaum(2 units) returned error: %v", err)
	}
	if !(imps2[0].Value > imps3[0].Value) {
		t.Fatalf("importance must rise after removal: %g !> %g", imps2[0].Value, imps3[0].Value)
	}
	// For n identical units in parallel the importance is (1-R)^(n-1):
	// 1 - (1 - (1-R)^(n-1)).
	if math.Abs(imps3[0].Value-math.Pow(1-0.8, 2)) > 1e-9 {
		t.Fatalf("I(3-unit) = %g, want %g", imps3[0].Value, math.Pow(1-0.8, 2))
	}
}

func TestBirnbaumForcedConsistency(t *testing.T) {
	// With a single unit the importance is exactly 1: forcing it working
	// gives R=1, forcing it failed gives R=0.
	root := unit(1.8, 4500)
	imps, err := Birnbaum(root, 1000)
	if err != nil {
		t.Fatalf("Birnbaum returned error: %v", err)
	}
	if len(imps) != 1 || math.Abs(imps[0].Value-1) > 1e-12 {
		t.Fatalf("single-unit importance = %v, want 1", imps)
	}
}
