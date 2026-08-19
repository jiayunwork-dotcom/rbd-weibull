package cli

import (
	"fmt"
	"io"

	"rbd-weibull/internal/reli"
)

// WriteResult prints the evaluation result as a small text report: the
// mission time, system reliability, hazard approximation, numerical MTTF
// with its integration grid and the per-leaf Birnbaum importances.
func WriteResult(w io.Writer, res reli.Result, source string) {
	if source != "" {
		fmt.Fprintf(w, "diagram: %s\n", source)
	}
	fmt.Fprintf(w, "mission time t        %.6g\n", res.MissionTime)
	fmt.Fprintf(w, "system R(t)           %.10f\n", res.Reliability)
	fmt.Fprintf(w, "system hazard(t)      %.6e\n", res.Hazard)
	fmt.Fprintf(w, "system MTTF           %.6f   (trapezoid %d steps over [0, %.6g])\n",
		res.MTTF, res.MTTFGrid, res.MTTFUpper)
	fmt.Fprintf(w, "unit leaves           %d\n", res.UnitCount)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Birnbaum importance:")
	if len(res.Importances) == 0 {
		fmt.Fprintln(w, "  (no unit leaves)")
		return
	}
	for _, imp := range res.Importances {
		label := imp.Name
		if imp.Leaf.Node.Name != "" {
			label = imp.Leaf.Node.Name
		}
		fmt.Fprintf(w, "  %-16s %10.6f\n", label, imp.Value)
	}
}

// WriteCSV prints the same importances in a headerless CSV form, which is
// convenient for piping into other tools.
func WriteCSV(w io.Writer, res reli.Result) {
	for _, imp := range res.Importances {
		label := imp.Name
		if imp.Leaf.Node.Name != "" {
			label = imp.Leaf.Node.Name
		}
		fmt.Fprintf(w, "%s,%.10f\n", label, imp.Value)
	}
}
