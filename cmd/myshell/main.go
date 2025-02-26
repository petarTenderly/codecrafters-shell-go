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
		command := parseCmd(rawCmd)
		command.exec()
	}

}

type Command struct {
	Name   string
	Args   []string
	Output string
}

func (command Command) Print(s string) {
	if command.Output == "" {
		fmt.Printf(s)
		return
	}
	f, err := os.Create(command.Output)
	if err != nil {
		fmt.Println("error creating file")
	}
	defer f.Close()

	_, err = f.WriteString(s)
	if err != nil {
		fmt.Println("error writing to file")
	}
}

func (command Command) exec() {
	switch command.Name {
	case exitCmd:
		if command.Args[0] == "0" {
			os.Exit(0)
		}
		command.Print("exit: status code must be 0\n")
	case echoCmd:
		command.Print(fmt.Sprintln(strings.Join(command.Args, " ")))
	case typeCmd:
		if slices.Contains(builtIns, command.Args[0]) {
			command.Print(fmt.Sprintf("%s is a shell builtin\n", command.Args[0]))
		} else {
			if path, err := exec.LookPath(command.Args[0]); err == nil {
				command.Print(fmt.Sprintf("%s is %s\n", command.Args[0], path))
			} else {
				command.Print(fmt.Sprintf("%s: not found\n", command.Args[0]))
			}
		}
	case pwdCmd:
		dir, err := os.Getwd()
		if err != nil {
			command.Print("error reading working directory\n")
		}
		fmt.Println(dir)
	case cdCmd:
		dir := command.Args[0]
		if command.Args[0] == "~" {
			dir = os.Getenv("HOME")
		}
		err := os.Chdir(dir)
		if err != nil {
			command.Print(fmt.Sprintf("cd: %s: No such file or directory\n", command.Args[0]))
		}
	default:
		for i, arg := range command.Args {
			command.Args[i] = strings.ReplaceAll(arg, "\n", "\\n")
		}
		c := exec.Command(command.Name, command.Args...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if command.Output != "" {
			f, err := os.Create(command.Output)
			if err != nil {
				command.Print("error creating file")
			}
			defer f.Close()
			c.Stdout = f
		}
		err := c.Run()
		if err != nil {
			command.Print(fmt.Sprintf("%s: command not found\n", command.Name))
		}
	}

}

func parseCmd(rawCmd string) Command {
	cmd := strings.TrimSpace(rawCmd)
	parts := resolveArguments(cmd)
	arguments := parts[1:]
	output := ""
	if len(arguments) > 2 {
		if arguments[len(arguments)-2] == ">" || arguments[len(arguments)-2] == "1>" {
			output = arguments[len(arguments)-1]
			arguments = arguments[:len(arguments)-2]
		}
	}

	return Command{
		Name:   parts[0],
		Args:   arguments,
		Output: output,
	}
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
