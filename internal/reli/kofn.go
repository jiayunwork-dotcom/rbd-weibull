package reli

import (
	"fmt"

	"rbd-weibull/internal/model"
)

// foldKofn combines children with a "k of n" voting rule: the block
// survives while at least k of its n children survive.
//
// When every child is identically distributed (their reliabilities agree
// within a small tolerance at time t) the survival probability is the
// binomial tail sum over j = k..n of C(n,j) R^j (1-R)^(n-j), which is
// closed-form and cheap for any n.
//
// Otherwise the evaluator enumerates all 2^n working/failed states and
// sums the probability of the states with at least k working children.
// Because that enumeration grows exponentially, n is capped at
// MaxEnumerateSubblocks and larger non-identical blocks are rejected with
// an error instead of hanging the process.
func foldKofn(n *model.Node, t float64) (float64, error) {
	rs, err := foldMany(n, t)
	if err != nil {
		return 0, err
	}
	k := *n.K
	if k > len(rs) {
		return 0, fmt.Errorf("%s: %w", model.Describe(n), model.ErrKGreaterN)
	}
	if AllEqual(rs, 1e-12) {
		return BinomAtLeast(k, len(rs), rs[0]), nil
	}
	if len(rs) > MaxEnumerateSubblocks {
		if err := dropTooMany(fmt.Errorf("%s: n=%d > %d: %w",
			model.Describe(n), len(rs), MaxEnumerateSubblocks, model.ErrTooManySubblocks)); err != nil {
			return 0, err
		}
	}
	return enumerateAtLeast(k, rs), nil
}

// enumerateAtLeast computes the probability that at least k of the given
// independent reliabilities survive by walking every state of a bit mask
// and accumulating the product of R_i for working bits and (1-R_i) for
// failed bits.
func enumerateAtLeast(k int, rs []float64) float64 {
	n := len(rs)
	total := 0.0
	limit := 1 << n
	for mask := 0; mask < limit; mask++ {
		working := 0
		p := 1.0
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				working++
				p *= rs[i]
			} else {
				p *= 1 - rs[i]
			}
		}
		if working >= k {
			total += p
		}
	}
	return total
}
