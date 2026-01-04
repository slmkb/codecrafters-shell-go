package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func (repl replConfig) findExecutablePath(target string) (string, error) {
	for _, path := range strings.Split(repl.env["PATH"], ":") {
		filePath := filepath.Clean(filepath.Join(path, target))
		if fi, err := os.Stat(filePath); err == nil {
			isExecutable := fi.Mode().Perm()&0111 != 0
			if isExecutable {
				// fmt.Printf("%s is %s\n", arg1, filePath)
				return filePath, nil
			}
		}
	}
	return "", errors.New("not found")
}
