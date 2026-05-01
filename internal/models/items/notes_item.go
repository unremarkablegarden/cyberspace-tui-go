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

// NoteItem implements list.Item for notes
type NoteItem struct{ Note entities.Note }

func (n NoteItem) FilterValue() string { return n.Note.Content }
func (n NoteItem) Title() string       { return n.Note.Content }
func (n NoteItem) Description() string { return TimeAgo(n.Note.CreatedAt) }

// NoteDelegate renders note items in the list
type NoteDelegate struct{}

func (d NoteDelegate) Height() int                               { return 3 }
func (d NoteDelegate) Spacing() int                              { return 0 }
func (d NoteDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d NoteDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	switch it := item.(type) {
	case NoteItem:
		isSelected := index == m.Index()

		content := strings.TrimSpace(it.Note.Content)
		content = strings.ReplaceAll(content, "\n", " ")
		if len(content) > 72 {
			content = content[:72] + "…"
		}

		date := TimeAgo(it.Note.CreatedAt)
		var meta string
		if len(it.Note.Topics) > 0 {
			meta = date + "  [" + strings.Join(it.Note.Topics, "] [") + "]"
		} else {
			meta = date
		}

		var contentLine, metaLine string
		if isSelected {
			contentLine = styles.Bright.Render("▸ " + content)
			metaLine = styles.Dim.Render("  " + meta)
		} else {
			contentLine = styles.Normal.Render("  " + content)
			metaLine = styles.Dim.Render("  " + meta)
		}
		fmt.Fprintf(w, "%s\n%s\n", contentLine, metaLine)

	case LoadMoreItem:
		fmt.Fprint(w, zone.Mark("load-more-notes", styles.Dim.Render("  ▼ load more")))
	}
}
