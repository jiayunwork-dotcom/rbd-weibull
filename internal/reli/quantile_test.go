package reli

import (
	"math"
	"testing"
)

func TestSystemQuantileUnitClosedForm(t *testing.T) {
	// For a single unit the bisection must reproduce the closed-form
	// Weibull quantile t = eta (-ln r)^(1/beta).
	root := unit(1.5, 100)
	got, err := SystemQuantile(root, 0.5, 0, 60)
	if err != nil {
		t.Fatalf("SystemQuantile returned error: %v", err)
	}
	want := math.Pow(-math.Log(0.5), 1/1.5) * 100
	if math.Abs(got-want)/want > 1e-4 {
		t.Fatalf("quantile = %g, closed form %g", got, want)
	}
}

func TestReliabilityLife(t *testing.T) {
	root := unit(1.0, 100)
	life, err := ReliabilityLife(root)
	if err != nil {
		t.Fatalf("ReliabilityLife returned error: %v", err)
	}
	// Exponential: R(t)=0.9 at t = -ln(0.9) * eta.
	want := -math.Log(0.9) * 100
	if math.Abs(life-want)/want > 1e-4 {
		t.Fatalf("life = %g, want %g", life, want)
	}
	b50, err := B50(root)
	if err != nil {
		t.Fatalf("B50 returned error: %v", err)
	}
	if math.Abs(b50-100*math.Ln2)/100 > 1e-4 {
		t.Fatalf("B50 = %g, want %g", b50, 100*math.Ln2)
	}
}

func TestSystemQuantileBracket(t *testing.T) {
	// The bisection bracket must satisfy R(0)=1 > r > R(horizon).
	root := unit(2.0, 50)
	r := 0.9
	q, err := SystemQuantile(root, r, 0, 0)
	if err != nil {
		t.Fatalf("SystemQuantile returned error: %v", err)
	}
	got := Reliability(q, 2.0, 50)
	if math.Abs(got-r) > 1e-6 {
		t.Fatalf("R(quantile) = %g, want %g", got, r)
	}
}
