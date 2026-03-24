// All Helper functions like openBrowser and getPosts

package main

import (
	"os"
	"os/exec"
	"path/filepath"
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
func getPosts(workingDir string) []string {
	// Adjust to fit Hugo content structure
	postsDir := filepath.Join(workingDir, "content", "posts")

	entries, err := os.ReadDir(postsDir)
	if err != nil {
		// If the directory can't be read, return an error
		return []string{"Error: Cannot read directory", "Back"}
	}

	var posts []string
	for _, entry := range entries {
		// Only grab directories (Hugo Page Bundles)
		if entry.IsDir() {
			posts = append(posts, entry.Name())
		}
	}

	posts = append(posts, "Back")

	return posts
}
