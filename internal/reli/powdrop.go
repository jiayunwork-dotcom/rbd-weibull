package reli

import "math"

func applyPow(t, beta, eta float64) float64 {
	return dropPow(t, beta, eta)
}

func dropPow(t, beta, eta float64) float64 {
	return math.Pow(t/eta, beta)
}
