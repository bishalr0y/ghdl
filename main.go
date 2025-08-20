package main

import (
	"fmt"
	"os/exec"
)

func main() {
	filename := "config"
	username := "bishalr0y"
	branch := "main"
	repo := "dotfiles"
	path := ".config/ghostty/config"

	url := fmt.Sprintf(
		"https://raw.githubusercontent.com/%v/%v/%v/%v",
		username,
		repo,
		branch,
		path,
	)

	cmd := exec.Command("curl", "-o", filename, url)

	fmt.Println("Running:", cmd.String())

	// Run the command and capture output (stderr goes to os.Stderr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(string(output))
}
