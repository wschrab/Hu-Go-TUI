// Core data. Model struct, sessionState constants, and initialModel() function

package main

import (
	"os/exec"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Define app states
type sessionState int

const (
	mainMenuView sessionState = iota
	hugoMenuView
	postsMenuView
	settingsMenuView
	postListView   // All posts viewable under view posts
	postActionView // View options for what to do after selecting an existing post
	newPostView    // View new post options (title, date)
	deleteConfirmView
)

const cols int = 3

type model struct {
	state        sessionState
	cursor       int
	choices      map[sessionState][]string
	postMap		 map[string]string // Map of post names to their content. Maps "display title" to "folder name"
	selectedPost string // Remembers what post is currently open
	textInput    textinput.Model
	serverCmd    *exec.Cmd // Track whether server is running
	workingDir   string    // Track current working directory
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "my-awesome-project"
	ti.CharLimit = 50
	ti.Width = 30

	// TODO: Replace this with os.Getwd() or a directory picker later!
	// For now, hardcode it for local testing.
	testDir := `C:\Users\wdsch\second-portfolio`

	return model{
		state: mainMenuView,
		choices: map[sessionState][]string{
			mainMenuView:     {"Hugo Controls", "Manage Posts", "Settings", "Quit"},
			hugoMenuView:     {"Start Local Server", "Start Local Server (Draft Mode)", "Open localhost:1313", "Build to public", "Back"},
			postsMenuView:    {"New Post", "View Posts", "Back"},
			settingsMenuView: {"Set Editor", "Set Browser", "Back"},
			postActionView:   {"Edit index.md", "View Media", "Delete", "Back"},
		},
		textInput:  ti,
		workingDir: testDir,
		postMap: make(map[string]string),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
