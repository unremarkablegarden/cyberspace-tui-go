package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/external/api"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/messages"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/ui"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// EditProfileKeyMap defines keybindings for the edit profile screen
type EditProfileKeyMap struct {
	SwitchField key.Binding
	Save        key.Binding
	Cancel      key.Binding
}

// NewEditProfileKeyMap returns default edit profile keybindings
func NewEditProfileKeyMap() EditProfileKeyMap {
	return EditProfileKeyMap{
		SwitchField: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next field"),
		),
		Save: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "save"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
	}
}

func (k EditProfileKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.SwitchField, k.Save, k.Cancel}
}

func (k EditProfileKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.SwitchField, k.Save, k.Cancel}}
}

var editFieldLabels = []string{
	"display name",
	"bio",
	"website name",
	"website url",
	"location",
}

var editFieldLimits = []int{64, 127, 64, 2048, 64}

// EditProfileModel is the edit profile form screen
type EditProfileModel struct {
	client  *api.Client
	user    entities.User
	fields  []textinput.Model
	focused int
	saving  bool
	err     error
	width   int
	height  int
	keys    EditProfileKeyMap
	help    help.Model
	spinner spinner.Model
}

// NewEditProfileModel creates a new edit profile screen pre-populated with user data
func NewEditProfileModel(client *api.Client, user entities.User) EditProfileModel {
	values := []string{
		user.DisplayName,
		user.Bio,
		user.WebsiteName,
		user.WebsiteURL,
		user.LocationName,
	}

	fields := make([]textinput.Model, len(editFieldLabels))
	for i := range fields {
		ti := textinput.New()
		ti.CharLimit = editFieldLimits[i]
		ti.Width = 50
		ti.SetValue(values[i])
		ti.PromptStyle = lipgloss.NewStyle().Foreground(styles.ColorPrimary)
		ti.TextStyle = lipgloss.NewStyle().Foreground(styles.ColorNormal)
		fields[i] = ti
	}
	fields[0].Focus()

	h := help.New()
	h.Styles = styles.HelpStyles()

	return EditProfileModel{
		client:  client,
		user:    user,
		fields:  fields,
		focused: 0,
		keys:    NewEditProfileKeyMap(),
		help:    h,
		spinner: ui.NewSpinner(),
	}
}

func (m EditProfileModel) Init() tea.Cmd {
	return m.fields[0].Focus()
}

func (m EditProfileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.saving {
			return m, nil
		}
		switch {
		case key.Matches(msg, m.keys.Cancel):
			return m, func() tea.Msg { return messages.SwitchToProfile{Username: m.user.Username} }
		case key.Matches(msg, m.keys.Save):
			m.saving = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.saveProfile())
		case key.Matches(msg, m.keys.SwitchField):
			m.fields[m.focused].Blur()
			m.focused = (m.focused + 1) % len(m.fields)
			return m, m.fields[m.focused].Focus()
		default:
			var cmd tea.Cmd
			m.fields[m.focused], cmd = m.fields[m.focused].Update(msg)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		for i := range m.fields {
			m.fields[i].Width = msg.Width - 20
		}

	case spinner.TickMsg:
		if m.saving {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case messages.EditProfileSaveSuccessMsg:
		return m, func() tea.Msg { return messages.SwitchToProfile{Username: m.user.Username} }

	case messages.EditProfileSaveSuccessErrMsg:
		m.saving = false
		m.err = msg.Err

	case ThemeChangedMsg:
		m.spinner.Style = styles.Spinner
		m.help.Styles = styles.HelpStyles()
		for i := range m.fields {
			m.fields[i].PromptStyle = lipgloss.NewStyle().Foreground(styles.ColorPrimary)
			m.fields[i].TextStyle = lipgloss.NewStyle().Foreground(styles.ColorNormal)
		}
	}

	return m, nil
}

func (m EditProfileModel) saveProfile() tea.Cmd {
	return func() tea.Msg {
		req := api.UpdateProfileRequest{
			DisplayName:  m.fields[0].Value(),
			Bio:          m.fields[1].Value(),
			WebsiteName:  m.fields[2].Value(),
			WebsiteURL:   m.fields[3].Value(),
			LocationName: m.fields[4].Value(),
		}
		if err := m.client.UpdateProfile(req); err != nil {
			return messages.EditProfileSaveSuccessErrMsg{Err: err}
		}
		return messages.EditProfileSaveSuccessMsg{}
	}
}

func (m EditProfileModel) View() string {
	w, _ := ui.SafeDimensions(m.width, m.height)

	borderStyle := lipgloss.NewStyle().Foreground(styles.ColorBright)
	titleStyle := lipgloss.NewStyle().Foreground(styles.ColorBright).Bold(true)
	dimStyle := styles.Dim

	innerWidth := max(w-6, 40)
	var b strings.Builder
	b.WriteString(ui.RenderHeader("▓▒░ EDIT PROFILE ░▒▓", w))
	b.WriteString("\n")

	title := "PROFILE DATA"
	if m.saving {
		title = "TRANSMITTING..."
	}
	dashesLen := max(innerWidth-len(title)-2, 1)

	b.WriteString(borderStyle.Render("╭─ ") + titleStyle.Render(title) + borderStyle.Render(" "+strings.Repeat("─", dashesLen)+"╮"))
	b.WriteString("\n")

	for i, label := range editFieldLabels {
		labelStr := label + ": "
		if m.focused == i {
			b.WriteString(borderStyle.Render("│ ") + styles.Bright.Render(labelStr) + m.fields[i].View())
		} else {
			b.WriteString(borderStyle.Render("│ ") + dimStyle.Render(labelStr) + m.fields[i].View())
		}
		b.WriteString("\n")
		b.WriteString(borderStyle.Render("│ ") + dimStyle.Render(fmt.Sprintf("%d / %d chars", len(m.fields[i].Value()), editFieldLimits[i])))
		b.WriteString("\n")
		if i < len(editFieldLabels)-1 {
			b.WriteString(borderStyle.Render("│"))
			b.WriteString("\n")
		}
	}

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
func (m *EditProfileModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	for i := range m.fields {
		m.fields[i].Width = width - 20
	}
}
