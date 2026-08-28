package settingsmodels

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/messages"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/items"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/keymaps"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/ui"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

type SettingsIndexModel struct {
	list     list.Model
	keybinds *keymaps.AppKeybinds
	help     help.Model
	width    int
	height   int
}

func NewSettingsIndexModel(keymap *keymaps.AppKeybinds) *SettingsIndexModel {
	li := []list.Item{
		items.SettingsIndexItem{Name: "Themes", Field: items.SettingsIndexTheme},
		items.SettingsIndexItem{Name: "Keybinds", Field: items.SettingsIndexKeybind},
	}

	l := list.New(li, items.SettingsIndexDelegate{}, 0, 0)
	items.ConfigList(&l)

	h := help.New()
	h.Styles = styles.HelpStyles()

	return &SettingsIndexModel{
		list:     l,
		keybinds: keymap,
		help:     h,
	}
}

func (sim *SettingsIndexModel) Init() tea.Cmd {
	return nil
}

func (sim *SettingsIndexModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case sim.keybinds.GlobalKeybinds.Quit:
			return sim, tea.Quit
		case sim.keybinds.GlobalKeybinds.Back:
			return sim, func() tea.Msg { return messages.SwitchToMenu{} }
		case sim.keybinds.GlobalKeybinds.Help:
			sim.help.ShowAll = !sim.help.ShowAll
			return sim, nil
		case sim.keybinds.GlobalKeybinds.Open:
			if item := sim.list.SelectedItem(); item != nil {
				var openSetting tea.Msg

				switch item.(items.SettingsIndexItem).Field {
				case items.SettingsIndexTheme:
					openSetting = messages.SwitchToThemeSwitcher{}
				case items.SettingsIndexKeybind:
					openSetting = messages.SwitchToSettings{Setting: uint8(items.SettingsIndexKeybind)}
				}

				return sim, func() tea.Msg { return openSetting }
			}
		}

	case tea.WindowSizeMsg:
		sim.width = msg.Width
		sim.height = msg.Height
		// Reserve space for our custom header (2 lines) and footer (2 lines)
		sim.list.SetSize(msg.Width, msg.Height-4)
	}

	var cmd tea.Cmd
	sim.list, cmd = sim.list.Update(msg)
	return sim, cmd
}

func (sim *SettingsIndexModel) View() string {
	w, _ := ui.SafeDimensions(sim.width, sim.height)

	var b strings.Builder
	b.WriteString(ui.RenderHeader("▓▒░ 𝓢є††เภﻮร ░▒▓", w))
	b.WriteString(sim.list.View())
	b.WriteString("\n")
	b.WriteString(
		ui.RenderFooterWithList(
			sim.help.View(sim.keybinds.BookmarksHelpKeys()), // TEMPORAL
			sim.list.Paginator.View(),
			w,
		))

	return b.String()
}
