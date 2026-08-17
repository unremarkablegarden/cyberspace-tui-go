package keymaps

import "github.com/charmbracelet/bubbles/key"

// BookmarksKeyMap defines keybindings for the bookmarks view.
type BookmarksKeyMap struct {
	Remove string `json:"remove"`
}

// NewDefaultBookmarksKeyMap returns the default bookmarks keybindings.
func NewDefaultBookmarksKeyMap() BookmarksKeyMap {
	return BookmarksKeyMap{
		Remove: "d",
	}
}

func (ak *AppKeybinds) BookmarksHelpKeys() helpKeybinds {
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
				ToKeybind(ak.BookmarksKeybinds.Remove, "remove"),
			},
		},
	}
}
