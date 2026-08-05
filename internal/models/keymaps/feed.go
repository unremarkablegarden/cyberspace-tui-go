package keymaps

import "github.com/charmbracelet/bubbles/key"

// FeedKeyMap defines keybindings for the feed view.
type FeedKeyMap struct {
	Logout  string `json:"logout"`
	Profile string `json:"profile"`
}

// NewDefaultFeedKeyMap returns the default feed keybindings.
func NewDefaultFeedKeyMap() FeedKeyMap {
	return FeedKeyMap{
		Logout:  "L",
		Profile: "p",
	}
}

func (ak *AppKeybinds) FeedHelpKeys() helpKeybinds {
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
			},
			{
				ToKeybind(ak.GlobalKeybinds.Down, "down"),
				ToKeybind(ak.GlobalKeybinds.Refresh, "refresh"),
				ToKeybind(ak.GlobalKeybinds.Help, "help"),
			},
			{
				ToKeybind("[CTRL+K]", "menu"),
				ToKeybind(ak.FeedKeybinds.Profile, "profile"),
				ToKeybind(ak.FeedKeybinds.Logout, "logout"),
			},
		},
	}
}
