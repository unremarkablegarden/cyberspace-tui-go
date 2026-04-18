package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/unremarkablegarden/cyberspace-tui-go/api"
	"github.com/unremarkablegarden/cyberspace-tui-go/models"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// NotificationsLoadedMsg is sent when notifications are fetched
type NotificationsLoadedMsg struct {
	Notifications []models.Notification
	Cursor        string
}

// MoreNotificationsLoadedMsg is sent when more notifications are loaded
type MoreNotificationsLoadedMsg struct {
	Notifications []models.Notification
	Cursor        string
}

// NotificationsErrorMsg is sent when fetching notifications fails
type NotificationsErrorMsg struct{ Err error }

// OpenPostFromNotificationMsg is sent when opening a post from a notification
type OpenPostFromNotificationMsg struct{ PostID string }

// BackFromNotificationsMsg is sent when navigating back from notifications
type BackFromNotificationsMsg struct{}

// NotificationsModel is the notifications screen
type NotificationsModel struct {
	list        list.Model
	loading     bool
	loadingMore bool
	spinner     spinner.Model
	err         error
	client      *api.Client
	nextCursor  string
	hasMore     bool
	width       int
	height      int
	keys        NotificationsKeyMap
	help        help.Model
}

// NewNotificationsModel creates a new notifications screen
func NewNotificationsModel(baseURL, idToken string) NotificationsModel {
	delegate := NotificationDelegate{}
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
		client:  api.NewClient(baseURL, idToken),
		spinner: NewSpinner(),
		loading: true,
		hasMore: true,
		keys:    NewNotificationsKeyMap(),
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
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case key.Matches(msg, m.keys.Back):
			return m, func() tea.Msg { return BackFromNotificationsMsg{} }
		case key.Matches(msg, m.keys.Refresh):
			m.loading = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.fetchNotifications())
		case key.Matches(msg, m.keys.MarkAllRead):
			return m, m.markAllRead()
		case key.Matches(msg, m.keys.Open):
			switch it := m.list.SelectedItem().(type) {
			case NotificationItem:
				n := it.Notification
				var cmds []tea.Cmd
				if !n.Read {
					cmds = append(cmds, m.markRead(n.ID))
				}
				if n.PostID != "" {
					postID := n.PostID
					cmds = append(cmds, func() tea.Msg {
						return OpenPostFromNotificationMsg{PostID: postID}
					})
				}
				return m, tea.Batch(cmds...)
			case LoadMoreItem:
				if !m.loadingMore {
					m.loadingMore = true
					return m, tea.Batch(m.spinner.Tick, m.fetchMoreNotifications())
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-4)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case NotificationsLoadedMsg:
		m.loading = false
		m.err = nil
		items := notificationsToItems(msg.Notifications)
		m.nextCursor = msg.Cursor
		m.hasMore = msg.Cursor != ""
		if m.hasMore {
			items = append(items, LoadMoreItem{})
		}
		cmd := m.list.SetItems(items)
		return m, cmd

	case MoreNotificationsLoadedMsg:
		m.loadingMore = false
		m.nextCursor = msg.Cursor
		m.hasMore = msg.Cursor != ""
		var items []list.Item
		for _, existing := range m.list.Items() {
			if _, ok := existing.(LoadMoreItem); ok {
				continue
			}
			items = append(items, existing)
		}
		for _, n := range msg.Notifications {
			items = append(items, NotificationItem{Notification: n})
		}
		if m.hasMore {
			items = append(items, LoadMoreItem{})
		}
		cmd := m.list.SetItems(items)
		return m, cmd

	case NotificationsErrorMsg:
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
	w, h := SafeDimensions(m.width, m.height)

	if m.loading {
		return m.renderLoadingScreen(w, h)
	}

	if m.err != nil {
		return m.renderErrorScreen(w, h)
	}

	var b strings.Builder
	b.WriteString(m.renderHeader(w))
	b.WriteString(m.list.View())
	b.WriteString("\n")
	b.WriteString(m.renderFooter(w))

	_ = h
	return b.String()
}

func (m NotificationsModel) renderHeader(width int) string {
	unread := 0
	for _, item := range m.list.Items() {
		if ni, ok := item.(NotificationItem); ok && !ni.Notification.Read {
			unread++
		}
	}
	title := "▓▒░ NOTIFICATIONS ░▒▓"
	if unread > 0 {
		title = "▓▒░ NOTIFICATIONS ░▒▓" + styles.Bright.Render(" ["+itoa(unread)+" unread]")
	}
	return RenderHeader(title, width)
}

func (m NotificationsModel) renderFooter(width int) string {
	helpView := m.help.View(m.keys)
	paginatorView := m.list.Paginator.View()

	helpWidth := lipgloss.Width(helpView)
	paginatorWidth := lipgloss.Width(paginatorView)

	dividerWidth := width - helpWidth - paginatorWidth - 2
	if dividerWidth < 1 {
		dividerWidth = 1
	}

	return helpView + " " + styles.Divider(dividerWidth) + " " + paginatorView
}

func (m NotificationsModel) renderLoadingScreen(width, height int) string {
	loadingBox := styles.DataBox("SCANNING SIGNALS",
		"\n"+
			"  "+m.spinner.View()+styles.Normal.Render(" Loading notifications...")+"\n"+
			"\n"+
			"  "+styles.Dim.Render("Intercepting incoming transmissions...")+"\n",
		50)
	return FullScreen(loadingBox, width, height, lipgloss.Center, lipgloss.Center)
}

func (m NotificationsModel) renderErrorScreen(width, height int) string {
	errorBox := styles.AlertBox(m.err.Error(), "error", 50) +
		"\n\n" +
		styles.Dim.Render("Press [r] to retry, [esc] to go back")
	return FullScreen(errorBox, width, height, lipgloss.Center, lipgloss.Center)
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
			return NotificationsErrorMsg{Err: err}
		}
		return NotificationsLoadedMsg{Notifications: notifs, Cursor: cursor}
	}
}

func (m NotificationsModel) fetchMoreNotifications() tea.Cmd {
	return func() tea.Msg {
		notifs, cursor, err := m.client.FetchMoreNotifications(30, m.nextCursor)
		if err != nil {
			return NotificationsErrorMsg{Err: err}
		}
		return MoreNotificationsLoadedMsg{Notifications: notifs, Cursor: cursor}
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

func notificationsToItems(notifs []models.Notification) []list.Item {
	items := make([]list.Item, len(notifs))
	for i, n := range notifs {
		items[i] = NotificationItem{Notification: n}
	}
	return items
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
