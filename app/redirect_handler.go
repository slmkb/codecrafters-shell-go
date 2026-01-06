package main

import (
	"os"
)

func (repl *replConfig) redirectHandler(streamNumber byte, target string) error {
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	// fmt.Println(streamNumber, target)
	switch streamNumber {
	case '1':
		repl.stdout = f
	case '2':
		repl.stderr = f
	}
	return nil
}
