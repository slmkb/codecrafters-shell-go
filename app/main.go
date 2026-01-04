package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	repl()
}

type cliCommand struct {
	name     string
	callback func([]string) error
}

func repl() {

	commands := make(map[string]cliCommand)

	commands["exit"] = cliCommand{
		name:     "exit",
		callback: commandExit,
	}

	commands["echo"] = cliCommand{
		name:     "echo",
		callback: commandEcho,
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("$ ")
		if !scanner.Scan() {
			break
		}
		commandLine := scanner.Text()
		commandName := strings.Fields(commandLine)[0]
		commandArgs := strings.Fields(commandLine)[1:]
		if v, ok := commands[commandName]; ok {
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

func commandExit([]string) error {
	os.Exit(0)
	return nil
}

func commandEcho(args []string) error {
	fmt.Println(strings.Join(args, " "))
	return nil
}
