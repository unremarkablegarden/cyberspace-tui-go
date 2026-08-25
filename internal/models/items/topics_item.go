package items

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

type TopicItem struct{ Topic entities.Topic }

func (t TopicItem) FilterValue() string { return t.Topic.Name }
func (t TopicItem) Title() string       { return "[" + t.Topic.Name + "]" }
func (t TopicItem) Description() string {
	return fmt.Sprintf("%d posts", t.Topic.PostCount)
}

type TopicDelegate struct{}

func (d TopicDelegate) Height() int                             { return 3 }
func (d TopicDelegate) Spacing() int                            { return 0 }
func (d TopicDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d TopicDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	it, ok := item.(TopicItem)
	if !ok {
		return
	}
	selected := index == m.Index()
	width := m.Width()
	height := d.Height()
	fmt.Fprint(w, renderTopicCard(it.Topic, width, height, selected))
}

func renderTopicCard(t entities.Topic, width int, height int, selected bool) string {
	innerWidth := max(width-4, 76)
	innerHeight := max(height-2, 1)

	tag := "[" + t.Name + "]"
	count := fmt.Sprintf("%d posts", t.PostCount)

	spacing := max(innerWidth-len(tag)-len(count), 1)

	line := styles.Bright.Render(tag) + strings.Repeat(" ", spacing) + styles.Dim.Render(count)
	return BuildCardBox(line, innerWidth, innerHeight, selected)
}
