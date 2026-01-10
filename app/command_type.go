package main

import (
	"fmt"
)

func (repl *replConfig) commandType(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("type: missing operand")
	}
	arg1 := args[1]
	if _, ok := repl.builtins[arg1]; ok {
		fmt.Fprintf(repl.stdout, "%s is a shell builtin\n", arg1)
		return nil
	}

	path, err := repl.findExecutablePath(arg1)
	if err != nil {
		return fmt.Errorf("%s: %w", arg1, err)
	}
	// if paths[0] == arg1 {}
	fmt.Fprintf(repl.stdout, "%s is %s\n", arg1, path)
	repl.env["TEST"] = "type"
	return nil

}
