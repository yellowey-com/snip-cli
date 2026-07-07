package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yellowey-com/snip-cli/pkg/storage"
)

type item struct {
	snippet storage.Snippet
}

func (i item) Title() string       { return i.snippet.Description }
func (i item) Description() string { return i.snippet.Command }
func (i item) FilterValue() string { return i.snippet.Description + " " + i.snippet.Command }

type Model struct {
	list     list.Model
	Selected string
}

func NewModel(snippets []storage.Snippet) Model {
	var items []list.Item
	for _, snip := range snippets {
		items = append(items, item{snippet: snip})
	}

	l := list.New(items, list.NewDefaultDelegate(), 80, 20)
	l.Title = "Select as snippet"

	return Model{list: l}
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
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return "\n" + m.list.View()
}
