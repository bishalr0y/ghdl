package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const helpMessage = `
ghdl - A tool to download files and directories from GitHub.

Usage:
  ghdl <github_url> <output_directory>

Example (file):
  ghdl https://github.com/username/repo/blob/main/path/to/file.ext ./output_dir

Example (directory):
  ghdl https://github.com/username/repo/tree/main/path/to/dir ./output_dir
`

// GitHubContent represents a file or directory in a GitHub repository.
type GitHubContent struct {
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
	Name        string `json:"name"`
	Path        string `json:"path"`
}

func main() {
	args := os.Args

	if len(args) != 3 || (len(args) > 1 && args[1] == "help") {
		fmt.Printf(helpMessage)
		return
	}

	githubURL := args[1]
	outputDir := args[2]

	// Validate the URL
	owner, repo, path, err := parseGitHubURL(githubURL)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Println("Error creating output directory:", err)
		os.Exit(1)
	}

	// Get the API URL for the content
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", owner, repo, path)

	// Download the content
	if err := downloadContent(apiURL, outputDir); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Content downloaded successfully!")
}

func parseGitHubURL(rawURL string) (owner, repo, path string, err error) {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", "", "", fmt.Errorf("invalid URL: %w", err)
	}

	if u.Scheme != "https" {
		return "", "", "", fmt.Errorf("invalid URL scheme: must be https")
	}

	if u.Host != "github.com" {
		return "", "", "", fmt.Errorf("invalid URL: must be a github.com URL")
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 4 {
		return "", "", "", fmt.Errorf("invalid GitHub URL format")
	}

	owner = parts[0]
	repo = parts[1]

	pathIndex := 4
	if len(parts) > 4 && (parts[2] == "tree" || parts[2] == "blob") {
		pathIndex = 4
	}
	path = strings.Join(parts[pathIndex:], "/")

	return owner, repo, path, nil
}

func downloadContent(apiURL, outputDir string) error {
	resp, err := http.Get(apiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not get content: %s", resp.Status)
	}

	var contents []GitHubContent
	if err := json.NewDecoder(resp.Body).Decode(&contents); err != nil {
		// If it's not a JSON array, it might be a single file
		var content GitHubContent
		// We need to "rewind" the body to read it again
		resp.Body.Close()
		resp, err = http.Get(apiURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
			return fmt.Errorf("error decoding github api response: %w", err)
		}
		contents = []GitHubContent{content}
	}

	for _, content := range contents {
		outputPath := filepath.Join(outputDir, content.Name)
		if content.Type == "file" {
			if err := downloadFile(outputPath, content.DownloadURL); err != nil {
				return fmt.Errorf("failed to download file %s: %w", content.Name, err)
			}
		} else if content.Type == "dir" {
			if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", outputPath, err)
			}
			dirAPIURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", getOwnerFromAPIURL(apiURL), getRepoFromAPIURL(apiURL), content.Path)
			if err := downloadContent(dirAPIURL, outputPath); err != nil {
				return err
			}
		}
	}

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

func getOwnerFromAPIURL(apiURL string) string {
	parts := strings.Split(apiURL, "/")
	if len(parts) > 4 {
		return parts[4]
	}
	return ""
}

func getRepoFromAPIURL(apiURL string) string {
	parts := strings.Split(apiURL, "/")
	if len(parts) > 5 {
		return parts[5]
	}
	return ""
}

