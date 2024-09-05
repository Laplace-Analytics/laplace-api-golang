package utilities

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindModuleRoot finds the relative path to the root of the Go module.
func FindModuleRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			// go.mod found
			return currentDir, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// Reached the root directory
			return "", fmt.Errorf("go.mod file not found")
		}

		currentDir = parentDir
	}
}
