package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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
)

var allCmds = []string{exitCmd, echoCmd, typeCmd, pwdCmd}

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
			if slices.Contains(allCmds, args[0]) {
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
				fmt.Println("pwd: error getting current directory")
			} else {
				fmt.Println(dir)
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
	cmdParts := strings.Split(cmd, " ")
	command := cmdParts[0]
	args := cmdParts[1:]
	return command, args
}
