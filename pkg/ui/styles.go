package ui

import "github.com/charmbracelet/lipgloss"

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))

	modalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("170")).
			Padding(1, 2)

	modalTitleStyle = lipgloss.NewStyle().Bold(true).MarginBottom(1)
	fieldLabelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).MarginTop(1)
	dangerStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
)
