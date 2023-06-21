package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// checks if ref is a tag
func IsRefTag(repoDir, ref string) bool {
	cmd := exec.Command("git", "show-ref", "--verify", "refs/tags/"+ref)
	cmd.Dir = repoDir
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// Check if the exit status is 0 or 1
			if exitError.ExitCode() == 0 {
				return true
			}
		}
		return false
	}
	return true
}

// checks if ref is a branch
func IsRefBranch(repoDir, ref string) bool {
	cmd := exec.Command("git", "show-ref", "--verify", "refs/heads/"+ref)
	cmd.Dir = repoDir
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// Check if the exit status is 0
			if exitError.ExitCode() == 0 {
				return true
			}
		}
		return false
	}
	return true
}

// FindRootDir finds the root directory of the repository
func FindRootDir() (string, error) {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Traverse up the directory tree until a go.mod file is found
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd, nil
		}

		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}

	return "", fmt.Errorf("go.mod file not found")
}
