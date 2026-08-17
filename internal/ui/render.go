package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"

	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

var SpView string

type LoadingScreenTexts struct {
	TitleText    string
	SubtitleText string
	BottomText   string
}

func RenderFooterWithList(helpView string, paginatorView string, width int) string {
	helpWidth := lipgloss.Width(helpView)
	paginatorWidth := lipgloss.Width(paginatorView)

	dividerWidth := max(width-helpWidth-paginatorWidth-2, 1)

	return helpView + " " + styles.Divider(dividerWidth) + " " + paginatorView
}

// RenderHeader renders a centered title bar with block-fill sides.
func RenderHeader(title string, width int) string {
	titleRendered := styles.Title.Render(title)
	titleWidth := lipgloss.Width(titleRendered)

	barStyle := lipgloss.NewStyle().Foreground(styles.ColorBright)

	barWidth := max((width-titleWidth)/2, 0)
	leftBar := barStyle.Render(strings.Repeat("█", barWidth))

	rightBarWidth := max(width-titleWidth-barWidth, 0)
	rightBar := barStyle.Render(strings.Repeat("█", rightBarWidth))

	return leftBar + titleRendered + rightBar + "\n"
}

func RenderLoadingScreen(
	lst LoadingScreenTexts,
	spinnerModel spinner.Model,
	width, height int,
) string {
	loadingBox := styles.DataBox(lst.TitleText,
		"\n"+
			"  "+spinnerModel.View()+styles.Normal.Render(lst.SubtitleText)+"\n"+
			"\n"+
			"  "+styles.Dim.Render(lst.BottomText)+"\n",
		50)
	return FullScreen(loadingBox, width, height, lipgloss.Center, lipgloss.Center)
}

func RenderErrorScreen(err error, width, height int) string {
	errorBox := styles.AlertBox(err.Error(), "error", 50) +
		"\n\n" +
		styles.Dim.Render("Press [r] to retry, [esc] to go back")
	return FullScreen(errorBox, width, height, lipgloss.Center, lipgloss.Center)
}
