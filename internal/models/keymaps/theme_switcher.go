package keymaps

import "github.com/charmbracelet/bubbles/key"

// ThemeSwitcherKeyMap defines keybindings for the theme switcher.
type ThemeSwitcherKeyMap struct{}

// NewDefaultThemeSwitcherKeyMap returns the default theme switcher keybindings.
func NewDefaultThemeSwitcherKeyMap() ThemeSwitcherKeyMap {
	return ThemeSwitcherKeyMap{}
}

func (ak *AppKeybinds) ThemeSwitcherHelpKeys() helpKeybinds {
	return helpKeybinds{
		short: []key.Binding{
			ToKeybind(ak.GlobalKeybinds.Up, "up"),
			ToKeybind(ak.GlobalKeybinds.Down, "down"),
			ToKeybind("[CTRL+K]", "menu"),
			ToKeybind(ak.GlobalKeybinds.Help, "help"),
			ToKeybind(ak.GlobalKeybinds.Quit, "quit"),
		},
		full: [][]key.Binding{
			{
				ToKeybind(ak.GlobalKeybinds.Up, "up"),
				ToKeybind(ak.GlobalKeybinds.Open, "open"),
				ToKeybind(ak.GlobalKeybinds.Quit, "quit"),
				ToKeybind(ak.GlobalKeybinds.Back, "back"),
			},
			{
				ToKeybind(ak.GlobalKeybinds.Down, "down"),
				ToKeybind(ak.GlobalKeybinds.Help, "help"),
				ToKeybind("[CTRL+K]", "menu"),
			},
		},
	}
}
