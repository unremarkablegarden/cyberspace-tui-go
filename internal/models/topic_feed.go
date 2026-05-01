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

// TopicFeedModel is the topic-filtered post feed screen
type TopicFeedModel struct {
	topic       entities.Topic
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
	keys        TopicFeedKeyMap
	help        help.Model
}

func NewTopicFeedModel(client *api.Client, topic entities.Topic) TopicFeedModel {
	delegate := PostDelegate{}
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
		spinner: NewSpinner(),
		loading: true,
		hasMore: true,
		keys:    NewTopicFeedKeyMap(),
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
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case key.Matches(msg, m.keys.Back):
			return m, func() tea.Msg { return messages.SwitchToTopics{} }
		case key.Matches(msg, m.keys.Refresh):
			m.loading = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.fetchPosts())
		case key.Matches(msg, m.keys.Profile):
			if item, ok := m.list.SelectedItem().(PostItem); ok {
				username := item.Post.AuthorUsername
				return m, func() tea.Msg {
					return messages.SwitchToProfile{
						Username:    username,
						BackMessage: messages.SwitchToTopicFeed{Topic: m.topic},
					}
				}
			}
		case key.Matches(msg, m.keys.Open):
			switch it := m.list.SelectedItem().(type) {
			case PostItem:
				post := it.Post
				return m, func() tea.Msg {
					return messages.SwitchToPostDetail{
						Post:        post,
						BackMessage: messages.SwitchToTopicFeed{Topic: m.topic},
					}
				}
			case LoadMoreItem:
				if !m.loadingMore {
					m.loadingMore = true
					return m, tea.Batch(m.spinner.Tick, m.fetchMorePosts())
				}
			}
		}

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionRelease && !m.loading {
			for _, item := range m.list.Items() {
				if pi, ok := item.(PostItem); ok {
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

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case messages.TopicPostsLoadedMsg:
		m.loading = false
		m.nextCursor = msg.Cursor
		m.hasMore = msg.Cursor != ""
		m.err = nil

		var items []list.Item
		if msg.IsAdditional {
			for _, existing := range m.list.Items() {
				if _, ok := existing.(LoadMoreItem); ok {
					continue
				}
				items = append(items, existing)
			}
			for _, p := range msg.Posts {
				items = append(items, PostItem{Post: p})
			}
		} else {
			items = postsToItems(msg.Posts)
		}

		if m.hasMore {
			items = append(items, LoadMoreItem{})
		}
		cmd := m.list.SetItems(items)
		return m, cmd

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
	w, h := SafeDimensions(m.width, m.height)

	if m.loading {
		loadingBox := styles.DataBox("FILTERING FEED",
			"\n"+
				"  "+m.spinner.View()+styles.Normal.Render(" Loading ["+m.topic.Name+"]...")+"\n"+
				"\n"+
				"  "+styles.Dim.Render("Filtering transmissions by topic...")+"\n",
			50)
		return FullScreen(loadingBox, w, h, lipgloss.Center, lipgloss.Center)
	}

	if m.err != nil {
		errorBox := styles.AlertBox(m.err.Error(), "error", 50) +
			"\n\n" +
			styles.Dim.Render("Press [r] to retry, [esc] to go back")
		return FullScreen(errorBox, w, h, lipgloss.Center, lipgloss.Center)
	}

	var b strings.Builder
	b.WriteString(RenderHeader("▓▒░ ["+m.topic.Name+"] ░▒▓", w))
	b.WriteString(m.list.View())
	b.WriteString("\n")
	b.WriteString(m.renderFooter(w))

	_ = h
	return b.String()
}

func (m TopicFeedModel) renderFooter(width int) string {
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
