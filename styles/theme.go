package styles

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ═══════════════════════════════════════════════════════════════════════════════
// 80s AMBER PHOSPHOR CRT PALETTE
// ═══════════════════════════════════════════════════════════════════════════════
//
// Monochrome amber terminal aesthetic - like classic 80s monitors
// All colors are variations of amber/orange phosphor glow

var (
	// Amber phosphor tones (bright to dim)
	ColorBright  = lipgloss.Color("220") // Bright amber (highlighted text)
	ColorPrimary = lipgloss.Color("214") // Standard amber (main text)
	ColorNormal  = lipgloss.Color("178") // Normal amber (content)
	ColorDim     = lipgloss.Color("172") // Dim amber (secondary)
	ColorMuted   = lipgloss.Color("136") // Muted amber (subtle/disabled)
	ColorDark    = lipgloss.Color("94")  // Dark amber (very subtle)

	// Semantic aliases (all amber, different intensities)
	ColorSecondary = ColorDim    // Secondary elements
	ColorAccent    = ColorBright // Accent/highlight
	ColorWarning   = ColorBright // Warnings (bright amber)
	ColorContent   = ColorNormal // Content text
	ColorText      = ColorNormal // General text

	// Status colors (still amber-based)
	ColorError   = lipgloss.Color("166") // Darker orange-red for errors
	ColorSuccess = ColorBright           // Bright amber for success
	ColorInfo    = ColorPrimary          // Standard amber for info

	// Background/UI colors
	ColorBgDark    = lipgloss.Color("232") // Near black
	ColorBgSelect  = lipgloss.Color("52")  // Dark amber/brown selection
	ColorBorder    = ColorDim              // Borders in dim amber
	ColorHighlight = lipgloss.Color("0")   // Black (for inverse video)
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

// LogoSmall is a smaller version for tighter spaces
var LogoSmall = `
▄████▄▓██   ██▓ ▄▄▄▄   ▓█████  ██▀███    ██████  ██▓███   ▄▄▄       ▄████▄  ▓█████
▒██▀ ▀█ ▒██  ██▒▓█████▄ ▓█   ▀ ▓██ ▒ ██▒▒██    ▒ ▓██░  ██▒▒████▄    ▒██▀ ▀█  ▓█   ▀
▒▓█    ▄ ▒██ ██░▒██▒ ▄██▒███   ▓██ ░▄█ ▒░ ▓██▄   ▓██░ ██▓▒▒██  ▀█▄  ▒▓█    ▄ ▒███
▒▓▓▄ ▄██▒░ ▐██▓░▒██░█▀  ▒▓█  ▄ ▒██▀▀█▄    ▒   ██▒▒██▄█▓▒ ▒░██▄▄▄▄██ ▒▓▓▄ ▄██▒▒▓█  ▄
▒ ▓███▀ ░░ ██▒▓░░▓█  ▀█▓░▒████▒░██▓ ▒██▒▒██████▒▒▒██▒ ░  ░ ▓█   ▓██▒▒ ▓███▀ ░░▒████▒
░ ░▒ ▒  ░ ██▒▒▒ ░▒▓███▀▒░░ ▒░ ░░ ▒▓ ░▒▓░▒ ▒▓▒ ▒ ░▒▓▒░ ░  ░ ▒▒   ▓▒█░░ ░▒ ▒  ░░░ ▒░ ░`

// LogoMini for very small terminals
var LogoMini = `╔═╗╦ ╦╔╗ ╔═╗╦═╗╔═╗╔═╗╔═╗╔═╗╔═╗
║  ╚╦╝╠╩╗║╣ ╠╦╝╚═╗╠═╝╠═╣║  ║╣
╚═╝ ╩ ╚═╝╚═╝╩╚═╚═╝╩  ╩ ╩╚═╝╚═╝`

// SystemBanner for headers
var SystemBanner = `┌─────────────────────────────────────────────────────────────────────────────────┐
│  ▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀  │
│  ██████╗██╗   ██╗██████╗ ███████╗██████╗ ███████╗██████╗  █████╗  ██████╗███████╗│
│  ██╔════╝╚██╗ ██╔╝██╔══██╗██╔════╝██╔══██╗██╔════╝██╔══██╗██╔══██╗██╔════╝██╔════╝│
│  ╚═════╝   ╚═╝   ╚═════╝ ╚══════╝╚═╝  ╚═╝╚══════╝╚═╝     ╚═╝  ╚═╝ ╚═════╝╚══════╝│
│  ▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄  │
└─────────────────────────────────────────────────────────────────────────────────┘`

// ═══════════════════════════════════════════════════════════════════════════════
// TEXT STYLES - AMBER PHOSPHOR
// ═══════════════════════════════════════════════════════════════════════════════

var (
	// Title styles - bright and bold for maximum phosphor glow
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorBright)

	TitleGlow = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorBright).
			Background(ColorBgDark)

	// User/Author styles - bright amber, bold
	Username = lipgloss.NewStyle().
			Foreground(ColorBright).
			Bold(true)

	// Time/date styles - dimmed phosphor
	Timestamp = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// Content text - normal amber glow
	Content = lipgloss.NewStyle().
		Foreground(ColorNormal)

	// Stats and metadata - dim
	Stats = lipgloss.NewStyle().
		Foreground(ColorDim)

	// Topics/tags - normal amber
	Topic = lipgloss.NewStyle().
		Foreground(ColorNormal)

	// Help text - muted/dim
	Help = lipgloss.NewStyle().
		Foreground(ColorMuted)

	// Labels - primary amber, bold
	Label = lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true)

	// Error messages - still visible but amber-ish
	Error = lipgloss.NewStyle().
		Foreground(ColorError).
		Bold(true)

	// Success messages - bright amber
	Success = lipgloss.NewStyle().
		Foreground(ColorBright).
		Bold(true)

	// Warning messages - bright amber bold
	Warning = lipgloss.NewStyle().
		Foreground(ColorBright).
		Bold(true)

	// Dim/subtle text - very low phosphor
	Dim = lipgloss.NewStyle().
		Foreground(ColorMuted)

	// Bright/highlighted text - full phosphor glow
	Bright = lipgloss.NewStyle().
		Foreground(ColorBright).
		Bold(true)

	// Normal text - standard amber
	Normal = lipgloss.NewStyle().
		Foreground(ColorNormal)

	// Dark/very subtle - barely visible phosphor
	Dark = lipgloss.NewStyle().
		Foreground(ColorDark)
)

// ═══════════════════════════════════════════════════════════════════════════════
// LAYOUT STYLES - AMBER PHOSPHOR
// ═══════════════════════════════════════════════════════════════════════════════

var (
	// Header style - bright amber bold
	Header = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorBright).
		Padding(0, 1)

	// Double-line border box (retro terminal style) - dim amber borders
	Box = lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(ColorDim).
		Padding(1, 2)

	// Single-line border for inner panels
	BoxSingle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorMuted).
			Padding(0, 1)

	// Inverse video selection (classic CRT style) - amber on black
	SelectedItem = lipgloss.NewStyle().
			Background(ColorPrimary).
			Foreground(ColorHighlight).
			Bold(true)

	// Subtle selection with border
	SelectedItemBorder = lipgloss.NewStyle().
				Background(ColorBgSelect).
				BorderLeft(true).
				BorderStyle(lipgloss.ThickBorder()).
				BorderForeground(ColorBright)

	// Footer - dim
	Footer = lipgloss.NewStyle().
		Foreground(ColorMuted).
		Padding(0, 1)

	// Modal dialog (for popups/alerts) - bright borders
	Modal = lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(ColorBright).
		Padding(1, 2)

	// Status bar - inverse amber
	StatusBar = lipgloss.NewStyle().
			Background(ColorPrimary).
			Foreground(ColorHighlight).
			Bold(true).
			Padding(0, 1)

	// Function key style - inverse for key number
	FnKey = lipgloss.NewStyle().
		Background(ColorDim).
		Foreground(ColorHighlight).
		Bold(true)

	FnLabel = lipgloss.NewStyle().
		Background(ColorBgSelect).
		Foreground(ColorNormal)

	// Scan line effect (decorative) - very dim
	ScanLine = lipgloss.NewStyle().
			Foreground(ColorDark)
)

// Spinner style - primary amber
var Spinner = lipgloss.NewStyle().Foreground(ColorPrimary)

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

// DoubleDivider returns a double-line horizontal divider (normal amber)
func DoubleDivider(width int) string {
	if width < 1 {
		width = 80
	}
	return lipgloss.NewStyle().
		Foreground(ColorDim).
		Render(strings.Repeat("═", width))
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

// GlitchDivider creates a "glitchy" divider for sci-fi effect (amber gradient)
func GlitchDivider(width int) string {
	if width < 1 {
		width = 80
	}
	pattern := "▓▒░"
	var result strings.Builder
	for i := 0; i < width; i++ {
		result.WriteByte(pattern[i%len(pattern)])
	}
	return lipgloss.NewStyle().
		Foreground(ColorDim).
		Render(result.String())
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
	lines := strings.Split(content, "\n")
	var middle strings.Builder
	for _, line := range lines {
		lineWidth := lipgloss.Width(line)
		padding := innerWidth - lineWidth
		if padding < 0 {
			padding = 0
		}
		middle.WriteString(borderStyle.Render("║"))
		middle.WriteString(line)
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
	lines := strings.Split(content, "\n")
	var middle strings.Builder
	for _, line := range lines {
		lineWidth := lipgloss.Width(line)
		padding := innerWidth - lineWidth
		if padding < 0 {
			padding = 0
		}
		middle.WriteString(borderStyle.Render("│"))
		middle.WriteString(line)
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

// ProgressBar creates an ASCII progress bar (amber phosphor style)
func ProgressBar(percent float64, width int) string {
	if width < 10 {
		width = 10
	}
	barWidth := width - 7

	filled := int(float64(barWidth) * percent)
	if filled > barWidth {
		filled = barWidth
	}
	empty := barWidth - filled

	// Dim brackets, bright fill, dark empty
	bar := lipgloss.NewStyle().Foreground(ColorDim).Render("[") +
		lipgloss.NewStyle().Foreground(ColorBright).Render(strings.Repeat("█", filled)) +
		lipgloss.NewStyle().Foreground(ColorDark).Render(strings.Repeat("░", empty)) +
		lipgloss.NewStyle().Foreground(ColorDim).Render("]")

	pct := lipgloss.NewStyle().Foreground(ColorNormal).Render(
		fmt.Sprintf("%3d%%", int(percent*100)),
	)

	return bar + " " + pct
}

// ProgressBarSimple creates a simple block progress bar
func ProgressBarSimple(percent float64, width int) string {
	filled := int(float64(width) * percent)
	if filled > width {
		filled = width
	}
	empty := width - filled

	// Bright fill, dark empty
	return lipgloss.NewStyle().Foreground(ColorPrimary).Render(strings.Repeat("█", filled)) +
		lipgloss.NewStyle().Foreground(ColorDark).Render(strings.Repeat("░", empty))
}

// StatusBarSegment creates a highlighted status bar segment
func StatusBarSegment(label, value string) string {
	return StatusBar.Render(label+":") +
		lipgloss.NewStyle().
			Background(ColorBgSelect).
			Foreground(ColorContent).
			Padding(0, 1).
			Render(value)
}

// MenuItem creates a menu item, optionally selected
func MenuItem(text string, selected bool, width int) string {
	textWidth := lipgloss.Width(text)
	if textWidth < width {
		text += strings.Repeat(" ", width-textWidth)
	}

	if selected {
		return SelectedItem.Render(text)
	}
	return lipgloss.NewStyle().Foreground(ColorContent).Render(text)
}

// TableHeader creates a table header row with underline (amber phosphor)
func TableHeader(columns []string, widths []int) string {
	var parts []string
	for i, col := range columns {
		w := 10
		if i < len(widths) {
			w = widths[i]
		}
		if len(col) > w {
			col = col[:w]
		} else {
			col += strings.Repeat(" ", w-len(col))
		}
		parts = append(parts, col)
	}

	// Bright amber for headers
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorBright).
		Render(strings.Join(parts, " "))

	totalWidth := 0
	for _, w := range widths {
		totalWidth += w + 1
	}
	// Dim amber for underline
	underline := lipgloss.NewStyle().
		Foreground(ColorMuted).
		Render(strings.Repeat("─", totalWidth))

	return header + "\n" + underline
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

	top := borderStyle.Render("┌─ ") + titleStyle.Render(icon) + borderStyle.Render(" " + strings.Repeat("─", innerWidth-len(icon)-2) + "┐")
	mid := borderStyle.Render("│ ") + lipgloss.NewStyle().Foreground(ColorNormal).Render(message) + strings.Repeat(" ", innerWidth-len(message)) + borderStyle.Render(" │")
	bottom := borderStyle.Render("└" + strings.Repeat("─", innerWidth+2) + "┘")

	return top + "\n" + mid + "\n" + bottom
}

// Blinker returns a blinking cursor character (for animation)
func Blinker(on bool) string {
	if on {
		return lipgloss.NewStyle().Foreground(ColorBright).Render("█")
	}
	return " "
}

// DataField renders a label: value pair (amber phosphor style)
func DataField(label, value string) string {
	return lipgloss.NewStyle().Foreground(ColorDim).Render(label+": ") +
		lipgloss.NewStyle().Foreground(ColorBright).Render(value)
}
