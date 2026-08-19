package reli

import (
	"math"
	"testing"
)

func TestFitExactUnit(t *testing.T) {
	// Sampling a single unit and fitting an equivalent Weibull must
	// recover the original beta and eta to high precision.
	beta, eta := 1.7, 250.0
	pts := make([]Point, 0, 50)
	for i := 0; i < 50; i++ {
		tt := 1.0 + float64(i)*8.0
		pts = append(pts, Point{T: tt, R: Reliability(tt, beta, eta)})
	}
	fit, err := FitWeibull(pts)
	if err != nil {
		t.Fatalf("FitWeibull returned error: %v", err)
	}
	if math.Abs(fit.Beta-beta)/beta > 1e-9 {
		t.Fatalf("fit beta = %g, want %g", fit.Beta, beta)
	}
	if math.Abs(fit.Eta-eta)/eta > 1e-9 {
		t.Fatalf("fit eta = %g, want %g", fit.Eta, eta)
	}
}

func TestFitDegenerateInput(t *testing.T) {
	// Samples with R out of (0,1) are skipped; two identical samples are
	// still not enough to fit two parameters... but all-degenerate input
	// must error cleanly.
	_, err := FitWeibull([]Point{{T: 1, R: 0}, {T: 2, R: 1}, {T: 3, R: 0}})
	if err == nil {
		t.Fatal("FitWeibull accepted degenerate samples")
	}
}

func TestFitExponentialRate(t *testing.T) {
	// An exponential (beta=1) system sampled at integer times must fit to
	// beta close to 1 and eta close to the reciprocal rate.
	rate := 0.01
	pts := make([]Point, 0, 30)
	for i := 1; i <= 30; i++ {
		pts = append(pts, Point{T: float64(i), R: math.Exp(-rate * float64(i))})
	}
	fit, err := FitWeibull(pts)
	if err != nil {
		t.Fatalf("FitWeibull returned error: %v", err)
	}
	if math.Abs(fit.Beta-1) > 1e-3 {
		t.Fatalf("fit beta = %g, want ~1", fit.Beta)
	}
	if math.Abs(fit.Eta-1/rate)/(1/rate) > 1e-3 {
		t.Fatalf("fit eta = %g, want %g", fit.Eta, 1/rate)
	}
}
