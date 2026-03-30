package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/euklides/cyberspace-cli/api"
	"github.com/euklides/cyberspace-cli/models"
	"github.com/euklides/cyberspace-cli/styles"
)

// PostDetailLoadedMsg is sent when post and replies are loaded
type PostDetailLoadedMsg struct {
	Post    models.Post
	Replies []models.Reply
}

// PostDetailErrorMsg is sent when loading fails
type PostDetailErrorMsg struct {
	Err error
}

// BackToFeedMsg is sent when user wants to go back
type BackToFeedMsg struct{}

// PostDetailModel is the post detail screen
type PostDetailModel struct {
	post    models.Post
	replies []models.Reply
	scroll  int
	loading bool
	spinner spinner.Model
	err     error
	client  *api.Client
	postID  string
	width   int
	height  int
}

// NewPostDetailModel creates a new post detail screen
func NewPostDetailModel(baseURL, idToken, postID string) PostDetailModel {
	return PostDetailModel{
		client:  api.NewClient(baseURL, idToken),
		postID:  postID,
		spinner: NewSpinner(),
		loading: true,
	}
}

// NewPostDetailModelWithPost creates a detail screen with post already loaded
func NewPostDetailModelWithPost(baseURL, idToken string, post models.Post) PostDetailModel {
	return PostDetailModel{
		client:  api.NewClient(baseURL, idToken),
		postID:  post.ID,
		post:    post,
		spinner: NewSpinner(),
		loading: true,
	}
}

func (m PostDetailModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchPostAndReplies())
}

func (m PostDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc", "b", "backspace":
			return m, func() tea.Msg { return BackToFeedMsg{} }
		case "j", "down":
			m.scroll++
			m.clampScroll()
		case "k", "up":
			m.scroll--
			m.clampScroll()
		case "g":
			m.scroll = 0
		case "G":
			m.scroll = m.maxScroll()
		case "r":
			m.loading = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.fetchPostAndReplies())
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.clampScroll()

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case PostDetailLoadedMsg:
		m.loading = false
		m.post = msg.Post
		m.replies = msg.Replies

	case PostDetailErrorMsg:
		m.loading = false
		m.err = msg.Err
	}

	return m, nil
}

func (m PostDetailModel) View() string {
	w, h := SafeDimensions(m.width, m.height)

	if m.loading && m.post.ID == "" {
		return m.renderLoadingScreen(w, h)
	}

	if m.err != nil {
		return m.renderErrorScreen(w, h)
	}

	var b strings.Builder

	// Header
	b.WriteString(m.renderHeader(w))

	// Build content
	content := m.buildContent(w)
	lines := strings.Split(content, "\n")

	// Apply scroll
	visibleLines := h - 5 // header + footer
	start := m.scroll
	end := Min(start+visibleLines, len(lines))

	for i := start; i < end && i < len(lines); i++ {
		b.WriteString(lines[i])
		b.WriteString("\n")
	}

	// Calculate padding needed
	renderedLines := end - start
	paddingNeeded := visibleLines - renderedLines
	for i := 0; i < paddingNeeded; i++ {
		b.WriteString("\n")
	}

	// Footer with scroll indicator
	b.WriteString(m.renderFooter(w, start+1, end, len(lines)))

	return b.String()
}

func (m PostDetailModel) renderHeader(width int) string {
	// Title centered with yellow bar across full width
	title := "▓▒░ ENTRY VIEWER ░▒▓"
	titleRendered := styles.Title.Render(title)
	titleWidth := lipgloss.Width(titleRendered)

	// Calculate bar widths for centering
	barWidth := (width - titleWidth) / 2
	if barWidth < 0 {
		barWidth = 0
	}
	rightBarWidth := width - titleWidth - barWidth
	if rightBarWidth < 0 {
		rightBarWidth = 0
	}

	// Yellow bar on left + centered title + yellow bar on right
	barStyle := lipgloss.NewStyle().Foreground(styles.ColorBright)
	leftBar := barStyle.Render(strings.Repeat("█", barWidth))
	rightBar := barStyle.Render(strings.Repeat("█", rightBarWidth))

	return leftBar + titleRendered + rightBar + "\n\n"
}

func (m PostDetailModel) renderFooter(width, start, end, total int) string {
	var b strings.Builder

	// Simple divider
	b.WriteString(styles.Divider(width))
	b.WriteString("\n")

	// Scroll position on left
	scrollInfo := styles.StatusBarSegment("LINE", fmt.Sprintf("%d-%d/%d", start, end, total))

	// Navigation hints on right
	navHint := styles.Dim.Render("[j/k] Scroll  [g/G] Top/Bottom  [r] Refresh  [ESC] Back")

	scrollWidth := lipgloss.Width(scrollInfo)
	navWidth := lipgloss.Width(navHint)
	spacing := width - scrollWidth - navWidth
	if spacing < 1 {
		spacing = 1
	}

	b.WriteString(scrollInfo)
	b.WriteString(strings.Repeat(" ", spacing))
	b.WriteString(navHint)

	return b.String()
}

func (m PostDetailModel) renderLoadingScreen(width, height int) string {
	loadingBox := styles.DataBox("DECRYPTING TRANSMISSION",
		"\n"+
			"  "+m.spinner.View()+" Accessing secured data...\n"+
			"\n"+
			"  "+styles.ProgressBarSimple(0.5, 30)+"\n"+
			"\n"+
			"  "+styles.Dim.Render("Decoding neural patterns...")+"\n",
		50)

	return FullScreen(loadingBox, width, height, lipgloss.Center, lipgloss.Center)
}

func (m PostDetailModel) renderErrorScreen(width, height int) string {
	errorBox := styles.AlertBox(m.err.Error(), "error", 50) +
		"\n\n" +
		styles.Dim.Render("Press [ESC] to return to feed, [r] to retry")

	return FullScreen(errorBox, width, height, lipgloss.Center, lipgloss.Center)
}

func (m PostDetailModel) buildContent(width int) string {
	var b strings.Builder

	contentWidth := width - 2
	if contentWidth < 40 {
		contentWidth = 78
	}

	// Metadata box (narrower - about 50 chars)
	metaWidth := 50
	metaContent := fmt.Sprintf("@%s\n%s", m.post.AuthorUsername, TimeAgo(m.post.CreatedAt))
	b.WriteString(renderBox("POST INFO", metaContent, metaWidth))
	b.WriteString("\n\n")

	// Message content box (full width with vertical lines) - strip markdown but keep link text
	cleanContent := StripMarkdownKeepNewlines(m.post.Content)
	b.WriteString(renderBox("MESSAGE", cleanContent, contentWidth))
	b.WriteString("\n\n")

	// Stats line
	replyWord := "replies"
	if m.post.RepliesCount == 1 {
		replyWord = "reply"
	}
	saveWord := "saves"
	if m.post.BookmarksCount == 1 {
		saveWord = "save"
	}
	statsLine := fmt.Sprintf("%d %s · %d %s", m.post.RepliesCount, replyWord, m.post.BookmarksCount, saveWord)

	// Topics
	if len(m.post.Topics) > 0 {
		tags := make([]string, len(m.post.Topics))
		for i, t := range m.post.Topics {
			tags[i] = "[" + t + "]"
		}
		statsLine += "  " + strings.Join(tags, " ")
	}
	b.WriteString(styles.Dim.Render(statsLine))
	b.WriteString("\n\n")

	// Replies section
	if m.loading {
		b.WriteString(m.spinner.View() + " Loading replies...")
	} else if len(m.replies) == 0 {
		b.WriteString(renderBox("REPLIES", "No replies yet", 40))
	} else {
		// Build replies content
		var repliesContent strings.Builder
		for i, reply := range m.replies {
			repliesContent.WriteString(styles.Username.Render("@" + reply.AuthorUsername))
			repliesContent.WriteString(styles.Dim.Render(" · " + TimeAgo(reply.CreatedAt)))
			repliesContent.WriteString("\n")
			repliesContent.WriteString(StripMarkdownKeepNewlines(reply.Content))
			if i < len(m.replies)-1 {
				repliesContent.WriteString("\n\n")
			}
		}
		replyTitle := fmt.Sprintf("REPLIES [%d]", len(m.replies))
		b.WriteString(renderBox(replyTitle, repliesContent.String(), contentWidth))
	}

	return b.String()
}

// renderBox renders content in a box with title
func renderBox(title, content string, width int) string {
	borderStyle := lipgloss.NewStyle().Foreground(styles.ColorDim)
	titleStyle := lipgloss.NewStyle().Foreground(styles.ColorBright).Bold(true)

	innerWidth := width - 4 // Account for "│ " and " │"
	if innerWidth < 10 {
		innerWidth = 10
	}

	// Top border with title: ┌─ TITLE ─────────────────┐
	// Structure: ┌ (1) + ─ (1) + space (1) + title + space (1) + dashes + ┐ (1)
	titleRendered := titleStyle.Render(title)
	titleVisualLen := lipgloss.Width(title)
	remainingDashes := width - 5 - titleVisualLen // 5 = ┌ + ─ + space + space + ┐
	if remainingDashes < 1 {
		remainingDashes = 1
	}
	top := borderStyle.Render("┌─ ") + titleRendered + borderStyle.Render(" "+strings.Repeat("─", remainingDashes)+"┐")

	// Bottom border
	bottom := borderStyle.Render("└" + strings.Repeat("─", width-2) + "┘")

	// Wrap content - split by newlines and wrap each line
	var middle strings.Builder
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		// Wrap long lines
		wrappedLines := wrapText(line, innerWidth)
		for _, wl := range wrappedLines {
			lineWidth := lipgloss.Width(wl)
			padding := innerWidth - lineWidth
			if padding < 0 {
				padding = 0
			}
			middle.WriteString(borderStyle.Render("│"))
			middle.WriteString(" ")
			middle.WriteString(wl)
			middle.WriteString(strings.Repeat(" ", padding))
			middle.WriteString(" ")
			middle.WriteString(borderStyle.Render("│"))
			middle.WriteString("\n")
		}
	}

	return top + "\n" + middle.String() + bottom
}

// wrapText wraps text to fit within width
func wrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	currentLine := words[0]
	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	lines = append(lines, currentLine)

	return lines
}

func (m PostDetailModel) renderReply(r models.Reply, width int) string {
	var b strings.Builder

	// Reply header
	header := fmt.Sprintf("  %s %s",
		styles.Username.Render("@"+r.AuthorUsername),
		styles.Timestamp.Render("· "+TimeAgo(r.CreatedAt)),
	)
	b.WriteString(header)
	b.WriteString("\n")

	// Reply content
	contentStyle := lipgloss.NewStyle().
		Foreground(styles.ColorContent).
		Width(width).
		PaddingLeft(4)

	b.WriteString(contentStyle.Render(r.Content))
	b.WriteString("\n")

	return b.String()
}

func (m PostDetailModel) fetchPostAndReplies() tea.Cmd {
	return func() tea.Msg {
		post := m.post
		if post.ID == "" {
			p, err := m.client.FetchPost(m.postID)
			if err != nil {
				return PostDetailErrorMsg{Err: err}
			}
			post = *p
		}

		replies, err := m.client.FetchReplies(m.postID)
		if err != nil {
			return PostDetailErrorMsg{Err: err}
		}

		return PostDetailLoadedMsg{Post: post, Replies: replies}
	}
}

func (m *PostDetailModel) clampScroll() {
	m.scroll = Clamp(m.scroll, 0, m.maxScroll())
}

// SetSize updates the view dimensions
func (m *PostDetailModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m PostDetailModel) maxScroll() int {
	w, _ := SafeDimensions(m.width, m.height)
	content := m.buildContent(w)
	lines := strings.Split(content, "\n")
	height := m.height
	if height < 10 {
		height = 24
	}
	visibleLines := height - 5
	maxScroll := len(lines) - visibleLines
	if maxScroll < 0 {
		return 0
	}
	return maxScroll
}
