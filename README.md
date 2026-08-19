# rbd-weibull

Reliability block diagram evaluator for nested series / parallel / k-of-n
structures of Weibull components. Given a diagram in JSON and a mission
time t, it prints the system reliability R(t), an approximation of the
system hazard rate, the Birnbaum importance of every unit leaf and the
numerical system MTTF.

Leaf components carry a Weibull shape `beta` and scale `eta`, with
`R(t) = exp(-(t/eta)^beta)`; structural nodes are evaluated at the same t
as their children, never at an MTTF substituted for a reliability. Series
nodes multiply, parallel nodes use `1 - prod(1 - R_i)`, and k-of-n nodes
use the binomial tail sum when children are identically distributed or an
exact 2^n state enumeration otherwise (capped at 12 sub-blocks, beyond
which the diagram is rejected). Input with unknown node types, empty
structural nodes, k > n, or non-positive beta/eta/time is rejected.

Birnbaum importance is `I_i(t) = R(t | i works) - R(t | i fails)`, each
side obtained from a full system evaluation of a cloned tree with the
leaf's reliability forced to 1 or 0.

## Build

```bash
go build .
go test ./...
```

## Usage

```bash
go run . eval example/pump2.json --t 1000
go run . eval example/kofn.json --t 2000
cat example/pump2.json | go run . eval --t 1000
go run . eval example/series3.json -t 500
```

`eval` reads a JSON file or stdin and prints R(t), hazard, MTTF and the
per-leaf importances. Exit codes: 0 success, 1 runtime/input error, 2
usage error.

## Numerical settings

The system MTTF is a trapezoid integral of R(t) over `[0, tMax]` with
`tMax = 8 * max(eta)` in the diagram and 2000 steps by default; at the
horizon the slowest component has decayed to at most `exp(-8^beta)`, so
truncation is below printed precision. The single-component closed form
`MTTF = eta * Gamma(1 + 1/beta)` is used as the analytic reference in
tests. The hazard is a forward difference of `ln R(t)` with a relative
step `1e-4 * t`.

## Example

`example/pump2.json` models two pumps in parallel upstream of a valve in
series. At t = 1000 h the system reliability is about 0.9953.
