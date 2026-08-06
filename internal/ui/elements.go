package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"

	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// NewSpinner creates a sci-fi styled spinner
func NewSpinner() spinner.Model {
	s := spinner.New()
	// Use a sci-fi looking spinner
	s.Spinner = spinner.Spinner{
		Frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		FPS:    time.Millisecond * 80,
	}
	s.Style = styles.Spinner
	return s
}
