package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

const (
	exitCmd = "exit"
	echoCmd = "echo"
	typeCmd = "type"
)

var allCmds = []string{exitCmd, echoCmd, typeCmd}

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

				path := os.Getenv("PATH")
				executableDirs := strings.Split(path, ":")
				isFound := false

				for _, executable := range executableDirs {
					dirEntries, err := os.ReadDir(executable)
					if err != nil {
						continue
					}
					for _, dirEntry := range dirEntries {
						if dirEntry.Name() == args[0] {
							fmt.Printf("%s is %s/%s\n", args[0], executable, args[0])
							isFound = true
							return
						}
					}
					if isFound {
						break
					}
				}
				if !isFound {
					fmt.Printf("%s: not found\n", args[0])
				}
			}
		default:
			fmt.Printf("%s: command not found\n", strings.TrimSpace(rawCmd))
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
