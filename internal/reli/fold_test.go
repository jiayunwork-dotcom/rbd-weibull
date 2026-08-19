package reli

import (
	"math"
	"testing"

	"rbd-weibull/internal/model"
)

func unit(b, e float64) *model.Node {
	return &model.Node{Type: model.TypeUnit, Beta: &b, Eta: &e}
}

func TestSeriesFold(t *testing.T) {
	// exp(-1) each, three in series: exp(-3).
	root := &model.Node{
		Type:   model.TypeSeries,
		Blocks: []*model.Node{unit(1, 1), unit(1, 1), unit(1, 1)},
	}
	r, err := System(root, 1)
	if err != nil {
		t.Fatalf("System returned error: %v", err)
	}
	if math.Abs(r-math.Exp(-3)) > 1e-12 {
		t.Fatalf("series R = %g, want %g", r, math.Exp(-3))
	}
}

func TestParallelFoldClosedForm(t *testing.T) {
	// Two identically distributed units in parallel must equal
	// 1 - (1-R)^2 at the shared mission time.
	root := &model.Node{
		Type:   model.TypeParallel,
		Blocks: []*model.Node{unit(1.8, 4500), unit(1.8, 4500)},
	}
	r, err := System(root, 1000)
	if err != nil {
		t.Fatalf("System returned error: %v", err)
	}
	one := Reliability(1000, 1.8, 4500)
	want := 1 - (1-one)*(1-one)
	if math.Abs(r-want) > 1e-12 {
		t.Fatalf("parallel R = %g, want closed form %g", r, want)
	}
}

func TestPump2DiagramManual(t *testing.T) {
	// Two pumps in parallel, then a valve in series; hand computation:
	// R_pump = exp(-(1000/4500)^1.8), R_par = 1-(1-R_pump)^2,
	// R_valve = exp(-(1000/20000)^2.5), R_sys = R_par * R_valve.
	root := &model.Node{
		Type: model.TypeSeries,
		Blocks: []*model.Node{
			{
				Type:   model.TypeParallel,
				Blocks: []*model.Node{unit(1.8, 4500), unit(1.8, 4500)},
			},
			unit(2.5, 20000),
		},
	}
	r, err := System(root, 1000)
	if err != nil {
		t.Fatalf("System returned error: %v", err)
	}
	rp := Reliability(1000, 1.8, 4500)
	rv := Reliability(1000, 2.5, 20000)
	want := (1 - (1-rp)*(1-rp)) * rv
	if math.Abs(r-want) > 1e-12 {
		t.Fatalf("system R = %g, want %g", r, want)
	}
	if r < 0.995 || r > 0.996 {
		t.Fatalf("system R = %g, want ~0.9953", r)
	}
}

func TestForcedLeafReliability(t *testing.T) {
	// A series of two units, the first forced to 1: system R equals the
	// second unit's R. Forced to 0: system R is 0.
	root := &model.Node{
		Type:   model.TypeSeries,
		Blocks: []*model.Node{unit(1, 2), unit(1, 4)},
	}
	if err := model.SetFixed(root, []int{0}, 1); err != nil {
		t.Fatalf("SetFixed returned error: %v", err)
	}
	r, err := System(root, 1)
	if err != nil {
		t.Fatalf("System returned error: %v", err)
	}
	want := Reliability(1, 1, 4)
	if math.Abs(r-want) > 1e-12 {
		t.Fatalf("forced-working R = %g, want %g", r, want)
	}
	if err := model.SetFixed(root, []int{0}, 0); err != nil {
		t.Fatalf("SetFixed returned error: %v", err)
	}
	r, err = System(root, 1)
	if err != nil {
		t.Fatalf("System returned error: %v", err)
	}
	if r != 0 {
		t.Fatalf("forced-failed R = %g, want 0", r)
	}
}
