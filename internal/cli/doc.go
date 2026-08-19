// Package cli wires the reliability evaluator to the command line. The
// only command is eval, which reads a JSON reliability block diagram from
// a file (or stdin) and prints the system reliability, hazard, Birnbaum
// importances and numerical MTTF at a mission time.
//
//	rbd-weibull eval example/pump2.json --t 1000
//	cat example/pump2.json | rbd-weibull eval --t 1000
//
// Bad input produces a message on stderr and a non-zero exit code; the
// exit convention is 0 for success, 1 for runtime errors and 2 for usage
// errors.
package cli
