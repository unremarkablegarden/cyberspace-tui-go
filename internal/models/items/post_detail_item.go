package items

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// ═══════════════════════════════════════════════════════════════════════════════
// POST ITEM — wraps entities.Post for the list bubble
// ═══════════════════════════════════════════════════════════════════════════════

// PostItem wraps a Post so it satisfies list.Item.
type PostItem struct {
	Post entities.Post
}

func (p PostItem) FilterValue() string {
	return p.Post.AuthorUsername + " " + p.Post.Content + " " + strings.Join(p.Post.Topics, " ")
}

func (p PostItem) Title() string { return "@" + p.Post.AuthorUsername }
func (p PostItem) Description() string {
	desc, _ := Truncate(StripMarkdown(p.Post.Content), 80)
	return desc
}

// ═══════════════════════════════════════════════════════════════════════════════
// POST DELEGATE — custom rendering for list items
// ═══════════════════════════════════════════════════════════════════════════════

// PostDelegate renders post items as styled cards.
type PostDelegate struct{}

func (d PostDelegate) Height() int  { return 7 }
func (d PostDelegate) Spacing() int { return 0 }

func (d PostDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d PostDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	selected := index == m.Index()
	width := m.Width()

	switch it := item.(type) {
	case PostItem:
		card := renderPostCard(&it.Post, nil, selected, width)
		fmt.Fprint(w, zone.Mark(it.Post.ID, card))
	case LoadMoreItem:
		card := renderLoadMoreCard(selected, width)
		fmt.Fprint(w, zone.Mark("load-more", card))
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// CARD RENDERING
// ═══════════════════════════════════════════════════════════════════════════════

func renderPostCard(p *entities.Post, b *entities.Bookmark, selected bool, width int) string {
	innerWidth := width - 4
	if innerWidth < 20 {
		innerWidth = 76
	}

	// Username on left, time + stats on right
	username := "@" + p.AuthorUsername
	rightStats := getPostStatsString(p, b)
	usernameWidth, statsWidth := len(username), len(rightStats)
	headerSpacing := max(innerWidth-usernameWidth-statsWidth, 1)
	headerLine := styles.Username.Render(username) +
		strings.Repeat(" ", headerSpacing) +
		styles.Dim.Render(rightStats)

	content, isTruncated := Truncate(StripMarkdown(p.Content), innerWidth*2-3)
	var tagsLine string
	if isTruncated {
		tagsLine += styles.Bright.Render("[load more] · ")
	}

	tagsLine += getTagLineString(p.Topics)

	var boxContent strings.Builder
	boxContent.WriteString(headerLine)
	boxContent.WriteString("\n")
	boxContent.WriteString(content)
	if tagsLine != "" {
		boxContent.WriteString("\n")
		boxContent.WriteString(tagsLine)
	}

	return BuildCardBox(boxContent.String(), innerWidth, selected)
}

func getPostStatsString(post *entities.Post, bookmark *entities.Bookmark) string {
	replyWord := "replies"
	if post.RepliesCount == 1 {
		replyWord = "reply"
	}

	if bookmark != nil {
		return fmt.Sprintf("%d %s · %s · saved %s",
			post.RepliesCount, replyWord,
			TimeAgo(post.CreatedAt),
			TimeAgo(bookmark.CreatedAt))
	}

	saveWord := "saves"
	if post.BookmarksCount == 1 {
		saveWord = "save"
	}
	return fmt.Sprintf("%d %s · %d %s · %s",
		post.RepliesCount, replyWord,
		post.BookmarksCount, saveWord,
		TimeAgo(post.CreatedAt))

}

func getTagLineString(topics []string) string {
	if len(topics) > 0 {
		tags := make([]string, len(topics))
		for i, t := range topics {
			tags[i] = "[" + t + "]"
		}
		return styles.Dim.Render(strings.Join(tags, " "))
	}

	return ""
}

func PostsToItems(posts []entities.Post) []list.Item {
	items := make([]list.Item, len(posts))
	for i, p := range posts {
		items[i] = PostItem{Post: p}
	}
	return items
}
