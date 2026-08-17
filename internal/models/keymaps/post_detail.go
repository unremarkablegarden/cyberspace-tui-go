package keymaps

import "github.com/charmbracelet/bubbles/key"

// PostDetailKeyMap defines keybindings for the post detail view.
type PostDetailKeyMap struct {
	Reply   string `json:"reply"`
	Send    string `json:"send"`
	Save    string `json:"save"`
	Delete  string `json:"delete"`
	Profile string `json:"profile"`
}

// NewDefaultPostDetailKeyMap returns the default post detail keybindings.
func NewDefaultPostDetailKeyMap() PostDetailKeyMap {
	return PostDetailKeyMap{
		Reply:   "c",
		Send:    "ctrl+s",
		Save:    "s",
		Delete:  "D",
		Profile: "p",
	}
}

func (ak *AppKeybinds) PostDetailHelpKeys() helpKeybinds {
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
				ToKeybind(ak.GlobalKeybinds.Refresh, "refresh"),
				ToKeybind(ak.GlobalKeybinds.Help, "help"),
			},
			{
				ToKeybind("[CTRL+K]", "menu"),
				ToKeybind(ak.PostDetailKeybinds.Reply, "reply"),
				ToKeybind(ak.PostDetailKeybinds.Send, "send"),
				ToKeybind(ak.PostDetailKeybinds.Save, "save"),
			},
			{
				ToKeybind(ak.PostDetailKeybinds.Delete, "delete"),
				ToKeybind(ak.PostDetailKeybinds.Profile, "profile"),
			},
		},
	}
}
