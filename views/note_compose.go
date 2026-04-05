package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/unremarkablegarden/cyberspace-tui-go/api"
	"github.com/unremarkablegarden/cyberspace-tui-go/models"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// NoteComposeDoneMsg is sent when the note compose screen closes
type NoteComposeDoneMsg struct{ Saved bool }

type noteSaveSuccessMsg struct{}
type noteSaveErrorMsg struct{ Err error }

// NoteComposeModel is the note create/edit screen
type NoteComposeModel struct {
	client      *api.Client
	note        models.Note // empty ID = new note
	isEdit      bool
	content     textarea.Model
	topicsInput textinput.Model
	focused     string // "content" or "topics"
	saving      bool
	err         error
	width       int
	height      int
	keys        NoteComposeKeyMap
	help        help.Model
	spinner     spinner.Model
}

// NewNoteComposeModel creates a note compose/edit screen
func NewNoteComposeModel(baseURL, idToken string, note models.Note, isEdit bool) NoteComposeModel {
	ta := textarea.New()
	ta.Placeholder = "Write your note..."
	ta.SetHeight(10)
	ta.SetWidth(60)
	ta.CharLimit = 32768
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle().Background(styles.ColorBgSelect)
	ta.FocusedStyle.Base = lipgloss.NewStyle().Foreground(styles.ColorNormal)
	ta.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(styles.ColorMuted)
	ta.FocusedStyle.EndOfBuffer = lipgloss.NewStyle().Foreground(styles.ColorDark)
	ta.BlurredStyle = ta.FocusedStyle
	if isEdit {
		ta.SetValue(note.Content)
	}
	ta.Focus()

	ti := textinput.New()
	ti.Placeholder = "linux, music, art"
	ti.CharLimit = 64
	ti.Width = 40
	ti.PromptStyle = lipgloss.NewStyle().Foreground(styles.ColorPrimary)
	ti.TextStyle = lipgloss.NewStyle().Foreground(styles.ColorNormal)
	if isEdit && len(note.Topics) > 0 {
		ti.SetValue(strings.Join(note.Topics, ", "))
	}

	h := help.New()
	h.Styles = styles.HelpStyles()

	return NoteComposeModel{
		client:      api.NewClient(baseURL, idToken),
		note:        note,
		isEdit:      isEdit,
		content:     ta,
		topicsInput: ti,
		focused:     "content",
		keys:        NewNoteComposeKeyMap(),
		help:        h,
		spinner:     NewSpinner(),
	}
}

func (m NoteComposeModel) Init() tea.Cmd {
	return m.content.Focus()
}

func (m NoteComposeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.saving {
			return m, nil
		}
		switch {
		case msg.String() == "esc":
			return m, func() tea.Msg { return NoteComposeDoneMsg{Saved: false} }
		case msg.String() == "ctrl+s":
			content := strings.TrimSpace(m.content.Value())
			if content == "" {
				return m, nil
			}
			m.saving = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.saveNote(content, m.parseTopics()))
		case msg.String() == "tab":
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
		if m.saving {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case noteSaveSuccessMsg:
		return m, func() tea.Msg { return NoteComposeDoneMsg{Saved: true} }

	case noteSaveErrorMsg:
		m.saving = false
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

func (m NoteComposeModel) parseTopics() []string {
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

func (m NoteComposeModel) saveNote(content string, topics []string) tea.Cmd {
	return func() tea.Msg {
		var err error
		if m.isEdit {
			err = m.client.UpdateNote(m.note.ID, content, topics)
		} else {
			_, err = m.client.CreateNote(content, topics)
		}
		if err != nil {
			return noteSaveErrorMsg{Err: err}
		}
		return noteSaveSuccessMsg{}
	}
}

func (m NoteComposeModel) View() string {
	w, _ := SafeDimensions(m.width, m.height)

	borderStyle := lipgloss.NewStyle().Foreground(styles.ColorBright)
	titleStyle := lipgloss.NewStyle().Foreground(styles.ColorBright).Bold(true)
	dimStyle := styles.Dim

	innerWidth := w - 6
	if innerWidth < 40 {
		innerWidth = 40
	}

	var b strings.Builder

	header := "▓▒░ NEW NOTE ░▒▓"
	if m.isEdit {
		header = "▓▒░ EDIT NOTE ░▒▓"
	}
	b.WriteString(RenderHeader(header, w))
	b.WriteString("\n")

	title := "COMPOSE"
	if m.isEdit {
		title = "EDITING"
	}
	if m.saving {
		title = "SAVING..."
	}
	dashesLen := innerWidth - len(title) - 2
	if dashesLen < 1 {
		dashesLen = 1
	}
	b.WriteString(borderStyle.Render("╭─ ") + titleStyle.Render(title) + borderStyle.Render(" "+strings.Repeat("─", dashesLen)+"╮"))
	b.WriteString("\n")

	for _, line := range strings.Split(m.content.View(), "\n") {
		b.WriteString(borderStyle.Render("│ "))
		b.WriteString(line)
		b.WriteString("\n")
	}

	b.WriteString(borderStyle.Render("│ "))
	b.WriteString(dimStyle.Render(fmt.Sprintf("%d / 32768", len(m.content.Value()))))
	b.WriteString("\n")

	b.WriteString(borderStyle.Render("├" + strings.Repeat("─", innerWidth+2) + "┤"))
	b.WriteString("\n")

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
func (m *NoteComposeModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.content.SetWidth(width - 8)
	m.topicsInput.Width = width - 20
}
