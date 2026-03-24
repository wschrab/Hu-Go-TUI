// All lipgloss variables (appStyle, cardStyle, etc)

package main

import "github.com/charmbracelet/lipgloss"

// Define our styling
var (
	// Base app styling
	appStyle   = lipgloss.NewStyle().Padding(1, 2)
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1).
			MarginBottom(1)

	// Card styling
	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238")). // Dark gray border
			Padding(0, 2).
			Width(22).
			Height(4).
			Align(lipgloss.Center).        // Horizontal center
			AlignVertical(lipgloss.Center) // Vertical center

	selectedCardStyle = cardStyle.
				BorderForeground(lipgloss.Color("#FF5F87")). // Pink highlight border
				Foreground(lipgloss.Color("#FF5F87"))        // Pink text
)
