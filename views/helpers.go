package views

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"

	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
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

// Truncate shortens a string to max visual width with ellipsis
func Truncate(s string, max int) string {
	if lipgloss.Width(s) <= max {
		return s
	}
	if max <= 3 {
		max = 3
	}
	target := max - 3
	width := 0
	for i, r := range s {
		w := runewidth.RuneWidth(r)
		if width+w > target {
			return s[:i] + "..."
		}
		width += w
	}
	return s + "..."
}

var (
	reLink       = regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
	reBold       = regexp.MustCompile(`\*\*(.+?)\*\*`)
	reBoldUndsc  = regexp.MustCompile(`__(.+?)__`)
	reItalic     = regexp.MustCompile(`\*(.+?)\*`)
	reItalUndsc  = regexp.MustCompile(`\b_(.+?)_\b`)
	reCode       = regexp.MustCompile("`([^`]+)`")
	reHeading    = regexp.MustCompile(`(?m)^#{1,6}\s+`)
	reCodeBlock  = regexp.MustCompile("(?s)```[a-z]*\n?(.*?)```")
)

// stripMarkdownCommon applies shared markdown stripping rules
func stripMarkdownCommon(s string) string {
	s = reCodeBlock.ReplaceAllString(s, "$1")
	s = reLink.ReplaceAllString(s, "$1")
	s = reBold.ReplaceAllString(s, "$1")
	s = reBoldUndsc.ReplaceAllString(s, "$1")
	s = reItalic.ReplaceAllString(s, "$1")
	s = reItalUndsc.ReplaceAllString(s, "$1")
	s = reCode.ReplaceAllString(s, "$1")
	s = reHeading.ReplaceAllString(s, "")
	// Replace HTML entities
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")
	s = ReplaceEmojis(s)
	return s
}

// StripMarkdown removes basic markdown formatting for plain text display (single line)
func StripMarkdown(s string) string {
	s = stripMarkdownCommon(s)
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.TrimSpace(s)
}

// StripMarkdownKeepNewlines removes markdown formatting but preserves line breaks
func StripMarkdownKeepNewlines(s string) string {
	s = stripMarkdownCommon(s)
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

// CenterText centers text within a given width
func CenterText(text string, width int) string {
	textWidth := lipgloss.Width(text)
	if textWidth >= width {
		return text
	}
	padding := (width - textWidth) / 2
	return strings.Repeat(" ", padding) + text
}

// styleLines applies a lipgloss style to each line of a multi-line string.
// This is needed because lipgloss.Render on a multi-line string only colors the first line.
func styleLines(s string, style lipgloss.Style) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = style.Render(line)
	}
	return strings.Join(lines, "\n")
}

// PadRight pads text to a given width
func PadRight(text string, width int) string {
	textWidth := lipgloss.Width(text)
	if textWidth >= width {
		return text
	}
	return text + strings.Repeat(" ", width-textWidth)
}

// ═══════════════════════════════════════════════════════════════════════════════
// EMOJI REPLACEMENT
// ═══════════════════════════════════════════════════════════════════════════════

// Plain Unicode symbols to replace emojis with — matching the web app aesthetic
var emojiReplacements = []rune{
	// Block elements
	'▀', '▄', '█', '▉', '▊', '▋', '▌', '▍', '▎', '▏', '▐', '░', '▒', '▓',
	// Mathematical
	'÷', '≠', '∑', '∏', '∫', '√', '∞', '∂', '∇', '∆', '∝', '∠',
	// Box-drawing
	'┼', '║', '╔', '╗', '╠', '╣', '╦', '╬',
	// Misc symbols
	'§', '¶', '†', '‡',
	'♠', '♣', '♥', '♦', '◊', '○', '●', '◐', '◑',
	'■', '□', '▲', '△', '▼', '▽', '◆', '◇',
	'★', '☆', '✦', '✧', '✩', '✪', '✫', '✬', '✭', '✮',
	'✱', '✲', '✳', '✴', '✵', '✶', '✷', '✸', '✹', '✺', '✻', '✼', '✽', '✾', '✿',
	'❀', '❁', '❂', '❃', '❄', '❅', '❆', '❇', '❈', '❉', '❊', '❋',
	'❍', '❏', '❐', '❑', '❒', '❖',
	'❡', '❢', '❣', '❤', '❥', '❦', '❧',
	// Geometric shapes
	'◧', '◨', '◩', '◪', '◫', '◬', '◭', '◮', '◯',
	'◰', '◱', '◲', '◳', '◴', '◵', '◶', '◷',
	'◸', '◹', '◺', '◻', '◼', '◽', '◾', '◿',
	// Braille patterns
	'⣀', '⣤', '⣶', '⣿', '⡀', '⡤', '⡶', '⡿', '⢀', '⢤', '⢶', '⢿',
}

// isEmoji returns true if the rune is in a common emoji Unicode range.
func isEmoji(r rune) bool {
	return unicode.Is(unicode.So, r) && (false ||
		(r >= 0x1F300 && r <= 0x1F9FF) || // Misc symbols, emoticons, supplemental
		(r >= 0x1FA00 && r <= 0x1FAFF) || // Symbols extended-A
		(r >= 0x2600 && r <= 0x26FF) || // Misc symbols
		(r >= 0x2700 && r <= 0x27BF) || // Dingbats
		(r >= 0x1F1E0 && r <= 0x1F1FF) || // Regional indicators (flags)
		(r >= 0xFE00 && r <= 0xFE0F) || // Variation selectors
		(r >= 0x200D && r <= 0x200D) || // ZWJ
		(r >= 0xE0020 && r <= 0xE007F)) // Tags
}

// ReplaceEmojis replaces emoji characters with deterministic plain Unicode symbols.
// Each unique emoji codepoint always maps to the same replacement, so the result
// is stable across re-renders.
func ReplaceEmojis(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i, r := range s {
		if isEmoji(r) {
			// Deterministic: hash the rune value and position to pick a stable replacement
			idx := (int(r)*31 + i*7) % len(emojiReplacements)
			if idx < 0 {
				idx = -idx
			}
			b.WriteRune(emojiReplacements[idx])
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
