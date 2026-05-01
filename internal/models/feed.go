package models

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/external/api"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/messages"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// FeedModel is the post feed screen
type FeedModel struct {
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
	keys        FeedKeyMap
	help        help.Model
}

// NewFeedModel creates a new feed screen
func NewFeedModel(client *api.Client) FeedModel {
	// Create list with custom delegate
	delegate := PostDelegate{}
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.SetShowTitle(false)
	l.SetShowFilter(false)
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	l.Styles = styles.ListStyles()

	// Pagination dots: half blocks, bright for active, dim for inactive
	l.Paginator.ActiveDot = styles.Bright.Render("▄")
	l.Paginator.InactiveDot = styles.Dark.Render("▄")

	// Disable list's built-in quit — we handle it ourselves
	l.KeyMap.Quit.SetEnabled(false)
	// Disable ForceQuit too (ctrl+c handled by main.go)
	l.KeyMap.ForceQuit.SetEnabled(false)

	h := help.New()
	h.Styles = styles.HelpStyles()

	return FeedModel{
		list:    l,
		client:  client,
		spinner: NewSpinner(),
		loading: true,
		hasMore: true,
		keys:    NewFeedKeyMap(),
		help:    h,
	}
}

func (m FeedModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchPosts())
}

func (m FeedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Don't process keys during initial load
		if m.loading {
			return m, nil
		}

		switch {
		// temporal access while menu is implemented
		case msg.String() == "t":
			return m, func() tea.Msg { return messages.SwitchToThemeSwitcher{} }
		case msg.String() == "esc":
			// esc never quits — swallow it on the feed screen
			return m, nil
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case key.Matches(msg, m.keys.Refresh):
			m.loading = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.fetchPosts())
		case key.Matches(msg, m.keys.Logout):
			return m, func() tea.Msg { return messages.LogoutMsg{} }
		case key.Matches(msg, m.keys.NewPost):
			return m, func() tea.Msg { return messages.SwitchToCompose{} }
		case key.Matches(msg, m.keys.Bookmarks):
			return m, func() tea.Msg { return messages.SwitchToBookmarks{} }
		case key.Matches(msg, m.keys.Notifications):
			return m, func() tea.Msg { return messages.SwitchToNotifications{} }
		case key.Matches(msg, m.keys.Topics):
			return m, func() tea.Msg { return messages.SwitchToTopics{} }
		case key.Matches(msg, m.keys.Notes):
			return m, func() tea.Msg { return messages.SwitchToNotes{} }
		case key.Matches(msg, m.keys.Profile):
			if item, ok := m.list.SelectedItem().(PostItem); ok {
				username := item.Post.AuthorUsername
				return m, func() tea.Msg { return messages.SwitchToProfile{Username: username} }
			}
		case key.Matches(msg, m.keys.Open):
			if item := m.list.SelectedItem(); item != nil {
				switch it := item.(type) {
				case PostItem:
					return m, func() tea.Msg {
						return messages.SwitchToPostDetail{Post: it.Post}
					}
				case LoadMoreItem:
					if !m.loadingMore {
						m.loadingMore = true
						return m, tea.Batch(m.spinner.Tick, m.fetchMorePosts())
					}
				}
			}
		}

	case messages.RefreshFeedMsg:
		m.loading = true
		m.err = nil
		return m, tea.Batch(m.spinner.Tick, m.fetchPosts())

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionRelease && !m.loading {
			// Check if a post card was clicked
			for _, item := range m.list.Items() {
				if pi, ok := item.(PostItem); ok {
					if zone.Get(pi.Post.ID).InBounds(msg) {
						post := pi.Post
						return m, func() tea.Msg {
							return messages.SwitchToPostDetail{Post: post}
						}
					}
				}
			}
			// Check load more
			if zone.Get("load-more").InBounds(msg) && m.hasMore && !m.loadingMore {
				m.loadingMore = true
				return m, tea.Batch(m.spinner.Tick, m.fetchMorePosts())
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Reserve space for our custom header (2 lines) and footer (2 lines)
		m.list.SetSize(msg.Width, msg.Height-4)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case messages.FeedLoadedMsg:
		m.loading = false
		m.err = nil
		m.nextCursor = msg.Cursor
		m.hasMore = msg.Cursor != ""

		items := buildListItems(
			&m.list, msg.IsAdditional, m.hasMore,
			postsToItems(msg.Posts),
		)
		return m, m.list.SetItems(items)

	case messages.FeedErrorMsg:
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

	// Forward all other messages to the list
	if !m.loading {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m FeedModel) View() string {
	w, h := SafeDimensions(m.width, m.height)

	if m.loading {
		return m.renderLoadingScreen(w, h)
	}

	if m.err != nil {
		return m.renderErrorScreen(w, h)
	}

	var b strings.Builder

	// Header: centered title with blocks on each side
	b.WriteString(m.renderHeader(w))

	// List content (title disabled, we render our own header)
	b.WriteString(m.list.View())

	// Footer: divider with paginator inline on the right
	b.WriteString("\n")
	b.WriteString(m.renderFooter(w))

	_ = h // height managed by list.SetSize
	return b.String()
}

func (m FeedModel) renderHeader(width int) string {
	return RenderHeader("▓▒░ ᑕ¥βєяรקค¢є ░▒▓", width)
}

func (m FeedModel) renderFooter(width int) string {
	helpView := m.help.View(m.keys)
	paginatorView := m.list.Paginator.View()

	helpWidth := lipgloss.Width(helpView)
	paginatorWidth := lipgloss.Width(paginatorView)

	// help ────── paginator
	dividerWidth := width - helpWidth - paginatorWidth - 2
	if dividerWidth < 1 {
		dividerWidth = 1
	}

	return helpView + " " + styles.Divider(dividerWidth) + " " + paginatorView
}

func (m FeedModel) renderLoadingScreen(width, height int) string {
	var b strings.Builder

	loadingBox := styles.DataBox("ESTABLISHING CONNECTION",
		"\n"+
			"  "+m.spinner.View()+styles.Normal.Render(" Synchronizing with network...")+"\n"+
			"\n"+
			"  "+styles.Dim.Render("Please wait while we access the datastream")+"\n",
		50)

	return FullScreen(b.String()+loadingBox, width, height, lipgloss.Center, lipgloss.Center)
}

func (m FeedModel) renderErrorScreen(width, height int) string {
	errorBox := styles.AlertBox(m.err.Error(), "error", 50) +
		"\n\n" +
		styles.Dim.Render("Press [r] to retry connection, [q] to disconnect")

	return FullScreen(errorBox, width, height, lipgloss.Center, lipgloss.Center)
}

// SetSize updates the view dimensions
func (m *FeedModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width, height-4)
}

func (m FeedModel) fetchPosts() tea.Cmd {
	return func() tea.Msg {
		posts, cursor, err := m.client.FetchPosts(30)
		if err != nil {
			return messages.FeedErrorMsg{Err: err}
		}
		return messages.FeedLoadedMsg{Posts: posts, Cursor: cursor}
	}
}

func (m FeedModel) fetchMorePosts() tea.Cmd {
	return func() tea.Msg {
		posts, cursor, err := m.client.FetchMorePosts(30, m.nextCursor)
		if err != nil {
			return messages.FeedErrorMsg{Err: err}
		}
		return messages.FeedLoadedMsg{Posts: posts, Cursor: cursor, IsAdditional: true}
	}
}
