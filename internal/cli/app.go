package cli

import (
	"fmt"
	"io"
)

const usage = `rbd-weibull — reliability block diagram evaluator

evaluates a nested series / parallel / k-of-n diagram of Weibull components
at a mission time: system reliability, hazard approximation, Birnbaum
importance per leaf and numerical system MTTF.

usage:
  rbd-weibull eval [flags] [file]
    file            path to the diagram JSON (default: read stdin)
    -t, --t <time>  mission time, must be positive (default 1000)
  rbd-weibull curve --to <t1> [--from <t0>] [--n <count>] [file]
    sample R(t) on a regular grid and print "t R" per line
  rbd-weibull describe [file]
    print the validated diagram as an indented tree

input JSON:
  { "type": "series"|"parallel"|"kofn", "k": 2, "blocks": [ ... ] }
  unit:   { "type": "unit", "beta": 1.8, "eta": 4500 }
  kofn:   { "type": "kofn", "k": 2, "blocks": [ unit, unit, unit ] }

examples:
  rbd-weibull eval example/pump2.json --t 1000
  cat example/pump2.json | rbd-weibull eval --t 1000
  rbd-weibull eval example/kofn.json --t 2000
  rbd-weibull curve example/pump2.json --from 0 --to 80000 --n 40
`

// Run dispatches the first argument as a subcommand and returns a process
// exit code. Unknown commands and missing commands print the usage text to
// stderr and exit 2.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) < 1 {
		fmt.Fprint(stderr, usage)
		return 2
	}
	switch args[0] {
	case "eval":
		return RunEval(args[1:], stdout, stderr)
	case "curve":
		return RunCurve(args[1:], stdout, stderr)
	case "describe":
		return RunDescribe(args[1:], stdout, stderr)
	case "help", "-h", "--help":
		fmt.Fprint(stdout, usage)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command %q\n\n%s\n", args[0], usage)
		return 2
	}
}
