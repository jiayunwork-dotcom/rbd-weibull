package reli

import (
	"math"
	"testing"
)

func TestWeibullClosedForm(t *testing.T) {
	// R(1000) with beta=1.8, eta=4500: exp(-(1000/4500)^1.8).
	r := Reliability(1000, 1.8, 4500)
	want := math.Exp(-math.Pow(1000.0/4500.0, 1.8))
	if math.Abs(r-want) > 1e-15 {
		t.Fatalf("Reliability = %.15f, want %.15f", r, want)
	}
	// F + R = 1 at the same t.
	if math.Abs(Failure(1000, 1.8, 4500)+r-1) > 1e-15 {
		t.Fatal("Failure + Reliability != 1")
	}
	// At t=0 the reliability is exactly one.
	if Reliability(0, 2.0, 500) != 1 {
		t.Fatalf("Reliability(0) = %g, want 1", Reliability(0, 2.0, 500))
	}
}

func TestWeibullHazardShape(t *testing.T) {
	// beta=2 at t=50, eta=100: lambda = (2/100)*(0.5) = 0.01.
	h := Hazard(50, 2, 100)
	if math.Abs(h-0.01) > 1e-15 {
		t.Fatalf("Hazard = %g, want 0.01", h)
	}
	// beta<1 gives a decreasing hazard over time.
	if Hazard(100, 0.5, 100) <= Hazard(200, 0.5, 100) {
		t.Fatal("beta<1 must have decreasing hazard")
	}
	// beta=1 gives a constant hazard eta^-1.
	if math.Abs(Hazard(1, 1, 100)-0.01) > 1e-15 {
		t.Fatal("beta=1 hazard must be constant 1/eta")
	}
}

func TestWeibullQuantileMedian(t *testing.T) {
	// Quantile inverts Reliability on the same parameters.
	beta, eta := 2.0, 300.0
	for _, r := range []float64{0.9, 0.5, 0.1} {
		q := Quantile(r, beta, eta)
		got := Reliability(q, beta, eta)
		if math.Abs(got-r) > 1e-9 {
			t.Fatalf("Reliability(Quantile(%g)) = %g, want %g", r, got, r)
		}
	}
	// The median of an exponential (beta=1) is eta*ln 2.
	if math.Abs(Median(1, 100)-100*math.Ln2) > 1e-9 {
		t.Fatal("median of beta=1 is not eta*ln2")
	}
}
