package main

import (
	"bufio"
	"fmt"
	"github.com/golang-collections/collections/stack"
	"golang.org/x/term"
	"io"
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
	pwdCmd  = "pwd"
	cdCmd   = "cd"
)

var builtIns = []string{exitCmd, echoCmd, typeCmd, pwdCmd, cdCmd}

func isTab(input string) bool {
	return strings.Contains(input, "\t")
}

func main() {

	// Uncomment this block to pass the first stage
	for true {
		fmt.Fprint(os.Stdout, "\r$ ")

		r := bufio.NewReader(os.Stdin)
		input := readInput(r)
		command := parseCmd(input)
		command.exec()
	}
}

func readInput(rd io.Reader) (input string) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	r := bufio.NewReader(rd)
loop:
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			fmt.Println(err)
			continue
		}
		switch c {
		case '\x03': // Ctrl+C
			os.Exit(0)
		case '\r', '\n': // Enter
			fmt.Fprint(os.Stdout, "\r\n")
			break loop
		case '\x7F': // Backspace
			if length := len(input); length > 0 {
				input = input[:length-1]
				fmt.Fprint(os.Stdout, "\b \b")
			}
		case '\t': // Tab
			suffix := autocomplete(input)
			if suffix != "" {
				input += suffix + " "
				fmt.Fprint(os.Stdout, suffix+" ")
			}
		default:
			input += string(c)
			fmt.Fprint(os.Stdout, string(c))
		}
	}
	return
}

func autocomplete(input string) string {
	for _, cmd := range builtIns {
		if strings.HasPrefix(cmd, input) {
			return cmd[len(input):]
		}
	}
	return ""
}

func parseCmd(rawCmd string) Command {
	cmd := strings.TrimSpace(rawCmd)
	parts := resolveArguments(cmd)

	return NewCommand(parts)
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
