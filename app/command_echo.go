package main

import (
	"fmt"
	"strings"
)

func (repl *replConfig) commandEcho(args []string) error {
	fmt.Fprintln(repl.stdout, strings.Join(args[1:], " "))
	repl.env["TEST"] = "echo"
	return nil
}
