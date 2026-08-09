package models

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/external/api"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/external/cache"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/messages"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/items"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/keymaps"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/ui"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// FeedModel is the post feed screen
type FeedModel struct {
	list        list.Model
	loading     bool
	loadingMore bool
	spinner     *spinner.Model
	err         error
	client      *api.Client
	cache       cache.ICache
	nextCursor  string
	hasMore     bool
	width       int
	height      int
	keys        keymaps.AppKeybinds
	help        help.Model
}

// NewFeedModel creates a new feed screen
func NewFeedModel(client *api.Client, cache cache.ICache, keymap keymaps.AppKeybinds, sp *spinner.Model) FeedModel {
	// Create list with custom delegate
	delegate := items.PostDelegate{}
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
		cache:   cache,
		spinner: sp,
		loading: true,
		hasMore: true,
		keys:    keymap,
		help:    h,
	}
}

func (m FeedModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchPosts(false))
}

func (m FeedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Don't process keys during initial load
		if m.loading {
			return m, nil
		}

		switch msg.String() {
		case m.keys.GlobalKeybinds.Quit:
			return m, tea.Quit
		case m.keys.GlobalKeybinds.Help:
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case m.keys.GlobalKeybinds.Refresh:
			m.loading = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.fetchPosts(true))
		case m.keys.FeedKeybinds.Logout:
			return m, func() tea.Msg { return messages.LogoutMsg{} }
		case m.keys.FeedKeybinds.Profile:
			if item, ok := m.list.SelectedItem().(items.PostItem); ok {
				username := item.Post.AuthorUsername
				return m, func() tea.Msg { return messages.SwitchToProfile{Username: username} }
			}
		case m.keys.GlobalKeybinds.Open:
			if item := m.list.SelectedItem(); item != nil {
				switch it := item.(type) {
				case items.PostItem:
					return m, func() tea.Msg {
						return messages.SwitchToPostDetail{Post: it.Post}
					}
				case items.LoadMoreItem:
					if !m.loadingMore {
						m.loadingMore = true
						return m, tea.Batch(m.spinner.Tick, m.fetchMorePosts())
					}
				}
			}
		}

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionRelease && !m.loading {
			// Check if a post card was clicked
			for _, item := range m.list.Items() {
				if pi, ok := item.(items.PostItem); ok {
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

	case messages.FeedLoadedMsg:
		m.loading = false
		m.loadingMore = false
		m.err = nil
		m.nextCursor = msg.Cursor
		m.hasMore = msg.Cursor != ""

		items := items.BuildListItems(
			&m.list, msg.IsAdditional, m.hasMore,
			items.PostsToItems(msg.Posts),
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
	w, h := ui.SafeDimensions(m.width, m.height)

	if m.loading {
		return ui.RenderLoadingScreen(
			ui.LoadingScreenTexts{
				TitleText:    "ESTABLISHING CONNECTION",
				SubtitleText: " Synchronizing with network...",
				BottomText:   "Please wait while we access the datastream",
			},
			*m.spinner,
			w, h,
		)
	}

	if m.err != nil {
		return ui.RenderErrorScreen(m.err, w, h)
	}

	var b strings.Builder

	// Header: centered title with blocks on each side
	b.WriteString(ui.RenderHeader("▓▒░ ᑕ¥βєяรקค¢є ░▒▓", w))

	// List content (title disabled, we render our own header)
	b.WriteString(m.list.View())

	// Footer: divider with paginator inline on the right
	b.WriteString("\n")
	b.WriteString(
		ui.RenderFooterWithList(
			m.help.View(m.keys.FeedHelpKeys()),
			m.list.Paginator.View(),
			w,
		))

	return b.String()
}

// SetSize updates the view dimensions
func (m *FeedModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width, height-4)
}

// Network stuff

func (m FeedModel) fetchPosts(isRefresh bool) tea.Cmd {
	return func() tea.Msg {
		if isRefresh {
			return m.postsFromAPI()
		}

		cachedPosts, postsFound := m.cache.Get(cache.DefaultFeedCacheKey + "posts")
		cachedCursor, cursorFound := m.cache.Get(cache.DefaultFeedCacheKey + "cursor")

		if !postsFound || !cursorFound {
			return m.postsFromAPI()
		}

		posts, ok := cachedPosts.([]entities.Post)
		if !ok {
			return m.postsFromAPI()
		}

		cursor, ok := cachedCursor.(string)
		if !ok {
			return m.postsFromAPI()
		}

		return messages.FeedLoadedMsg{
			Posts:  posts,
			Cursor: cursor,
		}
	}
}

func (m FeedModel) postsFromAPI() tea.Msg {
	posts, cursor, err := m.syncPostsFromAPI()
	if err != nil {
		return messages.FeedErrorMsg{Err: err}
	}

	return messages.FeedLoadedMsg{Posts: posts, Cursor: cursor}
}

func (m FeedModel) fetchMorePosts() tea.Cmd {
	return func() tea.Msg {
		var posts []entities.Post

		if p, found := m.cache.Get(cache.DefaultFeedCacheKey + "posts"); !found {
			apiPosts, apiCursor, err := m.syncPostsFromAPI()
			if err != nil {
				return messages.FeedErrorMsg{Err: err}
			}

			m.nextCursor = apiCursor
			posts = apiPosts
		} else {
			posts = p.([]entities.Post)
		}

		additionalPosts, cursor, err := m.client.FetchMorePosts(30, m.nextCursor)
		if err != nil {
			return messages.FeedErrorMsg{Err: err}
		}

		m.cache.Set(cache.DefaultFeedCacheKey+"posts", append(posts, additionalPosts...), 0)
		m.cache.Set(cache.DefaultFeedCacheKey+"cursor", cursor, 0)

		return messages.FeedLoadedMsg{Posts: additionalPosts, Cursor: cursor, IsAdditional: true}
	}
}

func (m FeedModel) syncPostsFromAPI() ([]entities.Post, string, error) {
	posts, cursor, err := m.client.FetchPosts(30)
	if err != nil {
		return []entities.Post{}, "", err
	}

	m.cache.Set(cache.DefaultFeedCacheKey+"posts", posts, 0)
	m.cache.Set(cache.DefaultFeedCacheKey+"cursor", cursor, 0)

	return posts, cursor, nil
}
