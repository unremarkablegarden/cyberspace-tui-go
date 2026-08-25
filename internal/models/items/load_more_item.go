package items

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const LoadMoreInnerHeight = 1 // Single string
const LoadMoreContentMessage = "▼ LOAD MORE POSTS ▼"

// ═══════════════════════════════════════════════════════════════════════════════
// LOAD MORE ITEM — sentinel item at the bottom of the list
// ═══════════════════════════════════════════════════════════════════════════════

// LoadMoreItem is a sentinel list item for triggering pagination.
type LoadMoreItem struct{}

func (l LoadMoreItem) FilterValue() string { return "" }
func (l LoadMoreItem) Title() string       { return "LOAD MORE" }
func (l LoadMoreItem) Description() string { return "" }

func renderLoadMoreCard(selected bool, width int) string {
	innerWidth := max(width-4, 76)

	contentWidth := lipgloss.Width(LoadMoreContentMessage)
	padding := max((innerWidth-contentWidth)/2, 0)
	centeredContent := strings.Repeat(" ", padding) + LoadMoreContentMessage

	return BuildCardBox(centeredContent, innerWidth, LoadMoreInnerHeight, selected)
}
