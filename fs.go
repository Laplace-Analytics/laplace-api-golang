package laplace

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func findModuleRoot(t testing.TB) (string, error) {
	t.Helper()
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
