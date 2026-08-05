package keymaps

import "github.com/charmbracelet/bubbles/key"

// NotesKeyMap defines keybindings for the notes list view.
type NotesKeyMap struct {
	New    string `json:"new"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
}

func NewDefaultNotesKeyMap() NotesKeyMap {
	return NotesKeyMap{
		New:    "n",
		Edit:   "e",
		Delete: "d",
	}
}

func (ak *AppKeybinds) NotesHelpKeys() helpKeybinds {
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
				ToKeybind(ak.NotesKeybinds.New, "new"),
				ToKeybind(ak.NotesKeybinds.Edit, "edit"),
				ToKeybind(ak.NotesKeybinds.Delete, "delete"),
			},
		},
	}
}
