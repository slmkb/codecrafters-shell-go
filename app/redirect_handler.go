package main

import (
	"os"
)

type redirectTarget struct {
	target   string
	options  int
	targetFD byte
}

func (repl *replConfig) redirectHandler(rt redirectTarget) error {
	f, err := os.OpenFile(rt.target, os.O_CREATE|os.O_WRONLY|rt.options, 0644)
	if err != nil {
		return err
	}
	// fmt.Println(streamNumber, target)
	switch rt.targetFD {
	case 1:
		repl.stdout = f
	case 2:
		repl.stderr = f
	}
	return nil
}
