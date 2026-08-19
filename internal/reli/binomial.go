package reli

import "math"

// Binom returns the binomial coefficient C(n, k) computed by an iterative
// product that avoids factorial overflow for the moderate n the evaluator
// handles (i.i.d. k-of-n blocks can be large).
func Binom(n, k int) float64 {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	c := 1.0
	for i := 1; i <= k; i++ {
		c = c * float64(n-k+i) / float64(i)
	}
	return c
}

// BinomAtLeast returns the binomial tail probability that at least k of n
// independent trials, each succeeding with probability p, succeed:
//
//	sum_{j=k}^{n} C(n,j) p^j (1-p)^{n-j}
//
// The sum is walked from j=0 with each term derived from the previous one
// by a single multiply and divide, so the whole tail costs O(n) work
// regardless of k.
func BinomAtLeast(k, n int, p float64) float64 {
	if k <= 0 {
		return 1
	}
	if k > n {
		return 0
	}
	if p <= 0 {
		return 0
	}
	if p >= 1 {
		return 1
	}
	q := 1 - p
	term := math.Pow(q, float64(n)) // C(n,0) p^0 q^n
	coef := 1.0
	tail := 0.0
	for j := 0; j <= n; j++ {
		if j >= k {
			tail += coef * term
		}
		if j < n {
			coef = coef * float64(n-j) / float64(j+1)
			term = term * p / q
		}
	}
	return tail
}
