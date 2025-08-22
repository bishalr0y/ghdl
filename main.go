package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const helpMessage = `
ghdl - A simple tool to download a file from a raw GitHub URL.

Usage:
  ghdl <output_filename> <github_url>

Example:
  ghdl my_file.txt https://raw.githubusercontent.com/username/repo/branch/path/to/file.ext
`

func main() {
	args := os.Args

	if len(args) != 3 || (len(args) > 1 && args[1] == "help") {
		fmt.Printf(helpMessage)
		return
	}

	outputFilename := args[1]
	rawURL := args[2]

	// Validate the URL
	if err := validateURL(rawURL); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Validate the output location
	if err := validateOutputLocation(outputFilename); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Download the file
	if err := downloadFile(outputFilename, rawURL); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("File downloaded successfully!")
}

func validateURL(rawURL string) error {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if u.Scheme != "https" {
		return fmt.Errorf("invalid URL scheme: must be https")
	}

	if !strings.HasPrefix(u.Host, "raw.githubusercontent.com") {
		return fmt.Errorf("invalid URL: must be a raw GitHub URL")
	}

	return nil
}

func validateOutputLocation(outputFilename string) error {
	dir := filepath.Dir(outputFilename)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("output directory does not exist: %s", dir)
	}

	// Check for write permissions by trying to create a temporary file
	tmpfile, err := os.CreateTemp(dir, "ghdl-")
	if err != nil {
		return fmt.Errorf("no write permissions for output directory: %s", dir)
	}
	tmpfile.Close()
	os.Remove(tmpfile.Name())

	return nil
}

func downloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not download file: %s", resp.Status)
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

