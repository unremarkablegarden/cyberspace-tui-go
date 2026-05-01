package models

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/external/api"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/messages"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/items"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// ─── Topic list item ────────────────────────────────────────────────────────
type TopicItem struct{ Topic entities.Topic }

func (t TopicItem) FilterValue() string { return t.Topic.Name }
func (t TopicItem) Title() string       { return "[" + t.Topic.Name + "]" }
func (t TopicItem) Description() string {
	return fmt.Sprintf("%d posts", t.Topic.PostCount)
}

type TopicDelegate struct{}

func (d TopicDelegate) Height() int                             { return 3 }
func (d TopicDelegate) Spacing() int                            { return 0 }
func (d TopicDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d TopicDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	it, ok := item.(TopicItem)
	if !ok {
		return
	}
	selected := index == m.Index()
	width := m.Width()
	fmt.Fprint(w, renderTopicCard(it.Topic, selected, width))
}

func renderTopicCard(t entities.Topic, selected bool, width int) string {
	innerWidth := width - 4
	if innerWidth < 20 {
		innerWidth = 76
	}

	tag := "[" + t.Name + "]"
	count := fmt.Sprintf("%d posts", t.PostCount)

	spacing := innerWidth - len(tag) - len(count)
	if spacing < 1 {
		spacing = 1
	}

	line := styles.Bright.Render(tag) + strings.Repeat(" ", spacing) + styles.Dim.Render(count)
	return items.BuildCardBox(line, innerWidth, selected)
}

// ─── TopicsModel ────────────────────────────────────────────────────────────

type TopicsModel struct {
	list    list.Model
	loading bool
	spinner spinner.Model
	err     error
	client  *api.Client
	width   int
	height  int
	keys    TopicsKeyMap
	help    help.Model
}

func NewTopicsModel(client *api.Client) TopicsModel {
	delegate := TopicDelegate{}
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
		spinner: items.NewSpinner(),
		loading: true,
		keys:    NewTopicsKeyMap(),
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
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case key.Matches(msg, m.keys.Back):
			return m, func() tea.Msg { return messages.SwitchToFeed{} }
		case key.Matches(msg, m.keys.Open):
			if it, ok := m.list.SelectedItem().(TopicItem); ok {
				topic := it.Topic
				return m, func() tea.Msg { return messages.SwitchToTopicFeed{Topic: topic} }
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

	case messages.TopicsLoadedMsg:
		m.loading = false
		localItems := make([]list.Item, len(msg.Topics))
		for i, t := range msg.Topics {
			localItems[i] = TopicItem{Topic: t}
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
	w, h := items.SafeDimensions(m.width, m.height)

	if m.loading {
		loadingBox := styles.DataBox("INDEXING TOPICS",
			"\n"+
				"  "+m.spinner.View()+styles.Normal.Render(" Loading topics...")+"\n"+
				"\n"+
				"  "+styles.Dim.Render("Scanning the datastream...")+"\n",
			50)
		return items.FullScreen(loadingBox, w, h, lipgloss.Center, lipgloss.Center)
	}

	if m.err != nil {
		errorBox := styles.AlertBox(m.err.Error(), "error", 50) +
			"\n\n" +
			styles.Dim.Render("Press [esc] to go back")
		return items.FullScreen(errorBox, w, h, lipgloss.Center, lipgloss.Center)
	}

	var b strings.Builder
	b.WriteString(items.RenderHeader("▓▒░ TOPICS ░▒▓", w))
	b.WriteString(m.list.View())
	b.WriteString("\n")

	helpView := m.help.View(m.keys)
	paginatorView := m.list.Paginator.View()
	helpWidth := lipgloss.Width(helpView)
	paginatorWidth := lipgloss.Width(paginatorView)
	dividerWidth := w - helpWidth - paginatorWidth - 2
	if dividerWidth < 1 {
		dividerWidth = 1
	}
	b.WriteString(helpView + " " + styles.Divider(dividerWidth) + " " + paginatorView)

	_ = h
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
