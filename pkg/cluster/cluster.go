package cluster

import (
	"errors"
)

func AddNode(node string) error {
	if node == "" {
		return errors.New("node name cannot be empty")
	}
	return nil
}

func RemoveNode(node string) error {
	if node == "" {
		return errors.New("node name cannot be empty")
	}
	return nil
}

func ListNodes() ([]string, error) {
	return []string{"node1", "node2", "node3"}, nil
}
