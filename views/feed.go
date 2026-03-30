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

// PostsLoadedMsg is sent when posts are fetched
type PostsLoadedMsg struct {
	Posts  []models.Post
	Cursor string
}

// MorePostsLoadedMsg is sent when more posts are loaded
type MorePostsLoadedMsg struct {
	Posts  []models.Post
	Cursor string
}

// PostsErrorMsg is sent when fetching posts fails
type PostsErrorMsg struct {
	Err error
}

// OpenPostMsg is sent when user wants to view a post
type OpenPostMsg struct {
	Post models.Post
}

// FeedModel is the post feed screen
type FeedModel struct {
	posts       []models.Post
	cursor      int
	scrollOff   int
	loading     bool
	loadingMore bool
	spinner     spinner.Model
	err         error
	client      *api.Client
	nextCursor  string
	hasMore     bool
	width       int
	height      int
}

// NewFeedModel creates a new feed screen
func NewFeedModel(baseURL, idToken string) FeedModel {
	return FeedModel{
		client:  api.NewClient(baseURL, idToken),
		spinner: NewSpinner(),
		loading: true,
		hasMore: true,
	}
}

func (m FeedModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchPosts())
}

func (m FeedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "j", "down":
			maxCursor := len(m.posts) - 1
			if m.hasMore {
				maxCursor = len(m.posts) // load more button
			}
			if m.cursor < maxCursor {
				m.cursor++
				m.adjustScroll()
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
				m.adjustScroll()
			}
		case "g":
			m.cursor = 0
			m.scrollOff = 0
		case "G":
			if m.hasMore {
				m.cursor = len(m.posts) // Go to load more button
			} else {
				m.cursor = len(m.posts) - 1
			}
			m.adjustScroll()
		case "r":
			m.loading = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.fetchPosts())
		case "enter":
			if len(m.posts) > 0 {
				if m.cursor == len(m.posts) && m.hasMore {
					// Load more button
					if !m.loadingMore {
						m.loadingMore = true
						return m, tea.Batch(m.spinner.Tick, m.fetchMorePosts())
					}
				} else if m.cursor < len(m.posts) {
					return m, func() tea.Msg {
						return OpenPostMsg{Post: m.posts[m.cursor]}
					}
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.adjustScroll()

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case PostsLoadedMsg:
		m.loading = false
		m.posts = msg.Posts
		m.cursor = 0
		m.scrollOff = 0
		m.nextCursor = msg.Cursor
		m.hasMore = msg.Cursor != ""

	case MorePostsLoadedMsg:
		m.loadingMore = false
		m.posts = append(m.posts, msg.Posts...)
		m.nextCursor = msg.Cursor
		m.hasMore = msg.Cursor != ""

	case PostsErrorMsg:
		m.loading = false
		m.loadingMore = false
		m.err = msg.Err
	}

	return m, nil
}

func (m FeedModel) View() string {
	w, h := SafeDimensions(m.width, m.height)

	if m.loading {
		return m.renderLoadingScreen(w, h)
	}

	if m.err != nil {
		return m.renderErrorScreen(w, h)
	}

	if len(m.posts) == 0 {
		return m.renderEmptyScreen(w, h)
	}

	var b strings.Builder

	// Header with sci-fi styling
	header := m.renderHeader(w)
	b.WriteString(header)

	// Calculate visible items (accounting for header and footer)
	totalItems := len(m.posts)
	if m.hasMore {
		totalItems++ // +1 for load more button
	}
	visibleItems := m.visiblePostCount()

	// Render posts and load more button
	var postLines int
	for i := m.scrollOff; i < totalItems && i < m.scrollOff+visibleItems; i++ {
		var item string
		if i < len(m.posts) {
			item = m.renderPost(m.posts[i], i == m.cursor, w)
		} else {
			// Load more button
			item = m.renderLoadMoreButton(i == m.cursor, w)
		}
		b.WriteString(item)
		b.WriteString("\n")
		postLines += strings.Count(item, "\n") + 1
	}

	// Footer area
	headerLines := 3
	footerLines := 2 // divider + status/help line
	usedLines := headerLines + postLines
	paddingNeeded := h - usedLines - footerLines

	// Add blank lines to push footer to bottom
	for i := 0; i < paddingNeeded; i++ {
		b.WriteString("\n")
	}

	// Status bar (includes help text on same line)
	b.WriteString(m.renderStatusBar(w))

	return b.String()
}

func (m FeedModel) renderHeader(width int) string {
	// Title centered with yellow bar across full width
	title := "▓▒░ ᑕ¥βєяรקค¢є ░▒▓"
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

func (m FeedModel) renderStatusBar(width int) string {
	// Left side: position (show "LOAD MORE" if on that button)
	var leftSide string
	if m.cursor == len(m.posts) {
		leftSide = styles.StatusBarSegment("ACTION", "LOAD MORE")
	} else {
		leftSide = styles.StatusBarSegment("POST", fmt.Sprintf("%d/%d", m.cursor+1, len(m.posts)))
	}

	// Right side: help text
	rightSide := styles.Dim.Render("[j/k] Navigate  [ENTER] Open  [r] Refresh  [q] Quit")

	// Calculate spacing
	leftWidth := lipgloss.Width(leftSide)
	rightWidth := lipgloss.Width(rightSide)
	spacing := width - leftWidth - rightWidth
	if spacing < 1 {
		spacing = 1
	}

	return styles.Divider(width) + "\n" +
		leftSide + strings.Repeat(" ", spacing) + rightSide
}

func (m FeedModel) renderLoadingScreen(width, height int) string {
	var b strings.Builder

	// Centered loading animation
	loadingBox := styles.DataBox("ESTABLISHING CONNECTION",
		"\n"+
		"  "+m.spinner.View()+" Synchronizing with network...\n"+
		"\n"+
		"  "+styles.ProgressBarSimple(0.3, 30)+"\n"+
		"\n"+
		"  "+styles.Dim.Render("Please wait while we access the datastream")+"\n",
		50)

	return FullScreen(b.String()+loadingBox, width, height, lipgloss.Center, lipgloss.Center)
}

func (m FeedModel) renderErrorScreen(width, height int) string {
	errorBox := styles.AlertBox(m.err.Error(), "error", 50) +
		"\n\n" +
		styles.Dim.Render("Press [r] to retry connection, [q] to disconnect")

	return FullScreen(errorBox, width, height, lipgloss.Center, lipgloss.Center)
}

func (m FeedModel) renderEmptyScreen(width, height int) string {
	emptyBox := styles.DataBox("NO DATA",
		"\n"+
		"  "+styles.Dim.Render("No transmissions detected in this sector")+"\n"+
		"\n"+
		"  "+styles.Help.Render("Press [r] to rescan frequency")+"\n",
		45)

	return FullScreen(emptyBox, width, height, lipgloss.Center, lipgloss.Center)
}

func (m FeedModel) visiblePostCount() int {
	visibleHeight := m.height - 5 // header + footer
	postHeight := 6               // average post height (borders + content + gap)
	count := visibleHeight / postHeight
	if count < 1 {
		return 1
	}
	return count
}

func (m *FeedModel) adjustScroll() {
	visiblePosts := m.visiblePostCount()

	if m.cursor >= m.scrollOff+visiblePosts {
		m.scrollOff = m.cursor - visiblePosts + 1
	}
	if m.cursor < m.scrollOff {
		m.scrollOff = m.cursor
	}
}

func (m FeedModel) renderPost(p models.Post, selected bool, width int) string {
	innerWidth := width - 4 // Account for box borders and padding
	if innerWidth < 20 {
		innerWidth = 76
	}

	// Username on left, time + stats on right
	username := "@" + p.AuthorUsername

	// Format stats: "2 replies · 1 save" or "0 replies · 0 saves"
	replyWord := "replies"
	if p.RepliesCount == 1 {
		replyWord = "reply"
	}
	saveWord := "saves"
	if p.BookmarksCount == 1 {
		saveWord = "save"
	}
	rightStats := fmt.Sprintf("%d %s · %d %s · %s",
		p.RepliesCount, replyWord,
		p.BookmarksCount, saveWord,
		TimeAgo(p.CreatedAt))

	// Calculate spacing for header line
	usernameWidth := len(username)
	statsWidth := len(rightStats)
	headerSpacing := innerWidth - usernameWidth - statsWidth
	if headerSpacing < 1 {
		headerSpacing = 1
	}

	// Build header line
	headerLine := styles.Username.Render(username) +
		strings.Repeat(" ", headerSpacing) +
		styles.Dim.Render(rightStats)

	// Content (limited to ~2 lines worth for preview) - strip markdown for performance
	content := Truncate(StripMarkdown(p.Content), innerWidth*2-3)

	// Topics with [brackets]
	var tagsLine string
	if len(p.Topics) > 0 {
		tags := make([]string, len(p.Topics))
		for i, t := range p.Topics {
			tags[i] = "[" + t + "]"
		}
		tagsLine = styles.Dim.Render(strings.Join(tags, " "))
	}

	// Build box content
	var boxContent strings.Builder
	boxContent.WriteString(headerLine)
	boxContent.WriteString("\n")
	boxContent.WriteString(content)
	if tagsLine != "" {
		boxContent.WriteString("\n")
		boxContent.WriteString(tagsLine)
	}

	// Choose border style based on selection
	if selected {
		// Bright border for selected
		return renderPostBox(boxContent.String(), innerWidth, true)
	}

	// Dim border for unselected
	return renderPostBox(boxContent.String(), innerWidth, false)
}

// renderPostBox renders a post inside a box with max lines limit
func renderPostBox(content string, width int, selected bool) string {
	var borderColor lipgloss.Color
	if selected {
		borderColor = styles.ColorBright
	} else {
		borderColor = styles.ColorDim
	}

	borderStyle := lipgloss.NewStyle().Foreground(borderColor)

	// Top border
	top := borderStyle.Render("┌" + strings.Repeat("─", width) + "┐")

	// Bottom border
	bottom := borderStyle.Render("└" + strings.Repeat("─", width) + "┘")

	innerWidth := width - 2 // Account for padding spaces

	// Wrap content lines, but limit total lines to keep consistent height
	lines := strings.Split(content, "\n")
	var middle strings.Builder
	totalLines := 0
	maxLines := 4 // header + 2 content lines + tags

	for _, line := range lines {
		if totalLines >= maxLines {
			break
		}
		// Wrap long lines to fit within the box
		wrappedLines := wrapText(line, innerWidth)
		for _, wl := range wrappedLines {
			if totalLines >= maxLines {
				break
			}
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
			totalLines++
		}
	}

	return top + "\n" + middle.String() + bottom
}

// SetSize updates the view dimensions
func (m *FeedModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m FeedModel) fetchPosts() tea.Cmd {
	return func() tea.Msg {
		posts, cursor, err := m.client.FetchPosts(30)
		if err != nil {
			return PostsErrorMsg{Err: err}
		}
		return PostsLoadedMsg{Posts: posts, Cursor: cursor}
	}
}

func (m FeedModel) fetchMorePosts() tea.Cmd {
	return func() tea.Msg {
		posts, cursor, err := m.client.FetchMorePosts(30, m.nextCursor)
		if err != nil {
			return PostsErrorMsg{Err: err}
		}
		return MorePostsLoadedMsg{Posts: posts, Cursor: cursor}
	}
}

func (m FeedModel) renderLoadMoreButton(selected bool, width int) string {
	innerWidth := width - 4
	if innerWidth < 20 {
		innerWidth = 76
	}

	var content string
	if m.loadingMore {
		content = m.spinner.View() + " Loading more posts..."
	} else {
		content = "▼ LOAD MORE POSTS ▼"
	}

	// Center the content
	contentWidth := lipgloss.Width(content)
	padding := (innerWidth - contentWidth) / 2
	if padding < 0 {
		padding = 0
	}
	centeredContent := strings.Repeat(" ", padding) + content

	return renderPostBox(centeredContent, innerWidth, selected)
}
