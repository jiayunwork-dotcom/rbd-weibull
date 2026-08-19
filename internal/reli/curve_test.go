package reli

import (
	"math"
	"testing"
)

func TestCurveGrid(t *testing.T) {
	root := unit(1.8, 4500)
	pts, err := Curve(root, 0, 4500, 5)
	if err != nil {
		t.Fatalf("Curve returned error: %v", err)
	}
	if len(pts) != 5 {
		t.Fatalf("Curve returned %d points, want 5", len(pts))
	}
	if pts[0].T != 0 || pts[4].T != 4500 {
		t.Fatalf("grid endpoints = %g..%g, want 0..4500", pts[0].T, pts[4].T)
	}
	// R decreases strictly along the curve.
	for i := 1; i < len(pts); i++ {
		if pts[i].R >= pts[i-1].R {
			t.Fatalf("curve not strictly decreasing at %d: %g >= %g", i, pts[i].R, pts[i-1].R)
		}
	}
}

func TestCurveAreaMatchesMTTF(t *testing.T) {
	// Integrating the sampled R(t) curve with the trapezoid rule must
	// agree with SystemMTTF on the same diagram.
	root := unit(1.5, 100)
	pts, err := Curve(root, 0, 800, 4000)
	if err != nil {
		t.Fatalf("Curve returned error: %v", err)
	}
	area := pts.Area()
	mttf, err := SystemMTTF(root, 800, 4000)
	if err != nil {
		t.Fatalf("SystemMTTF returned error: %v", err)
	}
	if math.Abs(area-mttf)/mttf > 1e-9 {
		t.Fatalf("curve area = %g, MTTF = %g", area, mttf)
	}
	if area < 1 || area > 200 {
		t.Fatalf("curve area out of plausible MTTF range: %g", area)
	}
}
