// View() function that draws the screen

package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	// If we are typing text input, bypass normal grid rendering
	if m.state == newPostView {
		return appStyle.Render(
			fmt.Sprintf(
				"Enter new project name:\n\n%s\n\n%s",
				m.textInput.View(),
				lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("(esc to cancel, enter to submit)"),
			),
		)
	}

	var s string

	// Renter Title based on state
	switch m.state {
	case mainMenuView:
		s += titleStyle.Render("Hugo TUI Manager") + "\n"
	case hugoMenuView:
		s += titleStyle.Render("Hugo Controls") + "\n"
	case postsMenuView:
		s += titleStyle.Render("Posts") + "\n"
	case settingsMenuView:
		s += titleStyle.Render("Settings") + "\n"
	case postListView:
		s += titleStyle.Render("Current Posts") + "\n"
	case postActionView:
		s += titleStyle.Render("Managing: "+m.selectedPost) + "\n"
	}

	var rows []string
	var currentRow []string

	// 1. Render the card by itself
	for i, choice := range m.choices[m.state] {
		var card string
		if i == m.cursor {
			card = selectedCardStyle.Render(choice)
		} else {
			card = cardStyle.Render(choice)
		}
		currentRow = append(currentRow, card)

		// 2. If the row is full or on the last item join horizontal items
		// and start a new row
		if len(currentRow) == cols {
			// Join Horizontal stitches strings side-by-side
			joinedRow := lipgloss.JoinHorizontal(lipgloss.Top, currentRow...)
			rows = append(rows, joinedRow)
			currentRow = nil
		}
	}

	if len(currentRow) > 0 {
		joinedRow := lipgloss.JoinHorizontal(lipgloss.Top, currentRow...)
		rows = append(rows, joinedRow)
	}

	// 3. Stack all the rows vertically
	grid := lipgloss.JoinVertical(lipgloss.Left, rows...)
	s += grid + "\n"

	s += lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("\nPress q to quit.")
	return appStyle.Render(s)
}
