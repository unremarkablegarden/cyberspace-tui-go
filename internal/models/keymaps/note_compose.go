package keymaps

import "github.com/charmbracelet/bubbles/key"

// NoteComposeKeyMap defines keybindings for the note compose/edit screen.
type NoteComposeKeyMap struct {
	SwitchField string `json:"switch_field"`
	Save        string `json:"save"`
	Cancel      string `json:"cancel"`
}

func NewDefaultNoteComposeKeyMap() NoteComposeKeyMap {
	return NoteComposeKeyMap{
		SwitchField: "tab",
		Save:        "ctrl+s",
		Cancel:      "esc",
	}
}

func (ak *AppKeybinds) NoteComposeHelpKeys() helpKeybinds {
	return helpKeybinds{
		short: []key.Binding{
			ToKeybind("[CTRL+K]", "menu"),
			ToKeybind(ak.NoteComposeKeybinds.SwitchField, "switch field"),
			ToKeybind(ak.NoteComposeKeybinds.Save, "save"),
			ToKeybind(ak.NoteComposeKeybinds.Cancel, "cancel"),
		},
		full: [][]key.Binding{
			{
				ToKeybind("[CTRL+K]", "menu"),
				ToKeybind(ak.NoteComposeKeybinds.SwitchField, "switch field"),
				ToKeybind(ak.NoteComposeKeybinds.Save, "save"),
				ToKeybind(ak.NoteComposeKeybinds.Cancel, "cancel"),
			},
		},
	}
}
