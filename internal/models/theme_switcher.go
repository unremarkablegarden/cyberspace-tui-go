package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/messages"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/keymaps"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/ui"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// ThemeChangedMsg is sent when the user selects a new theme
type ThemeChangedMsg struct {
	ThemeKey string
}

// ThemeSwitcherModel is the popup modal for switching themes
type ThemeSwitcherModel struct {
	themes        []styles.ThemeDefinition
	cursor        int
	width         int
	height        int
	keys          keymaps.AppKeybinds
	help          help.Model
	originalTheme string // theme active when switcher was opened, for reverting on ESC
}

// NewThemeSwitcherModel creates a new theme switcher
func NewThemeSwitcherModel(keymap keymaps.AppKeybinds) ThemeSwitcherModel {
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
		keys:          keymap,
		help:          h,
		originalTheme: current,
	}
}

func (m ThemeSwitcherModel) Init() tea.Cmd {
	return nil
}

func (m ThemeSwitcherModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case m.keys.GlobalKeybinds.Back:
			// Revert to original theme
			_ = styles.ApplyTheme(m.originalTheme)
			return m, func() tea.Msg { return messages.SwitchToSettings{} }
		case m.keys.GlobalKeybinds.Down:
			if m.cursor < len(m.themes)-1 {
				m.cursor++
				m.previewTheme()
			}
		case m.keys.GlobalKeybinds.Up:
			if m.cursor > 0 {
				m.cursor--
				m.previewTheme()
			}
		case m.keys.GlobalKeybinds.Open:
			if m.cursor < len(m.themes) {
				// Theme is already applied via preview, just confirm it
				selected := m.themes[m.cursor]
				return m, func() tea.Msg {
					return messages.ThemeChangedMsg{ThemeKey: selected.Key}
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
	w, h := ui.SafeDimensions(m.width, m.height)

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
	content.WriteString(m.help.View(m.keys.ThemeSwitcherHelpKeys()))

	// Wrap in a titled box
	boxWidth := max(min(60, w-6), 40)

	box := styles.TitledBox("THEME SELECTOR", content.String(), boxWidth)

	return ui.FullScreen(box, w, h, lipgloss.Center, lipgloss.Center)
}

func (m *ThemeSwitcherModel) previewTheme() {
	if m.cursor < len(m.themes) {
		_ = styles.ApplyTheme(m.themes[m.cursor].Key)
		m.help.Styles = styles.HelpStyles()
	}
}
