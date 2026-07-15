package ui

import "github.com/charmbracelet/lipgloss"

var (
	itemStyle = lipgloss.NewStyle().PaddingLeft(2)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(magenta)

	modalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(magenta).
			Padding(1, 2)

	modalTitleStyle = lipgloss.NewStyle().
			Bold(true).
			MarginBottom(1)

	fieldLabelStyle = lipgloss.NewStyle().Foreground(overlay)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f38ba8")).
			MarginTop(1)

	dangerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f38ba8")).
			Bold(true)
)
