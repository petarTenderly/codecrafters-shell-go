package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

const (
	exitCmd = "exit"
	echoCmd = "echo"
	typeCmd = "type"
	pwdCmd  = "pwd"
	cdCmd   = "cd"
)

type Command struct {
	Name            string
	Args            []string
	Output          *os.File
	ErrorOutput     *os.File
	HandlerRegistry map[string]cmdHandler
	Shell           *Shell
}

func NewCommand(parts []string, shell *Shell) Command {
	const exeMod = 0666

	output := os.Stdout
	errorOutput := os.Stderr
	arguments := parts[1:]

	stdOutCmds := []string{">", "1>", ">>", "1>>"}
	stdErrCmds := []string{"2>", "2>>"}

	if len(parts) > 2 {
		redirectPosition := len(arguments) - 2
		if slices.Contains(stdOutCmds, arguments[redirectPosition]) || slices.Contains(stdErrCmds, arguments[redirectPosition]) {
			fileName := arguments[len(arguments)-1]
			openFlag := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
			if arguments[redirectPosition] == ">>" || arguments[redirectPosition] == "1>>" || arguments[redirectPosition] == "2>>" {
				openFlag = os.O_APPEND | os.O_WRONLY | os.O_CREATE
			}
			outputFile, _ := os.OpenFile(fileName, openFlag, exeMod)
			if arguments[redirectPosition] == "2>" || arguments[redirectPosition] == "2>>" {
				errorOutput = outputFile
			} else {
				output = outputFile
			}

			arguments = arguments[:redirectPosition]
		}

	}

	handleRegistry := map[string]cmdHandler{
		exitCmd: handleExit,
		echoCmd: handleEcho,
		typeCmd: handleType,
		pwdCmd:  handlePwd,
		cdCmd:   handleCd,
	}

	return Command{
		Name:            parts[0],
		Args:            arguments,
		Output:          output,
		ErrorOutput:     errorOutput,
		HandlerRegistry: handleRegistry,
		Shell:           shell,
	}
}

type cmdHandler func(command Command)

func (command Command) exec() {
	if handler, ok := command.HandlerRegistry[command.Name]; ok {
		handler(command)
	} else {
		handleExec(command)
	}
}

func handleExec(command Command) {
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

func handleCd(command Command) {
	dir := command.Args[0]
	if command.Args[0] == "~" {
		dir = os.Getenv("HOME")
	}
	err := os.Chdir(dir)
	if err != nil {
		fmt.Fprintf(command.Output, fmt.Sprintf("cd: %s: No such file or directory\n", command.Args[0]))
	}
}

func handlePwd(command Command) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(command.Output, "error reading working directory\n")
	}
	fmt.Println(dir)
}

func handleType(command Command) {
	if slices.Contains(command.Shell.builtIns, command.Args[0]) {
		fmt.Fprintf(command.Output, fmt.Sprintf("%s is a shell builtin\n", command.Args[0]))
	} else {
		if path, err := exec.LookPath(command.Args[0]); err == nil {
			fmt.Fprintf(command.Output, fmt.Sprintf("%s is %s\n", command.Args[0], path))
		} else {
			fmt.Fprintf(command.Output, fmt.Sprintf("%s: not found\n", command.Args[0]))
		}
	}
}

func handleEcho(command Command) {
	fmt.Fprintf(command.Output, fmt.Sprintln(strings.Join(command.Args, " ")))
}

func handleExit(command Command) {
	if command.Args[0] == "0" {
		os.Exit(0)
	}
	fmt.Fprintf(command.Output, "exit: status code must be 0\n")
}
