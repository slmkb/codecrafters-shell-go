package main

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"golang.org/x/term"
)

type cliCommand struct {
	name     string
	callback func([]string) error
}

type replConfig struct {
	builtins map[string]cliCommand
	env      map[string]string
	stdin    *os.File
	stdout   *os.File
	stderr   *os.File
}

func run() error {

	repl := newRepl()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	b := make([]byte, 1)
	var line strings.Builder
	pressedTab := false

	for {
		fmt.Print("\r\x1b[K$ " + line.String())
		n, err := os.Stdin.Read(b)
		if err != nil {
			return err
		}
		if n != 1 {
			return errors.New("received EOF")
		}
		switch b[0] {
		case '\t':
			var token string
			fields := strings.Fields(line.String())
			lastIndex := len(fields) - 1
			if lastIndex >= 0 {
				token = fields[lastIndex]
			}
			candidates := repl.autoComplete(token)

			// fmt.Printf("\n%#v\n", candidates)
			switch len(candidates) {
			case 0:
				fmt.Print("\x07")
			case 1:
				line.Reset()
				var command string
				for k := range candidates {
					command = k
				}
				newLine := strings.Join(fields[:lastIndex], " ") + command + " "
				line.WriteString(newLine)
			default:
				fmt.Print("\x07")
				if pressedTab {
					fmt.Print("\n\r\x1b[K")
					sortedCandidates := slices.Sorted(maps.Keys(candidates))
					fmt.Println(strings.Join(sortedCandidates, "  "))
				}
				pressedTab = false
			}
			pressedTab = true
			continue
		case '\r', '\n':
			term.Restore(int(os.Stdin.Fd()), oldState)
			repl.stderr = os.Stderr
			repl.stdout = os.Stdout
			fmt.Println()
			pressedTab = false
		default:
			line.WriteByte(b[0])
			pressedTab = false
			continue
		}

		commandLine := strings.TrimSpace(line.String())
		line.Reset()
		if commandLine == "" {
			if _, err := term.MakeRaw(int(os.Stdin.Fd())); err != nil {
				return err
			}
			continue
		}

		commandArgs, err := repl.argParser(commandLine)
		commandName := commandArgs[0]
		if err != nil {
			fmt.Fprintln(repl.stderr, err)
			if _, err := term.MakeRaw(int(os.Stdin.Fd())); err != nil {
				return err
			}
			continue
		}

		if cmd, ok := repl.builtins[commandName]; ok {
			err := cmd.callback(commandArgs)
			if err != nil {
				if errors.Is(err, errExit) {
					return err
				}
				fmt.Fprintln(repl.stderr, err)
			}
		} else {
			_, err := repl.findExecutablePath(commandName)
			if err != nil {
				fmt.Fprintf(repl.stderr, "%s: command %s\n", commandName, err)
				if _, err := term.MakeRaw(int(os.Stdin.Fd())); err != nil {
					return err
				}
				continue
			}
			if err := repl.executeExternal(commandName, commandArgs[1:]); err != nil {
				// return err
			}
		}
		if _, err := term.MakeRaw(int(os.Stdin.Fd())); err != nil {
			return err
		}
	}

}

func newRepl() *replConfig {
	var repl replConfig

	repl.stdin = os.Stdin
	repl.stdout = os.Stdout
	repl.stderr = os.Stderr

	repl.env = make(map[string]string)
	for _, token := range os.Environ() {
		parts := strings.SplitN(token, "=", 2)
		if len(parts) < 2 {
			continue
		}
		k := parts[0]
		v := parts[1]
		repl.env[k] = v
	}

	repl.builtins = make(map[string]cliCommand)
	repl.builtins["exit"] = cliCommand{
		name:     "exit",
		callback: repl.commandExit,
	}

	repl.builtins["echo"] = cliCommand{
		name:     "echo",
		callback: repl.commandEcho,
	}

	repl.builtins["type"] = cliCommand{
		name:     "type",
		callback: repl.commandType,
	}

	repl.builtins["pwd"] = cliCommand{
		name:     "pwd",
		callback: repl.commandPwd,
	}

	repl.builtins["cd"] = cliCommand{
		name:     "cd",
		callback: repl.commandCd,
	}

	return &repl
}
