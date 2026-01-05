package main

import (
	"errors"
	"fmt"
	"os"
)

func (repl *replConfig) commandCd(args []string) error {
	if len(args) < 2 {
		return nil
	}
	targetDir := args[1]
	fileInfo, err := os.Stat(targetDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("cd: %s: No such file or directory", targetDir)
		}
		return fmt.Errorf("cd: %w", err)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("Error: directory not found")
	}

	err = os.Chdir(targetDir)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			return fmt.Errorf("cd: %s: Permission denied", targetDir)
		}
		return fmt.Errorf("cd: %w", err)
	}

	repl.env["PWD"] = targetDir

	return nil
}
