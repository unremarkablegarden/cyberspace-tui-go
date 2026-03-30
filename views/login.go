package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/euklides/cyberspace-cli/api"
	"github.com/euklides/cyberspace-cli/styles"
)

// LoginSuccessMsg is sent when login succeeds
type LoginSuccessMsg struct {
	IDToken      string
	RefreshToken string
	RTDBToken    string
}

// LoginErrorMsg is sent when login fails
type LoginErrorMsg struct {
	Err error
}

// LoginModel is the login screen
type LoginModel struct {
	emailInput    textinput.Model
	passwordInput textinput.Model
	focusIndex    int
	loading       bool
	err           error
	baseURL       string
	width         int
	height        int
}

// NewLoginModel creates a new login screen
func NewLoginModel(baseURL string) LoginModel {
	ei := textinput.New()
	ei.Placeholder = "user@network.net"
	ei.Focus()
	ei.CharLimit = 64
	ei.Width = 30
	ei.PromptStyle = lipgloss.NewStyle().Foreground(styles.ColorPrimary)
	ei.TextStyle = lipgloss.NewStyle().Foreground(styles.ColorBright)

	pi := textinput.New()
	pi.Placeholder = "••••••••••••"
	pi.EchoMode = textinput.EchoPassword
	pi.EchoCharacter = '•'
	pi.CharLimit = 64
	pi.Width = 30
	pi.PromptStyle = lipgloss.NewStyle().Foreground(styles.ColorPrimary)
	pi.TextStyle = lipgloss.NewStyle().Foreground(styles.ColorBright)

	return LoginModel{
		emailInput:    ei,
		passwordInput: pi,
		focusIndex:    0,
		baseURL:       baseURL,
	}
}

func (m LoginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab", "down":
			m.focusIndex = (m.focusIndex + 1) % 2
			return m, m.updateFocus()
		case "shift+tab", "up":
			m.focusIndex = (m.focusIndex - 1 + 2) % 2
			return m, m.updateFocus()
		case "enter":
			if m.focusIndex == 1 || (m.emailInput.Value() != "" && m.passwordInput.Value() != "") {
				m.loading = true
				m.err = nil
				return m, m.attemptLogin()
			}
			if m.focusIndex == 0 {
				m.focusIndex = 1
				return m, m.updateFocus()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case LoginSuccessMsg:
		m.loading = false
		return m, nil

	case LoginErrorMsg:
		m.loading = false
		m.err = msg.Err
		return m, nil
	}

	// Update the focused input
	var cmd tea.Cmd
	if m.focusIndex == 0 {
		m.emailInput, cmd = m.emailInput.Update(msg)
	} else {
		m.passwordInput, cmd = m.passwordInput.Update(msg)
	}

	return m, cmd
}

func (m LoginModel) View() string {
	w, h := SafeDimensions(m.width, m.height)
	var b strings.Builder

	// Top scan line effect
	b.WriteString(styles.ScanLineDivider(w))
	b.WriteString("\n\n")

	// ASCII Logo
	logo := styles.RenderLogo(w)
	b.WriteString(logo)
	b.WriteString("\n\n")

	// Subtitle
	subtitle := lipgloss.NewStyle().
		Foreground(styles.ColorMuted).
		Render("══════════════ NEURAL NETWORK INTERFACE v2.049 ══════════════")
	b.WriteString(lipgloss.PlaceHorizontal(w, lipgloss.Center, subtitle))
	b.WriteString("\n\n")

	// Build login form content
	var form strings.Builder

	// System status line
	form.WriteString(styles.SystemPrompt("AUTHENTICATION REQUIRED"))
	form.WriteString("\n\n")

	// Email field
	emailLabel := styles.Label.Render("├─ IDENTITY ")
	if m.focusIndex == 0 {
		emailLabel = styles.Bright.Render("├─ IDENTITY ")
	}
	form.WriteString(emailLabel)
	form.WriteString("\n")
	form.WriteString("│  ")
	form.WriteString(m.emailInput.View())
	form.WriteString("\n│\n")

	// Password field
	passLabel := styles.Label.Render("├─ ACCESS KEY ")
	if m.focusIndex == 1 {
		passLabel = styles.Bright.Render("├─ ACCESS KEY ")
	}
	form.WriteString(passLabel)
	form.WriteString("\n")
	form.WriteString("│  ")
	form.WriteString(m.passwordInput.View())
	form.WriteString("\n│\n")

	// Status line
	form.WriteString("└─ STATUS: ")
	if m.loading {
		form.WriteString(styles.Warning.Render("■ AUTHENTICATING..."))
	} else if m.err != nil {
		form.WriteString(styles.Error.Render("✖ ACCESS DENIED: " + m.err.Error()))
	} else {
		form.WriteString(styles.Success.Render("● AWAITING CREDENTIALS"))
	}
	form.WriteString("\n")

	// Wrap in titled box
	boxWidth := 60
	if w < 65 {
		boxWidth = w - 6
	}
	box := styles.TitledBox("SECURE LOGIN TERMINAL", form.String(), boxWidth)

	// Center the box
	b.WriteString(lipgloss.PlaceHorizontal(w, lipgloss.Center, box))
	b.WriteString("\n\n")

	// Help text
	help := styles.Dim.Render("[ TAB ] Navigate  [ ENTER ] Connect  [ ESC ] Disconnect")
	b.WriteString(lipgloss.PlaceHorizontal(w, lipgloss.Center, help))
	b.WriteString("\n\n")

	// Bottom scan line
	b.WriteString(styles.ScanLineDivider(w))

	// Function key bar
	fnKeys := map[string]string{
		"1":  "Help",
		"10": "Exit",
	}
	fnBar := styles.FnKeyBar(fnKeys, w)

	// Calculate vertical centering
	content := b.String()
	contentLines := strings.Count(content, "\n") + 1
	topPad := (h - contentLines - 1) / 2
	if topPad < 0 {
		topPad = 0
	}

	var result strings.Builder
	for i := 0; i < topPad; i++ {
		result.WriteString("\n")
	}
	result.WriteString(content)

	// Fill remaining space and add fn bar
	currentLines := topPad + contentLines
	for i := currentLines; i < h-1; i++ {
		result.WriteString("\n")
	}
	result.WriteString(fnBar)

	return result.String()
}

// SetSize updates the view dimensions
func (m *LoginModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *LoginModel) updateFocus() tea.Cmd {
	if m.focusIndex == 0 {
		m.passwordInput.Blur()
		return m.emailInput.Focus()
	}
	m.emailInput.Blur()
	return m.passwordInput.Focus()
}

func (m LoginModel) attemptLogin() tea.Cmd {
	return func() tea.Msg {
		resp, err := api.SignIn(m.emailInput.Value(), m.passwordInput.Value(), m.baseURL)
		if err != nil {
			return LoginErrorMsg{Err: err}
		}

		return LoginSuccessMsg{
			IDToken:      resp.IDToken,
			RefreshToken: resp.RefreshToken,
			RTDBToken:    resp.RTDBToken,
		}
	}
}
