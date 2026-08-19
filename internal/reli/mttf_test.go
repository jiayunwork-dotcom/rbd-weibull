package reli

import (
	"math"
	"testing"

	"rbd-weibull/internal/model"
)

func TestUnitMTTFClosedForm(t *testing.T) {
	// MTTF = eta * Gamma(1 + 1/beta); beta=1 gives exactly eta.
	if math.Abs(UnitMTTF(1, 100)-100) > 1e-9 {
		t.Fatalf("UnitMTTF(1,100) = %g, want 100", UnitMTTF(1, 100))
	}
	// beta=2, eta=100: 100 * Gamma(1.5) = 100 * 0.8862269...
	want := 100 * math.Gamma(1.5)
	if math.Abs(UnitMTTF(2, 100)-want) > 1e-9 {
		t.Fatalf("UnitMTTF(2,100) = %g, want %g", UnitMTTF(2, 100), want)
	}
}

func TestSystemMTTFUnitMatch(t *testing.T) {
	// The numerical system integrator on a single unit must reproduce the
	// closed-form unit MTTF within tolerance.
	root := unit(1.5, 100)
	got, err := SystemMTTF(root, 0, 2000)
	if err != nil {
		t.Fatalf("SystemMTTF returned error: %v", err)
	}
	want := UnitMTTF(1.5, 100)
	if math.Abs(got-want)/want > 1e-3 {
		t.Fatalf("numerical MTTF = %g, closed form %g (rel err %.2e)", got, want, math.Abs(got-want)/want)
	}
}

func TestSystemMTTFExplicitHorizon(t *testing.T) {
	// An explicit horizon behaves like the default one for the same
	// diagram: both integrate R(t) and agree to high precision.
	root := unit(2.0, 50)
	withDefault, err := SystemMTTF(root, 0, 1000)
	if err != nil {
		t.Fatalf("SystemMTTF(default) returned error: %v", err)
	}
	withExplicit, err := SystemMTTF(root, 400, 1000)
	if err != nil {
		t.Fatalf("SystemMTTF(explicit) returned error: %v", err)
	}
	if math.Abs(withDefault-withExplicit)/withDefault > 1e-6 {
		t.Fatalf("default/explicit MTTF disagree: %g vs %g", withDefault, withExplicit)
	}
}

func TestMTTFHorizonScalesWithEta(t *testing.T) {
	root := &model.Node{
		Type: model.TypeSeries,
		Blocks: []*model.Node{
			unit(1.8, 4500),
			unit(2.5, 20000),
		},
	}
	h := MTTFHorizon(root)
	if math.Abs(h-160000) > 1e-9 {
		t.Fatalf("horizon = %g, want 8 * 20000 = 160000", h)
	}
}
