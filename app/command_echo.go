package main

import (
	"fmt"
	"strings"
)

func (repl replConfig) commandEcho(args []string) error {
	fmt.Println(strings.Join(args, " "))
	return nil
}
