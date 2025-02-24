package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

const (
	exitCmd = "exit"
	echoCmd = "echo"
	typeCmd = "type"
	pwdCmd  = "pwd"
	cdCmd   = "cd"
)

var builtIns = []string{exitCmd, echoCmd, typeCmd, pwdCmd, cdCmd}

func main() {
	// Uncomment this block to pass the first stage
	for true {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		rawCmd, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		command, args := parseCmd(rawCmd)

		switch command {
		case exitCmd:
			if args[0] == "0" {
				os.Exit(0)
			}
			fmt.Println("exit: status code must be 0")
		case echoCmd:
			fmt.Println(strings.Join(args, " "))
		case typeCmd:
			if slices.Contains(builtIns, args[0]) {
				fmt.Printf("%s is a shell builtin\n", args[0])
			} else {
				if path, err := exec.LookPath(args[0]); err == nil {
					fmt.Printf("%s is %s\n", args[0], path)
				} else {
					fmt.Printf("%s: not found\n", args[0])
				}
			}
		case pwdCmd:
			dir, err := os.Getwd()
			if err != nil {
				fmt.Println("error reading working directory")
				return
			}
			fmt.Println(dir)
		case cdCmd:
			dir := args[0]
			if args[0] == "~" {
				dir = os.Getenv("HOME")
			}
			err := os.Chdir(dir)
			if err != nil {
				fmt.Printf("cd: %s: No such file or directory\n", args[0])
			}
		default:
			c := exec.Command(command, args...)
			c.Stderr = os.Stderr
			c.Stdout = os.Stdout
			err := c.Run()
			if err != nil {
				fmt.Printf("%s: command not found\n", command)
			}
		}

	}

}

func parseCmd(rawCmd string) (string, []string) {
	cmd := strings.TrimSpace(rawCmd)
	cmd = strings.ReplaceAll(cmd, "''", "")
	cmd = strings.ReplaceAll(cmd, "\"\"", "")
	//cmdParts := strings.Split(cmd, " ")
	//
	//command := cmdParts[0]
	//args := cmdParts[1:]
	//return command, args

	// Regular expression to capture the command and arguments
	re := regexp.MustCompile(`(\w+)(?:\s+((?:'[^']*'|"[^"]*"|\S+)(?:\s+(?:'[^']*'|"[^"]*"|\S+))*))?`) // Match the input string
	matches := re.FindStringSubmatch(cmd)

	// Extract the command and arguments
	command := matches[1]
	argumentsString := matches[2]

	// Split the arguments into a list
	arguments := parseArguments(argumentsString)

	return command, arguments
}

func parseArguments(argumentsString string) []string {
	// Use a regex to match either single-quoted strings or non-whitespace sequences
	re := regexp.MustCompile(`'[^']*'|"[^"]*"|\S+`)
	args := re.FindAllString(argumentsString, -1)
	realArgs := make([]string, len(args))
	for i, arg := range args {
		// Remove the single quotes from the matched strings
		realArgs[i] = strings.Trim(arg, "'")
		realArgs[i] = strings.Trim(realArgs[i], "\"")
	}
	return realArgs
}
