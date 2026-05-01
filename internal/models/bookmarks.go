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

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/external/api"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/messages"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// BookmarksModel is the bookmarks list screen
type BookmarksModel struct {
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
	keys        BookmarksKeyMap
	help        help.Model
}

// NewBookmarksModel creates a new bookmarks screen
func NewBookmarksModel(client *api.Client) BookmarksModel {
	delegate := BookmarkDelegate{}
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

	return BookmarksModel{
		list:    l,
		client:  client,
		spinner: NewSpinner(),
		loading: true,
		hasMore: true,
		keys:    NewBookmarksKeyMap(),
		help:    h,
	}
}

func (m BookmarksModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchBookmarks())
}

func (m BookmarksModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, func() tea.Msg { return messages.SwitchToFeed{} }
		case key.Matches(msg, m.keys.Refresh):
			m.loading = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.fetchBookmarks())
		case key.Matches(msg, m.keys.Remove):
			if item, ok := m.list.SelectedItem().(BookmarkItem); ok && item.Bookmark.ID != "" {
				return m, m.deleteBookmark(item.Bookmark.ID)
			}
		case key.Matches(msg, m.keys.Open):
			switch it := m.list.SelectedItem().(type) {
			case BookmarkItem:
				post := it.Bookmark.Post
				return m, func() tea.Msg { return messages.SwitchToPost{Post: post} }
			case LoadMoreItem:
				if !m.loadingMore {
					m.loadingMore = true
					return m, tea.Batch(m.spinner.Tick, m.fetchMoreBookmarks())
				}
			}
		}

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionRelease && !m.loading {
			for _, item := range m.list.Items() {
				if bi, ok := item.(BookmarkItem); ok {
					if zone.Get(bi.Bookmark.Post.ID).InBounds(msg) {
						post := bi.Bookmark.Post
						return m, func() tea.Msg { return messages.SwitchToPost{Post: post} }
					}
				}
			}
			if zone.Get("load-more-bookmarks").InBounds(msg) && m.hasMore && !m.loadingMore {
				m.loadingMore = true
				return m, tea.Batch(m.spinner.Tick, m.fetchMoreBookmarks())
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

	case messages.BookmarksLoadedMsg:
		m.loading = false
		m.err = nil
		m.nextCursor = msg.Cursor
		m.hasMore = msg.Cursor != ""

		var items []list.Item
		if msg.IsAdditional {
			for _, existing := range m.list.Items() {
				if _, ok := existing.(LoadMoreItem); ok {
					continue
				}
				items = append(items, existing)
			}
			for _, b := range msg.Bookmarks {
				items = append(items, BookmarkItem{Bookmark: b})
			}
		} else {
			items = bookmarksToItems(msg.Bookmarks)
		}

		if m.hasMore {
			items = append(items, LoadMoreItem{})
		}
		cmd := m.list.SetItems(items)
		return m, cmd

	case messages.BookmarksLoadedErrMsg:
		m.loading = false
		m.loadingMore = false
		m.err = msg.Err

	case messages.BookmarkRemovedMsg:
		var items []list.Item
		for _, existing := range m.list.Items() {
			if bi, ok := existing.(BookmarkItem); ok {
				if bi.Bookmark.ID == msg.BookmarkID {
					continue
				}
			}
			items = append(items, existing)
		}
		cmd := m.list.SetItems(items)
		return m, cmd

	case messages.BookmarkRemovedErrMsg:
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

func (m BookmarksModel) View() string {
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

func (m BookmarksModel) renderHeader(width int) string {
	return RenderHeader("▓▒░ BOOKMARKS ░▒▓", width)
}

func (m BookmarksModel) renderFooter(width int) string {
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

func (m BookmarksModel) renderLoadingScreen(width, height int) string {
	loadingBox := styles.DataBox("RETRIEVING SAVED DATA",
		"\n"+
			"  "+m.spinner.View()+styles.Normal.Render(" Loading bookmarks...")+"\n"+
			"\n"+
			"  "+styles.Dim.Render("Accessing your saved transmissions...")+"\n",
		50)
	return FullScreen(loadingBox, width, height, lipgloss.Center, lipgloss.Center)
}

func (m BookmarksModel) renderErrorScreen(width, height int) string {
	errorBox := styles.AlertBox(m.err.Error(), "error", 50) +
		"\n\n" +
		styles.Dim.Render("Press [r] to retry, [esc] to go back")
	return FullScreen(errorBox, width, height, lipgloss.Center, lipgloss.Center)
}

// SetSize updates the view dimensions
func (m *BookmarksModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width, height-4)
}

func (m BookmarksModel) fetchBookmarks() tea.Cmd {
	return func() tea.Msg {
		bookmarks, cursor, err := m.client.FetchBookmarks(20)
		if err != nil {
			return messages.BookmarksLoadedErrMsg{Err: err}
		}
		return messages.BookmarksLoadedMsg{Bookmarks: bookmarks, Cursor: cursor}
	}
}

func (m BookmarksModel) fetchMoreBookmarks() tea.Cmd {
	return func() tea.Msg {
		bookmarks, cursor, err := m.client.FetchMoreBookmarks(20, m.nextCursor)
		if err != nil {
			return messages.BookmarksLoadedErrMsg{Err: err}
		}
		return messages.BookmarksLoadedMsg{Bookmarks: bookmarks, Cursor: cursor, IsAdditional: true}
	}
}

func (m BookmarksModel) deleteBookmark(bookmarkID string) tea.Cmd {
	return func() tea.Msg {
		if err := m.client.DeleteBookmark(bookmarkID); err != nil {
			return messages.BookmarkRemovedErrMsg{Err: err}
		}
		return messages.BookmarkRemovedMsg{BookmarkID: bookmarkID}
	}
}

func bookmarksToItems(bookmarks []entities.Bookmark) []list.Item {
	items := make([]list.Item, 0, len(bookmarks))
	for _, b := range bookmarks {
		if !b.Post.Deleted {
			items = append(items, BookmarkItem{Bookmark: b})
		}
	}
	return items
}
