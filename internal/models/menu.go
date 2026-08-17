package models

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/messages"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/items"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/keymaps"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/ui"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

type MenuModel struct {
	list     list.Model
	keybinds keymaps.AppKeybinds
	help     help.Model
	width    int
	height   int
}

func NewMenuModel(keymap keymaps.AppKeybinds) MenuModel {
	li := []list.Item{
		items.MenuItem{
			Name: "Feed", Field: items.MenuSectionFeed,
			Keybind: keymaps.ToKeybind(keymap.MenuKeybinds.SectionFeed, "Feed"),
		},
		items.MenuItem{
			Name: "Write", Field: items.MenuSectionCompose,
			Keybind: keymaps.ToKeybind(keymap.MenuKeybinds.SectionCompose, "Write"),
		},
		items.MenuItem{
			Name: "Notifications", Field: items.MenuSectionNotifications,
			Keybind: keymaps.ToKeybind(keymap.MenuKeybinds.SectionNotifications, "Notifications"),
		},
		items.MenuItem{
			Name: "Profile", Field: items.MenuSectionOwnProfile,
			Keybind: keymaps.ToKeybind(keymap.MenuKeybinds.SectionOwnProfile, "Profile"),
		},
		items.MenuItem{
			Name: "Journal", Field: items.MenuSectionNotes,
			Keybind: keymaps.ToKeybind(keymap.MenuKeybinds.SectionNotes, "Journal"),
		},
		items.MenuItem{
			Name: "Topics", Field: items.MenuSectionTopics,
			Keybind: keymaps.ToKeybind(keymap.MenuKeybinds.SectionTopics, "Topics"),
		},
		items.MenuItem{
			Name: "Bookmarks", Field: items.MenuSectionBookmarks,
			Keybind: keymaps.ToKeybind(keymap.MenuKeybinds.SectionBookmarks, "Bookmarks"),
		},
		items.MenuItem{
			Name: "Settings", Field: items.MenuSectionSettings,
			Keybind: keymaps.ToKeybind(keymap.MenuKeybinds.SectionSettings, "Settings"),
		},
	}

	l := list.New(li, items.MenuDelegate{}, 0, 0)
	l.SetShowTitle(false)
	l.SetShowFilter(false)
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	l.Styles = styles.ListStyles()

	h := help.New()
	h.Styles = styles.HelpStyles()

	return MenuModel{
		list:     l,
		keybinds: keymap,
		help:     h,
	}
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() == m.keybinds.GlobalKeybinds.Open:
			if item := m.list.SelectedItem(); item != nil {
				var view tea.Msg
				switch item.(items.MenuItem).Field {
				case items.MenuSectionFeed:
					view = messages.SwitchToFeed{}
				case items.MenuSectionCompose:
					view = messages.SwitchToCompose{}
				case items.MenuSectionNotifications:
					view = messages.SwitchToNotifications{}
				case items.MenuSectionOwnProfile:
					view = messages.SwitchToOwnProfile{}
				case items.MenuSectionNotes:
					view = messages.SwitchToNotes{}
				case items.MenuSectionTopics:
					view = messages.SwitchToTopics{}
				case items.MenuSectionBookmarks:
					view = messages.SwitchToBookmarks{}
				case items.MenuSectionSettings:
					view = messages.SwitchToThemeSwitcher{} // Tmp setting
				default:
					view = messages.SwitchToFeed{}
				}

				return m, func() tea.Msg { return view }
			}

			// Shorcuts
		case key.Matches(msg, keymaps.ToKeybind(m.keybinds.MenuKeybinds.SectionFeed, "Feed")):
			return m, func() tea.Msg { return messages.SwitchToFeed{} }
		case key.Matches(msg, keymaps.ToKeybind(m.keybinds.MenuKeybinds.SectionCompose, "Write")):
			return m, func() tea.Msg { return messages.SwitchToCompose{} }
		case key.Matches(msg, keymaps.ToKeybind(m.keybinds.MenuKeybinds.SectionNotifications, "Notifications")):
			return m, func() tea.Msg { return messages.SwitchToNotifications{} }
		case key.Matches(msg, keymaps.ToKeybind(m.keybinds.MenuKeybinds.SectionOwnProfile, "Profile")):
			return m, func() tea.Msg { return messages.SwitchToOwnProfile{} }
		case key.Matches(msg, keymaps.ToKeybind(m.keybinds.MenuKeybinds.SectionNotes, "Journal")):
			return m, func() tea.Msg { return messages.SwitchToNotes{} }
		case key.Matches(msg, keymaps.ToKeybind(m.keybinds.MenuKeybinds.SectionTopics, "Topics")):
			return m, func() tea.Msg { return messages.SwitchToTopics{} }
		case key.Matches(msg, keymaps.ToKeybind(m.keybinds.MenuKeybinds.SectionBookmarks, "Bookmarks")):
			return m, func() tea.Msg { return messages.SwitchToBookmarks{} }
		case key.Matches(msg, keymaps.ToKeybind(m.keybinds.MenuKeybinds.SectionSettings, "Settings")):
			return m, func() tea.Msg { return messages.SwitchToThemeSwitcher{} }

			// Help
		case key.Matches(msg, keymaps.ToKeybind(m.keybinds.GlobalKeybinds.Help, "?")):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Reserve space for our custom header (2 lines) and footer (2 lines)
		m.list.SetSize(msg.Width, msg.Height-4)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m MenuModel) View() string {
	w, _ := ui.SafeDimensions(m.width, m.height)

	var b strings.Builder

	// Header: centered title with blocks on each side
	b.WriteString(ui.RenderHeader("▓▒░ ᗰєภย ░▒▓", w))

	// List content (title disabled, we render our own header)
	b.WriteString(m.list.View())

	// Footer: divider with paginator inline on the right
	b.WriteString("\n")
	b.WriteString(
		ui.RenderFooterWithList(
			m.help.View(m.keybinds.MenuHelpKeys()),
			m.list.Paginator.View(),
			w,
		))

	return b.String()
}
