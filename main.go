package main

import (
	"fmt"
	"os"
	"os/exec"
)

const helpMessage = `
ghdl - A simple tool to download a file from a raw GitHub URL.

Usage:
  ghdl <output_filename> <github_url>

Example:
  ghdl my_file.txt https://raw.githubusercontent.com/username/repo/branch/path/to/file.ext
`

// curl -o filename.ext https://raw.githubusercontent.com/username/repo/branch/path/to/file.ext

func main() {
	args := os.Args

	if len(args) != 3 || (len(args) > 1 && args[1] == "help") {
		fmt.Printf(helpMessage)
		return
	}

	filename := args[1]
	url := args[2]

	cmd := exec.Command("curl", "-o", filename, url)

	fmt.Println("Running:", cmd.String())

	// Run the command and capture output (stderr goes to os.Stderr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(string(output))
}

