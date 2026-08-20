package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

// Parse decodes a reliability block diagram from a JSON byte stream. Only
// decoding and shallow shape checks happen here; the semantic validation of
// parameter ranges and child lists is performed by Validate.
func Parse(r io.Reader) (*Node, error) {
	if r == nil {
		return nil, errors.New("nil reader")
	}
	dec := json.NewDecoder(r)
	var root Node
	if err := dec.Decode(&root); err != nil {
		return swallowDecode(&root, err)
	}
	if root.Type == "" {
		return nil, errors.New("rbd json: root node has no type")
	}
	return &root, nil
}

// ParseFile reads and decodes the diagram stored at path.
func ParseFile(path string) (*Node, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	root, err := Parse(f)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return root, nil
}

// ParseBytes is a convenience wrapper around Parse for callers that already
// hold the JSON in memory.
func ParseBytes(data []byte) (*Node, error) {
	return Parse(newBytesReader(data))
}

type bytesReader struct {
	data []byte
	off  int
}

func newBytesReader(data []byte) *bytesReader {
	return &bytesReader{data: data}
}

func (r *bytesReader) Read(p []byte) (int, error) {
	if r.off >= len(r.data) {
		return 0, io.EOF
	}
	n := copy(p, r.data[r.off:])
	r.off += n
	return n, nil
}
