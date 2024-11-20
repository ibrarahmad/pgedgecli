package spock

import "errors"

func CreateSet(set string) error {
	if set == "" {
		return errors.New("replication set name cannot be empty")
	}
	// Simulate creating a replication set
	return nil
}

func DropSet(set string) error {
	if set == "" {
		return errors.New("replication set name cannot be empty")
	}
	// Simulate dropping a replication set
	return nil
}

func CheckSet(set string) (string, error) {
	if set == "" {
		return "", errors.New("replication set name cannot be empty")
	}
	// Simulate checking replication set
	return "active", nil
}

