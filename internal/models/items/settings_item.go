package items

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

type SettingsIndex uint8

const (
	SettingsIndexTheme SettingsIndex = iota
)

type SettingsIndexItem struct {
	Name  string
	Field SettingsIndex
}

func (sii SettingsIndexItem) FilterValue() string { return sii.Name }

type SettingsIndexDelegate struct{}

func (sid SettingsIndexDelegate) Height() int                             { return 3 }
func (sid SettingsIndexDelegate) Spacing() int                            { return 0 }
func (sid SettingsIndexDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (sid SettingsIndexDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	it, ok := item.(SettingsIndexItem)
	if !ok {
		return
	}

	selected := index == m.Index()
	width := m.Width()
	height := sid.Height()

	fmt.Fprint(w, renderSettingsIndexItem(it.Name, width, height, selected))
}

func renderSettingsIndexItem(name string, width int, height int, selected bool) string {
	innerWidth := max(width-4, 76)
	innerHeight := max(height-2, 1)

	// line := styles.Bright.Render(name)
	return BuildCardBox(name, innerWidth, innerHeight, selected)
}
