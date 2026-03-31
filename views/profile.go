package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"

	"github.com/unremarkablegarden/cyberspace-tui-go/api"
	"github.com/unremarkablegarden/cyberspace-tui-go/models"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// OpenProfileMsg is sent when the user wants to view a profile
type OpenProfileMsg struct{ Username string }

// ProfileLoadedMsg is sent when a user's profile and posts are fetched
type ProfileLoadedMsg struct {
	User   models.User
	Posts  []models.Post
	Cursor string
}

// MoreProfilePostsLoadedMsg is sent when more posts are loaded
type MoreProfilePostsLoadedMsg struct {
	Posts  []models.Post
	Cursor string
}

// ProfileErrorMsg is sent when loading a profile fails
type ProfileErrorMsg struct{ Err error }

// BackFromProfileMsg is sent when navigating back from a profile
type BackFromProfileMsg struct{}

// profileHeaderHeight is the number of lines reserved for the profile info section
const profileHeaderHeight = 9

// ProfileModel is the user profile screen
type ProfileModel struct {
	username    string
	user        models.User
	list        list.Model
	loading     bool
	loadingMore bool
	spinner     spinner.Model
	err         error
	client      *api.Client
	nextCursor  string
	hasMore     bool
	width       int
	height      int
	keys        ProfileKeyMap
	help        help.Model
}

// NewProfileModel creates a new profile screen for the given username
func NewProfileModel(baseURL, idToken, username string) ProfileModel {
	delegate := PostDelegate{}
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.SetShowTitle(false)
	l.SetShowFilter(false)
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	l.Styles = styles.ListStyles()
	l.Paginator.ActiveDot = styles.Bright.Render("▄")
	l.Paginator.InactiveDot = styles.Dark.Render("▄")
	l.KeyMap.Quit.SetEnabled(false)
	l.KeyMap.ForceQuit.SetEnabled(false)

	h := help.New()
	h.Styles = styles.HelpStyles()

	return ProfileModel{
		username: username,
		list:     l,
		client:   api.NewClient(baseURL, idToken),
		spinner:  NewSpinner(),
		loading:  true,
		keys:     NewProfileKeyMap(),
		help:     h,
	}
}

func (m ProfileModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchProfile())
}

func (m ProfileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case key.Matches(msg, m.keys.Back):
			return m, func() tea.Msg { return BackFromProfileMsg{} }
		case key.Matches(msg, m.keys.Refresh):
			m.loading = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.fetchProfile())
		case key.Matches(msg, m.keys.Open):
			switch it := m.list.SelectedItem().(type) {
			case PostItem:
				post := it.Post
				return m, func() tea.Msg { return OpenPostMsg{Post: post} }
			case LoadMoreItem:
				if !m.loadingMore {
					m.loadingMore = true
					return m, tea.Batch(m.spinner.Tick, m.fetchMorePosts())
				}
			}
		}

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionRelease && !m.loading {
			for _, item := range m.list.Items() {
				if pi, ok := item.(PostItem); ok {
					if zone.Get(pi.Post.ID).InBounds(msg) {
						post := pi.Post
						return m, func() tea.Msg { return OpenPostMsg{Post: post} }
					}
				}
			}
			if zone.Get("load-more-profile").InBounds(msg) && m.hasMore && !m.loadingMore {
				m.loadingMore = true
				return m, tea.Batch(m.spinner.Tick, m.fetchMorePosts())
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		listHeight := msg.Height - profileHeaderHeight - 4
		if listHeight < 1 {
			listHeight = 1
		}
		m.list.SetSize(msg.Width, listHeight)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case ProfileLoadedMsg:
		m.loading = false
		m.err = nil
		m.user = msg.User
		items := postsToItems(msg.Posts)
		m.nextCursor = msg.Cursor
		m.hasMore = msg.Cursor != ""
		if m.hasMore {
			items = append(items, LoadMoreItem{})
		}
		cmd := m.list.SetItems(items)
		return m, cmd

	case MoreProfilePostsLoadedMsg:
		m.loadingMore = false
		m.nextCursor = msg.Cursor
		m.hasMore = msg.Cursor != ""
		var items []list.Item
		for _, existing := range m.list.Items() {
			if _, ok := existing.(LoadMoreItem); ok {
				continue
			}
			items = append(items, existing)
		}
		for _, p := range msg.Posts {
			items = append(items, PostItem{Post: p})
		}
		if m.hasMore {
			items = append(items, LoadMoreItem{})
		}
		cmd := m.list.SetItems(items)
		return m, cmd

	case ProfileErrorMsg:
		m.loading = false
		m.loadingMore = false
		m.err = msg.Err

	case ThemeChangedMsg:
		m.spinner.Style = styles.Spinner
		m.list.Styles = styles.ListStyles()
		m.help.Styles = styles.HelpStyles()
		m.list.Paginator.ActiveDot = styles.Bright.Render("▄")
		m.list.Paginator.InactiveDot = styles.Dark.Render("▄")
	}

	if !m.loading {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m ProfileModel) View() string {
	w, h := SafeDimensions(m.width, m.height)

	if m.loading {
		loadingBox := styles.DataBox("ACCESSING USER DATA",
			"\n"+
				"  "+m.spinner.View()+styles.Normal.Render(" Loading @"+m.username+"...")+"\n"+
				"\n"+
				"  "+styles.Dim.Render("Retrieving profile data...")+"\n",
			50)
		return FullScreen(loadingBox, w, h, lipgloss.Center, lipgloss.Center)
	}

	if m.err != nil {
		errorBox := styles.AlertBox(m.err.Error(), "error", 50) +
			"\n\n" +
			styles.Dim.Render("Press [esc] to go back, [r] to retry")
		return FullScreen(errorBox, w, h, lipgloss.Center, lipgloss.Center)
	}

	var b strings.Builder
	b.WriteString(RenderHeader("▓▒░ PROFILE ░▒▓", w))
	b.WriteString(m.renderProfileInfo(w))
	b.WriteString(m.list.View())
	b.WriteString("\n")
	b.WriteString(m.renderFooter(w))

	return b.String()
}

func (m ProfileModel) renderProfileInfo(width int) string {
	borderStyle := lipgloss.NewStyle().Foreground(styles.ColorDim)
	innerWidth := width - 4
	if innerWidth < 40 {
		innerWidth = 40
	}

	var content strings.Builder

	// Username + display name
	name := styles.Username.Render("@" + m.user.Username)
	if m.user.DisplayName != "" {
		name += styles.Normal.Render("  " + m.user.DisplayName)
	}
	content.WriteString(name + "\n")

	// Member since
	if !m.user.CreatedAt.IsZero() {
		content.WriteString(styles.Dim.Render("joined " + m.user.CreatedAt.Format("Jan 2006")) + "\n")
	}

	// Bio
	if m.user.Bio != "" {
		content.WriteString("\n")
		content.WriteString(styles.Normal.Render(m.user.Bio) + "\n")
	}

	// Website
	if m.user.WebsiteName != "" || m.user.WebsiteURL != "" {
		label := m.user.WebsiteName
		if label == "" {
			label = m.user.WebsiteURL
		}
		content.WriteString(styles.Dim.Render("⬡ " + label) + "\n")
	}

	// Location
	if m.user.LocationName != "" {
		content.WriteString(styles.Dim.Render("⌖ " + m.user.LocationName) + "\n")
	}

	// Posts header
	postCount := len(m.list.Items())
	postsLabel := fmt.Sprintf("POSTS [%d]", postCount)
	if m.hasMore {
		postsLabel = "POSTS"
	}

	top := borderStyle.Render("╭─ ") +
		lipgloss.NewStyle().Foreground(styles.ColorBright).Bold(true).Render(postsLabel) +
		borderStyle.Render(" "+strings.Repeat("─", innerWidth-len(postsLabel)-2)+"╮")

	// Render content lines inside box
	var mid strings.Builder
	for _, line := range strings.Split(strings.TrimRight(content.String(), "\n"), "\n") {
		wrappedLines := wrapText(line, innerWidth)
		for _, wl := range wrappedLines {
			lineWidth := lipgloss.Width(wl)
			pad := innerWidth - lineWidth
			if pad < 0 {
				pad = 0
			}
			mid.WriteString(borderStyle.Render("│") + " " + wl + strings.Repeat(" ", pad) + " " + borderStyle.Render("│") + "\n")
		}
	}

	return top + "\n" + mid.String()
}

func (m ProfileModel) renderFooter(width int) string {
	helpView := m.help.View(m.keys)
	paginatorView := m.list.Paginator.View()

	helpWidth := lipgloss.Width(helpView)
	paginatorWidth := lipgloss.Width(paginatorView)

	dividerWidth := width - helpWidth - paginatorWidth - 2
	if dividerWidth < 1 {
		dividerWidth = 1
	}

	return helpView + " " + styles.Divider(dividerWidth) + " " + paginatorView
}

// SetSize updates the view dimensions
func (m *ProfileModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	listHeight := height - profileHeaderHeight - 4
	if listHeight < 1 {
		listHeight = 1
	}
	m.list.SetSize(width, listHeight)
}

func (m ProfileModel) fetchProfile() tea.Cmd {
	return func() tea.Msg {
		user, err := m.client.FetchUser(m.username)
		if err != nil {
			return ProfileErrorMsg{Err: err}
		}
		posts, cursor, err := m.client.FetchUserPosts(m.username, 20)
		if err != nil {
			return ProfileErrorMsg{Err: err}
		}
		return ProfileLoadedMsg{User: *user, Posts: posts, Cursor: cursor}
	}
}

func (m ProfileModel) fetchMorePosts() tea.Cmd {
	return func() tea.Msg {
		posts, cursor, err := m.client.FetchMoreUserPosts(m.username, 20, m.nextCursor)
		if err != nil {
			return ProfileErrorMsg{Err: err}
		}
		return MoreProfilePostsLoadedMsg{Posts: posts, Cursor: cursor}
	}
}
