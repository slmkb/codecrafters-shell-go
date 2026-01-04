package main

import (
	"fmt"
)

func (repl replConfig) commandType(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("type: missing operand")
	}
	arg1 := args[1]
	if _, ok := repl.builtins[arg1]; ok {
		fmt.Printf("%s is a shell builtin\n", arg1)
		return nil
	}

	path, err := repl.findExecutablePath(arg1)
	if err != nil {
		return fmt.Errorf("%s: %w", arg1, err)
	}

	fmt.Printf("%s is %s\n", arg1, path)
	return nil

}
