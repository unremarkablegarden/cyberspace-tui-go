package keymaps

import "github.com/charmbracelet/bubbles/key"

// ProfileKeyMap defines keybindings for the profile view.
type ProfileKeyMap struct {
	Follow      string `json:"follow"`
	EditProfile string `json:"edit_profile"`
}

// NewDefaultProfileKeyMap returns the default profile keybindings.
func NewDefaultProfileKeyMap() ProfileKeyMap {
	return ProfileKeyMap{
		Follow:      "f",
		EditProfile: "e",
	}
}

// ShortHelp returns the short help bindings.
func (k ProfileKeyMap) ShortHelp() []key.Binding {
	// return []key.Binding{k.Up, k.Open, k.Follow, k.Back, k.Help, k.Quit}
	return []key.Binding{}
}

// FullHelp returns the full help bindings grouped in columns.
func (k ProfileKeyMap) FullHelp() [][]key.Binding {
	// return [][]key.Binding{
	// 	{k.Up, k.Down, k.Top, k.Bottom},
	// 	{k.Open, k.Follow, k.EditProfile, k.Refresh},
	// 	{k.Back, k.Help, k.Quit},
	// }
	return [][]key.Binding{}
}
