package main

import (
	"errors"
	"fmt"
	"os"
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

func run() {

	repl := newRepl()

	/* scanner ==========================
	scanner := bufio.NewScanner(os.Stdin)
	for {
		repl.stderr = os.Stderr
		repl.stdout = os.Stdout
		fmt.Print("$ ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
	====================================*/
	oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	b := make([]byte, 1)
	var line strings.Builder

	for {
		// oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
		fmt.Print("\r$ ", line.String())
		os.Stdin.Read(b)
		switch b[0] {
		case '\t':
			// fmt.Println("TAB")
			newLine := repl.autoComplete(line.String())
			line.Reset()
			line.WriteString(newLine)
			continue
		case '\r':
			fmt.Print("\r")
			term.Restore(int(os.Stdin.Fd()), oldState)
			fmt.Println()
			// fmt.Print("\r\x1b[K")
		default:
			line.WriteByte(b[0])
			continue
		}

		commandLine := strings.TrimSpace(line.String())
		line.Reset()
		if commandLine == "" {
			oldState, _ = term.MakeRaw(int(os.Stdin.Fd()))
			continue
		}

		fields := strings.Fields(commandLine)
		commandName := fields[0]
		commandArgs, err := repl.argParser(commandLine)
		// fmt.Println(commandLine)
		// fmt.Println(fields)
		// fmt.Println(commandName)
		// fmt.Println(commandArgs)
		if err != nil {
			fmt.Fprintln(repl.stderr, err)
			oldState, _ = term.MakeRaw(int(os.Stdin.Fd()))
			continue
		}

		if cmd, ok := repl.builtins[commandName]; ok {
			err := cmd.callback(commandArgs)
			if err != nil {
				if errors.Is(err, errExit) {
					return
				}
				fmt.Fprintln(repl.stderr, err)
			}
		} else {
			_, err := repl.findExecutablePath(commandName)
			if err != nil {
				fmt.Fprintf(repl.stderr, "%s: command %s\n", commandName, err)
				oldState, _ = term.MakeRaw(int(os.Stdin.Fd()))
				continue
			}
			err = repl.executeExternal(commandName, commandArgs[1:])
			if err != nil {
				fmt.Println(err)
			}
		}
		oldState, _ = term.MakeRaw(int(os.Stdin.Fd()))
		// defer term.Restore(int(os.Stdin.Fd()), oldState)
	}

	/* scanner
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(repl.stderr, "Error reading input:", err)
		os.Exit(1)
	}
	*/
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
