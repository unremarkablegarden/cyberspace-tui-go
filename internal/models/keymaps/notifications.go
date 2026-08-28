package keymaps

import "github.com/charmbracelet/bubbles/key"

// NotificationsKeyMap defines keybindings for the notifications view.
type NotificationsKeyMap struct {
	MarkAllRead string `json:"mark_all_read"`
}

func (nkm *NotificationsKeyMap) GetMetadata() []KeybindMetadata {
	return []KeybindMetadata{
		{
			ID:    "notifications.mark_all_read",
			Name:  "Notifications - Mark All Read",
			Value: nkm.MarkAllRead,
		},
	}
}

func (nkm *NotificationsKeyMap) Update(field string, value string) {
	switch field {
	case "mark_all_read":
		nkm.MarkAllRead = value
	}
}

// NewDefaultNotificationsKeyMap returns the default notifications keybindings.
func NewDefaultNotificationsKeyMap() *NotificationsKeyMap {
	return &NotificationsKeyMap{
		MarkAllRead: "a",
	}
}

func (ak *AppKeybinds) NotificationsHelpKeys() helpKeybinds {
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
				ToKeybind(ak.NotificationsKeybinds.MarkAllRead, "mark all read"),
			},
		},
	}
}
