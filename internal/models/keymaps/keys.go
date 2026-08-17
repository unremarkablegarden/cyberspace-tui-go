package keymaps

import (
	"github.com/charmbracelet/bubbles/key"
)

type AppKeybinds struct {
	BookmarksKeybinds     BookmarksKeyMap     `json:"bookmarks"`
	FeedKeybinds          FeedKeyMap          `json:"feed"`
	GlobalKeybinds        GlobalKeyMaps       `json:"global"`
	MenuKeybinds          MenuKeymaps         `json:"menu"`
	NotesKeybinds         NotesKeyMap         `json:"notes"`
	NoteComposeKeybinds   NoteComposeKeyMap   `json:"note_compose"`
	NotificationsKeybinds NotificationsKeyMap `json:"notifications"`
	PostDetailKeybinds    PostDetailKeyMap    `json:"post_detail"`
	ProfileKeybinds       ProfileKeyMap       `json:"profile"`
	ThemeSwitcherKeybinds ThemeSwitcherKeyMap `json:"theme_switcher"`
	TopicsKeybinds        TopicsKeyMap        `json:"topics"`
	TopicsFeedKeybinds    TopicFeedKeyMap     `json:"topics_feed"`
}

type helpKeybinds struct {
	short []key.Binding
	full  [][]key.Binding
}

func NewDefaultAppKeymaps() AppKeybinds {
	return AppKeybinds{
		BookmarksKeybinds:     NewDefaultBookmarksKeyMap(),
		FeedKeybinds:          NewDefaultFeedKeyMap(),
		GlobalKeybinds:        NewDefaultGlobalKeyMaps(),
		MenuKeybinds:          NewDefaultMenuKeyMap(),
		NotesKeybinds:         NewDefaultNotesKeyMap(),
		NoteComposeKeybinds:   NewDefaultNoteComposeKeyMap(),
		NotificationsKeybinds: NewDefaultNotificationsKeyMap(),
		PostDetailKeybinds:    NewDefaultPostDetailKeyMap(),
		ProfileKeybinds:       NewDefaultProfileKeyMap(),
		ThemeSwitcherKeybinds: NewDefaultThemeSwitcherKeyMap(),
		TopicsKeybinds:        NewDefaultTopicsKeyMap(),
		TopicsFeedKeybinds:    NewDefaultTopicFeedKeyMap(),
	}

}

func (hk helpKeybinds) ShortHelp() []key.Binding {
	return hk.short
}

func (hk helpKeybinds) FullHelp() [][]key.Binding {
	return hk.full
}

func ToKeybind(keyS string, helpText string) key.Binding {
	return key.NewBinding(
		key.WithKeys(keyS),
		key.WithHelp(keyS, helpText),
	)
}

// ═══════════════════════════════════════════════════════════════════════════════
// LOGIN KEY MAP
// ═══════════════════════════════════════════════════════════════════════════════

// LoginKeyMap defines keybindings for the login view.
type LoginKeyMap struct {
	NextField key.Binding
	PrevField key.Binding
	Submit    key.Binding
}

// NewLoginKeyMap returns the default login keybindings.
func NewLoginKeyMap() LoginKeyMap {
	return LoginKeyMap{
		NextField: key.NewBinding(
			key.WithKeys("tab", "down"),
			key.WithHelp("tab", "next field"),
		),
		PrevField: key.NewBinding(
			key.WithKeys("shift+tab", "up"),
			key.WithHelp("shift+tab", "prev field"),
		),
		Submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "connect"),
		),
	}
}

// ShortHelp returns the short help bindings.
func (k LoginKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NextField, k.Submit}
}

// FullHelp returns the full help bindings grouped in columns.
func (k LoginKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextField, k.PrevField, k.Submit},
	}
}
