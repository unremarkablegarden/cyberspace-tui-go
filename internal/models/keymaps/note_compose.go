package keymaps

import "github.com/charmbracelet/bubbles/key"

// NoteComposeKeyMap defines keybindings for the note compose/edit screen.
type NoteComposeKeyMap struct {
	SwitchField string `json:"switch_field"`
	Save        string `json:"save"`
	Cancel      string `json:"cancel"`
}

func (nckm *NoteComposeKeyMap) GetMetadata() []KeybindMetadata {
	return []KeybindMetadata{
		{
			ID:    "note_compose.switch_field",
			Name:  "Note Compose - Switch Field",
			Value: nckm.SwitchField,
		},
		{
			ID:    "note_compose.save",
			Name:  "Note Compose - Save",
			Value: nckm.Save,
		},
		{
			ID:    "note_compose.cancel",
			Name:  "Note Compose - Cancel",
			Value: nckm.Cancel,
		},
	}
}
func (nckm *NoteComposeKeyMap) Update(field string, value string) {
	switch field {
	case "switch_field":
		nckm.SwitchField = value
	case "save":
		nckm.Save = value
	case "cancel":
		nckm.Cancel = value
	}
}

func NewDefaultNoteComposeKeyMap() *NoteComposeKeyMap {
	return &NoteComposeKeyMap{
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
