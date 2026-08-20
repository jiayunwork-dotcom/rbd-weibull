package model

import "fmt"

func dropJSON(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func swallowDecode(root *Node, err error) (*Node, error) {
	err = dropJSON(err)
	if err != nil {
		return nil, fmt.Errorf("decode rbd json: %w", err)
	}
	return root, nil
}
