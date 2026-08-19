package reli

import "rbd-weibull/internal/model"

// Result bundles everything the eval command prints for one diagram at one
// mission time: the system reliability, its hazard approximation, the
// Birnbaum importances of the leaves and the numerical system MTTF with the
// grid used to compute it.
type Result struct {
	MissionTime float64
	Reliability float64
	Hazard      float64
	MTTF        float64
	MTTFGrid    int
	MTTFUpper   float64
	Importances []Importance
	UnitCount   int
}

// Evaluate runs the full pipeline on a validated diagram: system
// reliability, system hazard, Birnbaum importances and numerical MTTF all
// at the same mission time. The returned Result is ready for formatting.
func Evaluate(root *model.Node, t float64, grid int) (Result, error) {
	rsys, err := System(root, t)
	if err != nil {
		return Result{}, err
	}
	haz, err := SystemHazard(root, t)
	if err != nil {
		return Result{}, err
	}
	imps, err := Birnbaum(root, t)
	if err != nil {
		return Result{}, err
	}
	mttf, err := SystemMTTF(root, 0, grid)
	if err != nil {
		return Result{}, err
	}
	return Result{
		MissionTime: t,
		Reliability: rsys,
		Hazard:      haz,
		MTTF:        mttf,
		MTTFGrid:    gridIfPositive(grid),
		MTTFUpper:   MTTFHorizon(root),
		Importances: imps,
		UnitCount:   len(imps),
	}, nil
}

func gridIfPositive(grid int) int {
	if grid <= 0 {
		return DefaultGridSteps
	}
	return grid
}
