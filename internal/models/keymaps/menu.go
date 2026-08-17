package keymaps

import "github.com/charmbracelet/bubbles/key"

type MenuKeymaps struct {
	SectionFeed          string `json:"section_feed"`
	SectionCompose       string `json:"section_compose"`
	SectionNotifications string `json:"section_notifications"`
	SectionOwnProfile    string `json:"section_own_profile"`
	SectionNotes         string `json:"section_notes"`
	SectionTopics        string `json:"section_topics"`
	SectionBookmarks     string `json:"section_bookmarks"`
	SectionSettings      string `json:"section_settings"`
}

func NewDefaultMenuKeyMap() MenuKeymaps {
	return MenuKeymaps{
		SectionFeed:          "1",
		SectionCompose:       "2",
		SectionNotifications: "3",
		SectionOwnProfile:    "4",
		SectionNotes:         "5",
		SectionTopics:        "6",
		SectionBookmarks:     "7",
		SectionSettings:      "8",
	}
}

func (ak *AppKeybinds) MenuHelpKeys() helpKeybinds {
	return helpKeybinds{
		short: []key.Binding{
			ToKeybind(ak.GlobalKeybinds.Up, "up"),
			ToKeybind(ak.GlobalKeybinds.Down, "down"),
			ToKeybind(ak.GlobalKeybinds.Help, "help"),
			ToKeybind(ak.GlobalKeybinds.Quit, "quit"),
		},
		full: [][]key.Binding{
			{
				ToKeybind(ak.GlobalKeybinds.Up, "up"),
				ToKeybind(ak.GlobalKeybinds.Open, "open"),
				ToKeybind(ak.GlobalKeybinds.Quit, "quit"),
			},
			{
				ToKeybind(ak.GlobalKeybinds.Down, "down"),
				ToKeybind(ak.GlobalKeybinds.Help, "help"),
			},
			{
				ToKeybind(ak.MenuKeybinds.SectionFeed, "feed"),
				ToKeybind(ak.MenuKeybinds.SectionCompose, "compose"),
				ToKeybind(ak.MenuKeybinds.SectionNotifications, "notifications"),
				ToKeybind(ak.MenuKeybinds.SectionOwnProfile, "own profile"),
			},
			{
				ToKeybind(ak.MenuKeybinds.SectionNotes, "notes"),
				ToKeybind(ak.MenuKeybinds.SectionTopics, "topics"),
				ToKeybind(ak.MenuKeybinds.SectionBookmarks, "bookmarks"),
				ToKeybind(ak.MenuKeybinds.SectionSettings, "settings"),
			},
		},
	}
}
