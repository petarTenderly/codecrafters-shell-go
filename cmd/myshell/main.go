package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

const (
	exitCmd = "exit"
	echoCmd = "echo"
)

func main() {
	// Uncomment this block to pass the first stage
	for true {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		rawCmd, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		cmd := strings.TrimSpace(rawCmd)
		command, arg, found := strings.Cut(cmd, " ")
		if !found {
			fmt.Printf("%s: command not found\n", cmd)
			continue
		}
		switch command {
		case exitCmd:
			if arg == "0" {
				os.Exit(0)
			}
			fmt.Println("exit: status code must be 0")
		case echoCmd:
			fmt.Println(arg)
		default:
			fmt.Printf("%s: command not found\n", cmd)
		}

	}

}
