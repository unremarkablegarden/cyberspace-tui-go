package keymaps

import "github.com/charmbracelet/bubbles/key"

// NotesKeyMap defines keybindings for the notes list view.
type NotesKeyMap struct {
	New    string `json:"new"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
}

func (nkm *NotesKeyMap) GetMetadata() []KeybindMetadata {
	return []KeybindMetadata{
		{
			ID:    "notes.new",
			Name:  "Notes - New",
			Value: nkm.New,
		},
		{
			ID:    "notes.edit",
			Name:  "Notes - Edit",
			Value: nkm.Edit,
		},
		{
			ID:    "notes.delete",
			Name:  "Notes - Delete",
			Value: nkm.Delete,
		},
	}
}

func (nkm *NotesKeyMap) Update(field string, value string) {
	switch field {
	case "new":
		nkm.New = value
	case "edit":
		nkm.Edit = value
	case "delete":
		nkm.Delete = value
	}
}

func NewDefaultNotesKeyMap() *NotesKeyMap {
	return &NotesKeyMap{
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
