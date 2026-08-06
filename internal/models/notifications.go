package models

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/external/api"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/messages"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/items"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/keymaps"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/ui"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// NotificationsModel is the notifications screen
type NotificationsModel struct {
	list        list.Model
	loading     bool
	loadingMore bool
	spinner     *spinner.Model
	err         error
	client      *api.Client
	nextCursor  string
	hasMore     bool
	width       int
	height      int
	keys        keymaps.AppKeybinds
	help        help.Model
}

// NewNotificationsModel creates a new notifications screen
func NewNotificationsModel(client *api.Client, keymap keymaps.AppKeybinds, sp *spinner.Model) NotificationsModel {
	delegate := items.NotificationDelegate{}
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.SetShowTitle(false)
	l.SetShowFilter(false)
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	l.Styles = styles.ListStyles()
	l.Paginator.ActiveDot = styles.Bright.Render("▄")
	l.Paginator.InactiveDot = styles.Dark.Render("▄")
	l.KeyMap.Quit.SetEnabled(false)
	l.KeyMap.ForceQuit.SetEnabled(false)

	h := help.New()
	h.Styles = styles.HelpStyles()

	return NotificationsModel{
		list:    l,
		client:  client,
		spinner: sp,
		loading: true,
		hasMore: true,
		keys:    keymap,
		help:    h,
	}
}

func (m NotificationsModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchNotifications())
}

func (m NotificationsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}
		switch msg.String() {
		case m.keys.GlobalKeybinds.Quit:
			return m, tea.Quit
		case m.keys.GlobalKeybinds.Help:
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case m.keys.GlobalKeybinds.Back:
			return m, func() tea.Msg { return messages.SwitchToFeed{} }
		case m.keys.GlobalKeybinds.Refresh:
			m.loading = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.fetchNotifications())
		case m.keys.GlobalKeybinds.Open:
			switch it := m.list.SelectedItem().(type) {
			case items.NotificationItem:
				n := it.Notification
				var cmds []tea.Cmd
				if !n.Read {
					cmds = append(cmds, m.markRead(n.ID))
				}

				if n.Metadata.ReplyID != "" {
					if p, postErr := m.client.FetchPostByReplyID(n.Metadata.ReplyID); postErr == nil {
						cmds = append(cmds, func() tea.Msg {
							return messages.SwitchToPostDetail{
								Post:        *p,
								BackMessage: messages.SwitchToNotifications{},
							}
						})
					}
				}
				return m, tea.Batch(cmds...)
			case items.LoadMoreItem:
				if !m.loadingMore {
					m.loadingMore = true
					return m, tea.Batch(m.spinner.Tick, m.fetchMoreNotifications())
				}
			}
		case m.keys.NotificationsKeybinds.MarkAllRead:
			return m, m.markAllRead()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-4)

	case messages.NotificationsLoadedMsg:
		m.loading = false
		m.loadingMore = false
		m.nextCursor = msg.Cursor
		m.hasMore = msg.Cursor != ""
		m.err = nil

		localItems := items.BuildListItems(
			&m.list, msg.IsAdditional, m.hasMore,
			items.NotificationsToItems(msg.Notifications),
		)
		return m, m.list.SetItems(localItems)

	case messages.NotificationsErrorMsg:
		m.loading = false
		m.loadingMore = false
		m.err = msg.Err

	case ThemeChangedMsg:
		m.spinner.Style = styles.Spinner
		m.list.Styles = styles.ListStyles()
		m.help.Styles = styles.HelpStyles()
		m.list.Paginator.ActiveDot = styles.Bright.Render("▄")
		m.list.Paginator.InactiveDot = styles.Dark.Render("▄")
	}

	if !m.loading {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m NotificationsModel) View() string {
	w, h := ui.SafeDimensions(m.width, m.height)

	if m.loading {
		return ui.RenderLoadingScreen(
			ui.LoadingScreenTexts{
				TitleText:    "SCANNING SIGNALS",
				SubtitleText: " Loading notifications...",
				BottomText:   "Intercepting incoming transmissions...",
			},
			*m.spinner,
			w, h,
		)
	}

	if m.err != nil {
		return ui.RenderErrorScreen(m.err, w, h)
	}

	var b strings.Builder
	b.WriteString(ui.RenderHeader("▓▒░ NOTIFICATIONS ░▒▓", w))
	b.WriteString(m.list.View())
	b.WriteString("\n")
	b.WriteString(
		ui.RenderFooterWithList(
			m.help.View(m.keys.NotificationsHelpKeys()),
			m.list.Paginator.View(),
			w,
		))

	_ = h
	return b.String()
}

// SetSize updates the view dimensions
func (m *NotificationsModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width, height-4)
}

func (m NotificationsModel) fetchNotifications() tea.Cmd {
	return func() tea.Msg {
		notifs, cursor, err := m.client.FetchNotifications(30)
		if err != nil {
			return messages.NotificationsErrorMsg{Err: err}
		}
		return messages.NotificationsLoadedMsg{Notifications: notifs, Cursor: cursor}
	}
}

func (m NotificationsModel) fetchMoreNotifications() tea.Cmd {
	return func() tea.Msg {
		notifs, cursor, err := m.client.FetchMoreNotifications(30, m.nextCursor)
		if err != nil {
			return messages.NotificationsErrorMsg{Err: err}
		}
		return messages.NotificationsLoadedMsg{Notifications: notifs, Cursor: cursor, IsAdditional: true}
	}
}

func (m NotificationsModel) markRead(notificationID string) tea.Cmd {
	return func() tea.Msg {
		_ = m.client.MarkNotificationRead(notificationID)
		return nil
	}
}

func (m NotificationsModel) markAllRead() tea.Cmd {
	return func() tea.Msg {
		_ = m.client.MarkAllNotificationsRead()
		return nil
	}
}

// itoa converts an int to string without importing strconv everywhere
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}
