package model

import "fmt"

// Sentinel error classes shared by the whole package. Wrap them with
// additional context (node type, path, offending value) before returning so
// tests can match with errors.Is while messages stay informative.
var (
	// ErrUnknownType is reported when a node's type string is not one of
	// unit / series / parallel / kofn.
	ErrUnknownType = fmt.Errorf("unknown node type")

	// ErrEmptySeries is reported when a series or parallel node has an
	// empty block list.
	ErrEmptySeries = fmt.Errorf("structural node without child blocks")

	// ErrKGreaterN is reported when a kofn node has k > n.
	ErrKGreaterN = fmt.Errorf("k greater than number of blocks")

	// ErrKBelowOne is reported when a kofn node has k < 1.
	ErrKBelowOne = fmt.Errorf("k below one")

	// ErrNonPositive is reported when a Weibull parameter or mission time
	// is not strictly positive.
	ErrNonPositive = fmt.Errorf("parameter must be positive")

	// ErrTooManySubblocks is reported when an enumerate-only kofn node
	// exceeds the state enumeration limit.
	ErrTooManySubblocks = fmt.Errorf("too many sub-blocks to enumerate")
)

// PathError attaches the structural path (blocks[i].blocks[j]...) of the
// offending node to a sentinel error.
type PathError struct {
	Path string
	Err  error
}

func (e *PathError) Error() string {
	if e.Path == "" {
		return e.Err.Error()
	}
	return fmt.Sprintf("%s: %v", e.Path, e.Err)
}

func (e *PathError) Unwrap() error {
	return e.Err
}

// At wraps err with a node path description.
func At(path string, err error) error {
	if err == nil {
		return nil
	}
	return &PathError{Path: path, Err: err}
}
