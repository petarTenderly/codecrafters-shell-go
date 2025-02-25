package main

import (
	"bufio"
	"fmt"
	"github.com/golang-collections/collections/stack"
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
			for i, arg := range args {
				args[i] = strings.ReplaceAll(arg, "\n", "\\n")
			}
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
	//split := strings.SplitN(cmd, " ", 2)
	//Extract the command and arguments
	//command := split[0]
	//if len(split) == 1 {
	//	return command, []string{}
	//}
	//argumentsString := split[1]

	// Split the arguments into a list
	arguments := resolveArguments(cmd)

	return arguments[0], arguments[1:]
}

func resolveArguments(argument string) []string {
	sb := strings.Builder{}
	s := stack.New()
	argList := make([]string, 0)

	for i := 0; i < len(argument); {
		// strongest rule is to ignore special chars when we have quotes unless
		if s.Peek() == uint8(39) {
			if argument[i] == '\'' {
				s.Pop()
				i++
				continue
			}
			sb.WriteByte(argument[i])
			i++
			continue
			// logic when in single quotas
		}
		//double quotes
		if s.Peek() == uint8(34) {
			if argument[i] == '"' {
				s.Pop()
				i++
				continue
			}
			specRunes := []uint8{34, 92, 36}
			if argument[i] == '\\' && i+2 < len(argument) {
				if slices.Contains(specRunes, argument[i+1]) {
					sb.WriteByte(argument[i+1])
					i += 2
					continue
				}
			}
			sb.WriteByte(argument[i])
			i++
			continue
		}

		if argument[i] == '\'' || argument[i] == '"' {
			if s.Peek() == argument[i] {
				s.Pop()
			} else {
				s.Push(argument[i])
			}
			i++
			continue
		}

		// when we are in single column then we just append the char

		if s.Peek() == nil {
			// new line not in quotes, this is new argument
			if argument[i] == ' ' {
				argList = append(argList, sb.String())
				sb = strings.Builder{}
				i++
				// and we need to skip all lines
				for argument[i] == ' ' {
					i++
				}
				continue
			}

			if argument[i] == '\\' {
				// we need to skip this char and add next char
				i++
				if i < len(argument) {
					sb.WriteByte(argument[i])
				}
				i++
				continue
			}

			//we need to add special char if there is escape char
			sb.WriteByte(argument[i])
			i++

			continue
		}

	}

	// we need to add last argument if command is not finished with space
	if sb.Len() > 0 {
		argList = append(argList, sb.String())

	}
	return argList
}

//func main() {
//	fmt.Printf(resolveArguments("hello'script'\\n'world"))
//}
