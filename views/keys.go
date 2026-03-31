package views

import (
	"github.com/charmbracelet/bubbles/key"
)

// ═══════════════════════════════════════════════════════════════════════════════
// FEED KEY MAP
// ═══════════════════════════════════════════════════════════════════════════════

// FeedKeyMap defines keybindings for the feed view.
type FeedKeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Top     key.Binding
	Bottom  key.Binding
	Open    key.Binding
	Refresh key.Binding
	Theme   key.Binding
	Logout        key.Binding
	NewPost       key.Binding
	Bookmarks     key.Binding
	Notifications key.Binding
	Topics        key.Binding
	Profile       key.Binding
	Help          key.Binding
	Quit          key.Binding
}

// NewFeedKeyMap returns the default feed keybindings.
func NewFeedKeyMap() FeedKeyMap {
	return FeedKeyMap{
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k/↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j/↓", "down"),
		),
		Top: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "top"),
		),
		Bottom: key.NewBinding(
			key.WithKeys("G"),
			key.WithHelp("G", "bottom"),
		),
		Open: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "open"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		Theme: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "theme"),
		),
		Logout: key.NewBinding(
			key.WithKeys("L"),
			key.WithHelp("L", "logout"),
		),
		NewPost: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new post"),
		),
		Bookmarks: key.NewBinding(
			key.WithKeys("B"),
			key.WithHelp("B", "bookmarks"),
		),
		Notifications: key.NewBinding(
			key.WithKeys("N"),
			key.WithHelp("N", "notifications"),
		),
		Topics: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "topics"),
		),
		Profile: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "profile"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
	}
}

// ShortHelp returns the short help bindings.
func (k FeedKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Open, k.NewPost, k.Topics, k.Bookmarks, k.Help, k.Quit}
}

// FullHelp returns the full help bindings grouped in columns.
func (k FeedKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Top, k.Bottom},
		{k.Open, k.NewPost, k.Topics, k.Bookmarks, k.Notifications, k.Refresh},
		{k.Profile, k.Theme, k.Logout, k.Help, k.Quit},
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// POST DETAIL KEY MAP
// ═══════════════════════════════════════════════════════════════════════════════

// PostDetailKeyMap defines keybindings for the post detail view.
type PostDetailKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	HalfUp   key.Binding
	HalfDown key.Binding
	Top      key.Binding
	Bottom   key.Binding
	Back     key.Binding
	Refresh  key.Binding
	Reply   key.Binding
	Send    key.Binding
	Save    key.Binding
	Profile key.Binding
	Theme   key.Binding
	Help     key.Binding
	Quit     key.Binding
}

// NewPostDetailKeyMap returns the default post detail keybindings.
func NewPostDetailKeyMap() PostDetailKeyMap {
	return PostDetailKeyMap{
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k/↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j/↓", "down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup"),
			key.WithHelp("pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("f", "pgdown", " "),
			key.WithHelp("f/pgdn", "page down"),
		),
		HalfUp: key.NewBinding(
			key.WithKeys("u", "ctrl+u"),
			key.WithHelp("u", "half up"),
		),
		HalfDown: key.NewBinding(
			key.WithKeys("d", "ctrl+d"),
			key.WithHelp("d", "half down"),
		),
		Top: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "top"),
		),
		Bottom: key.NewBinding(
			key.WithKeys("G"),
			key.WithHelp("G", "bottom"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "backspace", "b"),
			key.WithHelp("esc", "back"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		Reply: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "reply"),
		),
		Send: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "send"),
		),
		Save: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "save"),
		),
		Profile: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "profile"),
		),
		Theme: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "theme"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
	}
}

// ShortHelp returns the short help bindings.
func (k PostDetailKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Back, k.Help, k.Quit}
}

// FullHelp returns the full help bindings grouped in columns.
func (k PostDetailKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.PageUp, k.PageDown, k.HalfUp, k.HalfDown},
		{k.Top, k.Bottom, k.Back, k.Refresh},
		{k.Reply, k.Save, k.Profile, k.Theme, k.Help, k.Quit},
	}
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

// ═══════════════════════════════════════════════════════════════════════════════
// BOOKMARKS KEY MAP
// ═══════════════════════════════════════════════════════════════════════════════

// BookmarksKeyMap defines keybindings for the bookmarks view.
type BookmarksKeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Top     key.Binding
	Bottom  key.Binding
	Open    key.Binding
	Remove  key.Binding
	Refresh key.Binding
	Back    key.Binding
	Theme   key.Binding
	Help    key.Binding
	Quit    key.Binding
}

// NewBookmarksKeyMap returns the default bookmarks keybindings.
func NewBookmarksKeyMap() BookmarksKeyMap {
	return BookmarksKeyMap{
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k/↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j/↓", "down"),
		),
		Top: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "top"),
		),
		Bottom: key.NewBinding(
			key.WithKeys("G"),
			key.WithHelp("G", "bottom"),
		),
		Open: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "open"),
		),
		Remove: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "remove"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "backspace", "b"),
			key.WithHelp("esc", "back"),
		),
		Theme: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "theme"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
	}
}

// ShortHelp returns the short help bindings.
func (k BookmarksKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Open, k.Remove, k.Back, k.Help, k.Quit}
}

// FullHelp returns the full help bindings grouped in columns.
func (k BookmarksKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Top, k.Bottom},
		{k.Open, k.Remove, k.Refresh},
		{k.Theme, k.Back, k.Help, k.Quit},
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// TOPICS KEY MAP
// ═══════════════════════════════════════════════════════════════════════════════

// TopicsKeyMap defines keybindings for the topics browser.
type TopicsKeyMap struct {
	Up   key.Binding
	Down key.Binding
	Open key.Binding
	Back key.Binding
	Help key.Binding
	Quit key.Binding
}

func NewTopicsKeyMap() TopicsKeyMap {
	return TopicsKeyMap{
		Up:   key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k/↑", "up")),
		Down: key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j/↓", "down")),
		Open: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "browse")),
		Back: key.NewBinding(key.WithKeys("esc", "backspace", "b"), key.WithHelp("esc", "back")),
		Help: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Quit: key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
	}
}

func (k TopicsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Open, k.Back, k.Quit}
}

func (k TopicsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Open, k.Back, k.Help, k.Quit}}
}

// ═══════════════════════════════════════════════════════════════════════════════
// TOPIC FEED KEY MAP
// ═══════════════════════════════════════════════════════════════════════════════

// TopicFeedKeyMap defines keybindings for the topic feed.
type TopicFeedKeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Top     key.Binding
	Bottom  key.Binding
	Open    key.Binding
	Profile key.Binding
	Refresh key.Binding
	Back    key.Binding
	Help    key.Binding
	Quit    key.Binding
}

func NewTopicFeedKeyMap() TopicFeedKeyMap {
	return TopicFeedKeyMap{
		Up:      key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k/↑", "up")),
		Down:    key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j/↓", "down")),
		Top:     key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "top")),
		Bottom:  key.NewBinding(key.WithKeys("G"), key.WithHelp("G", "bottom")),
		Open:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "open")),
		Profile: key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "profile")),
		Refresh: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
		Back:    key.NewBinding(key.WithKeys("esc", "backspace", "b"), key.WithHelp("esc", "back")),
		Help:    key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Quit:    key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
	}
}

func (k TopicFeedKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Open, k.Profile, k.Back, k.Help, k.Quit}
}

func (k TopicFeedKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Top, k.Bottom},
		{k.Open, k.Profile, k.Refresh, k.Back},
		{k.Help, k.Quit},
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// PROFILE KEY MAP
// ═══════════════════════════════════════════════════════════════════════════════

// ProfileKeyMap defines keybindings for the profile view.
type ProfileKeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Top     key.Binding
	Bottom  key.Binding
	Open    key.Binding
	Refresh key.Binding
	Back    key.Binding
	Help    key.Binding
	Quit    key.Binding
}

// NewProfileKeyMap returns the default profile keybindings.
func NewProfileKeyMap() ProfileKeyMap {
	return ProfileKeyMap{
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k/↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j/↓", "down"),
		),
		Top: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "top"),
		),
		Bottom: key.NewBinding(
			key.WithKeys("G"),
			key.WithHelp("G", "bottom"),
		),
		Open: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "open post"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "backspace", "b"),
			key.WithHelp("esc", "back"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
	}
}

// ShortHelp returns the short help bindings.
func (k ProfileKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Open, k.Back, k.Help, k.Quit}
}

// FullHelp returns the full help bindings grouped in columns.
func (k ProfileKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Top, k.Bottom},
		{k.Open, k.Refresh, k.Back},
		{k.Help, k.Quit},
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// NOTIFICATIONS KEY MAP
// ═══════════════════════════════════════════════════════════════════════════════

// NotificationsKeyMap defines keybindings for the notifications view.
type NotificationsKeyMap struct {
	Up          key.Binding
	Down        key.Binding
	Open        key.Binding
	MarkAllRead key.Binding
	Refresh     key.Binding
	Back        key.Binding
	Help        key.Binding
	Quit        key.Binding
}

// NewNotificationsKeyMap returns the default notifications keybindings.
func NewNotificationsKeyMap() NotificationsKeyMap {
	return NotificationsKeyMap{
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k/↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j/↓", "down"),
		),
		Open: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "open post"),
		),
		MarkAllRead: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "mark all read"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "backspace", "b"),
			key.WithHelp("esc", "back"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
	}
}

// ShortHelp returns the short help bindings.
func (k NotificationsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Open, k.MarkAllRead, k.Back, k.Help, k.Quit}
}

// FullHelp returns the full help bindings grouped in columns.
func (k NotificationsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Open, k.MarkAllRead, k.Refresh},
		{k.Back, k.Help, k.Quit},
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// THEME SWITCHER KEY MAP
// ═══════════════════════════════════════════════════════════════════════════════

// ThemeSwitcherKeyMap defines keybindings for the theme switcher.
type ThemeSwitcherKeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Apply key.Binding
	Close key.Binding
}

// NewThemeSwitcherKeyMap returns the default theme switcher keybindings.
func NewThemeSwitcherKeyMap() ThemeSwitcherKeyMap {
	return ThemeSwitcherKeyMap{
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k/↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j/↓", "down"),
		),
		Apply: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "apply"),
		),
		Close: key.NewBinding(
			key.WithKeys("esc", "t"),
			key.WithHelp("esc", "close"),
		),
	}
}

// ShortHelp returns the short help bindings.
func (k ThemeSwitcherKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Apply, k.Close}
}

// FullHelp returns the full help bindings grouped in columns.
func (k ThemeSwitcherKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Apply, k.Close},
	}
}
