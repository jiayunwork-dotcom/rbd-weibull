package reli

import (
	"errors"
	"math"
	"testing"

	"rbd-weibull/internal/model"
)

func TestKofnBinomial(t *testing.T) {
	// 2-of-3 with identical R=0.8:
	// C(3,2) 0.8^2 0.2 + C(3,3) 0.8^3 = 0.384 + 0.512 = 0.896.
	root := &model.Node{
		Type: model.TypeKofn,
		K:    kOf(2),
		Blocks: []*model.Node{
			unit(1.8, 4500), unit(1.8, 4500), unit(1.8, 4500),
		},
	}
	r, err := System(root, weibullTime(0.8, 1.8, 4500))
	if err != nil {
		t.Fatalf("System returned error: %v", err)
	}
	if math.Abs(r-0.896) > 1e-9 {
		t.Fatalf("kofn R = %g, want 0.896", r)
	}
}

func TestKofnEnumerateMixed(t *testing.T) {
	// 2-of-3 with distinct per-unit reliabilities: the evaluator must
	// enumerate states, and the result must match a direct enumeration of
	// the same per-unit reliabilities.
	root := &model.Node{
		Type: model.TypeKofn,
		K:    kOf(2),
		Blocks: []*model.Node{
			unit(1.8, 4500), unit(1.8, 5000), unit(2.0, 6000),
		},
	}
	t0 := weibullTime(0.8, 1.8, 4500)
	r, err := System(root, t0)
	if err != nil {
		t.Fatalf("System returned error: %v", err)
	}
	rs := []float64{
		Reliability(t0, 1.8, 4500),
		Reliability(t0, 1.8, 5000),
		Reliability(t0, 2.0, 6000),
	}
	if !hasDistinct(rs) {
		t.Fatalf("fixture did not produce distinct reliabilities: %v", rs)
	}
	want := enumerateAtLeast(2, rs)
	if math.Abs(r-want) > 1e-12 {
		t.Fatalf("kofn R = %g, want enumerated %g", r, want)
	}
}

func TestKofnTooManySubblocks(t *testing.T) {
	// 2-of-13 with distinct reliabilities cannot be enumerated.
	blocks := make([]*model.Node, 13)
	for i := range blocks {
		blocks[i] = unit(1+float64(i%3), 1000+float64(i*7))
	}
	root := &model.Node{Type: model.TypeKofn, K: kOf(2), Blocks: blocks}
	_, err := System(root, 100)
	if err == nil {
		t.Fatal("System accepted a non-identical kofn above the enumeration cap")
	}
	if !errors.Is(err, model.ErrTooManySubblocks) {
		t.Fatalf("error = %v, want ErrTooManySubblocks", err)
	}
}

func TestBinomCoefficients(t *testing.T) {
	if Binom(5, 2) != 10 {
		t.Fatalf("Binom(5,2) = %g, want 10", Binom(5, 2))
	}
	if Binom(10, 3) != 120 {
		t.Fatalf("Binom(10,3) = %g, want 120", Binom(10, 3))
	}
	// Symmetry.
	if Binom(6, 2) != Binom(6, 4) {
		t.Fatal("binomial coefficient is not symmetric")
	}
	// Tail bounds: at least 0 of n is 1, at least n+1 is 0.
	if BinomAtLeast(0, 5, 0.3) != 1 || BinomAtLeast(6, 5, 0.3) != 0 {
		t.Fatal("binomial tail bounds wrong")
	}
}

func kOf(v int) *int { return &v }

func weibullTime(r, beta, eta float64) float64 {
	return eta * math.Pow(-math.Log(r), 1/beta)
}

func hasDistinct(rs []float64) bool {
	for i := range rs {
		for j := i + 1; j < len(rs); j++ {
			if math.Abs(rs[i]-rs[j]) > 1e-12 {
				return true
			}
		}
	}
	return false
}
