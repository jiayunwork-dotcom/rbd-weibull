package reli

import (
	"fmt"

	"rbd-weibull/internal/model"
)

// Point is one sample of the system reliability curve at a mission time.
type Point struct {
	T float64
	R float64
}

// Points is a sampled reliability curve with trapezoid integration and
// endpoint accessors.
type Points []Point

// Curve samples the system reliability R(t) on a regular grid from t0 to
// t1 (inclusive) with n points. n smaller than two is treated as two; t0
// may be zero because the evaluator is only called at strictly positive
// times after the first sample.
func Curve(root *model.Node, t0, t1 float64, n int) (Points, error) {
	if root == nil {
		return nil, fmt.Errorf("nil diagram")
	}
	if n < 2 {
		n = 2
	}
	if t1 <= t0 {
		return nil, fmt.Errorf("curve horizon %g must exceed start %g", t1, t0)
	}
	out := make(Points, 0, n)
	for i := 0; i < n; i++ {
		t := t0 + (t1-t0)*float64(i)/float64(n-1)
		tEval := t
		if tEval <= 0 {
			tEval = 1e-9 // keep the evaluator strictly positive; keep the reported T as 0
		}
		r, err := System(root, tEval)
		if err != nil {
			return nil, err
		}
		out = append(out, Point{T: t, R: r})
	}
	return out, nil
}

// Area integrates the sampled curve with the trapezoid rule. With t0 = 0
// and a horizon that reaches negligible reliability this equals the system
// MTTF and provides an independent cross-check of SystemMTTF.
func (ps Points) Area() float64 {
	if len(ps) < 2 {
		return 0
	}
	area := 0.0
	for i := 1; i < len(ps); i++ {
		h := ps[i].T - ps[i-1].T
		area += 0.5 * (ps[i].R + ps[i-1].R) * h
	}
	return area
}

// Start and End return the first and last sample times.
func (ps Points) Start() float64 {
	if len(ps) == 0 {
		return 0
	}
	return ps[0].T
}

func (ps Points) End() float64 {
	if len(ps) == 0 {
		return 0
	}
	return ps[len(ps)-1].T
}
