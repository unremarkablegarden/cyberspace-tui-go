package keymaps

import "github.com/charmbracelet/bubbles/key"

// TopicFeedKeyMap defines keybindings for the topic feed.
type TopicFeedKeyMap struct {
	Profile string `json:"profile"`
}

func NewDefaultTopicFeedKeyMap() TopicFeedKeyMap {
	return TopicFeedKeyMap{
		Profile: "p",
	}
}

func (ak *AppKeybinds) TopicFeedHelpKeys() helpKeybinds {
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
				ToKeybind(ak.FeedKeybinds.Profile, "profile"),
				ToKeybind(ak.FeedKeybinds.Logout, "logout"),
			},
		},
	}
}
