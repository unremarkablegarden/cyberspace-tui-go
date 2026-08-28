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

func (pdkm *PostDetailKeyMap) GetMetadata() []KeybindMetadata {
	return []KeybindMetadata{
		{
			ID:    "post_detail.reply",
			Name:  "Post Detail - Reply",
			Value: pdkm.Reply,
		},
		{
			ID:    "post_detail.send",
			Name:  "Post Detail - Send",
			Value: pdkm.Send,
		},
		{
			ID:    "post_detail.save",
			Name:  "Post Detail - Save",
			Value: pdkm.Save,
		},
		{
			ID:    "post_detail.delete",
			Name:  "Post Detail - Delete",
			Value: pdkm.Delete,
		},
		{
			ID:    "post_detail.profile",
			Name:  "Post Detail - Profile",
			Value: pdkm.Profile,
		},
	}
}

func (pdkm *PostDetailKeyMap) Update(field string, value string) {
	switch field {
	case "reply":
		pdkm.Reply = value
	case "send":
		pdkm.Send = value
	case "save":
		pdkm.Save = value
	case "delete":
		pdkm.Delete = value
	case "profile":
		pdkm.Profile = value
	}
}

// NewDefaultPostDetailKeyMap returns the default post detail keybindings.
func NewDefaultPostDetailKeyMap() *PostDetailKeyMap {
	return &PostDetailKeyMap{
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
