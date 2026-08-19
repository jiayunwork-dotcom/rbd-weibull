package cli

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"rbd-weibull/internal/model"
	"rbd-weibull/internal/reli"
)

// options collects the flags understood by the eval command.
type options struct {
	time float64
	file string
}

// RunEval executes the eval subcommand. Flags may appear in any position
// relative to the file argument, which is what makes
//
//	rbd-weibull eval example/pump2.json --t 1000
//
// work without a flag-parser preamble.
func RunEval(args []string, stdout, stderr io.Writer) int {
	opt, err := parseEvalArgs(args)
	if err == errHelpRequested {
		fmt.Fprint(stdout, usage)
		return 0
	}
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n\n%s\n", err, usage)
		return 2
	}
	var root *model.Node
	if opt.file != "" {
		root, err = model.ParseFile(opt.file)
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
	res, err := reli.Evaluate(root, opt.time, 0)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	WriteResult(stdout, res, opt.file)
	return 0
}

func parseEvalArgs(args []string) (options, error) {
	opt := options{time: 1000}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-h" || arg == "--help":
			return opt, errHelpRequested
		case arg == "-t" || arg == "--t":
			if i+1 >= len(args) {
				return opt, fmt.Errorf("flag %s needs a value", arg)
			}
			i++
			v, err := strconv.ParseFloat(args[i], 64)
			if err != nil {
				return opt, fmt.Errorf("invalid time %q", args[i])
			}
			opt.time = v
		case len(arg) > 3 && (arg[:2] == "-t" || arg[:3] == "--t"):
			// -t=1000 or --t=1000
			body := arg
			eq := indexByte(body, '=')
			if eq < 0 {
				return opt, fmt.Errorf("unknown flag %q", arg)
			}
			v, err := strconv.ParseFloat(body[eq+1:], 64)
			if err != nil {
				return opt, fmt.Errorf("invalid time %q", body[eq+1:])
			}
			opt.time = v
		default:
			if opt.file == "" {
				opt.file = arg
			} else {
				return opt, fmt.Errorf("unexpected argument %q", arg)
			}
		}
	}
	if opt.time <= 0 {
		return opt, fmt.Errorf("mission time must be positive, got %g", opt.time)
	}
	return opt, nil
}

func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}

// errHelpRequested is a sentinel that is never used as a real error value
// beyond short-circuiting flag parsing; help is printed by the caller.
var errHelpRequested = fmt.Errorf("help requested")
