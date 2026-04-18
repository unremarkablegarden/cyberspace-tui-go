package views

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"

	"github.com/unremarkablegarden/cyberspace-tui-go/models"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// BookmarkItem wraps a Bookmark for the list bubble
type BookmarkItem struct {
	Bookmark models.Bookmark
}

func (b BookmarkItem) FilterValue() string {
	p := b.Bookmark.Post
	return p.AuthorUsername + " " + p.Content + " " + strings.Join(p.Topics, " ")
}

func (b BookmarkItem) Title() string       { return "@" + b.Bookmark.Post.AuthorUsername }
func (b BookmarkItem) Description() string { return Truncate(StripMarkdown(b.Bookmark.Post.Content), 80) }

// BookmarkDelegate renders bookmark items as styled cards
type BookmarkDelegate struct{}

func (d BookmarkDelegate) Height() int  { return 6 }
func (d BookmarkDelegate) Spacing() int { return 0 }

func (d BookmarkDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d BookmarkDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	selected := index == m.Index()
	width := m.Width()

	switch it := item.(type) {
	case BookmarkItem:
		card := renderBookmarkCard(it.Bookmark, selected, width)
		fmt.Fprint(w, zone.Mark(it.Bookmark.Post.ID, card))
	case LoadMoreItem:
		card := renderLoadMoreCard(selected, width)
		fmt.Fprint(w, zone.Mark("load-more-bookmarks", card))
	}
}

func renderBookmarkCard(b models.Bookmark, selected bool, width int) string {
	p := b.Post
	innerWidth := width - 4
	if innerWidth < 20 {
		innerWidth = 76
	}

	username := "@" + p.AuthorUsername

	replyWord := "replies"
	if p.RepliesCount == 1 {
		replyWord = "reply"
	}
	rightStats := fmt.Sprintf("%d %s · %s · saved %s",
		p.RepliesCount, replyWord,
		TimeAgo(p.CreatedAt),
		TimeAgo(b.CreatedAt))

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
