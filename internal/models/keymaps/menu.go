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

func (mkm *MenuKeymaps) GetMetadata() []KeybindMetadata {
	return []KeybindMetadata{
		{
			ID:    "menu.feed",
			Name:  "Menu - Feed",
			Value: mkm.SectionFeed,
		},
		{
			ID:    "menu.compose",
			Name:  "Menu - Compose",
			Value: mkm.SectionCompose,
		},
		{
			ID:    "menu.notifications",
			Name:  "Menu - Notifications",
			Value: mkm.SectionNotifications,
		},
		{
			ID:    "menu.profile",
			Name:  "Menu - Profile",
			Value: mkm.SectionOwnProfile,
		},
		{
			ID:    "menu.notes",
			Name:  "Menu - Notes",
			Value: mkm.SectionNotes,
		},
		{
			ID:    "menu.topics",
			Name:  "Menu - Topics",
			Value: mkm.SectionTopics,
		},
		{
			ID:    "menu.bookmarks",
			Name:  "Menu - Bookmarks",
			Value: mkm.SectionBookmarks,
		},
		{
			ID:    "menu.settings",
			Name:  "Menu - Settings",
			Value: mkm.SectionSettings,
		},
	}
}
func (mkm *MenuKeymaps) Update(field string, value string) {
	switch field {
	case "feed":
		mkm.SectionFeed = value
	case "compose":
		mkm.SectionCompose = value
	case "notifications":
		mkm.SectionNotifications = value
	case "profile":
		mkm.SectionOwnProfile = value
	case "notes":
		mkm.SectionNotes = value
	case "topics":
		mkm.SectionTopics = value
	case "bookmarks":
		mkm.SectionBookmarks = value
	case "settings":
		mkm.SectionSettings = value
	}
}

func NewDefaultMenuKeyMap() *MenuKeymaps {
	return &MenuKeymaps{
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
