package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	logoAccentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#82aaff"))
	logoDimStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#5c6370"))
	logoTextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#abb2bf"))
	logoBoxStyle    = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#2d3139")).
			Padding(1, 3).
			MaxWidth(50)
)

func RenderLogo() string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		logoAccentStyle.Bold(true).Render("> snip run reborn"),
		"",
		logoDimStyle.Render("# Delete lockfiles, reinstall and run"),
		logoTextStyle.Render("rm -rf node_modules package-lock.json &&"),
		logoTextStyle.Render("npm install && npm run dev"),
	)

	return logoBoxStyle.Render(content)
}
