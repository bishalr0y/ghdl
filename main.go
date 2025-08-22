package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// curl -o filename.ext https://raw.githubusercontent.com/username/repo/branch/path/to/file.ext

func main() {
	args := os.Args

	if len(args) == 1 || args[1] == "help" {
		fmt.Printf("Usage:\n\nghdl output_location github_url\n")
		return
	}

	urlSlice := strings.Split(args[2], "/")
	filename := urlSlice[len(urlSlice)-1]

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
