package ui

import (
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

type Model struct {
	list    list.Model
	dirPath string
	state   uiState

	descInput  textinput.Model
	cmdInput   textinput.Model
	focusIndex int
	formIsEdit bool

	targetItem  item
	targetIndex int

	help help.Model

	errMsg string

	width  int
	height int

	Selected string
	Execute  bool
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
	keyMap := GetKeys(m.state)

	helpView := m.help.View(keyMap)

	switch m.state {
	case stateForm:
		return m.renderModal(m.formView() + "\n\n" + helpView)
	case stateConfirm:
		return m.renderModal(m.confirmView() + "\n\n" + helpView)
	default:
		return "\n" + m.list.View() + "\n" + helpView
	}
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
