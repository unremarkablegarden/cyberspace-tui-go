package models

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/external/api"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/messages"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/items"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/keymaps"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/ui"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// BookmarksModel is the bookmarks list screen
type BookmarksModel struct {
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

// NewBookmarksModel creates a new bookmarks screen
func NewBookmarksModel(client *api.Client, keybinds keymaps.AppKeybinds, sp *spinner.Model) BookmarksModel {
	delegate := items.BookmarkDelegate{}
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
		spinner: sp,
		loading: true,
		hasMore: true,
		keys:    keybinds,
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
			return m, tea.Batch(m.spinner.Tick, m.fetchBookmarks())
		case m.keys.GlobalKeybinds.Open:
			switch it := m.list.SelectedItem().(type) {
			case items.BookmarkItem:
				post := it.Bookmark.Post
				return m, func() tea.Msg {
					return messages.SwitchToPostDetail{
						Post:        post,
						BackMessage: messages.SwitchToBookmarks{},
					}
				}
			case items.LoadMoreItem:
				if !m.loadingMore {
					m.loadingMore = true
					return m, tea.Batch(m.spinner.Tick, m.fetchMoreBookmarks())
				}
			}
		case m.keys.BookmarksKeybinds.Remove:
			if item, ok := m.list.SelectedItem().(items.BookmarkItem); ok && item.Bookmark.ID != "" {
				return m, m.deleteBookmark(item.Bookmark.ID)
			}
		}

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionRelease && !m.loading {
			for _, item := range m.list.Items() {
				if bi, ok := item.(items.BookmarkItem); ok {
					if zone.Get(bi.Bookmark.Post.ID).InBounds(msg) {
						post := bi.Bookmark.Post
						return m, func() tea.Msg {
							return messages.SwitchToPostDetail{
								Post:        post,
								BackMessage: messages.SwitchToBookmarks{},
							}
						}
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

	case messages.BookmarksLoadedMsg:
		m.loading = false
		m.loadingMore = false
		m.err = nil
		m.nextCursor = msg.Cursor
		m.hasMore = msg.Cursor != ""

		localItems := items.BuildListItems(
			&m.list, msg.IsAdditional, m.hasMore,
			items.BookmarksToItems(msg.Bookmarks),
		)
		return m, m.list.SetItems(localItems)

	case messages.BookmarksLoadedErrMsg:
		m.loading = false
		m.loadingMore = false
		m.err = msg.Err

	case messages.BookmarkRemovedMsg:
		var localItems []list.Item
		for _, existing := range m.list.Items() {
			if bi, ok := existing.(items.BookmarkItem); ok {
				if bi.Bookmark.ID == msg.BookmarkID {
					continue
				}
			}
			localItems = append(localItems, existing)
		}
		cmd := m.list.SetItems(localItems)
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
	w, h := ui.SafeDimensions(m.width, m.height)

	if m.loading {
		return ui.RenderLoadingScreen(
			ui.LoadingScreenTexts{
				TitleText:    "RETRIEVING SAVED DATA",
				SubtitleText: " Loading bookmarks...",
				BottomText:   "Accessing your saved transmissions...",
			},
			*m.spinner,
			w, h,
		)
	}

	if m.err != nil {
		return ui.RenderErrorScreen(m.err, w, h)
	}

	var b strings.Builder
	b.WriteString(ui.RenderHeader("▓▒░ BOOKMARKS ░▒▓", w))
	b.WriteString(m.list.View())
	b.WriteString("\n")
	b.WriteString(
		ui.RenderFooterWithList(
			m.help.View(m.keys.BookmarksHelpKeys()),
			m.list.Paginator.View(),
			w,
		))

	return b.String()
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
