package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yellowey-com/snip-cli/pkg/storage"
)

type item struct {
	snippet  storage.Snippet
	category string
}

func (i item) Title() string       { return i.snippet.Description }
func (i item) Description() string { return "[" + i.category + "] " + i.snippet.Command }
func (i item) FilterValue() string {
	return i.snippet.Description + " " + i.snippet.Command + " " + i.category
}

type Model struct {
	list     list.Model
	Selected string
	Execute  bool
}

func NewModel(items []list.Item) Model {
	l := list.New(items, list.NewDefaultDelegate(), 80, 20)
	l.Title = "Select a snippet"

	l.SetShowHelp(false)
	l.SetShowStatusBar(false)

	l.KeyMap.CursorUp = key.NewBinding(key.WithKeys("up", "k"))
	l.KeyMap.CursorDown = key.NewBinding(key.WithKeys("down", "j"))
	l.KeyMap.Filter = key.NewBinding(key.WithKeys("/"))
	l.KeyMap.ClearFilter = key.NewBinding(key.WithKeys("esc"))

	return Model{list: l}
}

func NewItem(snippet storage.Snippet, category string) list.Item {
	return item{
		snippet:  snippet,
		category: category,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			if selected, ok := m.list.SelectedItem().(item); ok {
				m.Selected = selected.snippet.Command
			}
			return m, tea.Quit

		case "ctrl+x":
			if selected, ok := m.list.SelectedItem().(item); ok {
				m.Selected = selected.snippet.Command
				m.Execute = false
			}
			return m, tea.Quit

		}

	case tea.WindowSizeMsg:
		h, v := msg.Width, msg.Height
		m.list.SetSize(h-4, v-4)

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return "\n" + m.list.View()
}
