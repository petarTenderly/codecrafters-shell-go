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

type Shell struct {
	executiveFiles  []string
	argumentHistory []string
	tabPressed      int
}

func NewShell() *Shell {
	executiveFiles := execList()
	return &Shell{
		executiveFiles:  executiveFiles,
		argumentHistory: make([]string, 0),
		tabPressed:      0,
	}
}

func (shell *Shell) StartLoop() {
	for {
		fmt.Fprint(os.Stdout, "\r$ ")

		r := bufio.NewReader(os.Stdin)
		input := shell.readInput(r)
		command := shell.parseCmd(input)
		shell.argumentHistory = append(shell.argumentHistory, command.Args...)
		command.exec()
	}
}

func (shell *Shell) autocomplete(input string) string {
	if strings.Contains(input, " ") {
		//autocomplete arguments
		lastArg := strings.Split(input, " ")[len(strings.Split(input, " "))-1]
		for _, arg := range shell.argumentHistory {
			if strings.HasPrefix(arg, lastArg) {
				return arg[len(input):] + " "
			}
		}
	} else {
		for _, cmd := range builtIns {
			if strings.HasPrefix(cmd, input) {
				return cmd[len(input):] + " "
			}
		}

		listOfMatches := make([]string, 0)
		for _, cmd := range shell.executiveFiles {
			if strings.HasPrefix(cmd, input) {
				listOfMatches = append(listOfMatches, cmd)
			}
		}
		if len(listOfMatches) == 1 {
			shell.tabPressed = 0
			return listOfMatches[0][len(input):] + " "
		}
		if len(listOfMatches) > 1 {
			if shell.tabPressed > 0 {
				fmt.Print("\a")
			}
			cursor := len(input)
		findCommonSuffix:
			for {
				if cursor >= len(listOfMatches[0]) {
					break findCommonSuffix
				}
				runeAtCursor := listOfMatches[0][cursor]
				for _, match := range listOfMatches {
					if cursor >= len(match) {
						break findCommonSuffix
					}
					if runeAtCursor != match[cursor] {
						break findCommonSuffix
					}

				}
				cursor++
			}

			if cursor > len(input) {
				shell.tabPressed = 0
				return listOfMatches[0][len(input):cursor]
			}
			shell.tabPressed++
			slices.Sort(listOfMatches)
			fmt.Printf("\r\n%s\n\r", strings.Join(listOfMatches, "  "))
			fmt.Print("$ ", input)

			return ""
		}
	}
	fmt.Print("\a")
	return ""
}

func (shell *Shell) readInput(rd io.Reader) (input string) {
	const (
		escRune       = '\x03'
		rRune         = '\r'
		nRune         = '\n'
		backspaceRune = '\x7F'
		tabRune       = '\t'
	)

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
		case escRune: // Ctrl+C
			os.Exit(0)
		case rRune, nRune: // Enter
			fmt.Fprint(os.Stdout, "\r\n")
			break loop
		case backspaceRune: // Backspace
			if length := len(input); length > 0 {
				input = input[:length-1]
				fmt.Fprint(os.Stdout, "\b \b")
			}
		case tabRune: // Tab
			suffix := shell.autocomplete(input)
			if suffix != "" {
				input += suffix
				fmt.Fprint(os.Stdout, suffix)
			}
		default:
			input += string(c)
			fmt.Fprint(os.Stdout, string(c))
		}
	}
	return
}

var tabPressed = 0

func (shell *Shell) parseCmd(rawCmd string) Command {
	cmd := strings.TrimSpace(rawCmd)
	parts := shell.resolveArguments(cmd)

	return NewCommand(parts)
}

func (shell *Shell) resolveArguments(argument string) []string {
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
