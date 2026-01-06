package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

	scanner := bufio.NewScanner(os.Stdin)
	for {
		repl.stderr = os.Stderr
		repl.stdout = os.Stdout
		fmt.Print("$ ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		commandName := fields[0]
		commandArgs, err := repl.argParser(line)
		if err != nil {
			fmt.Fprintln(repl.stderr, err)
		}

		if cmd, ok := repl.builtins[commandName]; ok {
			err := cmd.callback(commandArgs)
			// fmt.Fprintln(repl.stdout, "TEST")
			// fmt.Fprintln(repl.stdout, repl.env["TEST"])
			if err != nil {
				fmt.Fprintln(repl.stderr, err)
			}
		} else {
			_, err := repl.findExecutablePath(commandName)
			if err != nil {
				fmt.Fprintf(repl.stderr, "%s: command %s\n", commandName, err)
				continue
			}
			err = repl.executeExternal(commandName, commandArgs[1:])
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(repl.stderr, "Error reading input:", err)
		os.Exit(1)
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
