package views

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"

	"github.com/unremarkablegarden/cyberspace-tui-go/models"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// ═══════════════════════════════════════════════════════════════════════════════
// POST ITEM — wraps models.Post for the list bubble
// ═══════════════════════════════════════════════════════════════════════════════

// PostItem wraps a Post so it satisfies list.Item.
type PostItem struct {
	Post models.Post
}

func (p PostItem) FilterValue() string {
	return p.Post.AuthorUsername + " " + p.Post.Content + " " + strings.Join(p.Post.Topics, " ")
}

func (p PostItem) Title() string       { return "@" + p.Post.AuthorUsername }
func (p PostItem) Description() string { return Truncate(StripMarkdown(p.Post.Content), 80) }

// ═══════════════════════════════════════════════════════════════════════════════
// LOAD MORE ITEM — sentinel item at the bottom of the list
// ═══════════════════════════════════════════════════════════════════════════════

// LoadMoreItem is a sentinel list item for triggering pagination.
type LoadMoreItem struct{}

func (l LoadMoreItem) FilterValue() string { return "" }
func (l LoadMoreItem) Title() string       { return "LOAD MORE" }
func (l LoadMoreItem) Description() string { return "" }

// ═══════════════════════════════════════════════════════════════════════════════
// POST DELEGATE — custom rendering for list items
// ═══════════════════════════════════════════════════════════════════════════════

// PostDelegate renders post items as styled cards.
type PostDelegate struct{}

func (d PostDelegate) Height() int  { return 6 }
func (d PostDelegate) Spacing() int { return 0 }

func (d PostDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d PostDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	selected := index == m.Index()
	width := m.Width()

	switch it := item.(type) {
	case PostItem:
		card := renderPostCard(it.Post, selected, width)
		fmt.Fprint(w, zone.Mark(it.Post.ID, card))
	case LoadMoreItem:
		card := renderLoadMoreCard(selected, width)
		fmt.Fprint(w, zone.Mark("load-more", card))
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// CARD RENDERING
// ═══════════════════════════════════════════════════════════════════════════════

func renderPostCard(p models.Post, selected bool, width int) string {
	innerWidth := width - 4
	if innerWidth < 20 {
		innerWidth = 76
	}

	// Username on left, time + stats on right
	username := "@" + p.AuthorUsername

	replyWord := "replies"
	if p.RepliesCount == 1 {
		replyWord = "reply"
	}
	saveWord := "saves"
	if p.BookmarksCount == 1 {
		saveWord = "save"
	}
	rightStats := fmt.Sprintf("%d %s · %d %s · %s",
		p.RepliesCount, replyWord,
		p.BookmarksCount, saveWord,
		TimeAgo(p.CreatedAt))

	usernameWidth := len(username)
	statsWidth := len(rightStats)
	headerSpacing := innerWidth - usernameWidth - statsWidth
	if headerSpacing < 1 {
		headerSpacing = 1
	}

	headerLine := styles.Username.Render(username) +
		strings.Repeat(" ", headerSpacing) +
		styles.Dim.Render(rightStats)

	content := Truncate(StripMarkdown(p.Content), innerWidth*2-3)

	var tagsLine string
	if len(p.Topics) > 0 {
		tags := make([]string, len(p.Topics))
		for i, t := range p.Topics {
			tags[i] = "[" + t + "]"
		}
		tagsLine = styles.Dim.Render(strings.Join(tags, " "))
	}

	var boxContent strings.Builder
	boxContent.WriteString(headerLine)
	boxContent.WriteString("\n")
	boxContent.WriteString(content)
	if tagsLine != "" {
		boxContent.WriteString("\n")
		boxContent.WriteString(tagsLine)
	}

	return buildCardBox(boxContent.String(), innerWidth, selected)
}

func renderLoadMoreCard(selected bool, width int) string {
	innerWidth := width - 4
	if innerWidth < 20 {
		innerWidth = 76
	}

	content := "▼ LOAD MORE POSTS ▼"
	contentWidth := lipgloss.Width(content)
	padding := (innerWidth - contentWidth) / 2
	if padding < 0 {
		padding = 0
	}
	centeredContent := strings.Repeat(" ", padding) + content

	return buildCardBox(centeredContent, innerWidth, selected)
}

// buildCardBox renders content in a bordered card with rounded corners.
func buildCardBox(content string, width int, selected bool) string {
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
	maxLines := 4

	for _, line := range lines {
		if totalLines >= maxLines {
			break
		}
		wrappedLines := wrapText(line, innerWidth)
		for _, wl := range wrappedLines {
			if totalLines >= maxLines {
				break
			}
			styled := contentStyle.Render(wl)
			lineWidth := lipgloss.Width(styled)
			pad := innerWidth - lineWidth
			if pad < 0 {
				pad = 0
			}
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

	return top + "\n" + middle.String() + bottom
}
