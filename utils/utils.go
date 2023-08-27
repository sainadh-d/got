package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func FindRepoRoot(path string) (string, error) {
	curr, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get current directory", err)
		return "", err
	}

	// Check if the given path contains a .git directory
	_, err = os.Stat(filepath.Join(curr, ".git"))
	if err == nil {
		return curr, nil
	}

	// If not found, check the parent directory
	parent := filepath.Dir(path)
	if parent == path {
		// Reached the root directory, .git directory not found
		return "", errors.New("no .git directory found")
	}

	return FindRepoRoot(parent)
}
