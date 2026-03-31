package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/unremarkablegarden/cyberspace-tui-go/api"
	"github.com/unremarkablegarden/cyberspace-tui-go/models"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

const headerHeight = 2 // title bar + blank line
const footerHeight = 2 // divider + status line

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

// ReplyCreatedMsg is sent when a reply is successfully created
type ReplyCreatedMsg struct{ ReplyID string }

// ReplyErrorMsg is sent when creating a reply fails
type ReplyErrorMsg struct{ Err error }

// BookmarkAddedMsg is sent when a post is successfully bookmarked
type BookmarkAddedMsg struct{ BookmarkID string }

// BookmarkAddErrorMsg is sent when bookmarking a post fails
type BookmarkAddErrorMsg struct{ Err error }

const composeHeight = 6 // textarea height including border

// PostDetailModel is the post detail screen
type PostDetailModel struct {
	post         models.Post
	replies      []models.Reply
	loading      bool
	spinner      spinner.Model
	err          error
	client       *api.Client
	postID       string
	width        int
	height       int
	keys     PostDetailKeyMap
	help     help.Model
	viewport viewport.Model
	ready        bool // true once we've received a WindowSizeMsg
	replyInput   textarea.Model
	composing    bool
	replySending bool
	replyErr     error
	bookmarking  bool
	bookmarked   bool
	bookmarkErr  error
}

func newReplyTextarea() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "Type your reply..."
	ta.SetHeight(3)
	ta.SetWidth(60)
	ta.CharLimit = 32768
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle().Background(styles.ColorBgSelect)
	ta.FocusedStyle.Base = lipgloss.NewStyle().Foreground(styles.ColorNormal)
	ta.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(styles.ColorMuted)
	ta.FocusedStyle.EndOfBuffer = lipgloss.NewStyle().Foreground(styles.ColorDark)
	ta.BlurredStyle = ta.FocusedStyle
	ta.Blur()
	return ta
}

func newDetailViewport() viewport.Model {
	vp := viewport.New(0, 0)
	vp.MouseWheelEnabled = true
	// Override default keymap: remove 'b' from PageUp (we use it for "back")
	km := viewport.DefaultKeyMap()
	km.PageUp = key.NewBinding(key.WithKeys("pgup"))
	km.PageDown = key.NewBinding(key.WithKeys("f", "pgdown", " "))
	km.HalfPageUp = key.NewBinding(key.WithKeys("u", "ctrl+u"))
	km.HalfPageDown = key.NewBinding(key.WithKeys("d", "ctrl+d"))
	vp.KeyMap = km
	return vp
}

// NewPostDetailModel creates a new post detail screen
func NewPostDetailModel(baseURL, idToken, postID string) PostDetailModel {
	h := help.New()
	h.Styles = styles.HelpStyles()
	return PostDetailModel{
		client:     api.NewClient(baseURL, idToken),
		postID:     postID,
		spinner:    NewSpinner(),
		loading:    true,
		keys:       NewPostDetailKeyMap(),
		help:       h,
		viewport:   newDetailViewport(),
		replyInput: newReplyTextarea(),
	}
}

// NewPostDetailModelWithPost creates a detail screen with post already loaded
func NewPostDetailModelWithPost(baseURL, idToken string, post models.Post) PostDetailModel {
	h := help.New()
	h.Styles = styles.HelpStyles()
	vp := newDetailViewport()
	m := PostDetailModel{
		client:     api.NewClient(baseURL, idToken),
		postID:     post.ID,
		post:       post,
		spinner:    NewSpinner(),
		loading:    true,
		keys:       NewPostDetailKeyMap(),
		help:       h,
		viewport:   vp,
		replyInput: newReplyTextarea(),
	}
	// Pre-populate viewport so post shows immediately while replies load
	w, _ := SafeDimensions(0, 0)
	m.viewport.SetContent(m.buildContent(w))
	return m
}

func (m PostDetailModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchPostAndReplies())
}

func (m PostDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// If composing, route keys to textarea
		if m.composing {
			switch {
			case key.Matches(msg, m.keys.Send):
				// Send reply
				content := strings.TrimSpace(m.replyInput.Value())
				if content != "" && !m.replySending {
					m.replySending = true
					m.replyErr = nil
					return m, m.sendReply(content)
				}
				return m, nil
			case msg.String() == "esc":
				// Exit compose mode
				m.composing = false
				m.replyInput.Blur()
				m.resizeViewport()
				return m, nil
			default:
				// Forward to textarea
				var cmd tea.Cmd
				m.replyInput, cmd = m.replyInput.Update(msg)
				return m, cmd
			}
		}

		// Normal (non-compose) key handling
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case key.Matches(msg, m.keys.Back):
			return m, func() tea.Msg { return BackToFeedMsg{} }
		case key.Matches(msg, m.keys.Refresh):
			m.loading = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.fetchPostAndReplies())
		case key.Matches(msg, m.keys.Reply):
			m.composing = true
			m.replyErr = nil
			m.replyInput.SetWidth(m.width - 6)
			m.replyInput.Focus()
			m.resizeViewport()
			return m, m.replyInput.Focus()
		case key.Matches(msg, m.keys.Save):
			if !m.bookmarking && !m.bookmarked {
				m.bookmarking = true
				m.bookmarkErr = nil
				return m, m.addBookmark()
			}
		case key.Matches(msg, m.keys.Profile):
			username := m.post.AuthorUsername
			return m, func() tea.Msg { return OpenProfileMsg{Username: username} }
		}
		// Everything else (j/k, g/G, pgup/pgdn, etc.) falls through to viewport

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		vpHeight := msg.Height - headerHeight - footerHeight
		if vpHeight < 1 {
			vpHeight = 1
		}
		m.viewport.Width = msg.Width
		m.viewport.Height = vpHeight
		if !m.ready {
			m.ready = true
		}
		// Rebuild content at new width if we have data
		if m.post.ID != "" {
			m.viewport.SetContent(m.buildContent(msg.Width))
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		// Rebuild viewport content while loading so spinner animates
		if m.loading && m.post.ID != "" {
			w, _ := SafeDimensions(m.width, m.height)
			m.viewport.SetContent(m.buildContent(w))
		}
		return m, cmd

	case PostDetailLoadedMsg:
		m.loading = false
		m.post = msg.Post
		m.replies = msg.Replies
		w, _ := SafeDimensions(m.width, m.height)
		m.viewport.SetContent(m.buildContent(w))
		m.viewport.GotoTop()

	case PostDetailErrorMsg:
		m.loading = false
		m.err = msg.Err

	case ReplyCreatedMsg:
		m.replySending = false
		m.composing = false
		m.replyInput.Reset()
		m.replyInput.Blur()
		m.resizeViewport()
		// Re-fetch to show the new reply
		m.loading = true
		return m, tea.Batch(m.spinner.Tick, m.fetchPostAndReplies())

	case ReplyErrorMsg:
		m.replySending = false
		m.replyErr = msg.Err

	case BookmarkAddedMsg:
		m.bookmarking = false
		m.bookmarked = true
		m.post.BookmarksCount++
		w, _ := SafeDimensions(m.width, m.height)
		m.viewport.SetContent(m.buildContent(w))

	case BookmarkAddErrorMsg:
		m.bookmarking = false
		m.bookmarkErr = msg.Err

	case ThemeChangedMsg:
		m.spinner.Style = styles.Spinner
		m.help.Styles = styles.HelpStyles()
		m.replyInput.FocusedStyle.CursorLine = lipgloss.NewStyle().Background(styles.ColorBgSelect)
		m.replyInput.FocusedStyle.Base = lipgloss.NewStyle().Foreground(styles.ColorNormal)
		m.replyInput.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(styles.ColorMuted)
		m.replyInput.BlurredStyle = m.replyInput.FocusedStyle
		// Rebuild content with new theme colors
		if m.post.ID != "" {
			w, _ := SafeDimensions(m.width, m.height)
			m.viewport.SetContent(m.buildContent(w))
		}

	}

	// Forward to viewport for scroll handling
	var vpCmd tea.Cmd
	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, vpCmd)

	return m, tea.Batch(cmds...)
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

	// Viewport content
	b.WriteString(m.viewport.View())
	b.WriteString("\n")

	// Compose area (if active)
	if m.composing {
		b.WriteString(m.renderComposeArea(w))
		b.WriteString("\n")
	}

	// Footer
	b.WriteString(m.renderFooter(w))

	return b.String()
}

func (m PostDetailModel) renderHeader(width int) string {
	return RenderHeader("▓▒░ ENTRY VIEWER ░▒▓", width) + "\n"
}

func (m PostDetailModel) renderFooter(width int) string {
	navHint := m.help.View(m.keys)
	navWidth := lipgloss.Width(navHint)

	var status string
	if m.bookmarking {
		status = styles.Dim.Render(" [saving...]")
	} else if m.bookmarked {
		status = styles.Normal.Render(" [■ saved]")
	} else if m.bookmarkErr != nil {
		status = styles.Error.Render(" [save failed: " + m.bookmarkErr.Error() + "]")
	}
	statusWidth := lipgloss.Width(status)

	dividerWidth := width - navWidth - statusWidth - 1
	if dividerWidth < 1 {
		dividerWidth = 1
	}
	return styles.Divider(dividerWidth) + status + " " + navHint
}

func (m PostDetailModel) renderComposeArea(width int) string {
	borderStyle := lipgloss.NewStyle().Foreground(styles.ColorBright)
	titleStyle := lipgloss.NewStyle().Foreground(styles.ColorBright).Bold(true)

	title := "COMPOSE REPLY"
	if m.replySending {
		title = "SENDING..."
	}

	innerWidth := width - 4
	if innerWidth < 20 {
		innerWidth = 60
	}

	dashesLen := innerWidth - len(title) - 4
	if dashesLen < 1 {
		dashesLen = 1
	}

	top := borderStyle.Render("╭─ ") + titleStyle.Render(title) + borderStyle.Render(" "+strings.Repeat("─", dashesLen)+"╮")
	bottom := borderStyle.Render("╰" + strings.Repeat("─", innerWidth+2) + "╯")

	var mid strings.Builder
	if m.replyErr != nil {
		mid.WriteString(borderStyle.Render("│ "))
		mid.WriteString(styles.Error.Render("Error: " + m.replyErr.Error()))
		mid.WriteString("\n")
	}
	mid.WriteString(borderStyle.Render("│ "))
	// Render textarea lines within box
	taView := m.replyInput.View()
	taLines := strings.Split(taView, "\n")
	for i, line := range taLines {
		if i > 0 {
			mid.WriteString(borderStyle.Render("│ "))
		}
		mid.WriteString(line)
		if i < len(taLines)-1 {
			mid.WriteString("\n")
		}
	}
	mid.WriteString("\n")
	mid.WriteString(borderStyle.Render("│ "))
	mid.WriteString(styles.Dim.Render("[ctrl+s] send  [esc] cancel"))
	mid.WriteString("\n")

	return top + "\n" + mid.String() + bottom
}

func (m *PostDetailModel) resizeViewport() {
	vpHeight := m.height - headerHeight - footerHeight
	if m.composing {
		vpHeight -= composeHeight
	}
	if vpHeight < 1 {
		vpHeight = 1
	}
	m.viewport.Height = vpHeight
}

func (m PostDetailModel) addBookmark() tea.Cmd {
	return func() tea.Msg {
		id, err := m.client.CreateBookmark(m.postID)
		if err != nil {
			return BookmarkAddErrorMsg{Err: err}
		}
		return BookmarkAddedMsg{BookmarkID: id}
	}
}

func (m PostDetailModel) sendReply(content string) tea.Cmd {
	return func() tea.Msg {
		replyID, err := m.client.CreateReply(m.postID, content)
		if err != nil {
			return ReplyErrorMsg{Err: err}
		}
		return ReplyCreatedMsg{ReplyID: replyID}
	}
}

func (m PostDetailModel) renderLoadingScreen(width, height int) string {
	loadingBox := styles.DataBox("DECRYPTING TRANSMISSION",
		"\n"+
			"  "+m.spinner.View()+styles.Normal.Render(" Accessing secured data...")+"\n"+
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

	// Metadata box
	metaWidth := 50
	metaContent := styles.Username.Render("@"+m.post.AuthorUsername) + "\n" + styles.Dim.Render(TimeAgo(m.post.CreatedAt))
	b.WriteString(renderBox("POST INFO", metaContent, metaWidth))
	b.WriteString("\n\n")

	// Message content box
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
		b.WriteString(m.spinner.View() + styles.Normal.Render(" Loading replies..."))
	} else if len(m.replies) == 0 {
		b.WriteString(renderBox("REPLIES", "No replies yet", 40))
	} else {
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

	innerWidth := width - 4
	if innerWidth < 10 {
		innerWidth = 10
	}

	titleRendered := titleStyle.Render(title)
	titleVisualLen := lipgloss.Width(title)
	remainingDashes := width - 5 - titleVisualLen
	if remainingDashes < 1 {
		remainingDashes = 1
	}
	top := borderStyle.Render("╭─ ") + titleRendered + borderStyle.Render(" "+strings.Repeat("─", remainingDashes)+"╮")

	bottom := borderStyle.Render("╰" + strings.Repeat("─", width-2) + "╯")

	contentStyle := lipgloss.NewStyle().Foreground(styles.ColorNormal)

	var middle strings.Builder
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		wrappedLines := wrapText(line, innerWidth)
		for _, wl := range wrappedLines {
			// Apply theme foreground to each line
			styled := contentStyle.Render(wl)
			lineWidth := lipgloss.Width(styled)
			padding := innerWidth - lineWidth
			if padding < 0 {
				padding = 0
			}
			middle.WriteString(borderStyle.Render("│"))
			middle.WriteString(" ")
			middle.WriteString(styled)
			middle.WriteString(strings.Repeat(" ", padding))
			middle.WriteString(" ")
			middle.WriteString(borderStyle.Render("│"))
			middle.WriteString("\n")
		}
	}

	return top + "\n" + middle.String() + bottom
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

// Composing returns true when the reply textarea is active
func (m PostDetailModel) Composing() bool { return m.composing }

// SetSize updates the view dimensions
func (m *PostDetailModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	vpHeight := height - headerHeight - footerHeight
	if vpHeight < 1 {
		vpHeight = 1
	}
	m.viewport.Width = width
	m.viewport.Height = vpHeight
	m.ready = true
}
