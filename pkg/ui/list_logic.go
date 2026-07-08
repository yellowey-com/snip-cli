package ui

import (
	"io"

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
func (i item) FilterValue() string { return i.snippet.Description }

func NewItem(snippet storage.Snippet, category string) list.Item {
	return item{snippet: snippet, category: category}
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 2 }
func (d itemDelegate) Spacing() int                              { return 1 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}
	str := i.Title() + "\n" + i.Description()
	if index == m.Index() {
		str = selectedItemStyle.Render("> " + str)
	} else {
		str = itemStyle.Render("  " + str)
	}
	io.WriteString(w, str)
}
