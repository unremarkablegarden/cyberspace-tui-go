package models

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

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/external/api"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/messages"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/items"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/keymaps"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

const headerHeight = 2  // title bar + blank line
const footerHeight = 2  // divider + status line
const hintHeight = 1    // contextual actions hint line
const composeHeight = 6 // textarea height including border

// PostDetailModel is the post detail screen
type PostDetailModel struct {
	post             entities.Post
	replies          []entities.Reply
	loading          bool
	spinner          spinner.Model
	err              error
	client           *api.Client
	postID           string
	currentUsername  string
	width            int
	height           int
	keys             keymaps.AppKeybinds
	help             help.Model
	viewport         viewport.Model
	ready            bool // true once we've received a WindowSizeMsg
	replyInput       textarea.Model
	composing        bool
	replySending     bool
	replyErr         error
	bookmarking      bool
	bookmarked       bool
	bookmarkErr      error
	confirmingDelete bool
	deleting         bool
	deleteErr        error
	prevMsg          messages.PrevMessage
}

// NewPostDetailModel creates a detail screen with post already loaded
func NewPostDetailModel(
	client *api.Client,
	keymap keymaps.AppKeybinds,
	post entities.Post,
	currentUsername string,
	prevMsg messages.PrevMessage,
) PostDetailModel {
	h := help.New()
	h.Styles = styles.HelpStyles()
	vp := newDetailViewport()
	m := PostDetailModel{
		client:          client,
		postID:          post.ID,
		post:            post,
		currentUsername: currentUsername,
		spinner:         items.NewSpinner(),
		loading:         true,
		keys:            keymap,
		help:            h,
		viewport:        vp,
		replyInput:      newReplyTextarea(),
		prevMsg:         prevMsg,
	}
	// Pre-populate viewport so post shows immediately while replies load
	w, _ := items.SafeDimensions(0, 0)
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
			switch msg.String() {
			case m.keys.PostDetailKeybinds.Send:
				// Send reply
				content := strings.TrimSpace(m.replyInput.Value())
				if content != "" && !m.replySending {
					m.replySending = true
					m.replyErr = nil
					return m, m.sendReply(content)
				}
				return m, nil
			case "esc":
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

		// Confirm delete prompt
		if m.confirmingDelete {
			switch msg.String() {
			case "y", "Y":
				m.confirmingDelete = false
				m.deleting = true
				m.deleteErr = nil
				return m, tea.Batch(m.spinner.Tick, m.deletePost())
			case "n", "N", "esc":
				m.confirmingDelete = false
				return m, nil
			}
			return m, nil
		}

		// Normal (non-compose) key handling
		switch msg.String() {
		case m.keys.GlobalKeybinds.Quit:
			return m, tea.Quit
		case m.keys.GlobalKeybinds.Help:
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case m.keys.GlobalKeybinds.Back:
			return m, func() tea.Msg {
				if m.prevMsg != nil {
					return m.prevMsg
				}
				return messages.SwitchToFeed{}
			}
		case m.keys.GlobalKeybinds.Refresh:
			m.loading = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.fetchPostAndReplies())

		case m.keys.PostDetailKeybinds.Delete:
			if !m.deleting && m.currentUsername != "" && m.post.AuthorUsername == m.currentUsername {
				m.confirmingDelete = true
				return m, nil
			}
		case m.keys.PostDetailKeybinds.Reply:
			m.composing = true
			m.replyErr = nil
			m.replyInput.SetWidth(m.width - 6)
			m.replyInput.Focus()
			m.resizeViewport()
			return m, m.replyInput.Focus()
		case m.keys.PostDetailKeybinds.Save:
			if !m.bookmarking && !m.bookmarked {
				m.bookmarking = true
				m.bookmarkErr = nil
				return m, m.addBookmark()
			}
		case m.keys.PostDetailKeybinds.Profile:
			username := m.post.AuthorUsername
			return m, func() tea.Msg {
				return messages.SwitchToProfile{
					Username: username,
					BackMessage: messages.SwitchToPostDetail{
						Post:        m.post,
						BackMessage: m.prevMsg,
					},
				}
			}
		}
		// Everything else (j/k, g/G, pgup/pgdn, etc.) falls through to viewport

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		vpHeight := max(msg.Height-headerHeight-footerHeight-hintHeight, 1)
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
			w, _ := items.SafeDimensions(m.width, m.height)
			m.viewport.SetContent(m.buildContent(w))
		}
		return m, cmd

	case messages.PostDetailLoadedMsg:
		m.loading = false
		m.post = msg.Post
		m.replies = msg.Replies
		w, _ := items.SafeDimensions(m.width, m.height)
		m.viewport.SetContent(m.buildContent(w))
		m.viewport.GotoTop()

	case messages.PostDetailErrorMsg:
		m.loading = false
		m.err = msg.Err

	case messages.PostReplyCreatedMsg:
		m.replySending = false
		m.composing = false
		m.replyInput.Reset()
		m.replyInput.Blur()
		m.resizeViewport()
		// Re-fetch to show the new reply
		m.loading = true
		return m, tea.Batch(m.spinner.Tick, m.fetchPostAndReplies())

	case messages.PostReplyCreatedErrMsg:
		m.replySending = false
		m.replyErr = msg.Err

	case messages.PostBookmarkAddedMsg:
		m.bookmarking = false
		m.bookmarked = true
		m.post.BookmarksCount++
		w, _ := items.SafeDimensions(m.width, m.height)
		m.viewport.SetContent(m.buildContent(w))

	case messages.PostBookmarkAddedErrMsg:
		m.bookmarking = false
		m.bookmarkErr = msg.Err

	case messages.PostDeleteMsg:
		return m, func() tea.Msg { return messages.SwitchToFeed{} }

	case messages.PostDeleteErrMsg:
		m.deleting = false
		m.deleteErr = msg.Err

	case ThemeChangedMsg:
		m.spinner.Style = styles.Spinner
		m.help.Styles = styles.HelpStyles()
		m.replyInput.FocusedStyle.CursorLine = lipgloss.NewStyle().Background(styles.ColorBgSelect)
		m.replyInput.FocusedStyle.Base = lipgloss.NewStyle().Foreground(styles.ColorNormal)
		m.replyInput.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(styles.ColorMuted)
		m.replyInput.BlurredStyle = m.replyInput.FocusedStyle
		// Rebuild content with new theme colors
		if m.post.ID != "" {
			w, _ := items.SafeDimensions(m.width, m.height)
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
	w, h := items.SafeDimensions(m.width, m.height)

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

	// Contextual hint line
	b.WriteString(m.renderHints(w))
	b.WriteString("\n")

	// Footer
	b.WriteString(m.renderFooter(w))

	return b.String()
}

func (m PostDetailModel) renderHeader(width int) string {
	return items.RenderHeader("▓▒░ ENTRY VIEWER ░▒▓", width) + "\n"
}

func (m PostDetailModel) renderFooter(width int) string {
	navHint := m.help.View(m.keys.PostDetailHelpKeys())
	navWidth := lipgloss.Width(navHint)

	var status string
	if m.confirmingDelete {
		status = styles.Error.Render(" [delete post? y/n]")
	} else if m.deleting {
		status = styles.Dim.Render(" [deleting...]")
	} else if m.deleteErr != nil {
		status = styles.Error.Render(" [delete failed: " + m.deleteErr.Error() + "]")
	} else if m.bookmarking {
		status = styles.Dim.Render(" [saving...]")
	} else if m.bookmarked {
		status = styles.Normal.Render(" [■ saved]")
	} else if m.bookmarkErr != nil {
		status = styles.Error.Render(" [save failed: " + m.bookmarkErr.Error() + "]")
	}
	statusWidth := lipgloss.Width(status)

	dividerWidth := max(width-navWidth-statusWidth-1, 1)
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

	dashesLen := max(innerWidth-len(title)-4, 1)

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
	vpHeight := m.height - headerHeight - footerHeight - hintHeight
	if m.composing {
		vpHeight -= composeHeight
	}
	if vpHeight < 1 {
		vpHeight = 1
	}
	m.viewport.Height = vpHeight
}

func (m PostDetailModel) renderHints(width int) string {
	var hints []string

	if !m.composing {
		hints = append(hints, styles.Dim.Render("[c]")+styles.Normal.Render(" reply"))
		if m.bookmarked {
			hints = append(hints, styles.Success.Render("■ saved"))
		} else {
			hints = append(hints, styles.Dim.Render("[s]")+styles.Normal.Render(" save"))
		}
		hints = append(hints, styles.Dim.Render("[p]")+styles.Normal.Render(" @"+m.post.AuthorUsername))
		if m.currentUsername != "" && m.post.AuthorUsername == m.currentUsername {
			hints = append(hints, styles.Dim.Render("[D]")+styles.Error.Render(" delete"))
		}
		hints = append(hints, styles.Dim.Render("[b]")+styles.Normal.Render(" back"))
	}

	line := "  " + strings.Join(hints, styles.Dim.Render("  ·  "))
	lineWidth := lipgloss.Width(line)
	if lineWidth < width {
		line += strings.Repeat(" ", width-lineWidth)
	}
	return line
}

// replyNode is a node in the reply tree
type replyNode struct {
	Reply    entities.Reply
	Children []*replyNode
}

// buildReplyTree organises a flat reply list into a tree using ParentReplyID
func buildReplyTree(replies []entities.Reply) []*replyNode {
	nodes := make(map[string]*replyNode, len(replies))
	for i := range replies {
		nodes[replies[i].ID] = &replyNode{Reply: replies[i]}
	}
	var roots []*replyNode
	for _, r := range replies {
		node := nodes[r.ID]
		if r.ParentReplyID == "" {
			roots = append(roots, node)
		} else if parent, ok := nodes[r.ParentReplyID]; ok {
			parent.Children = append(parent.Children, node)
		} else {
			roots = append(roots, node)
		}
	}
	return roots
}

// renderReplyNode renders a reply and its children with indentation
func renderReplyNode(node *replyNode, depth int, isLast bool, contentWidth int) string {
	var b strings.Builder

	// Build indent prefix
	var prefix, childPrefix string
	if depth == 0 {
		prefix = ""
		childPrefix = ""
	} else {
		indent := strings.Repeat("│  ", depth-1)
		if isLast {
			prefix = indent + "└─ "
			childPrefix = indent + "   "
		} else {
			prefix = indent + "├─ "
			childPrefix = indent + "│  "
		}
	}

	// Header line: @username · time
	b.WriteString(styles.Dim.Render(prefix))
	b.WriteString(styles.Username.Render("@" + node.Reply.AuthorUsername))
	b.WriteString(styles.Dim.Render(" · " + items.TimeAgo(node.Reply.CreatedAt)))
	b.WriteString("\n")

	// Content lines with continuation indent
	content := items.StripMarkdownKeepNewlines(node.Reply.Content)
	for line := range strings.SplitSeq(content, "\n") {
		b.WriteString(styles.Dim.Render(childPrefix))
		b.WriteString(styles.Normal.Render(line))
		b.WriteString("\n")
	}

	// Render children
	for i, child := range node.Children {
		b.WriteString("\n")
		b.WriteString(renderReplyNode(child, depth+1, i == len(node.Children)-1, contentWidth))
	}

	return b.String()
}

func (m PostDetailModel) deletePost() tea.Cmd {
	return func() tea.Msg {
		if err := m.client.DeletePost(m.postID); err != nil {
			return messages.PostDeleteErrMsg{Err: err}
		}
		return messages.PostDeleteMsg{}
	}
}

func (m PostDetailModel) addBookmark() tea.Cmd {
	return func() tea.Msg {
		id, err := m.client.CreateBookmark(m.postID)
		if err != nil {
			return messages.PostBookmarkAddedErrMsg{Err: err}
		}
		return messages.PostBookmarkAddedMsg{BookmarkID: id}
	}
}

func (m PostDetailModel) sendReply(content string) tea.Cmd {
	return func() tea.Msg {
		replyID, err := m.client.CreateReply(m.postID, content)
		if err != nil {
			return messages.PostReplyCreatedErrMsg{Err: err}
		}
		return messages.PostReplyCreatedMsg{ReplyID: replyID}
	}
}

func (m PostDetailModel) renderLoadingScreen(width, height int) string {
	loadingBox := styles.DataBox("DECRYPTING TRANSMISSION",
		"\n"+
			"  "+m.spinner.View()+styles.Normal.Render(" Accessing secured data...")+"\n"+
			"\n"+
			"  "+styles.Dim.Render("Decoding neural patterns...")+"\n",
		50)

	return items.FullScreen(loadingBox, width, height, lipgloss.Center, lipgloss.Center)
}

func (m PostDetailModel) renderErrorScreen(width, height int) string {
	errorBox := styles.AlertBox(m.err.Error(), "error", 50) +
		"\n\n" +
		styles.Dim.Render("Press [ESC] to return to feed, [r] to retry")

	return items.FullScreen(errorBox, width, height, lipgloss.Center, lipgloss.Center)
}

func (m PostDetailModel) buildContent(width int) string {
	var b strings.Builder

	contentWidth := width - 2
	if contentWidth < 40 {
		contentWidth = 78
	}

	// Metadata box
	metaWidth := 50
	metaContent := styles.Username.Render("@"+m.post.AuthorUsername) + "\n" + styles.Dim.Render(items.TimeAgo(m.post.CreatedAt))
	b.WriteString(renderBox("POST INFO", metaContent, metaWidth))
	b.WriteString("\n\n")

	// Message content box
	cleanContent := items.StripMarkdownKeepNewlines(m.post.Content)
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
		roots := buildReplyTree(m.replies)
		var repliesContent strings.Builder
		for i, root := range roots {
			repliesContent.WriteString(renderReplyNode(root, 0, false, contentWidth))
			if i < len(roots)-1 {
				repliesContent.WriteString("\n")
			}
		}
		replyTitle := fmt.Sprintf("REPLIES [%d]", len(m.replies))
		b.WriteString(renderBox(replyTitle, strings.TrimRight(repliesContent.String(), "\n"), contentWidth))
	}

	return b.String()
}

// renderBox renders content in a box with title
func renderBox(title, content string, width int) string {
	borderStyle := lipgloss.NewStyle().Foreground(styles.ColorDim)
	titleStyle := lipgloss.NewStyle().Foreground(styles.ColorBright).Bold(true)

	innerWidth := max(width-4, 10)

	titleRendered := titleStyle.Render(title)
	titleVisualLen := lipgloss.Width(title)
	remainingDashes := max(width-5-titleVisualLen, 1)
	top := borderStyle.Render("╭─ ") + titleRendered + borderStyle.Render(" "+strings.Repeat("─", remainingDashes)+"╮")

	bottom := borderStyle.Render("╰" + strings.Repeat("─", width-2) + "╯")

	contentStyle := lipgloss.NewStyle().Foreground(styles.ColorNormal)

	var middle strings.Builder
	for line := range strings.SplitSeq(content, "\n") {
		wrappedLines := items.WrapText(line, innerWidth)
		for _, wl := range wrappedLines {
			// Apply theme foreground to each line
			styled := contentStyle.Render(wl)
			lineWidth := lipgloss.Width(styled)
			padding := max(innerWidth-lineWidth, 0)
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
		replies, err := m.client.FetchReplies(m.postID)
		if err != nil {
			return messages.PostDetailErrorMsg{Err: err}
		}

		return messages.PostDetailLoadedMsg{Post: post, Replies: replies}
	}
}

// Composing returns true when the reply textarea is active
func (m PostDetailModel) Composing() bool { return m.composing }

// SetSize updates the view dimensions
func (m *PostDetailModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	vpHeight := max(height-headerHeight-footerHeight-hintHeight, 1)
	m.viewport.Width = width
	m.viewport.Height = vpHeight
	m.ready = true
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
