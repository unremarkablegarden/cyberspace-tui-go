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
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/messages"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/items"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/keymaps"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/ui"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// TopicFeedModel is the topic-filtered post feed screen
type TopicFeedModel struct {
	topic       entities.Topic
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

func NewTopicFeedModel(client *api.Client, keymap keymaps.AppKeybinds, sp *spinner.Model, topic entities.Topic) TopicFeedModel {
	delegate := items.PostDelegate{}
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

	return TopicFeedModel{
		topic:   topic,
		list:    l,
		client:  client,
		spinner: sp,
		loading: true,
		hasMore: true,
		keys:    keymap,
		help:    h,
	}
}

func (m TopicFeedModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchPosts())
}

func (m TopicFeedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, func() tea.Msg { return messages.SwitchToTopics{} }
		case m.keys.GlobalKeybinds.Refresh:
			m.loading = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.fetchPosts())
		case m.keys.GlobalKeybinds.Open:
			switch it := m.list.SelectedItem().(type) {
			case items.PostItem:
				post := it.Post
				return m, func() tea.Msg {
					return messages.SwitchToPostDetail{
						Post:        post,
						BackMessage: messages.SwitchToTopicFeed{Topic: m.topic},
					}
				}
			case items.LoadMoreItem:
				if !m.loadingMore {
					m.loadingMore = true
					return m, tea.Batch(m.spinner.Tick, m.fetchMorePosts())
				}
			}
		case m.keys.TopicsFeedKeybinds.Profile:
			if item, ok := m.list.SelectedItem().(items.PostItem); ok {
				username := item.Post.AuthorUsername
				return m, func() tea.Msg {
					return messages.SwitchToProfile{
						Username:    username,
						BackMessage: messages.SwitchToTopicFeed{Topic: m.topic},
					}
				}
			}
		}

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionRelease && !m.loading {
			for _, item := range m.list.Items() {
				if pi, ok := item.(items.PostItem); ok {
					if zone.Get(pi.Post.ID).InBounds(msg) {
						post := pi.Post
						return m, func() tea.Msg {
							return messages.SwitchToPostDetail{
								Post:        post,
								BackMessage: messages.SwitchToTopicFeed{Topic: m.topic},
							}
						}
					}
				}
			}
			if zone.Get("load-more-topic").InBounds(msg) && m.hasMore && !m.loadingMore {
				m.loadingMore = true
				return m, tea.Batch(m.spinner.Tick, m.fetchMorePosts())
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-4)

	case messages.TopicPostsLoadedMsg:
		m.loading = false
		m.loadingMore = false
		m.nextCursor = msg.Cursor
		m.hasMore = msg.Cursor != ""
		m.err = nil

		localItems := items.BuildListItems(
			&m.list, msg.IsAdditional, m.hasMore,
			items.PostsToItems(msg.Posts),
		)
		return m, m.list.SetItems(localItems)

	case messages.TopicPostsLoadedErrMsg:
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

func (m TopicFeedModel) View() string {
	w, h := ui.SafeDimensions(m.width, m.height)

	if m.loading {
		return ui.RenderLoadingScreen(
			ui.LoadingScreenTexts{
				TitleText:    "FILTERING FEED",
				SubtitleText: " Loading [" + m.topic.Name + "]...",
				BottomText:   "Filtering transmissions by topic...",
			},
			*m.spinner,
			w, h,
		)
	}

	if m.err != nil {
		return ui.RenderErrorScreen(m.err, w, h)
	}

	var b strings.Builder
	b.WriteString(ui.RenderHeader("▓▒░ ["+m.topic.Name+"] ░▒▓", w))
	b.WriteString(m.list.View())
	b.WriteString("\n")
	b.WriteString(
		ui.RenderFooterWithList(
			m.help.View(m.keys.TopicFeedHelpKeys()),
			m.list.Paginator.View(),
			w,
		))

	return b.String()
}

func (m *TopicFeedModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width, height-4)
}

func (m TopicFeedModel) fetchPosts() tea.Cmd {
	return func() tea.Msg {
		posts, cursor, err := m.client.FetchTopicPosts(m.topic.Name, 30)
		if err != nil {
			return messages.TopicPostsLoadedErrMsg{Err: err}
		}
		return messages.TopicPostsLoadedMsg{Posts: posts, Cursor: cursor}
	}
}

func (m TopicFeedModel) fetchMorePosts() tea.Cmd {
	return func() tea.Msg {
		posts, cursor, err := m.client.FetchMoreTopicPosts(m.topic.Name, 30, m.nextCursor)
		if err != nil {
			return messages.TopicPostsLoadedErrMsg{Err: err}
		}
		return messages.TopicPostsLoadedMsg{Posts: posts, Cursor: cursor, IsAdditional: true}
	}
}
