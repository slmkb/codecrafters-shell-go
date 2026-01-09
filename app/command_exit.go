package main

import (
	"errors"
)

var errExit = errors.New("Exit")

func (repl replConfig) commandExit([]string) error {
	return errExit
}
