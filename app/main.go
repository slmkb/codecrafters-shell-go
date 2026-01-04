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
}

func run() {

	repl := newRepl()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("$ ")
		if !scanner.Scan() {
			break
		}
		commandLine := scanner.Text()
		commandName := strings.Fields(commandLine)[0]
		commandArgs := strings.Fields(commandLine)
		if v, ok := repl.builtins[commandName]; ok {
			err := v.callback(commandArgs)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("%s: command not found\n", commandName)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
		os.Exit(1)
	}

}

func (c replConfig) commandExit([]string) error {
	os.Exit(0)
	return nil
}

func (c replConfig) commandEcho(args []string) error {
	fmt.Println(strings.Join(args[1:], " ")) //arg[0] == commandName
	return nil
}

func (c replConfig) commandType(args []string) error {
	arg1 := args[1]
	if _, ok := c.builtins[arg1]; !ok {
		return fmt.Errorf("%s: not found", arg1)
	}
	fmt.Printf("%s is a shell builtin\n", arg1)
	return nil
}

func newRepl() *replConfig {
	var repl replConfig
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

	return &repl
}

func main() {
	run()
}
