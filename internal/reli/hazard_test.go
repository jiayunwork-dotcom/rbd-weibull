package reli

import (
	"math"
	"testing"

	"rbd-weibull/internal/model"
)

func TestSystemHazardApproximation(t *testing.T) {
	// A single Weibull unit's hazard is exact; the numeric forward
	// difference must land on the closed form.
	root := unit(2, 100)
	got, err := SystemHazard(root, 50)
	if err != nil {
		t.Fatalf("SystemHazard returned error: %v", err)
	}
	want := Hazard(50, 2, 100)
	if math.Abs(got-want)/want > 1e-3 {
		t.Fatalf("hazard = %g, want %g (rel err %.2e)", got, want, math.Abs(got-want)/want)
	}
}

func TestSeriesHazardSum(t *testing.T) {
	// For independent series components the system hazard is the sum of
	// the component hazards; the numeric approximation must track it.
	root := &model.Node{
		Type: model.TypeSeries,
		Blocks: []*model.Node{
			unit(1.8, 4500),
			unit(2.5, 20000),
		},
	}
	got, err := SystemHazard(root, 1000)
	if err != nil {
		t.Fatalf("SystemHazard returned error: %v", err)
	}
	want := Hazard(1000, 1.8, 4500) + Hazard(1000, 2.5, 20000)
	if math.Abs(got-want)/want > 5e-3 {
		t.Fatalf("series hazard = %g, want sum %g (rel err %.2e)", got, want, math.Abs(got-want)/want)
	}
}
