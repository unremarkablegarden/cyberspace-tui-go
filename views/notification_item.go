package views

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/unremarkablegarden/cyberspace-tui-go/models"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// NotificationItem wraps a Notification for the list bubble
type NotificationItem struct {
	Notification models.Notification
}

func (n NotificationItem) FilterValue() string {
	return n.Notification.ActorUsername + " " + n.Notification.Type
}

func (n NotificationItem) Title() string {
	return notificationSummary(n.Notification)
}

func (n NotificationItem) Description() string {
	return TimeAgo(n.Notification.CreatedAt)
}

// NotificationDelegate renders notification items as styled cards
type NotificationDelegate struct{}

func (d NotificationDelegate) Height() int  { return 4 }
func (d NotificationDelegate) Spacing() int { return 0 }

func (d NotificationDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d NotificationDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	selected := index == m.Index()
	width := m.Width()

	switch it := item.(type) {
	case NotificationItem:
		fmt.Fprint(w, renderNotificationCard(it.Notification, selected, width))
	case LoadMoreItem:
		fmt.Fprint(w, renderLoadMoreCard(selected, width))
	}
}

func renderNotificationCard(n models.Notification, selected bool, width int) string {
	innerWidth := width - 4
	if innerWidth < 20 {
		innerWidth = 76
	}

	unreadMark := "  "
	if !n.Read {
		unreadMark = styles.Bright.Render("● ")
	}

	summary := notificationSummary(n)
	timeStr := TimeAgo(n.CreatedAt)

	summaryWidth := lipgloss.Width(unreadMark) + lipgloss.Width(summary)
	timeWidth := len(timeStr)
	spacing := innerWidth - summaryWidth - timeWidth
	if spacing < 1 {
		spacing = 1
	}

	line1 := unreadMark + styles.Normal.Render(summary) + strings.Repeat(" ", spacing) + styles.Dim.Render(timeStr)

	var line2 string
	if n.PostID != "" {
		line2 = styles.Dim.Render("  → open post  [enter]")
	}

	var boxContent strings.Builder
	boxContent.WriteString(line1)
	if line2 != "" {
		boxContent.WriteString("\n")
		boxContent.WriteString(line2)
	}

	return buildCardBox(boxContent.String(), innerWidth, selected)
}

func notificationSummary(n models.Notification) string {
	actor := "@" + n.ActorUsername
	switch n.Type {
	case "reply":
		return actor + " replied to your post"
	case "bookmark":
		return actor + " saved your post"
	case "poke":
		return actor + " poked you"
	case "follow":
		return actor + " followed you"
	default:
		if n.ActorUsername != "" {
			return actor + " — " + n.Type
		}
		return n.Type
	}
}
