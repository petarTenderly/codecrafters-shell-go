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
	ExitCmd = "exit 0"
)

func main() {
	// Uncomment this block to pass the first stage
	for true {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		rawCmd, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		cmd := strings.TrimSpace(rawCmd)
		switch cmd {
		case ExitCmd:
			os.Exit(0)
		}

		fmt.Printf("%s: command not found\n", cmd)
	}

}
