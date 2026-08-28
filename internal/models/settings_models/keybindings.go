package settingsmodels

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/messages"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/items"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/keymaps"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/ui"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

type SettingsKeybindsModel struct {
	list      list.Model
	keybinds  *keymaps.AppKeybinds
	help      help.Model
	width     int
	height    int
	isEditing bool
	id        string
}

func NewSettingsKeybindsModel(keymap *keymaps.AppKeybinds) *SettingsKeybindsModel {
	metadataList := keymap.GetKeybindsMap()

	var li []list.Item
	for _, metaItem := range metadataList {
		li = append(li, items.SettingsKeybindsItem{
			ID:    metaItem.ID,
			Name:  metaItem.Name,
			Value: metaItem.Value,
		})
	}

	l := list.New(li, items.SettingsKeybindsDelegate{}, 0, 0)
	items.ConfigList(&l)

	h := help.New()
	h.Styles = styles.HelpStyles()

	return &SettingsKeybindsModel{
		keybinds:  keymap,
		list:      l,
		help:      h,
		isEditing: false,
	}
}

func (skm *SettingsKeybindsModel) Init() tea.Cmd {
	return nil
}

func (skm *SettingsKeybindsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// stuff to handle when isEditing
	if skm.isEditing {
		skm.isEditing = false
		keyMsg, ok := msg.(tea.KeyMsg)
		if !ok {
			return skm, nil
		}

		skm.keybinds.UpdateKeybind(skm.id, keyMsg.String())

		return skm, func() tea.Msg {
			return messages.SaveKeymaps{}
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case skm.keybinds.GlobalKeybinds.Quit:
			return skm, tea.Quit
		case skm.keybinds.GlobalKeybinds.Back:
			return skm, func() tea.Msg { return messages.SwitchToSettings{} }
		case skm.keybinds.GlobalKeybinds.Open:
			if item := skm.list.SelectedItem(); item != nil {
				id := item.(items.SettingsKeybindsItem).ID
				skm.isEditing = true
				skm.id = id
			}
		}
	case tea.WindowSizeMsg:
		skm.width = msg.Width
		skm.height = msg.Height
		// Reserve space for our custom header (2 lines) and footer (2 lines)
		skm.list.SetSize(msg.Width, msg.Height-4)
	}

	var cmd tea.Cmd
	skm.list, cmd = skm.list.Update(msg)
	return skm, cmd
}

func (skm *SettingsKeybindsModel) View() string {
	w, h := ui.SafeDimensions(skm.width, skm.height)

	var b strings.Builder

	if skm.isEditing {
		b.WriteString("Press a key")
		boxWidth := max(min(60, w-6), 40)
		box := styles.TitledBox("CHANGE THE KEYBIND", b.String(), boxWidth)

		return ui.FullScreen(box, w, h, lipgloss.Center, lipgloss.Center)
	}

	b.WriteString(ui.RenderHeader("▓▒░ 𝓢є††เภﻮร ░▒▓", w))
	b.WriteString(skm.list.View())
	b.WriteString("\n")
	return b.String()

}
