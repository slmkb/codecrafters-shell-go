package main

import (
	"errors"
	"fmt"
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

func (repl *replConfig) argParser(argLine string) ([]string, error) {
	var parsedArgs []string
	var builder strings.Builder

	inSingle := false
	inDouble := false
	redirectStdout := false
	var targetStdout string
	redirectStderr := false
	var targetStderr string
	var lastChar byte

	for i := 0; i < len(argLine); i++ {
		ch := argLine[i]

		switch ch {
		case '\'':
			if inDouble {
				builder.WriteByte(ch)
			} else {
				inSingle = !inSingle
			}
		case '"':
			if inSingle {
				builder.WriteByte(ch)
			} else {
				inDouble = !inDouble
			}
		case ' ':
			if inSingle || inDouble {
				builder.WriteByte(ch)
			} else if builder.Len() == 0 {
				continue
			} else if redirectStdout && len(targetStdout) == 0 {
				targetStdout = builder.String()
				builder.Reset()
			} else if redirectStderr && len(targetStderr) == 0 {
				targetStderr = builder.String()
				builder.Reset()
			} else {
				parsedArgs = append(parsedArgs, builder.String())
				builder.Reset()
			}
		case '>':
			if inSingle || inDouble {
				builder.WriteByte(ch)
			} else {
				if builder.Len() != 0 && lastChar == '2' {
					redirectStderr = true
				} else {
					redirectStdout = true
				}
				builder.Reset()
			}
		default:
			builder.WriteByte(ch)
		}
		lastChar = ch
	}

	if inSingle || inDouble {
		return nil, errors.New("error: unterminated quote")
	}

	if builder.Len() > 0 {
		if redirectStdout {
			if len(targetStdout) == 0 {
				targetStdout = builder.String()
				builder.Reset()
			}
			repl.redirectHandler('1', targetStdout)
		} else if redirectStderr {
			if len(targetStderr) == 0 {
				targetStdout = builder.String()
				builder.Reset()
			}
			repl.redirectHandler('2', targetStderr)
		}
		if builder.Len() > 0 {
			parsedArgs = append(parsedArgs, builder.String())
		}
	}

	fmt.Printf("argline: %q\n", argLine)
	fmt.Printf("parsed: %#v\n", parsedArgs)
	fmt.Printf("out: %q\n", targetStdout)
	fmt.Printf("err: %q\n", targetStderr)
	return parsedArgs, nil
}
