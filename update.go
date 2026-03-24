// Update() quitApp() and handleEnter() logic

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

// Update the model based on the input tea.Msg
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:

		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m.quitApp()
		}

		if m.state == newPostView {
			switch msg.String() {
			case "enter":
				// 1. Get the typed name
				postName := m.textInput.Value()

				if postName != "" {
					// TODO: 2. Create the post
					hugoCmd := exec.Command("hugo", "new", fmt.Sprintf("posts/%s/index.md", postName))

					hugoCmd.Dir = m.workingDir

					err := hugoCmd.Run()
					if err != nil {
						// TODO: Handle error
					}

					// 3. Reset input and return to the posts menu
					m.textInput.SetValue("")
					m.textInput.Blur()
					m.state = postsMenuView
					return m, nil
				}

				// Update text input component with keystrokes
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd

			case "esc":
				// Cancel and go back
				m.textInput.SetValue("")
				m.textInput.Blur()
				m.state = postsMenuView
				return m, nil
			}
			// Update the text input component with keystrokes
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

		if m.state == deleteConfirmView {
			switch msg.String() {
			case "y", "Y":
				// 1. Build the path to the folder
				targetDir := filepath.Join(m.workingDir, "content", "posts", m.selectedPost)

				// 2. Nuke the directory and everything inside it
				err := os.RemoveAll(targetDir)
				if err != nil {
					// TODO: Handle error
				}

				// 3. Go back to the post list and REFRESH the list
				// so the deleted post disappears from the screen
				m.state = postListView
				m.cursor = 0
				titles, pMap := getPosts(m.workingDir)
				m.choices[postListView] = titles
				m.postMap = pMap
				return m, nil

			case "n", "N", "esc", "enter", "backspace":
				// Cancel and go back to the action menu safely
				m.state = postActionView
				return m, nil
			}
			return m, nil

		}

		switch msg.String() {

		case "up", "k":
			if m.cursor-cols >= 0 {
				m.cursor -= cols // Jump backwards a full row
			} else {
				m.cursor = 0
			}

		case "down", "j":
			if m.cursor+cols < len(m.choices[m.state]) {
				m.cursor += cols // Jump forward a full row
			} else {
				m.cursor = len(m.choices[m.state]) - 1 // - 1 Jump to the end of the last row
			}

		case "left", "h":
			if m.cursor%cols != 0 && m.cursor > 0 {
				m.cursor-- // Move left (unless already at the start of a row)
			}

		case "right", "l":
			if (m.cursor+1)%cols != 0 && m.cursor < len(m.choices[m.state])-1 {
				m.cursor++ // Move right (prevent wrapping around the right edge)
			}

		case "enter", " ":
			return m.handleEnter()

		case "backspace":
			if m.state != mainMenuView {
				m.state = mainMenuView
				m.cursor = 0
			}
		}
	}
	return m, nil
}

// handleEnter processes what happens when a user selects a "button"
func (m model) handleEnter() (tea.Model, tea.Cmd) {
	selected := m.choices[m.state][m.cursor]

	// Routing from Main Menu to Sub Menu
	if m.state == mainMenuView {
		switch selected {
		case "Hugo Controls":
			m.state = hugoMenuView
			m.cursor = 0

		case "Manage Posts":
			m.state = postsMenuView
			m.cursor = 0

		case "Settings":
			m.state = settingsMenuView
			m.cursor = 0

		case "Quit":
			return m.quitApp()
		}
		return m, nil // TODO: if this return is reached something went wrong
		// it should not be possible to select a state from the
		// main menu view that isn't one of these
	}

	// Handle actions inside the Hugo Menu
	if m.state == hugoMenuView {
		switch selected {

		case "Open localhost:1313":
			err := openBrowser("http://localhost:1313")
			if err != nil {
				// TODO: deal with broswer not opening error
			}
			return m, nil

		case "Start Local Server":
			// Only start if a server isn't already running
			if m.serverCmd == nil {
				m.serverCmd = exec.Command("hugo", "server")
				m.serverCmd.Dir = m.workingDir // Point to your Hugo root

				err := m.serverCmd.Start() // Run in background
				if err == nil {            // Only change the button if it actually started
					// Update the button text to "Stop"
					m.choices[hugoMenuView][m.cursor] = "Stop Local Server"
					m.choices[hugoMenuView][m.cursor+1] = "Stop Local Server (Draft Mode)"
				}

				if err != nil {
					// TODO: Handle error
				}
			}
			return m, nil

		case "Stop Local Server": // <-- New case to catch the changed button
			if m.serverCmd != nil && m.serverCmd.Process != nil {
				m.serverCmd.Process.Kill()
				m.serverCmd = nil // Reset the command

				// Revert the button text back to "Start"
				m.choices[hugoMenuView][m.cursor] = "Start Local Server"
				m.choices[hugoMenuView][m.cursor+1] = "Start Local Server (Draft Mode)"
			}
			return m, nil

		case "Start Local Server (Draft Mode)":
			// Only start if a server isn't already running
			if m.serverCmd == nil {
				m.serverCmd = exec.Command("hugo", "server", "-D")
				m.serverCmd.Dir = m.workingDir // Point to your Hugo root

				err := m.serverCmd.Start() // Run in background
				if err == nil {
					m.choices[hugoMenuView][m.cursor] = "Stop Local Server (Draft Mode)"
					m.choices[hugoMenuView][m.cursor-1] = "Stop Local Server"
				}

				if err != nil {
					// TODO: Handle error
				}
			}
			return m, nil

		case "Stop Local Server (Draft Mode)":
			if m.serverCmd != nil && m.serverCmd.Process != nil {
				m.serverCmd.Process.Kill()
				m.serverCmd = nil

				m.choices[hugoMenuView][m.cursor] = "Start Local Server (Draft Mode)"
				m.choices[hugoMenuView][m.cursor-1] = "Start Local Server"
			}
			return m, nil
		}

	}

	if m.state == postsMenuView {
		switch selected {
		case "View Posts":
			m.state = postListView
			m.cursor = 0
			// Dynamically load the posts into the choices map
			titles, pMap := getPosts(m.workingDir)
			m.choices[postListView] = titles
			m.postMap = pMap
			return m, nil

		case "New Post":
			m.state = newPostView
			m.textInput.Focus() // Activate the blinking cursor
			return m, nil
		}
	}

	if m.state == postListView {
		if selected == "Back" {
			m.state = postsMenuView
			m.cursor = 0
			return m, nil
		} else {
			// If it's not "Back", the selection must be a project folder
			m.selectedPost = m.postMap[selected] // Save the folder name
			m.state = postActionView  //Move to action view
			m.cursor = 0
			return m, nil
		}
	}

	if m.state == postActionView {
		switch selected {
		case "Edit index.md":
			// 1. Build the path to the file
			targetFile := filepath.Join(m.workingDir, "content", "posts", m.selectedPost, "index.md")

			// 2. Use start command to open file
			// Note: empty string "" is quirky windows reuirement when using start with paths
			// that may contain spaces. (according to Gemini)
			cmd := exec.Command("cmd", "/c", "start", "", targetFile)

			// 3. Launch application
			err := cmd.Start()
			if err != nil {
				// TODO: Handle the error (show an error message)
			}
			return m, nil

		case "View Media":
			// 1. Build the path to the specific project folder
			targetDir := filepath.Join(m.workingDir, "content", "posts", m.selectedPost)

			// 2. Open the folder in Windows File Explorer
			cmd := exec.Command("cmd", "/c", "start", "", targetDir)

			err := cmd.Start()
			if err != nil {
				// TODO: Handle the error
			}
			return m, nil

		case "Delete":
			m.state = deleteConfirmView
			return m, nil
		case "Back":
			m.state = postListView
			m.cursor = 0
			return m, nil
		}
	}

	// Handle "Back" buttons in sub-menus
	if selected == "Back" {
		if m.state == postListView {
			m.state = postsMenuView
		} else {
			m.state = mainMenuView
		}
		m.cursor = 0
		return m, nil
	}

	// TODO: Add control within submenus

	return m, nil
}

func (m model) quitApp() (tea.Model, tea.Cmd) {
	// 1. Clean up the Hugo server if it's running
	if m.serverCmd != nil && m.serverCmd.Process != nil {
		m.serverCmd.Process.Kill()
	}

	return m, tea.Quit
}
