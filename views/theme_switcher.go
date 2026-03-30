package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// ThemeChangedMsg is sent when the user selects a new theme
type ThemeChangedMsg struct {
	ThemeKey string
}

// ThemeSwitcherClosedMsg is sent when the user closes the theme switcher
type ThemeSwitcherClosedMsg struct{}

// ThemeSwitcherModel is the popup modal for switching themes
type ThemeSwitcherModel struct {
	themes       []styles.ThemeDefinition
	cursor       int
	width        int
	height       int
	keys         ThemeSwitcherKeyMap
	help         help.Model
	originalTheme string // theme active when switcher was opened, for reverting on ESC
}

// NewThemeSwitcherModel creates a new theme switcher
func NewThemeSwitcherModel() ThemeSwitcherModel {
	themes := styles.ListThemes()
	current := styles.CurrentThemeName()

	// Set cursor to current theme
	cursor := 0
	for i, t := range themes {
		if t.Key == current {
			cursor = i
			break
		}
	}

	h := help.New()
	h.Styles = styles.HelpStyles()
	return ThemeSwitcherModel{
		themes:        themes,
		cursor:        cursor,
		keys:          NewThemeSwitcherKeyMap(),
		help:          h,
		originalTheme: current,
	}
}

func (m ThemeSwitcherModel) Init() tea.Cmd {
	return nil
}

func (m ThemeSwitcherModel) Update(msg tea.Msg) (ThemeSwitcherModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Close):
			// Revert to original theme
			_ = styles.ApplyTheme(m.originalTheme)
			return m, func() tea.Msg { return ThemeSwitcherClosedMsg{} }
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.themes)-1 {
				m.cursor++
				m.previewTheme()
			}
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
				m.previewTheme()
			}
		case key.Matches(msg, m.keys.Apply):
			if m.cursor < len(m.themes) {
				// Theme is already applied via preview, just confirm it
				selected := m.themes[m.cursor]
				return m, func() tea.Msg {
					return ThemeChangedMsg{ThemeKey: selected.Key}
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m ThemeSwitcherModel) View() string {
	w, h := SafeDimensions(m.width, m.height)

	var content strings.Builder

	content.WriteString(styles.SystemPrompt("SELECT DISPLAY THEME"))
	content.WriteString("\n\n")

	current := styles.CurrentThemeName()

	for i, theme := range m.themes {
		selected := i == m.cursor
		active := theme.Key == current

		// Build the line
		var line strings.Builder
		if selected {
			line.WriteString(styles.Bright.Render("▸ "))
		} else {
			line.WriteString("  ")
		}

		name := theme.Name
		if active {
			name += " ●"
		}

		if selected {
			line.WriteString(styles.SelectedItem.Render(fmt.Sprintf(" %-20s ", name)))
		} else {
			line.WriteString(styles.Normal.Render(fmt.Sprintf(" %-20s ", name)))
		}

		line.WriteString("  ")
		line.WriteString(styles.Dim.Render(theme.Description))

		content.WriteString(line.String())
		content.WriteString("\n")
	}

	content.WriteString("\n")
	content.WriteString(m.help.View(m.keys))

	// Wrap in a titled box
	boxWidth := 60
	if w < 65 {
		boxWidth = w - 6
	}
	if boxWidth < 40 {
		boxWidth = 40
	}

	box := styles.TitledBox("THEME SELECTOR", content.String(), boxWidth)

	return FullScreen(box, w, h, lipgloss.Center, lipgloss.Center)
}

func (m *ThemeSwitcherModel) previewTheme() {
	if m.cursor < len(m.themes) {
		_ = styles.ApplyTheme(m.themes[m.cursor].Key)
		m.help.Styles = styles.HelpStyles()
	}
}

// SetSize updates the view dimensions
func (m *ThemeSwitcherModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
