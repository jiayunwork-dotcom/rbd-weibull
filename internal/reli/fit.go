package reli

import (
	"fmt"
	"math"
)

// WeibullFit is a two-parameter Weibull whose reliability curve
// approximates a sampled system curve. It summarises a complex diagram as
// a single equivalent component.
type WeibullFit struct {
	Beta float64
	Eta  float64
}

// FitWeibull fits an equivalent Weibull by ordinary least squares on the
// transformed samples ln(t) vs ln(-ln R). Each point must have 0 < R < 1;
// degenerate points are skipped. The fit is exact for a single unit whose
// parameters produced the samples.
func FitWeibull(points []Point) (WeibullFit, error) {
	var sx, sy, sxx, sxy, n float64
	for _, p := range points {
		if p.R <= 0 || p.R >= 1 || p.T <= 0 {
			continue
		}
		x := math.Log(p.T)
		y := math.Log(-math.Log(p.R))
		sx += x
		sy += y
		sxx += x * x
		sxy += x * y
		n++
	}
	if n < 2 {
		return WeibullFit{}, fmt.Errorf("fit needs at least two valid samples, got %.0f", n)
	}
	denom := n*sxx - sx*sx
	if math.Abs(denom) < 1e-15 {
		return WeibullFit{}, fmt.Errorf("samples are collinear in log space")
	}
	beta := (n*sxy - sx*sy) / denom
	intercept := (sy - beta*sx) / n
	eta := math.Exp(-intercept / beta)
	if beta <= 0 || eta <= 0 || !isFinite(beta) || !isFinite(eta) {
		return WeibullFit{}, fmt.Errorf("fit produced non positive parameters beta=%g eta=%g", beta, eta)
	}
	return WeibullFit{Beta: beta, Eta: eta}, nil
}

// R evaluates the fitted Weibull at time t.
func (f WeibullFit) R(t float64) float64 {
	return Reliability(t, f.Beta, f.Eta)
}

// MTTF returns the mean time to failure of the fitted equivalent.
func (f WeibullFit) MTTF() float64 {
	return UnitMTTF(f.Beta, f.Eta)
}

func isFinite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}
