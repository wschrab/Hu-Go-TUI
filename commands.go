// All Helper functions like openBrowser and getPosts

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func openBrowser(url string) error {
	// "cmd" is the Windows command prompt
	// "/c" tells it to run the command and then close the prompt
	// "start" is the Windows command that opens URLs in the default browser
	cmd := exec.Command("cmd", "/c", "start", url)

	// We use Start() instead of Run() because Start() is "fire and forget".
	// It opens the browser and immediately returns control to your TUI without waiting.
	return cmd.Start()
}

// getPosts scans the Hugo content directory
func getPosts(workingDir string) ([]string, map[string]string) {
	// Adjust to fit Hugo content structure
	postsDir := filepath.Join(workingDir, "content", "posts")
	entries, err := os.ReadDir(postsDir)

	postMap := make(map[string]string) // Initialize our translation map

	if err != nil {
		// If the directory can't be read, return an error
		return []string{"Error: Cannot read directory", "Back"}, postMap
	}

	var posts []string
	for _, entry := range entries {
		// Only grab directories (Hugo Page Bundles)
		if entry.IsDir() {
			folderName := entry.Name()
			indexPath := filepath.Join(postsDir, folderName, "index.md")

			// Try to get real title from index.md
			title := extractTitle(indexPath)

			// Fallback logic if title is empty
			displayTitle := title
			if displayTitle == "" {
				displayTitle = fmt.Sprintf("%s (index.md missing title)", folderName)
			}

			posts = append(posts, displayTitle)
			postMap[displayTitle] = folderName
		}
	}

	posts = append(posts, "Back")
	return posts, postMap
}

// extractTitle reads a file and looks for Hugo's title front matter
func extractTitle(filePath string) string {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		// If the file can't be read, return an error
		return "Error: Cannot read file"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Look for standard Hugo title front matter
		if strings.HasPrefix(line, "title:") || strings.HasPrefix(line, "title =") {
			// Split the string at the colon or equals sign
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 1 {
				parts = strings.SplitN(line, "=", 2)
			}

			if len(parts) == 2 {
				// Clean up extra spaces and remove quotation marks
				return strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			}
		}
	}
	return "" //For if there is not a title in the provided file
}
