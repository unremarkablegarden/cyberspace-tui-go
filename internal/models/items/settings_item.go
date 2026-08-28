package items

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

type SettingsIndex uint8

// Do not change the order of first item
// since I want index to be the default
// value for this "enum"
const (
	SettingsIndexIndex SettingsIndex = iota
	SettingsIndexTheme
	SettingsIndexKeybind
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

/*-------------------------------------
FOR KEYBINDS
-------------------------------------*/

type SettingsKeybindsItem struct {
	ID    string
	Name  string
	Value string
}

func (ski SettingsKeybindsItem) FilterValue() string { return ski.ID }

type SettingsKeybindsDelegate struct{}

func (skd SettingsKeybindsDelegate) Height() int                             { return 3 }
func (skd SettingsKeybindsDelegate) Spacing() int                            { return 0 }
func (skd SettingsKeybindsDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (skd SettingsKeybindsDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	it, ok := item.(SettingsKeybindsItem)
	if !ok {
		return
	}

	selected := index == m.Index()
	width := m.Width()
	height := skd.Height()

	fmt.Fprint(w, renderSettingsKeybindsItem(it.Name, it.Value, width, height, selected))
}

func renderSettingsKeybindsItem(
	name string,
	value string,
	width int,
	height int,
	selected bool,
) string {
	innerWidth := max(width-4, 76)
	innerHeight := max(height-2, 1)

	content := name + ": [" + value + "]"
	// line := styles.Bright.Render(name)
	return BuildCardBox(content, innerWidth, innerHeight, selected)
}
