package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func execList() []string {
	allExecutables := make([]string, 0)
	// Get the PATH environment variable
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		fmt.Println("PATH environment variable is not set")
		return nil
	}

	// Split the PATH into individual directories
	pathDirs := strings.Split(pathEnv, string(os.PathListSeparator))

	// Iterate through each directory in PATH
	for _, dir := range pathDirs {
		// Open the directory
		d, err := os.Open(dir)
		if err != nil {
			// If the directory cannot be opened, skip it
			continue
		}
		defer d.Close()

		// Read all files in the directory
		files, err := d.Readdir(-1)
		if err != nil {
			continue
		}

		// Check each file to see if it is executable
		for _, file := range files {
			// Skip directories
			if file.IsDir() {
				continue
			}

			// Construct the full path to the file
			fullPath := filepath.Join(dir, file.Name())

			// Check if the file is executable
			if isExecutable(fullPath) {
				if slices.Contains(allExecutables, file.Name()) {
					continue
				}
				allExecutables = append(allExecutables, file.Name())
			}
		}
	}

	return allExecutables
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// Check if the file is executable by the user
	return info.Mode().Perm()&0111 != 0
}
