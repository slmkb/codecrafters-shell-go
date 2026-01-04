package main

import "os"

func (repl replConfig) commandExit([]string) error {
	os.Exit(0)
	return nil
}
