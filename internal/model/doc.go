// Package model defines the reliability block diagram tree: unit, series,
// parallel and k-of-n nodes, their JSON decoding and the structural checks
// that must pass before any reliability number is computed.
//
// A diagram is a rooted tree. Inner nodes (series, parallel, kofn) hold an
// ordered list of child blocks; leaf nodes (unit) carry the two Weibull
// parameters beta (shape) and eta (scale). The tree is built from JSON like:
//
//	{ "type": "series", "blocks": [
//	    { "type": "unit", "beta": 1.8, "eta": 4500 },
//	    { "type": "parallel", "blocks": [ ... ] }
//	] }
//
// Unknown node types, empty child lists, k > n on k-of-n nodes and non
// positive Weibull parameters are rejected by Validate with a typed error so
// callers can distinguish the failure class without string matching.
package model
