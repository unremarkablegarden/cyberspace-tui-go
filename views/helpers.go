package views

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"

	"github.com/euklides/cyberspace-cli/styles"
)

// TimeAgo formats a time as a relative string (e.g., "5m", "2h", "3d")
func TimeAgo(t time.Time) string {
	d := time.Since(t)

	switch {
	case d < time.Minute:
		return "now"
	case d < time.Hour:
		m := int(d.Minutes())
		return fmt.Sprintf("%dm", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		return fmt.Sprintf("%dh", h)
	case d < 7*24*time.Hour:
		days := int(d.Hours() / 24)
		return fmt.Sprintf("%dd", days)
	default:
		return t.Format("Jan 2")
	}
}

// Truncate shortens a string to max length with ellipsis
func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

// StripMarkdown removes basic markdown formatting for plain text display (single line)
func StripMarkdown(s string) string {
	// Convert markdown links [text](url) to just the text
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
	s = linkRegex.ReplaceAllString(s, "$1")

	// Remove other markdown formatting
	s = strings.ReplaceAll(s, "**", "")
	s = strings.ReplaceAll(s, "__", "")
	s = strings.ReplaceAll(s, "*", "")
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, "`", "")
	s = strings.ReplaceAll(s, "#", "")
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.TrimSpace(s)
}

// StripMarkdownKeepNewlines removes markdown formatting but preserves line breaks
func StripMarkdownKeepNewlines(s string) string {
	// Convert markdown links [text](url) to just the text
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
	s = linkRegex.ReplaceAllString(s, "$1")

	// Remove other markdown formatting
	s = strings.ReplaceAll(s, "**", "")
	s = strings.ReplaceAll(s, "__", "")
	s = strings.ReplaceAll(s, "*", "")
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, "`", "")
	s = strings.ReplaceAll(s, "#", "")
	return strings.TrimSpace(s)
}

// Min returns the smaller of two integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the larger of two integers
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Clamp restricts a value to a range
func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// SafeWidth returns a width that's safe for rendering (minimum 1)
func SafeWidth(width, defaultWidth int) int {
	if width < 1 {
		return defaultWidth
	}
	return width
}

// SafeDimensions returns width and height with sensible defaults
// Use this before WindowSizeMsg has been received
func SafeDimensions(width, height int) (int, int) {
	if width < 10 {
		width = 80
	}
	if height < 10 {
		height = 24
	}
	return width, height
}

// FullScreen renders content in a full-screen container
func FullScreen(content string, width, height int, hAlign, vAlign lipgloss.Position) string {
	w, h := SafeDimensions(width, height)

	// Count current lines
	lines := strings.Split(content, "\n")
	contentHeight := len(lines)

	// Calculate padding needed
	var topPad, bottomPad int
	if vAlign == lipgloss.Center {
		topPad = (h - contentHeight) / 2
		bottomPad = h - contentHeight - topPad
	} else if vAlign == lipgloss.Bottom {
		topPad = h - contentHeight
	} else { // Top
		bottomPad = h - contentHeight
	}

	// Ensure non-negative
	if topPad < 0 {
		topPad = 0
	}
	if bottomPad < 0 {
		bottomPad = 0
	}

	// Build padded content
	var b strings.Builder

	// Top padding
	for i := 0; i < topPad; i++ {
		b.WriteString(strings.Repeat(" ", w))
		b.WriteString("\n")
	}

	// Content - pad each line to full width
	for i, line := range lines {
		lineLen := lipgloss.Width(line)
		var padLeft, padRight int

		if hAlign == lipgloss.Center {
			padLeft = (w - lineLen) / 2
			padRight = w - lineLen - padLeft
		} else if hAlign == lipgloss.Right {
			padLeft = w - lineLen
		} else { // Left
			padRight = w - lineLen
		}

		if padLeft < 0 {
			padLeft = 0
		}
		if padRight < 0 {
			padRight = 0
		}

		b.WriteString(strings.Repeat(" ", padLeft))
		b.WriteString(line)
		b.WriteString(strings.Repeat(" ", padRight))
		if i < len(lines)-1 {
			b.WriteString("\n")
		}
	}

	// Bottom padding
	for i := 0; i < bottomPad; i++ {
		b.WriteString("\n")
		b.WriteString(strings.Repeat(" ", w))
	}

	return b.String()
}

// NewSpinner creates a sci-fi styled spinner
func NewSpinner() spinner.Model {
	s := spinner.New()
	// Use a sci-fi looking spinner
	s.Spinner = spinner.Spinner{
		Frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		FPS:    time.Millisecond * 80,
	}
	s.Style = styles.Spinner
	return s
}

// NewDataSpinner creates an alternative data-transfer style spinner
func NewDataSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Spinner{
		Frames: []string{"▓▒░", "░▓▒", "▒░▓"},
		FPS:    time.Millisecond * 150,
	}
	s.Style = styles.Spinner
	return s
}

// RenderHeader renders a sci-fi styled header with title and help text
func RenderHeader(title, help string, width int) string {
	var b strings.Builder

	// Title with sci-fi decoration
	titleStyled := styles.Title.Render("▓▒░ " + title + " ░▒▓")
	b.WriteString(titleStyled)

	if help != "" {
		b.WriteString("  ")
		b.WriteString(styles.Help.Render(help))
	}
	b.WriteString("\n")

	// Double divider
	dividerWidth := width
	if dividerWidth < 1 {
		dividerWidth = 80
	}
	b.WriteString(styles.DoubleDivider(dividerWidth))
	b.WriteString("\n")

	return b.String()
}

// RenderError renders a sci-fi styled error message centered on screen
func RenderError(err error, hint string, width, height int) string {
	errorBox := styles.AlertBox(err.Error(), "error", 50)
	if hint != "" {
		errorBox += "\n\n" + styles.Dim.Render(hint)
	}
	return FullScreen(errorBox, width, height, lipgloss.Center, lipgloss.Center)
}

// RenderLoading renders a sci-fi styled loading message centered on screen
func RenderLoading(s spinner.Model, message string, width, height int) string {
	loadingContent := styles.DataBox("PROCESSING",
		"\n"+
			"  "+s.View()+" "+message+"\n"+
			"\n"+
			"  "+styles.ProgressBarSimple(0.3, 25)+"\n"+
			"\n",
		45)
	return FullScreen(loadingContent, width, height, lipgloss.Center, lipgloss.Center)
}

// FormatTopics formats a topics array with sci-fi brackets
func FormatTopics(topics []string) string {
	if len(topics) == 0 {
		return ""
	}
	return styles.Topic.Render("⟨" + strings.Join(topics, "⟩ ⟨") + "⟩")
}

// FormatStats formats reply and bookmark counts with sci-fi icons
func FormatStats(replies, bookmarks int) string {
	return styles.Stats.Render(fmt.Sprintf("◈ %d  ◆ %d", replies, bookmarks))
}

// FormatAuthor formats username and timestamp in sci-fi style
func FormatAuthor(username string, createdAt time.Time) string {
	return fmt.Sprintf("%s %s",
		styles.Username.Render("@"+username),
		styles.Timestamp.Render("· "+TimeAgo(createdAt)))
}

// CenterText centers text within a given width
func CenterText(text string, width int) string {
	textWidth := lipgloss.Width(text)
	if textWidth >= width {
		return text
	}
	padding := (width - textWidth) / 2
	return strings.Repeat(" ", padding) + text
}

// PadRight pads text to a given width
func PadRight(text string, width int) string {
	textWidth := lipgloss.Width(text)
	if textWidth >= width {
		return text
	}
	return text + strings.Repeat(" ", width-textWidth)
}
