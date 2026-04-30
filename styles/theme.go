package styles

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// ═══════════════════════════════════════════════════════════════════════════════
// THEME SYSTEM
// ═══════════════════════════════════════════════════════════════════════════════

// ThemeColors defines the color palette for a theme
type ThemeColors struct {
	Bright    string `json:"bright"`
	Primary   string `json:"primary"`
	Normal    string `json:"normal"`
	Dim       string `json:"dim"`
	Muted     string `json:"muted"`
	Dark      string `json:"dark"`
	Error     string `json:"error"`
	BgDark    string `json:"bg_dark"`
	BgSelect  string `json:"bg_select"`
	Highlight string `json:"highlight"`
}

// ThemeDefinition is a complete theme loaded from JSON
type ThemeDefinition struct {
	Key         string      `json:"-"` // filename without extension, used for ApplyTheme
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Colors      ThemeColors `json:"colors"`
}

//go:embed themes/*.json
var themesFS embed.FS

var (
	currentTheme = "dark"
)

// InitThemes stores the embedded filesystem and applies the default theme styles.
func InitThemes(f fs.FS) {
	rebuildStyles()
}

// CurrentThemeName returns the name of the currently active theme
func CurrentThemeName() string {
	return currentTheme
}

// ListThemes returns all available theme definitions
func ListThemes() []ThemeDefinition {
	var themes []ThemeDefinition
	entries, err := fs.ReadDir(themesFS, "themes")
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		data, err := fs.ReadFile(themesFS, "themes/"+entry.Name())
		if err != nil {
			continue
		}
		var def ThemeDefinition
		if err := json.Unmarshal(data, &def); err != nil {
			continue
		}
		def.Key = strings.TrimSuffix(entry.Name(), ".json")
		themes = append(themes, def)
	}
	return themes
}

// ApplyTheme loads a theme by filename (without .json extension) and applies it
func ApplyTheme(name string) error {
	data, err := fs.ReadFile(themesFS, "themes/"+name+".json")
	if err != nil {
		return fmt.Errorf("theme %q not found: %w", name, err)
	}

	var def ThemeDefinition
	if err := json.Unmarshal(data, &def); err != nil {
		return fmt.Errorf("invalid theme %q: %w", name, err)
	}

	// Apply base colors
	ColorBright = lipgloss.Color(def.Colors.Bright)
	ColorPrimary = lipgloss.Color(def.Colors.Primary)
	ColorNormal = lipgloss.Color(def.Colors.Normal)
	ColorDim = lipgloss.Color(def.Colors.Dim)
	ColorMuted = lipgloss.Color(def.Colors.Muted)
	ColorDark = lipgloss.Color(def.Colors.Dark)
	ColorError = lipgloss.Color(def.Colors.Error)
	ColorBgDark = lipgloss.Color(def.Colors.BgDark)
	ColorBgSelect = lipgloss.Color(def.Colors.BgSelect)
	ColorHighlight = lipgloss.Color(def.Colors.Highlight)

	rebuildStyles()
	currentTheme = name
	return nil
}

// rebuildStyles reassigns all style vars from the current color vars.
func rebuildStyles() {
	Title = lipgloss.NewStyle().Bold(true).Foreground(ColorBright)
	Username = lipgloss.NewStyle().Foreground(ColorBright).Bold(true)
	Label = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	Error = lipgloss.NewStyle().Foreground(ColorError).Bold(true)
	Success = lipgloss.NewStyle().Foreground(ColorBright).Bold(true)
	Warning = lipgloss.NewStyle().Foreground(ColorBright).Bold(true)
	Dim = lipgloss.NewStyle().Foreground(ColorMuted)
	Bright = lipgloss.NewStyle().Foreground(ColorBright).Bold(true)
	Normal = lipgloss.NewStyle().Foreground(ColorNormal)
	Dark = lipgloss.NewStyle().Foreground(ColorDark)
	SelectedItem = lipgloss.NewStyle().Background(ColorPrimary).Foreground(ColorHighlight).Bold(true)
	FnKey = lipgloss.NewStyle().Background(ColorDim).Foreground(ColorHighlight).Bold(true)
	FnLabel = lipgloss.NewStyle().Background(ColorBgSelect).Foreground(ColorNormal)
	Spinner = lipgloss.NewStyle().Foreground(ColorPrimary)
}

// ═══════════════════════════════════════════════════════════════════════════════
// COLOR PALETTE — defaults for the "dark" theme; overwritten by ApplyTheme
// ═══════════════════════════════════════════════════════════════════════════════

var (
	ColorBright    = lipgloss.Color("229")
	ColorPrimary   = lipgloss.Color("223")
	ColorNormal    = lipgloss.Color("222")
	ColorDim       = lipgloss.Color("180")
	ColorMuted     = lipgloss.Color("137")
	ColorDark      = lipgloss.Color("94")
	ColorError     = lipgloss.Color("166")
	ColorBgDark    = lipgloss.Color("232")
	ColorBgSelect  = lipgloss.Color("236")
	ColorHighlight = lipgloss.Color("0")
)

// ═══════════════════════════════════════════════════════════════════════════════
// STYLES — set by rebuildStyles(), called from InitThemes() and ApplyTheme()
// ═══════════════════════════════════════════════════════════════════════════════

var (
	Title        lipgloss.Style
	Username     lipgloss.Style
	Label        lipgloss.Style
	Error        lipgloss.Style
	Success      lipgloss.Style
	Warning      lipgloss.Style
	Dim          lipgloss.Style
	Bright       lipgloss.Style
	Normal       lipgloss.Style
	Dark         lipgloss.Style
	SelectedItem lipgloss.Style
	FnKey        lipgloss.Style
	FnLabel      lipgloss.Style
	Spinner      lipgloss.Style
)

// ═══════════════════════════════════════════════════════════════════════════════
// ASCII ART & LOGOS
// ═══════════════════════════════════════════════════════════════════════════════

// Logo is the main CYBERSPACE ASCII art logo
var Logo = `
 ██████╗██╗   ██╗██████╗ ███████╗██████╗ ███████╗██████╗  █████╗  ██████╗███████╗
██╔════╝╚██╗ ██╔╝██╔══██╗██╔════╝██╔══██╗██╔════╝██╔══██╗██╔══██╗██╔════╝██╔════╝
██║      ╚████╔╝ ██████╔╝█████╗  ██████╔╝███████╗██████╔╝███████║██║     █████╗
██║       ╚██╔╝  ██╔══██╗██╔══╝  ██╔══██╗╚════██║██╔═══╝ ██╔══██║██║     ██╔══╝
╚██████╗   ██║   ██████╔╝███████╗██║  ██║███████║██║     ██║  ██║╚██████╗███████╗
 ╚═════╝   ╚═╝   ╚═════╝ ╚══════╝╚═╝  ╚═╝╚══════╝╚═╝     ╚═╝  ╚═╝ ╚═════╝╚══════╝`

// LogoMini for very small terminals
var LogoMini = `╔═╗╦ ╦╔╗ ╔═╗╦═╗╔═╗╔═╗╔═╗╔═╗╔═╗
║  ╚╦╝╠╩╗║╣ ╠╦╝╚═╗╠═╝╠═╣║  ║╣
╚═╝ ╩ ╚═╝╚═╝╩╚═╚═╝╩  ╩ ╩╚═╝╚═╝`

// ═══════════════════════════════════════════════════════════════════════════════
// HELPER FUNCTIONS
// ═══════════════════════════════════════════════════════════════════════════════

// RenderLogo renders the appropriate logo size based on terminal width
func RenderLogo(width int) string {
	// Bright amber for the logo - maximum phosphor glow
	logoStyle := lipgloss.NewStyle().Foreground(ColorBright).Bold(true)

	if width >= 85 {
		return logoStyle.Render(Logo)
	} else if width >= 60 {
		return logoStyle.Render(LogoMini)
	}
	return logoStyle.Render("[ CYBERSPACE ]")
}

// Divider returns a single-line horizontal divider (dim amber)
func Divider(width int) string {
	if width < 1 {
		width = 80
	}
	return lipgloss.NewStyle().
		Foreground(ColorMuted).
		Render(strings.Repeat("─", width))
}

// ScanLineDivider creates a decorative scan-line effect divider (dark amber)
func ScanLineDivider(width int) string {
	if width < 1 {
		width = 80
	}
	return lipgloss.NewStyle().
		Foreground(ColorDark).
		Render(strings.Repeat("░", width))
}

// TitledBox creates a retro-style box with title embedded in the top border
func TitledBox(title, content string, width int) string {
	if width < len(title)+6 {
		width = len(title) + 6
	}

	// Dim amber for borders, bright for title
	borderStyle := lipgloss.NewStyle().Foreground(ColorDim)
	titleStyle := lipgloss.NewStyle().Foreground(ColorBright).Bold(true)

	innerWidth := width - 2
	titleLen := len(title) + 4 // "[ " + title + " ]"
	leftPad := (innerWidth - titleLen) / 2
	rightPad := innerWidth - titleLen - leftPad

	// Build top border with embedded title
	top := borderStyle.Render("╔") +
		borderStyle.Render(strings.Repeat("═", leftPad)) +
		borderStyle.Render("[ ") +
		titleStyle.Render(title) +
		borderStyle.Render(" ]") +
		borderStyle.Render(strings.Repeat("═", rightPad)) +
		borderStyle.Render("╗")

	// Build bottom border
	bottom := borderStyle.Render("╚") +
		borderStyle.Render(strings.Repeat("═", innerWidth)) +
		borderStyle.Render("╝")

	// Wrap content lines with vertical borders
	contentStyle := lipgloss.NewStyle().Foreground(ColorNormal)
	lines := strings.Split(content, "\n")
	var middle strings.Builder
	for _, line := range lines {
		styled := contentStyle.Render(line)
		lineWidth := lipgloss.Width(styled)
		padding := innerWidth - lineWidth
		if padding < 0 {
			padding = 0
		}
		middle.WriteString(borderStyle.Render("║"))
		middle.WriteString(styled)
		middle.WriteString(strings.Repeat(" ", padding))
		middle.WriteString(borderStyle.Render("║"))
		middle.WriteString("\n")
	}

	return top + "\n" + middle.String() + bottom
}

// DataBox creates a sci-fi "data terminal" box
func DataBox(title, content string, width int) string {
	if width < len(title)+10 {
		width = len(title) + 10
	}

	// Muted amber for borders, bright for title
	borderStyle := lipgloss.NewStyle().Foreground(ColorMuted)
	titleStyle := lipgloss.NewStyle().Foreground(ColorBright).Bold(true)
	cornerStyle := lipgloss.NewStyle().Foreground(ColorDim)

	innerWidth := width - 2

	// Top border with title
	top := cornerStyle.Render("┌") +
		borderStyle.Render("──┤ ") +
		titleStyle.Render(title) +
		borderStyle.Render(" ├") +
		borderStyle.Render(strings.Repeat("─", innerWidth-len(title)-6)) +
		cornerStyle.Render("┐")

	// Bottom border
	bottom := cornerStyle.Render("└") +
		borderStyle.Render(strings.Repeat("─", innerWidth)) +
		cornerStyle.Render("┘")

	// Wrap content
	contentStyle := lipgloss.NewStyle().Foreground(ColorNormal)
	lines := strings.Split(content, "\n")
	var middle strings.Builder
	for _, line := range lines {
		styled := contentStyle.Render(line)
		lineWidth := lipgloss.Width(styled)
		padding := innerWidth - lineWidth
		if padding < 0 {
			padding = 0
		}
		middle.WriteString(borderStyle.Render("│"))
		middle.WriteString(styled)
		middle.WriteString(strings.Repeat(" ", padding))
		middle.WriteString(borderStyle.Render("│"))
		middle.WriteString("\n")
	}

	return top + "\n" + middle.String() + bottom
}

// FnKeyBar creates a function key bar
func FnKeyBar(keys map[string]string, width int) string {
	var parts []string
	order := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	for _, k := range order {
		if label, ok := keys[k]; ok {
			key := FnKey.Render(k)
			lbl := FnLabel.Render(label)
			parts = append(parts, key+lbl)
		}
	}

	bar := strings.Join(parts, " ")
	barWidth := lipgloss.Width(bar)
	if barWidth < width {
		bar += strings.Repeat(" ", width-barWidth)
	}

	return lipgloss.NewStyle().
		Background(ColorBgSelect).
		Width(width).
		Render(bar)
}

// SystemPrompt creates a system prompt prefix (amber phosphor)
func SystemPrompt(text string) string {
	prompt := lipgloss.NewStyle().
		Foreground(ColorBright).
		Bold(true).
		Render("▶ SYSTEM:")
	message := lipgloss.NewStyle().
		Foreground(ColorNormal).
		Render(" " + text)
	return prompt + message
}

// AlertBox creates a highlighted alert/warning box (amber themed)
func AlertBox(message string, alertType string, width int) string {
	var color lipgloss.Color
	var icon string

	// All alerts use amber tones, just different intensities
	switch alertType {
	case "error":
		color = ColorError // Slightly different orange for errors
		icon = "✖ ERROR"
	case "warning":
		color = ColorBright
		icon = "⚠ WARNING"
	case "success":
		color = ColorBright
		icon = "✔ SUCCESS"
	default:
		color = ColorNormal
		icon = "ℹ INFO"
	}

	borderStyle := lipgloss.NewStyle().Foreground(color)
	titleStyle := lipgloss.NewStyle().Foreground(color).Bold(true)

	innerWidth := width - 4
	if innerWidth < len(message) {
		innerWidth = len(message) + 2
	}

	top := borderStyle.Render("┌─ ") + titleStyle.Render(icon) + borderStyle.Render(" "+strings.Repeat("─", innerWidth-len(icon)-2)+"┐")
	mid := borderStyle.Render("│ ") + lipgloss.NewStyle().Foreground(ColorNormal).Render(message) + strings.Repeat(" ", innerWidth-len(message)) + borderStyle.Render(" │")
	bottom := borderStyle.Render("└" + strings.Repeat("─", innerWidth+2) + "┘")

	return top + "\n" + mid + "\n" + bottom
}

// ListStyles returns list.Styles themed to the current color palette.
func ListStyles() list.Styles {
	s := list.DefaultStyles()
	s.TitleBar = lipgloss.NewStyle().
		Background(ColorBgDark).
		Foreground(ColorBright).
		Padding(0, 1)
	s.Title = lipgloss.NewStyle().
		Foreground(ColorBright).
		Bold(true)
	s.Spinner = lipgloss.NewStyle().Foreground(ColorPrimary)
	s.FilterPrompt = lipgloss.NewStyle().Foreground(ColorBright)
	s.FilterCursor = lipgloss.NewStyle().Foreground(ColorBright)
	s.DefaultFilterCharacterMatch = lipgloss.NewStyle().Foreground(ColorBright).Bold(true)
	s.StatusBar = lipgloss.NewStyle().Foreground(ColorMuted)
	s.StatusEmpty = lipgloss.NewStyle().Foreground(ColorDim)
	s.StatusBarActiveFilter = lipgloss.NewStyle().Foreground(ColorBright)
	s.StatusBarFilterCount = lipgloss.NewStyle().Foreground(ColorMuted)
	s.NoItems = lipgloss.NewStyle().Foreground(ColorDim)
	s.PaginationStyle = lipgloss.NewStyle().PaddingLeft(2)
	s.HelpStyle = lipgloss.NewStyle().PaddingLeft(2)
	s.ActivePaginationDot = lipgloss.NewStyle().Foreground(ColorBright)
	s.InactivePaginationDot = lipgloss.NewStyle().Foreground(ColorDark)
	s.DividerDot = lipgloss.NewStyle().Foreground(ColorDark)
	return s
}

// HelpStyles returns help.Styles themed to the current color palette.
// Call after ApplyTheme or on init to keep help bubble in sync.
func HelpStyles() help.Styles {
	return help.Styles{
		ShortKey:       lipgloss.NewStyle().Foreground(ColorBright).Bold(true),
		ShortDesc:      lipgloss.NewStyle().Foreground(ColorMuted),
		ShortSeparator: lipgloss.NewStyle().Foreground(ColorDark),
		Ellipsis:       lipgloss.NewStyle().Foreground(ColorDark),
		FullKey:        lipgloss.NewStyle().Foreground(ColorBright).Bold(true),
		FullDesc:       lipgloss.NewStyle().Foreground(ColorMuted),
		FullSeparator:  lipgloss.NewStyle().Foreground(ColorDark),
	}
}
