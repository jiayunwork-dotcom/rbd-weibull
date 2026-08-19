package model

import "encoding/json"

// Marshal serialises the diagram back to compact JSON. The Fixed override
// is excluded because it carries analysis state, not diagram content.
func Marshal(root *Node) ([]byte, error) {
	return json.Marshal(root)
}

// MarshalIndent serialises the diagram with indentation, suitable for
// writing a normalised copy of an input file.
func MarshalIndent(root *Node) ([]byte, error) {
	return json.MarshalIndent(root, "", "  ")
}

// Normalise parses and re-serialises a diagram, which canonicalises the
// JSON (field order, spacing) while keeping the structure and parameters
// identical.
func Normalise(data []byte) ([]byte, error) {
	root, err := ParseBytes(data)
	if err != nil {
		return nil, err
	}
	if err := Validate(root); err != nil {
		return nil, err
	}
	return MarshalIndent(root)
}
