package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/ansi"
)

var (
	emptyAccent = lipgloss.Color("#82aaff")
	emptyDim    = lipgloss.Color("#5c6370")
	emptyBorder = lipgloss.Color("#2d3139")
	emptyText   = lipgloss.Color("#abb2bf")
	emptyWhite  = lipgloss.Color("#ffffff")

	emptyHeaderAccentStyle = lipgloss.NewStyle().Foreground(emptyAccent).Bold(true)
	emptyHeaderDimStyle    = lipgloss.NewStyle().Foreground(emptyDim)

	emptyTitleStyle   = lipgloss.NewStyle().Foreground(emptyWhite).Bold(true)
	emptyDescStyle    = lipgloss.NewStyle().Foreground(emptyText)
	emptyPathStyle    = lipgloss.NewStyle().Foreground(emptyAccent)
	emptyStepStyle    = lipgloss.NewStyle().Foreground(emptyText)
	emptyStepNumStyle = lipgloss.NewStyle().Foreground(emptyAccent).Bold(true)

	emptyDividerStyle   = lipgloss.NewStyle().Foreground(emptyBorder)
	emptyFooterStyle    = lipgloss.NewStyle().Foreground(emptyDim)
	emptyFooterKeyStyle = lipgloss.NewStyle().Foreground(emptyAccent)
)

const (
	splitMinWidth = 80
	leftPaneWidth = 32
	paneGap       = 4
)

func RenderEmptyState(width, height int) string {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	header := emptyHeaderRow(width)
	divider := emptyDividerStyle.Render(strings.Repeat("─", width))
	footer := emptyGlobalFooter(width)

	var content string
	if width >= splitMinWidth {
		content = emptySplitContent()
	} else {
		content = emptyStackedContent(width)
	}

	usedHeight := lipgloss.Height(header) + lipgloss.Height(divider)*2 + lipgloss.Height(footer)
	contentHeight := height - usedHeight
	if contentHeight < lipgloss.Height(content) {
		contentHeight = lipgloss.Height(content)
	}

	centered := lipgloss.Place(width, contentHeight, lipgloss.Center, lipgloss.Center, content)

	return lipgloss.JoinVertical(lipgloss.Left, header, divider, centered, divider, footer)
}

func emptyHeaderRow(width int) string {
	left := emptyHeaderAccentStyle.Render("snip") + emptyHeaderDimStyle.Render("  •  0 snippets")
	right := emptyHeaderDimStyle.Render("filter: ") + emptyHeaderAccentStyle.Render("all")

	gap := width - ansi.PrintableRuneWidth(left) - ansi.PrintableRuneWidth(right)
	if gap < 1 {
		return emptyHeaderAccentStyle.Render("snip")
	}
	return left + strings.Repeat(" ", gap) + right
}

func emptyQuickStart() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		emptyStepNumStyle.Render("1  ")+emptyStepStyle.Render("Create a .md file"),
		emptyStepNumStyle.Render("2  ")+emptyStepStyle.Render("Add a title and command"),
		emptyStepNumStyle.Render("3  ")+emptyStepStyle.Render("Reload and run snippets"),
	)
}

func emptySplitContent() string {
	left := lipgloss.JoinVertical(
		lipgloss.Left,
		emptyTitleStyle.Width(leftPaneWidth).Render("No snippets yet"),
		"",
		emptyDescStyle.Width(leftPaneWidth).Render("Add your first markdown snippet file to get started."),
		emptyPathStyle.Width(leftPaneWidth).Render("~/.config/snip/snippets/"),
		"",
		emptyQuickStart(),
	)

	logo := RenderLogo()

	paneHeight := lipgloss.Height(left)
	if lh := lipgloss.Height(logo); lh > paneHeight {
		paneHeight = lh
	}

	left = lipgloss.NewStyle().Width(leftPaneWidth).Height(paneHeight).Render(left)
	logo = lipgloss.NewStyle().Height(paneHeight).Render(logo)

	divider := strings.TrimSuffix(strings.Repeat("│\n", paneHeight), "\n")
	divider = emptyDividerStyle.Render(divider)

	gap := strings.Repeat(" ", paneGap)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, gap, divider, gap, logo)
}

func emptyStackedContent(width int) string {
	cardWidth := width - 8
	if cardWidth > 50 {
		cardWidth = 50
	}
	if cardWidth < 20 {
		cardWidth = 20
	}

	title := emptyTitleStyle.Width(cardWidth).Align(lipgloss.Center).Render("No snippets yet")
	desc := emptyDescStyle.Width(cardWidth).Align(lipgloss.Center).Render("Add your first markdown snippet file to get started.")
	path := emptyPathStyle.Width(cardWidth).Align(lipgloss.Center).Render("~/.config/snip/snippets/")
	steps := lipgloss.NewStyle().Width(cardWidth).Align(lipgloss.Center).Render(emptyQuickStart())

	return lipgloss.JoinVertical(lipgloss.Center, title, "", desc, path, "", steps)
}

func emptyGlobalFooter(width int) string {
	items := []string{
		emptyFooterKeyStyle.Render("a") + emptyFooterStyle.Render(" new"),
		emptyFooterKeyStyle.Render("enter") + emptyFooterStyle.Render(" copy"),
		emptyFooterKeyStyle.Render("x") + emptyFooterStyle.Render(" execute"),
		emptyFooterKeyStyle.Render("/") + emptyFooterStyle.Render(" filter"),
		emptyFooterKeyStyle.Render("q") + emptyFooterStyle.Render(" quit"),
	}
	line := wrapItems(items, emptyFooterStyle.Render("    •    "), width)
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(line)
}

func wrapItems(items []string, sep string, width int) string {
	if len(items) == 0 {
		return ""
	}

	var lines []string
	current := items[0]

	for _, item := range items[1:] {
		candidate := current + sep + item
		if ansi.PrintableRuneWidth(candidate) > width {
			lines = append(lines, current)
			current = item
			continue
		}
		current = candidate
	}
	lines = append(lines, current)

	return strings.Join(lines, "\n")
}
