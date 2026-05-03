package items

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"
)

// BookmarkItem wraps a Bookmark for the list bubble
type BookmarkItem struct {
	Bookmark entities.Bookmark
}

func (b BookmarkItem) FilterValue() string {
	p := b.Bookmark.Post
	return p.AuthorUsername + " " + p.Content + " " + strings.Join(p.Topics, " ")
}

func (b BookmarkItem) Title() string { return "@" + b.Bookmark.Post.AuthorUsername }
func (b BookmarkItem) Description() string {
	desc, _ := Truncate(StripMarkdown(b.Bookmark.Post.Content), 80)
	return desc
}

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
		card := renderPostCard(&it.Bookmark.Post, &it.Bookmark, selected, width)
		fmt.Fprint(w, zone.Mark(it.Bookmark.Post.ID, card))
	case LoadMoreItem:
		card := renderLoadMoreCard(selected, width)
		fmt.Fprint(w, zone.Mark("load-more-bookmarks", card))
	}
}

func BookmarksToItems(bookmarks []entities.Bookmark) []list.Item {
	items := make([]list.Item, 0, len(bookmarks))
	for _, b := range bookmarks {
		if !b.Post.Deleted {
			items = append(items, BookmarkItem{Bookmark: b})
		}
	}
	return items
}
