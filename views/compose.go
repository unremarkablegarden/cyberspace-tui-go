package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/unremarkablegarden/cyberspace-tui-go/api"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// PostCreatedMsg is sent when a post is successfully created
type PostCreatedMsg struct{ PostID string }

// PostCreateErrorMsg is sent when creating a post fails
type PostCreateErrorMsg struct{ Err error }

// ComposeBackMsg is sent when the user wants to leave the compose screen
type ComposeBackMsg struct{}

// ComposeKeyMap defines keybindings for the compose screen
type ComposeKeyMap struct {
	Send       key.Binding
	Cancel     key.Binding
	SwitchField key.Binding
}

// NewComposeKeyMap returns the default compose keybindings
func NewComposeKeyMap() ComposeKeyMap {
	return ComposeKeyMap{
		Send: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "send"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		SwitchField: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "switch field"),
		),
	}
}

// ShortHelp returns short help bindings
func (k ComposeKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.SwitchField, k.Send, k.Cancel}
}

// FullHelp returns full help bindings
func (k ComposeKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.SwitchField, k.Send, k.Cancel}}
}

// ComposeModel is the new post compose screen
type ComposeModel struct {
	client      *api.Client
	content     textarea.Model
	topicsInput textinput.Model
	focused     string // "content" or "topics"
	sending     bool
	err         error
	width       int
	height      int
	keys        ComposeKeyMap
	help        help.Model
	spinner     spinner.Model
}

// NewComposeModel creates a new compose screen
func NewComposeModel(baseURL, idToken string) ComposeModel {
	ta := textarea.New()
	ta.Placeholder = "What's on your mind?"
	ta.SetHeight(8)
	ta.SetWidth(60)
	ta.CharLimit = 32768
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle().Background(styles.ColorBgSelect)
	ta.FocusedStyle.Base = lipgloss.NewStyle().Foreground(styles.ColorNormal)
	ta.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(styles.ColorMuted)
	ta.FocusedStyle.EndOfBuffer = lipgloss.NewStyle().Foreground(styles.ColorDark)
	ta.BlurredStyle = ta.FocusedStyle
	ta.Focus()

	ti := textinput.New()
	ti.Placeholder = "linux, music, art"
	ti.CharLimit = 64
	ti.Width = 40
	ti.PromptStyle = lipgloss.NewStyle().Foreground(styles.ColorPrimary)
	ti.TextStyle = lipgloss.NewStyle().Foreground(styles.ColorNormal)

	h := help.New()
	h.Styles = styles.HelpStyles()

	return ComposeModel{
		client:      api.NewClient(baseURL, idToken),
		content:     ta,
		topicsInput: ti,
		focused:     "content",
		keys:        NewComposeKeyMap(),
		help:        h,
		spinner:     NewSpinner(),
	}
}

func (m ComposeModel) Init() tea.Cmd {
	return m.content.Focus()
}

func (m ComposeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.sending {
			return m, nil
		}
		switch {
		case key.Matches(msg, m.keys.Cancel):
			return m, func() tea.Msg { return ComposeBackMsg{} }
		case key.Matches(msg, m.keys.Send):
			content := strings.TrimSpace(m.content.Value())
			if content == "" {
				return m, nil
			}
			m.sending = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.sendPost(content, m.parseTopics()))
		case key.Matches(msg, m.keys.SwitchField):
			if m.focused == "content" {
				m.focused = "topics"
				m.content.Blur()
				return m, m.topicsInput.Focus()
			}
			m.focused = "content"
			m.topicsInput.Blur()
			return m, m.content.Focus()
		default:
			if m.focused == "content" {
				var cmd tea.Cmd
				m.content, cmd = m.content.Update(msg)
				return m, cmd
			}
			var cmd tea.Cmd
			m.topicsInput, cmd = m.topicsInput.Update(msg)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.content.SetWidth(msg.Width - 8)
		m.topicsInput.Width = msg.Width - 20

	case spinner.TickMsg:
		if m.sending {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case PostCreatedMsg:
		return m, func() tea.Msg { return ComposeBackMsg{} }

	case PostCreateErrorMsg:
		m.sending = false
		m.err = msg.Err

	case ThemeChangedMsg:
		m.spinner.Style = styles.Spinner
		m.help.Styles = styles.HelpStyles()
		m.content.FocusedStyle.CursorLine = lipgloss.NewStyle().Background(styles.ColorBgSelect)
		m.content.FocusedStyle.Base = lipgloss.NewStyle().Foreground(styles.ColorNormal)
		m.content.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(styles.ColorMuted)
		m.content.BlurredStyle = m.content.FocusedStyle
		m.topicsInput.PromptStyle = lipgloss.NewStyle().Foreground(styles.ColorPrimary)
		m.topicsInput.TextStyle = lipgloss.NewStyle().Foreground(styles.ColorNormal)
	}

	return m, nil
}

func (m ComposeModel) parseTopics() []string {
	raw := strings.TrimSpace(m.topicsInput.Value())
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	var topics []string
	for _, p := range parts {
		t := strings.ToLower(strings.TrimSpace(p))
		if t != "" {
			topics = append(topics, t)
		}
	}
	if len(topics) > 3 {
		topics = topics[:3]
	}
	return topics
}

func (m ComposeModel) sendPost(content string, topics []string) tea.Cmd {
	return func() tea.Msg {
		postID, err := m.client.CreatePost(content, topics)
		if err != nil {
			return PostCreateErrorMsg{Err: err}
		}
		return PostCreatedMsg{PostID: postID}
	}
}

func (m ComposeModel) View() string {
	w, _ := SafeDimensions(m.width, m.height)

	borderStyle := lipgloss.NewStyle().Foreground(styles.ColorBright)
	titleStyle := lipgloss.NewStyle().Foreground(styles.ColorBright).Bold(true)
	dimStyle := styles.Dim

	innerWidth := w - 6
	if innerWidth < 40 {
		innerWidth = 40
	}

	var b strings.Builder

	b.WriteString(RenderHeader("▓▒░ NEW POST ░▒▓", w))
	b.WriteString("\n")

	// Box title
	title := "COMPOSE"
	if m.sending {
		title = "TRANSMITTING..."
	}
	dashesLen := innerWidth - len(title) - 2
	if dashesLen < 1 {
		dashesLen = 1
	}
	b.WriteString(borderStyle.Render("╭─ ") + titleStyle.Render(title) + borderStyle.Render(" "+strings.Repeat("─", dashesLen)+"╮"))
	b.WriteString("\n")

	// Content textarea
	for _, line := range strings.Split(m.content.View(), "\n") {
		b.WriteString(borderStyle.Render("│ "))
		b.WriteString(line)
		b.WriteString("\n")
	}

	// Char count
	b.WriteString(borderStyle.Render("│ "))
	b.WriteString(dimStyle.Render(fmt.Sprintf("%d / 32768", len(m.content.Value()))))
	b.WriteString("\n")

	// Separator
	b.WriteString(borderStyle.Render("├" + strings.Repeat("─", innerWidth+2) + "┤"))
	b.WriteString("\n")

	// Topics field
	topicsLabel := "topics: "
	if m.focused == "topics" {
		b.WriteString(borderStyle.Render("│ ") + styles.Bright.Render(topicsLabel) + m.topicsInput.View())
	} else {
		b.WriteString(borderStyle.Render("│ ") + dimStyle.Render(topicsLabel) + m.topicsInput.View())
	}
	b.WriteString("\n")

	b.WriteString(borderStyle.Render("│ "))
	b.WriteString(dimStyle.Render("comma-separated, max 3, lowercase"))
	b.WriteString("\n")

	// Error line
	if m.err != nil {
		b.WriteString(borderStyle.Render("│ "))
		b.WriteString(styles.Error.Render("error: " + m.err.Error()))
		b.WriteString("\n")
	}

	b.WriteString(borderStyle.Render("╰" + strings.Repeat("─", innerWidth+2) + "╯"))
	b.WriteString("\n\n")

	b.WriteString(m.help.View(m.keys))

	return b.String()
}

// SetSize updates the view dimensions
func (m *ComposeModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.content.SetWidth(width - 8)
	m.topicsInput.Width = width - 20
}
