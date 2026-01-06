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

func (repl *replConfig) argParser(argLine string) ([]string, error) {
	var args []string
	var cur strings.Builder

	inSingle := false
	inDouble := false

	var redirFD byte
	var outTarget redirectTarget
	var errTarget redirectTarget
	expectTarget := false
	appendMode := false

	flushToken := func() {
		if cur.Len() == 0 {
			return
		}
		tok := cur.String()
		cur.Reset()

		if expectTarget {
			var redirTarget redirectTarget
			redirTarget.target = tok
			if appendMode {
				redirTarget.options = os.O_APPEND
			} else {
				redirTarget.options = os.O_TRUNC
			}
			// assign redirect target
			switch redirFD {
			case 1:
				redirTarget.targetFD = 1
				outTarget = redirTarget
			case 2:
				redirTarget.targetFD = 2
				errTarget = redirTarget
			}
			expectTarget = false
			appendMode = false
			redirFD = 0
			return
		}

		args = append(args, tok)
	}

	for i := 0; i < len(argLine); i++ {
		ch := argLine[i]

		switch ch {
		case '\'':
			if inDouble {
				cur.WriteByte(ch)
			} else {
				inSingle = !inSingle
			}

		case '"':
			if inSingle {
				cur.WriteByte(ch)
			} else {
				inDouble = !inDouble
			}

		case ' ':
			if inSingle || inDouble {
				cur.WriteByte(ch)
			} else {
				flushToken()
			}

		case '>':
			if inSingle || inDouble {
				cur.WriteByte(ch)
				continue
			}

			if i+1 < len(argLine) && argLine[i+1] == '>' {
				appendMode = true
				i++ // consume the second '>'
			}

			tok := cur.String()
			cur.Reset()

			switch tok {
			case "2":
				redirFD = 2
			case "1", "":
				redirFD = 1
			default:
				redirFD = 1
				args = append(args, tok)
			}
			expectTarget = true

		default:
			cur.WriteByte(ch)
		}
	}

	if inSingle || inDouble {
		return nil, errors.New("error: unterminated quote")
	}

	flushToken()

	if expectTarget {
		return nil, errors.New("error: missing redirect target")
	}

	if outTarget.target != "" {
		err := repl.redirectHandler(outTarget)
		if err != nil {
			return nil, err
		}
	}

	if errTarget.target != "" {
		err := repl.redirectHandler(errTarget)
		if err != nil {
			return nil, err
		}
	}

	// fmt.Println(redirectStdout, redirectStderr)
	// fmt.Printf("argline: %q\n", argLine)
	// fmt.Printf("parsed: %#v\n", args)
	// fmt.Printf("out: %q\n", outTarget)
	// fmt.Printf("err: %q\n", errTarget)
	return args, nil
}
