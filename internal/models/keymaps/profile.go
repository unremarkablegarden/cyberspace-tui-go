package keymaps

import "github.com/charmbracelet/bubbles/key"

// ProfileKeyMap defines keybindings for the profile view.
type ProfileKeyMap struct {
	Follow      string `json:"follow"`
	EditProfile string `json:"edit_profile"`
}

func (pkm *ProfileKeyMap) GetMetadata() []KeybindMetadata {
	return []KeybindMetadata{
		{
			ID:    "profile.follow",
			Name:  "Profile - Follow",
			Value: pkm.Follow,
		},
		{
			ID:    "profile.edit_profile",
			Name:  "Profile - Edit Profile",
			Value: pkm.EditProfile,
		},
	}
}

func (pkm *ProfileKeyMap) Update(field string, value string) {
	switch field {
	case "follow":
		pkm.Follow = value
	case "edit_profile":
		pkm.EditProfile = value
	}
}

// NewDefaultProfileKeyMap returns the default profile keybindings.
func NewDefaultProfileKeyMap() *ProfileKeyMap {
	return &ProfileKeyMap{
		Follow:      "f",
		EditProfile: "e",
	}
}

func (ak *AppKeybinds) ProfileHelpKeys() helpKeybinds {
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
				ToKeybind(ak.GlobalKeybinds.Help, "help"),
				ToKeybind("[CTRL+K]", "menu"),
			},
			{
				ToKeybind(ak.ProfileKeybinds.EditProfile, "edit_profile"),
				ToKeybind(ak.ProfileKeybinds.Follow, "follow"),
			},
		},
	}
}
