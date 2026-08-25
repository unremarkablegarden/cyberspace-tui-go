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

func (d NoteDelegate) Height() int                               { return 6 }
func (d NoteDelegate) Spacing() int                              { return 0 }
func (d NoteDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d NoteDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	switch it := item.(type) {
	case NoteItem:
		isSelected := index == m.Index()
		fmt.Fprint(w, renderNoteCard(it.Note, m.Width(), d.Height(), isSelected))

	case LoadMoreItem:
		fmt.Fprint(w, zone.Mark("load-more-notes", styles.Dim.Render("  ▼ load more")))
	}
}
func renderNoteCard(n entities.Note, width int, height int, selected bool) string {
	innerWidth := max(width-4, 76)
	innerHeight := max(height-2, 1)

	content := strings.TrimSpace(n.Content)
	content = strings.ReplaceAll(content, "\n", " ")
	if len(content) > 72 {
		content = content[:72] + "…"
	}

	meta := TimeAgo(n.CreatedAt) +
		"·" +
		getTagLineString(n.Topics)

	var contentLine, metaLine string
	if selected {
		contentLine = styles.Bright.Render("▸ " + content)
	} else {
		contentLine = styles.Normal.Render("  " + content)
	}

	metaLine = styles.Dim.Render("  " + meta)

	var boxContent strings.Builder
	boxContent.WriteString(contentLine)
	boxContent.WriteString("\n")
	boxContent.WriteString(metaLine)

	return BuildCardBox(boxContent.String(), innerWidth, innerHeight, selected)
}

func NotesToItems(notes []entities.Note) []list.Item {
	items := make([]list.Item, len(notes))
	for i, n := range notes {
		items[i] = NoteItem{Note: n}
	}
	return items
}
