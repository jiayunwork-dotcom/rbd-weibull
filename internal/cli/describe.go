package cli

import (
	"fmt"
	"io"
	"os"

	"rbd-weibull/internal/model"
)

// RunDescribe executes the describe subcommand: it parses and validates a
// diagram, then prints the tree with every node's kind and parameters.
func RunDescribe(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintf(stderr, "error: describe needs a diagram file\n\n%s\n", usage)
		return 2
	}
	var root *model.Node
	var err error
	if args[0] == "-" {
		root, err = model.Parse(os.Stdin)
	} else {
		root, err = model.ParseFile(args[0])
	}
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	if err := model.Validate(root); err != nil {
		fmt.Fprintf(stderr, "error: invalid diagram: %v\n", err)
		return 1
	}
	fmt.Fprint(stdout, model.TextTree(root))
	return 0
}
