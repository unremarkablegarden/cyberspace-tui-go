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
	Logout  key.Binding
	Help    key.Binding
	Quit    key.Binding
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
	return []key.Binding{k.Up, k.Open, k.Refresh, k.Theme, k.Help, k.Quit}
}

// FullHelp returns the full help bindings grouped in columns.
func (k FeedKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Top, k.Bottom},
		{k.Open, k.Refresh},
		{k.Theme, k.Logout, k.Help, k.Quit},
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
	Reply    key.Binding
	Send     key.Binding
	Theme    key.Binding
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
			key.WithDisabled(),
		),
		Send: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "send"),
			key.WithDisabled(),
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
		{k.Theme, k.Help, k.Quit},
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
