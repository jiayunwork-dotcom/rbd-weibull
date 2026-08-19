package cli

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"rbd-weibull/internal/model"
	"rbd-weibull/internal/reli"
)

// RunCurve executes the curve subcommand: it samples the system
// reliability on a regular time grid and prints one "t R" pair per line,
// which is convenient for plotting or for numeric cross-checks.
func RunCurve(args []string, stdout, stderr io.Writer) int {
	from, to, n, file, err := parseCurveArgs(args)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n\n%s\n", err, usage)
		return 2
	}
	var root *model.Node
	if file != "" {
		root, err = model.ParseFile(file)
	} else {
		root, err = model.Parse(os.Stdin)
	}
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	if err := model.Validate(root); err != nil {
		fmt.Fprintf(stderr, "error: invalid diagram: %v\n", err)
		return 1
	}
	pts, err := reli.Curve(root, from, to, n)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	for _, p := range pts {
		fmt.Fprintf(stdout, "%.6g %.10f\n", p.T, p.R)
	}
	return 0
}

func parseCurveArgs(args []string) (from, to float64, n int, file string, err error) {
	from, to, n = 0, 0, 0
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-h" || arg == "--help":
			return 0, 0, 0, "", errHelpRequested
		case arg == "--from":
			if i+1 >= len(args) {
				return 0, 0, 0, "", fmt.Errorf("flag --from needs a value")
			}
			i++
			from, err = strconv.ParseFloat(args[i], 64)
		case arg == "--to":
			if i+1 >= len(args) {
				return 0, 0, 0, "", fmt.Errorf("flag --to needs a value")
			}
			i++
			to, err = strconv.ParseFloat(args[i], 64)
		case arg == "--n":
			if i+1 >= len(args) {
				return 0, 0, 0, "", fmt.Errorf("flag --n needs a value")
			}
			i++
			var nv int64
			nv, err = strconv.ParseInt(args[i], 10, 32)
			n = int(nv)
		default:
			if file == "" {
				file = arg
			} else {
				return 0, 0, 0, "", fmt.Errorf("unexpected argument %q", arg)
			}
		}
		if err != nil {
			return 0, 0, 0, "", fmt.Errorf("invalid numeric flag value: %v", err)
		}
	}
	if to <= 0 {
		return 0, 0, 0, "", fmt.Errorf("curve horizon --to must be positive")
	}
	if n < 2 {
		return 0, 0, 0, "", fmt.Errorf("curve sample count --n must be at least 2")
	}
	return from, to, n, file, nil
}
