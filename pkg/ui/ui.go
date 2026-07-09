package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type uiState int

const (
	stateList uiState = iota
	stateForm
	stateConfirm
)

var (
	magenta = lipgloss.Color("#cba6f7")
	subtext = lipgloss.Color("#a6adc8")
	overlay = lipgloss.Color("#6c7086")
	surface = lipgloss.Color("#313244")

	appStyle = lipgloss.NewStyle().Padding(1, 2)

	headerTitleStyle = lipgloss.NewStyle().Foreground(magenta).Bold(true)
	headerCountStyle = lipgloss.NewStyle().Foreground(subtext)
	headerLineStyle  = lipgloss.NewStyle().Foreground(surface)
	footerStyle      = lipgloss.NewStyle().Foreground(overlay)
	footerKeyStyle   = lipgloss.NewStyle().Foreground(magenta)

	activeTitleStyle = lipgloss.NewStyle().Foreground(magenta).Bold(true)
	activeDescStyle  = lipgloss.NewStyle().Foreground(overlay)

	inactiveTitleStyle = lipgloss.NewStyle().Foreground(subtext)
	inactiveDescStyle  = lipgloss.NewStyle().Foreground(surface)

	activeBorder = lipgloss.NewStyle().
			Border(lipgloss.Border{Left: "┃"}, false, false, false, true).
			BorderForeground(magenta).
			PaddingLeft(1)

	inactiveBorder = lipgloss.NewStyle().PaddingLeft(2)
)

type Model struct {
	list        list.Model
	dirPath     string
	state       uiState
	descInput   textinput.Model
	cmdInput    textinput.Model
	focusIndex  int
	formIsEdit  bool
	targetItem  item
	targetIndex int
	help        help.Model
	errMsg      string
	width       int
	height      int
	Selected    string
	Execute     bool
}

func NewModel(items []list.Item, dirPath string, filterQuery string) Model {
	l := list.New(items, itemDelegate{}, 80, 20)
	l.Title = ""
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)

	l.KeyMap.CursorUp = key.NewBinding(key.WithKeys("up", "k"))
	l.KeyMap.CursorDown = key.NewBinding(key.WithKeys("down", "j"))
	l.KeyMap.Filter = key.NewBinding(key.WithKeys("/"))
	l.KeyMap.ClearFilter = key.NewBinding(key.WithKeys("esc"))

	if filterQuery != "" {
		l.SetFilterText(filterQuery)
	}

	h := help.New()
	h.ShowAll = true

	desc := textinput.New()
	desc.Placeholder = "Description"
	desc.CharLimit = 200
	desc.Width = 40

	cmdIn := textinput.New()
	cmdIn.Placeholder = "Command"
	cmdIn.CharLimit = 500
	cmdIn.Width = 40

	return Model{list: l, dirPath: dirPath, descInput: desc, cmdInput: cmdIn, help: h}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if wsm, ok := msg.(tea.WindowSizeMsg); ok {
		m.width, m.height = wsm.Width, wsm.Height
		m.help.Width = wsm.Width
		m.list.SetSize(wsm.Width-4, wsm.Height-4)
		return m, nil
	}

	switch m.state {
	case stateForm:
		return m.updateForm(msg)
	case stateConfirm:
		return m.updateConfirm(msg)
	default:
		return m.updateList(msg)
	}
}

func (m Model) View() string {
	switch m.state {
	case stateForm:
		return m.renderModal(m.formView())
	case stateConfirm:
		return m.renderModal(m.confirmView())
	default:
		return m.listView()
	}
}

func (m Model) listView() string {
	var s strings.Builder
	totalItems := len(m.list.VisibleItems())

	title := headerTitleStyle.Render("snip")
	count := headerCountStyle.Render(fmt.Sprintf("%d snippets", totalItems))
	filterMode := headerTitleStyle.Render("all")
	if m.list.FilterState() == list.Filtering {
		filterMode = headerTitleStyle.Render("searching")
	}
	filterStr := headerCountStyle.Render("filter: ") + filterMode

	headerWidth := m.width - 4
	if headerWidth < 40 {
		headerWidth = 60
	}

	leftHeader := fmt.Sprintf("%s  •  %s ", title, count)
	rightHeader := fmt.Sprintf(" %s", filterStr)
	lineWidth := headerWidth - lipgloss.Width(leftHeader) - lipgloss.Width(rightHeader)

	var lineStr string
	if lineWidth > 0 {
		lineStr = headerLineStyle.Render(strings.Repeat("─", lineWidth))
	}

	s.WriteString(leftHeader + lineStr + rightHeader + "\n\n")

	if totalItems == 0 {
		s.WriteString(inactiveTitleStyle.Render("  No snippets found.") + "\n")
	} else {
		cursor := m.list.Index()
		for i, rawItem := range m.list.VisibleItems() {
			si := rawItem.(item)

			var itemStr string
			if i == cursor {
				titleLine := activeTitleStyle.Render(si.snippet.Description)
				descLine := activeDescStyle.Render(fmt.Sprintf("[%s] %s", si.category, si.snippet.Command))
				itemStr = activeBorder.Render(fmt.Sprintf("%s\n%s", titleLine, descLine))
			} else {
				titleLine := inactiveTitleStyle.Render(si.snippet.Description)
				descLine := inactiveDescStyle.Render(fmt.Sprintf("[%s] %s", si.category, si.snippet.Command))
				itemStr = inactiveBorder.Render(fmt.Sprintf("%s\n%s", titleLine, descLine))
			}
			s.WriteString(itemStr + "\n\n")
		}
	}

	s.WriteString(headerLineStyle.Render(strings.Repeat("─", headerWidth)) + "\n")

	nav := fmt.Sprintf("%s navigate", footerKeyStyle.Render("↑/↓"))
	run := fmt.Sprintf("%s run", footerKeyStyle.Render("enter"))
	execute := fmt.Sprintf("%s execute", footerKeyStyle.Render("x"))
	filt := fmt.Sprintf("%s filter", footerKeyStyle.Render("/"))
	quit := fmt.Sprintf("%s quit", footerKeyStyle.Render("ctrl+c"))

	footer := fmt.Sprintf("%-20s | %-16s | %-16s | %-16s | %s", nav, run, execute, filt, quit)
	s.WriteString(footerStyle.Render(footer))

	return appStyle.Render(s.String())
}

func (m Model) renderModal(content string) string {
	box := modalStyle.Render(content)
	if m.width == 0 || m.height == 0 {
		return box
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

func (m Model) formView() string {
	title := "Add Snippet"
	if m.formIsEdit {
		title = "Edit Snippet"
	}

	var b strings.Builder
	b.WriteString(modalTitleStyle.Render(title))
	b.WriteString("\n")

	b.WriteString(fieldLabelStyle.Render("Description"))
	b.WriteString("\n")
	if m.formIsEdit {
		b.WriteString(m.targetItem.snippet.Description)
	} else {
		b.WriteString(m.descInput.View())
	}
	b.WriteString("\n\n")

	b.WriteString(fieldLabelStyle.Render("Command"))
	b.WriteString("\n")
	b.WriteString(m.cmdInput.View())

	if m.errMsg != "" {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render(m.errMsg))
	}

	b.WriteString("\n\n")
	b.WriteString(m.help.View(GetKeys(m.state)))

	return b.String()
}

func (m Model) confirmView() string {
	var b strings.Builder
	b.WriteString(modalTitleStyle.Render("Delete Snippet"))
	b.WriteString("\n")
	b.WriteString("Delete ")
	b.WriteString(dangerStyle.Render(m.targetItem.snippet.Description))
	b.WriteString("?")

	if m.errMsg != "" {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render(m.errMsg))
	}

	b.WriteString("\n\n")
	b.WriteString(m.help.View(GetKeys(m.state)))

	return b.String()
}
