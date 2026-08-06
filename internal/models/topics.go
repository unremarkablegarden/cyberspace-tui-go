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

type TopicsModel struct {
	list    list.Model
	loading bool
	spinner *spinner.Model
	err     error
	client  *api.Client
	width   int
	height  int
	keys    keymaps.AppKeybinds
	help    help.Model
}

func NewTopicsModel(client *api.Client, keymap keymaps.AppKeybinds, sp *spinner.Model) TopicsModel {
	delegate := items.TopicDelegate{}
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

	return TopicsModel{
		list:    l,
		client:  client,
		spinner: sp,
		loading: true,
		keys:    keymap,
		help:    h,
	}
}

func (m TopicsModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchTopics())
}

func (m TopicsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case m.keys.GlobalKeybinds.Open:
			if it, ok := m.list.SelectedItem().(items.TopicItem); ok {
				topic := it.Topic
				return m, func() tea.Msg { return messages.SwitchToTopicFeed{Topic: topic} }
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-4)

	case messages.TopicsLoadedMsg:
		m.loading = false
		localItems := make([]list.Item, len(msg.Topics))
		for i, t := range msg.Topics {
			localItems[i] = items.TopicItem{Topic: t}
		}
		cmd := m.list.SetItems(localItems)
		return m, cmd

	case messages.TopicsLoadedErrMsg:
		m.loading = false
		m.err = msg.Err

	case ThemeChangedMsg:
		m.spinner.Style = styles.Spinner
		m.list.Styles = styles.ListStyles()
		m.help.Styles = styles.HelpStyles()
	}

	if !m.loading {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m TopicsModel) View() string {
	w, h := ui.SafeDimensions(m.width, m.height)

	if m.loading {
		return ui.RenderLoadingScreen(
			ui.LoadingScreenTexts{
				TitleText:    "INDEXING TOPICS",
				SubtitleText: " Loading topics...",
				BottomText:   "Scanning the datastream...",
			},
			*m.spinner,
			w, h,
		)
	}

	if m.err != nil {
		return ui.RenderErrorScreen(m.err, w, h)
	}

	var b strings.Builder
	b.WriteString(ui.RenderHeader("▓▒░ TOPICS ░▒▓", w))
	b.WriteString(m.list.View())
	b.WriteString("\n")
	b.WriteString(
		ui.RenderFooterWithList(
			m.help.View(m.keys.TopicsHelpKeys()),
			m.list.Paginator.View(),
			w,
		))

	return b.String()
}

func (m *TopicsModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width, height-4)
}

func (m TopicsModel) fetchTopics() tea.Cmd {
	return func() tea.Msg {
		topics, err := m.client.FetchTopics()
		if err != nil {
			return messages.TopicsLoadedErrMsg{Err: err}
		}
		return messages.TopicsLoadedMsg{Topics: topics}
	}
}
