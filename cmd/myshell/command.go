package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

type Command struct {
	Name        string
	Args        []string
	Output      *os.File
	ErrorOutput *os.File
}

func NewCommand(parts []string) Command {
	output := os.Stdout
	errorOutput := os.Stderr
	arguments := parts[1:]
	if len(parts) > 2 {
		if arguments[len(arguments)-2] == ">" || arguments[len(arguments)-2] == "1>" {
			fileName := arguments[len(arguments)-1]
			var err error
			output, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
			if err != nil {
				fmt.Println("error creating file")
			}
			arguments = arguments[:len(arguments)-2]
		} else if arguments[len(arguments)-2] == "2>" {
			fileName := arguments[len(arguments)-1]
			var err error
			errorOutput, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
			if err != nil {
				fmt.Println("error creating file")
			}
			arguments = arguments[:len(arguments)-2]
		}
	}

	return Command{
		Name:        parts[0],
		Args:        arguments,
		Output:      output,
		ErrorOutput: errorOutput,
	}
}

func (command Command) exec() {
	switch command.Name {
	case exitCmd:
		if command.Args[0] == "0" {
			os.Exit(0)
		}
		fmt.Fprintf(command.Output, "exit: status code must be 0\n")
	case echoCmd:
		fmt.Fprintf(command.Output, fmt.Sprintln(strings.Join(command.Args, " ")))
	case typeCmd:
		if slices.Contains(builtIns, command.Args[0]) {
			fmt.Fprintf(command.Output, fmt.Sprintf("%s is a shell builtin\n", command.Args[0]))
		} else {
			if path, err := exec.LookPath(command.Args[0]); err == nil {
				fmt.Fprintf(command.Output, fmt.Sprintf("%s is %s\n", command.Args[0], path))
			} else {
				fmt.Fprintf(command.Output, fmt.Sprintf("%s: not found\n", command.Args[0]))
			}
		}
	case pwdCmd:
		dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(command.Output, "error reading working directory\n")
		}
		fmt.Println(dir)
	case cdCmd:
		dir := command.Args[0]
		if command.Args[0] == "~" {
			dir = os.Getenv("HOME")
		}
		err := os.Chdir(dir)
		if err != nil {
			fmt.Fprintf(command.Output, fmt.Sprintf("cd: %s: No such file or directory\n", command.Args[0]))
		}
	default:
		for i, arg := range command.Args {
			command.Args[i] = strings.ReplaceAll(arg, "\n", "\\n")
		}
		c := exec.Command(command.Name, command.Args...)
		c.Stdout = command.Output
		c.Stderr = command.ErrorOutput
		err := c.Run()
		if err != nil {
			if errors.Is(err, exec.ErrNotFound) {
				fmt.Fprintf(command.Output, "%s: command not found\n", command.Name)
			}
		}
	}

}
