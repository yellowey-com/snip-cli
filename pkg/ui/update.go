package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yellowey-com/snip-cli/pkg/storage"
)

func (m Model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok && m.list.FilterState() != list.Filtering {
		switch k.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			if s, ok := m.list.SelectedItem().(item); ok {
				m.Selected = s.snippet.Command
				m.Execute = false
			}
			return m, tea.Quit

		case "x":
			if s, ok := m.list.SelectedItem().(item); ok {
				m.Selected = s.snippet.Command
				m.Execute = true
			}
			return m, tea.Quit

		case "a":
			return m, m.startAdd()

		case "e":
			if s, ok := m.list.SelectedItem().(item); ok {
				return m, m.startEdit(s)
			}
			return m, nil

		case "d":
			if s, ok := m.list.SelectedItem().(item); ok {
				m.targetItem = s
				m.targetIndex = m.list.Index()
				m.state = stateConfirm
				m.errMsg = ""
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *Model) startAdd() tea.Cmd {
	m.state = stateForm
	m.formIsEdit = false
	m.errMsg = ""
	m.targetItem = item{}
	if s, ok := m.list.SelectedItem().(item); ok {
		m.targetItem.category = s.category
	}
	m.descInput.Reset()
	m.cmdInput.Reset()
	m.focusIndex = 0
	m.cmdInput.Blur()
	return m.descInput.Focus()
}

func (m *Model) startEdit(it item) tea.Cmd {
	m.state = stateForm
	m.formIsEdit = true
	m.errMsg = ""
	m.targetItem = it
	m.targetIndex = m.list.Index()
	m.descInput.SetValue(it.snippet.Description)
	m.descInput.Blur()
	m.cmdInput.SetValue(it.snippet.Command)
	m.focusIndex = 1
	return m.cmdInput.Focus()
}

func (m Model) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok {
		switch k.String() {
		case "esc":
			m.state = stateList
			m.errMsg = ""
			return m, nil

		case "tab", "shift+tab", "down", "up":
			if m.formIsEdit {
				return m, nil
			}
			var cmd tea.Cmd
			m.focusIndex = 1 - m.focusIndex
			if m.focusIndex == 0 {
				cmd = m.descInput.Focus()
				m.cmdInput.Blur()
			} else {
				cmd = m.cmdInput.Focus()
				m.descInput.Blur()
			}
			return m, cmd

		case "enter":
			return m.submitForm()
		}
	}

	var cmd tea.Cmd
	if m.focusIndex == 0 && !m.formIsEdit {
		m.descInput, cmd = m.descInput.Update(msg)
	} else {
		m.cmdInput, cmd = m.cmdInput.Update(msg)
	}
	return m, cmd
}

func (m Model) submitForm() (tea.Model, tea.Cmd) {
	command := strings.TrimSpace(m.cmdInput.Value())
	if command == "" {
		m.errMsg = "command is required"
		return m, nil
	}

	if m.formIsEdit {
		if err := storage.EditSnippet(m.dirPath, m.targetItem.category, m.targetItem.snippet.Description, command); err != nil {
			m.errMsg = err.Error()
			return m, nil
		}
		updated := item{
			snippet:  storage.Snippet{Description: m.targetItem.snippet.Description, Command: command},
			category: m.targetItem.category,
		}
		cmd := m.list.SetItem(m.targetIndex, updated)
		m.state = stateList
		m.errMsg = ""
		return m, cmd
	}

	desc := strings.TrimSpace(m.descInput.Value())
	if desc == "" {
		m.errMsg = "description is required"
		return m, nil
	}
	filename := m.targetItem.category
	if filename == "" {
		m.errMsg = "no snippet file to add to"
		return m, nil
	}
	if err := storage.AppendSnippet(m.dirPath, filename, desc, command); err != nil {
		m.errMsg = err.Error()
		return m, nil
	}
	newItem := item{snippet: storage.Snippet{Description: desc, Command: command}, category: filename}
	cmd := m.list.InsertItem(len(m.list.Items()), newItem)
	m.state = stateList
	m.errMsg = ""
	return m, cmd
}

func (m Model) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch k.String() {
	case "y", "enter":
		if err := storage.RemoveSnippet(m.dirPath, m.targetItem.category, m.targetItem.snippet.Description); err != nil {
			m.errMsg = err.Error()
			m.state = stateList
			return m, nil
		}
		m.list.RemoveItem(m.targetIndex)
		m.state = stateList
		m.errMsg = ""
		return m, nil

	case "n", "esc":
		m.state = stateList
		m.errMsg = ""
		return m, nil
	}
	return m, nil
}
