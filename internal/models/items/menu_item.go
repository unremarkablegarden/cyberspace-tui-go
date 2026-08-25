package items

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

type MenuSection uint8

const (
	MenuSectionFeed MenuSection = iota
	MenuSectionCompose
	MenuSectionNotifications
	MenuSectionOwnProfile
	MenuSectionNotes
	MenuSectionTopics
	MenuSectionBookmarks
	MenuSectionSettings
)

type MenuItem struct {
	Name    string
	Field   MenuSection
	Keybind key.Binding
}

func (mi MenuItem) FilterValue() string { return mi.Name }
func (mi MenuItem) Title() string       { return "[" + mi.Name + "]" }
func (mi MenuItem) Description() string { return "" }

type MenuDelegate struct{}

func (d MenuDelegate) Height() int                             { return 3 }
func (d MenuDelegate) Spacing() int                            { return 0 }
func (d MenuDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d MenuDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	it, ok := item.(MenuItem)
	if !ok {
		return
	}

	selected := index == m.Index()
	width := m.Width()
	height := d.Height()

	kbSlice := it.Keybind.Keys()
	kb := "N/A"
	if len(kbSlice) > 0 {
		kb = kbSlice[0]
	}

	fmt.Fprint(w, renderMenuItem(it.Name, kb, width, height, selected))
}

func renderMenuItem(name string, keybind string, width int, height int, selected bool) string {
	innerWidth := max(width-4, 76)
	innerHeight := max(height-2, 1)

	tag := "[" + keybind + "]" + " " + name

	line := styles.Bright.Render(tag)
	return BuildCardBox(line, innerWidth, innerHeight, selected)
}
