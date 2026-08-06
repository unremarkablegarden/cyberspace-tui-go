package ui

import "github.com/charmbracelet/lipgloss"

const (
	DefaultMinWidth  = 10
	DefaultSafeWidth = 80

	DefaultMinHeight  = 10
	DefaultSafeHeight = 24
)

// SafeDimensions returns width and height with sensible defaults
// Use this before WindowSizeMsg has been received
func SafeDimensions(width, height int) (int, int) {
	if width < DefaultMinWidth {
		width = DefaultSafeWidth
	}
	if height < DefaultMinHeight {
		height = DefaultSafeHeight
	}

	return width, height
}

// FullScreen renders content centered in a full-screen container.
func FullScreen(content string, width, height int, hAlign, vAlign lipgloss.Position) string {
	w, h := SafeDimensions(width, height)
	return lipgloss.Place(w, h, hAlign, vAlign, content)
}
