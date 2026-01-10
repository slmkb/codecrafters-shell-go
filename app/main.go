package main

import (
	"errors"
	"log"
)

func main() {
	err := run()
	if err != nil {
		if errors.Is(err, errExit) {
			return
		}
		log.Fatal(err)
	}
}
