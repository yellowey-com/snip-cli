package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	state uiState

	Up      key.Binding
	Down    key.Binding
	Select  key.Binding
	Execute key.Binding
	Add     key.Binding
	Edit    key.Binding
	Delete  key.Binding
	Filter  key.Binding
	Quit    key.Binding

	Next   key.Binding
	Submit key.Binding
	Cancel key.Binding

	Confirm key.Binding
	Deny    key.Binding
}

var baseKeyMap = keyMap{
	Up:      key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:    key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Select:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "copy")),
	Execute: key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "execute")),
	Add:     key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
	Edit:    key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
	Delete:  key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
	Filter:  key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
	Quit:    key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),

	Next:   key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next field")),
	Submit: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
	Cancel: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),

	Confirm: key.NewBinding(key.WithKeys("y", "enter"), key.WithHelp("y", "confirm")),
	Deny:    key.NewBinding(key.WithKeys("n", "esc"), key.WithHelp("n", "cancel")),
}

func GetKeys(state uiState) help.KeyMap {
	k := baseKeyMap
	k.state = state
	return k
}

func (k keyMap) ShortHelp() []key.Binding {
	switch k.state {
	case stateForm:
		return []key.Binding{k.Submit, k.Next, k.Cancel}
	case stateConfirm:
		return []key.Binding{k.Confirm, k.Deny}
	default:
		return []key.Binding{k.Select, k.Execute, k.Add, k.Edit, k.Delete, k.Filter, k.Quit}
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	switch k.state {
	case stateForm:
		return [][]key.Binding{{k.Submit, k.Cancel}, {k.Next}}
	case stateConfirm:
		return [][]key.Binding{{k.Confirm, k.Deny}}
	default:
		return [][]key.Binding{
			{k.Up, k.Down, k.Filter},
			{k.Add, k.Edit, k.Delete},
			{k.Select, k.Execute, k.Quit},
		}
	}
}
