package keymaps

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
)

type AppKeybinds struct {
	BookmarksKeybinds     *BookmarksKeyMap     `json:"bookmarks"`
	FeedKeybinds          *FeedKeyMap          `json:"feed"`
	GlobalKeybinds        *GlobalKeyMaps       `json:"global"`
	MenuKeybinds          *MenuKeymaps         `json:"menu"`
	NotesKeybinds         *NotesKeyMap         `json:"notes"`
	NoteComposeKeybinds   *NoteComposeKeyMap   `json:"note_compose"`
	NotificationsKeybinds *NotificationsKeyMap `json:"notifications"`
	PostDetailKeybinds    *PostDetailKeyMap    `json:"post_detail"`
	ProfileKeybinds       *ProfileKeyMap       `json:"profile"`
	ThemeSwitcherKeybinds *ThemeSwitcherKeyMap `json:"theme_switcher"`
	TopicsKeybinds        *TopicsKeyMap        `json:"topics"`
	TopicsFeedKeybinds    *TopicFeedKeyMap     `json:"topics_feed"`
}

type KeybindMetadata struct {
	ID    string
	Name  string
	Value string
}

func NewDefaultAppKeymaps() *AppKeybinds {
	return &AppKeybinds{
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

func (ak *AppKeybinds) GetKeybindsMap() []KeybindMetadata {
	var list []KeybindMetadata

	list = append(list, ak.GlobalKeybinds.GetMetadata()...)
	list = append(list, ak.FeedKeybinds.GetMetadata()...)
	list = append(list, ak.PostDetailKeybinds.GetMetadata()...)
	list = append(list, ak.MenuKeybinds.GetMetadata()...)
	list = append(list, ak.BookmarksKeybinds.GetMetadata()...)
	list = append(list, ak.NotesKeybinds.GetMetadata()...)
	list = append(list, ak.NoteComposeKeybinds.GetMetadata()...)
	list = append(list, ak.NotificationsKeybinds.GetMetadata()...)
	list = append(list, ak.ProfileKeybinds.GetMetadata()...)
	list = append(list, ak.TopicsFeedKeybinds.GetMetadata()...)

	return list
}

func (ak *AppKeybinds) UpdateKeybind(id string, value string) {
	splitted := strings.Split(id, ".")

	// not valid format for "parent.field"
	if len(splitted) != 2 {
		return
	}

	switch splitted[0] {
	case "global":
		ak.GlobalKeybinds.Update(splitted[1], value)
	case "bookmarks":
		ak.BookmarksKeybinds.Update(splitted[1], value)
	case "feed":
		ak.FeedKeybinds.Update(splitted[1], value)
	}
}

type helpKeybinds struct {
	short []key.Binding
	full  [][]key.Binding
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
