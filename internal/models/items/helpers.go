package items

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"

	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

var (
	reLink        = regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
	reBold        = regexp.MustCompile(`\*\*(.+?)\*\*`)
	reBoldUndsc   = regexp.MustCompile(`__(.+?)__`)
	reItalic      = regexp.MustCompile(`\*(.+?)\*`)
	reItalUndsc   = regexp.MustCompile(`\b_(.+?)_\b`)
	reCode        = regexp.MustCompile("`([^`]+)`")
	reHeading     = regexp.MustCompile(`(?m)^#{1,6}\s+`)
	reCodeBlock   = regexp.MustCompile("(?s)```[a-z]*\n?(.*?)```")
	reNbspLine    = regexp.MustCompile(`(?m)^[ \t]*&nbsp;[ \t]*$`)
	reNbspUniLine = regexp.MustCompile("(?m)^[ \t]*\u00A0[ \t]*$")
	reMultiBlank  = regexp.MustCompile(`\n{3,}`)
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
// returns a boolean to indicate if was truncated or not
func Truncate(s string, max int) (string, bool) {
	if lipgloss.Width(s) <= max {
		return s, false
	}
	if max <= 3 {
		max = 3
	}
	target := max - 3
	width := 0
	for i, r := range s {
		w := runewidth.RuneWidth(r)
		if width+w > target {
			return s[:i] + "...", true
		}
		width += w
	}
	return s + "...", true
}

// WrapText wraps text to fit within a visual width, splitting on word boundaries.
func WrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var currentLine string
	var currentWidth int

	for _, word := range words {
		wordWidth := lipgloss.Width(word)

		if currentWidth+1+wordWidth <= width {
			currentLine += " " + word
			currentWidth += 1 + wordWidth
			continue
		}

		lines = append(lines, currentLine)
		currentLine = word
		currentWidth = wordWidth
	}

	lines = append(lines, currentLine)

	return lines
}

// cleanContent removes &nbsp;-only lines and collapses excessive blank lines,
// matching the Nuxt4 web app's useMarkdownRenderer sanitization.
func cleanContent(s string) string {
	s = reNbspLine.ReplaceAllString(s, "")
	s = reNbspUniLine.ReplaceAllString(s, "")
	s = reMultiBlank.ReplaceAllString(s, "\n\n")
	return s
}

// stripMarkdownCommon applies shared markdown stripping rules
func stripMarkdownCommon(s string) string {
	s = cleanContent(s)
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

// BuildCardBox renders content in a bordered card with rounded corners.
func BuildCardBox(content string, width int, height int, selected bool) string {
	var borderColor lipgloss.Color
	if selected {
		borderColor = styles.ColorBright
	} else {
		borderColor = styles.ColorDim
	}

	borderStyle := lipgloss.NewStyle().Foreground(borderColor)
	contentStyle := lipgloss.NewStyle().Foreground(styles.ColorNormal)

	// Rounded corners
	top := borderStyle.Render("╭" + strings.Repeat("─", width) + "╮")
	bottom := borderStyle.Render("╰" + strings.Repeat("─", width) + "╯")

	innerWidth := width - 2

	lines := strings.Split(content, "\n")
	var middle strings.Builder
	totalLines := 0

	for _, line := range lines {
		if totalLines >= height {
			break
		}
		wrappedLines := WrapText(line, innerWidth)
		for _, wl := range wrappedLines {
			if totalLines >= height {
				break
			}
			styled := contentStyle.Render(wl)
			lineWidth := lipgloss.Width(styled)
			pad := max(innerWidth-lineWidth, 0)
			middle.WriteString(borderStyle.Render("│"))
			middle.WriteString(" ")
			middle.WriteString(styled)
			middle.WriteString(strings.Repeat(" ", pad))
			middle.WriteString(" ")
			middle.WriteString(borderStyle.Render("│"))
			middle.WriteString("\n")
			totalLines++
		}
	}

	for range height - totalLines {
		middle.WriteString(borderStyle.Render("│"))
		middle.WriteString(strings.Repeat(" ", width))
		middle.WriteString(borderStyle.Render("│"))
		middle.WriteString("\n")
	}

	return top + "\n" + middle.String() + bottom
}

// BuildListItems updates lists when added
func BuildListItems(
	l *list.Model,
	isAdditional bool,
	hasMore bool,
	newItems []list.Item,
) []list.Item {

	var base []list.Item

	if isAdditional {
		existing := l.Items()

		base = make([]list.Item, 0, len(existing))
		for _, it := range existing {
			if _, isLoadMore := it.(LoadMoreItem); !isLoadMore {
				base = append(base, it)
			}
		}
	}

	items := make([]list.Item, 0, len(base)+len(newItems)+1)
	items = append(items, base...)
	items = append(items, newItems...)

	if hasMore {
		items = append(items, LoadMoreItem{})
	}

	return items
}
