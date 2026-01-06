package main

import (
	"fmt"
	"os/exec"
)

func (repl replConfig) executeExternal(command string, args []string) error {
	cmd := exec.Command(command, args...)
	cmd.Stderr = repl.stderr
	cmd.Stdout = repl.stdout
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("execute external: %w", err)
	}
	return nil
}
