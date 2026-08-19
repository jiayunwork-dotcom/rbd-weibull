package reli

import (
	"fmt"

	"rbd-weibull/internal/model"
)

// LeafDetail is one row of the per-component diagnostic table: the
// component's own reliability and hazard at the mission time, its closed
// form MTTF and its Birnbaum importance in the system.
type LeafDetail struct {
	Name       string
	Path       string
	R          float64
	Hazard     float64
	MTTF       float64
	Importance float64
}

// LeafDetails builds the diagnostic table for every unit leaf at time t.
// The importance values come from Birnbaum, which already re-evaluates the
// system; the remaining fields are computed from the leaf's own Weibull
// parameters so the table shows what drives each importance.
func LeafDetails(root *model.Node, t float64) ([]LeafDetail, error) {
	imps, err := Birnbaum(root, t)
	if err != nil {
		return nil, err
	}
	out := make([]LeafDetail, 0, len(imps))
	for _, imp := range imps {
		leaf := imp.Leaf.Node
		if !leaf.IsUnit() || leaf.Beta == nil || leaf.Eta == nil {
			return nil, fmt.Errorf("leaf %s is not a parameterised unit", imp.Path)
		}
		out = append(out, LeafDetail{
			Name:       imp.Name,
			Path:       imp.Path,
			R:          Reliability(t, *leaf.Beta, *leaf.Eta),
			Hazard:     Hazard(t, *leaf.Beta, *leaf.Eta),
			MTTF:       UnitMTTF(*leaf.Beta, *leaf.Eta),
			Importance: imp.Value,
		})
	}
	return out, nil
}

// RankByImportance returns the leaf names ordered from most to least
// important at time t.
func RankByImportance(root *model.Node, t float64) ([]string, error) {
	imps, err := Birnbaum(root, t)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(imps))
	for _, imp := range imps {
		names = append(names, imp.Name)
	}
	// Insertion sort by descending importance; n is small in practice.
	for i := 1; i < len(names); i++ {
		for j := i; j > 0 && imps[j].Value > imps[j-1].Value; j-- {
			names[j], names[j-1] = names[j-1], names[j]
			imps[j], imps[j-1] = imps[j-1], imps[j]
		}
	}
	return names, nil
}

// FractionalImportance normalises the Birnbaum values so they sum to one;
// it is a rough guide to where redesign effort pays off.
func FractionalImportance(root *model.Node, t float64) ([]Importance, error) {
	imps, err := Birnbaum(root, t)
	if err != nil {
		return nil, err
	}
	total := 0.0
	for _, imp := range imps {
		total += imp.Value
	}
	if total <= 0 {
		return imps, nil
	}
	out := make([]Importance, len(imps))
	for i, imp := range imps {
		out[i] = imp
		out[i].Value = imp.Value / total
	}
	return out, nil
}
