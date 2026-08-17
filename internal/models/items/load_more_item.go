package items

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ═══════════════════════════════════════════════════════════════════════════════
// LOAD MORE ITEM — sentinel item at the bottom of the list
// ═══════════════════════════════════════════════════════════════════════════════

// LoadMoreItem is a sentinel list item for triggering pagination.
type LoadMoreItem struct{}

func (l LoadMoreItem) FilterValue() string { return "" }
func (l LoadMoreItem) Title() string       { return "LOAD MORE" }
func (l LoadMoreItem) Description() string { return "" }

func renderLoadMoreCard(selected bool, width int) string {
	innerWidth := width - 4
	if innerWidth < 20 {
		innerWidth = 76
	}

	content := "▼ LOAD MORE POSTS ▼"
	contentWidth := lipgloss.Width(content)
	padding := max((innerWidth-contentWidth)/2, 0)
	centeredContent := strings.Repeat(" ", padding) + content

	return BuildCardBox(centeredContent, innerWidth, selected)
}
