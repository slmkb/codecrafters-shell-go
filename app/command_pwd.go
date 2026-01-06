package main

import (
	"fmt"
	"os"
)

func (repl replConfig) commandPwd(args []string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Fprintln(repl.stdout, pwd)
	return nil
}
