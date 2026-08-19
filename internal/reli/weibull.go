package reli

import "math"

// Weibull reliability and hazard functions for a component with shape beta
// and scale eta. The same formulas are used for leaf nodes and for the
// closed-form checks in tests, so every layer of the evaluator shares one
// source of truth for R(t) = exp(-(t/eta)^beta).

// Reliability returns the probability that a Weibull component survives
// mission time t: exp(-(t/eta)^beta).
func Reliability(t, beta, eta float64) float64 {
	if t <= 0 {
		return 1
	}
	return math.Exp(-math.Pow(t/eta, beta))
}

// Failure returns the unreliability F(t) = 1 - R(t).
func Failure(t, beta, eta float64) float64 {
	return 1 - Reliability(t, beta, eta)
}

// Hazard returns the instantaneous failure rate
// lambda(t) = (beta/eta) * (t/eta)^(beta-1).
func Hazard(t, beta, eta float64) float64 {
	if t <= 0 {
		return 0
	}
	return (beta / eta) * math.Pow(t/eta, beta-1)
}

// Density returns the Weibull probability density f(t).
func Density(t, beta, eta float64) float64 {
	if t <= 0 {
		return 0
	}
	return Hazard(t, beta, eta) * Reliability(t, beta, eta)
}

// Quantile returns the time at which reliability reaches r (0 < r < 1):
// t = eta * (-ln r)^(1/beta).
func Quantile(r, beta, eta float64) float64 {
	if r <= 0 || r >= 1 {
		return math.Inf(1)
	}
	return eta * math.Pow(-math.Log(r), 1/beta)
}

// Median is the 50% reliability quantile.
func Median(beta, eta float64) float64 {
	return Quantile(0.5, beta, eta)
}

// ShapeLimit classifies a Weibull shape for diagnostics: below 1 is a
// decreasing hazard ("infant mortality"), equal to 1 is constant and above
// 1 is an increasing hazard ("wear out").
func ShapeLimit(beta float64) string {
	switch {
	case beta < 1:
		return "decreasing hazard"
	case beta > 1:
		return "increasing hazard"
	default:
		return "constant hazard"
	}
}
