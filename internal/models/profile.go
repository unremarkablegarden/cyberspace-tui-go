package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/external/api"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/messages"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/items"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/keymaps"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// profileHeaderHeight is the number of lines reserved for the profile info section
const profileHeaderHeight = 12

// ProfileModel is the user profile screen
type ProfileModel struct {
	username        string
	currentUsername string
	isOwnProfile    bool
	user            entities.User
	list            list.Model
	loading         bool
	loadingMore     bool
	spinner         spinner.Model
	err             error
	client          *api.Client
	nextCursor      string
	hasMore         bool
	width           int
	height          int
	keys            keymaps.AppKeybinds
	help            help.Model
	prevMsg         messages.PrevMessage
	// follow state
	isFollowing   bool
	followID      string
	followLoaded  bool
	followPending bool
}

// NewProfileModel creates a new profile screen for the given username
func NewProfileModel(
	client *api.Client,
	keymap keymaps.AppKeybinds,
	username,
	currentUsername string,
	prevMsg messages.PrevMessage,
) ProfileModel {
	delegate := items.PostDelegate{}
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

	isOwn := currentUsername != "" && username == currentUsername

	return ProfileModel{
		username:        username,
		currentUsername: currentUsername,
		isOwnProfile:    isOwn,
		list:            l,
		client:          client,
		spinner:         items.NewSpinner(),
		loading:         true,
		keys:            keymap,
		help:            h,
		prevMsg:         prevMsg,
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
		switch msg.String() {
		case m.keys.GlobalKeybinds.Open:
			switch it := m.list.SelectedItem().(type) {
			case items.PostItem:
				post := it.Post
				return m, func() tea.Msg {
					return messages.SwitchToPostDetail{
						Post:        post,
						BackMessage: messages.SwitchToProfile{Username: m.username},
					}
				}
			case items.LoadMoreItem:
				if !m.loadingMore {
					m.loadingMore = true
					return m, tea.Batch(m.spinner.Tick, m.fetchMorePosts())
				}
			}
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
			m.followLoaded = false
			return m, tea.Batch(m.spinner.Tick, m.fetchProfile())
		case m.keys.ProfileKeybinds.Follow:
			if !m.isOwnProfile && m.followLoaded && !m.followPending {
				m.followPending = true
				return m, tea.Batch(m.spinner.Tick, m.toggleFollow())
			}
		case m.keys.ProfileKeybinds.EditProfile:
			if m.isOwnProfile {
				user := m.user
				return m, func() tea.Msg { return messages.SwitchToEditProfile{User: user} }
			}
		}

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionRelease && !m.loading {
			for _, item := range m.list.Items() {
				if pi, ok := item.(items.PostItem); ok {
					if zone.Get(pi.Post.ID).InBounds(msg) {
						post := pi.Post
						return m, func() tea.Msg {
							return messages.SwitchToPostDetail{
								Post:        post,
								BackMessage: messages.SwitchToProfile{Username: m.username},
							}
						}
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
		listHeight := max(msg.Height-profileHeaderHeight-4, 1)
		m.list.SetSize(msg.Width, listHeight)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case messages.ProfileLoadedMsg:
		m.loading = false
		m.loadingMore = false
		m.nextCursor = msg.Cursor
		m.err = nil
		m.hasMore = msg.Cursor != ""
		m.user = msg.User

		localItems := items.BuildListItems(
			&m.list, msg.IsAdditional, m.hasMore,
			items.PostsToItems(msg.Posts),
		)
		cmd := m.list.SetItems(localItems)

		// Fetch follow status for other users
		if !m.isOwnProfile && !msg.IsAdditional {
			return m, tea.Batch(cmd, m.fetchFollowStatus())
		}
		return m, cmd

	case messages.ProfileLoadedErrMsg:
		m.loading = false
		m.loadingMore = false
		m.err = msg.Err

	case messages.ProfileFollowMsg:
		if msg.InitialLoading {
			m.followLoaded = true
			m.isFollowing = msg.IsFollowing
			m.followID = msg.FollowID
		} else {
			m.followPending = false
			m.isFollowing = msg.IsFollowing
			m.followID = msg.FollowID
		}

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
	w, h := items.SafeDimensions(m.width, m.height)

	if m.loading {
		loadingBox := styles.DataBox("ACCESSING USER DATA",
			"\n"+
				"  "+m.spinner.View()+styles.Normal.Render(" Loading @"+m.username+"...")+"\n"+
				"\n"+
				"  "+styles.Dim.Render("Retrieving profile data...")+"\n",
			50)
		return items.FullScreen(loadingBox, w, h, lipgloss.Center, lipgloss.Center)
	}

	if m.err != nil {
		errorBox := styles.AlertBox(m.err.Error(), "error", 50) +
			"\n\n" +
			styles.Dim.Render("Press [esc] to go back, [r] to retry")
		return items.FullScreen(errorBox, w, h, lipgloss.Center, lipgloss.Center)
	}

	var b strings.Builder
	b.WriteString(items.RenderHeader("▓▒░ PROFILE ░▒▓", w))
	b.WriteString(m.renderProfileInfo(w))
	b.WriteString(m.list.View())
	b.WriteString("\n")
	b.WriteString(m.renderFooter(w))

	return b.String()
}

func (m ProfileModel) renderProfileInfo(width int) string {
	borderStyle := lipgloss.NewStyle().Foreground(styles.ColorDim)
	innerWidth := max(width-4, 40)

	var content strings.Builder

	// Username + display name
	name := styles.Username.Render("@" + m.user.Username)
	if m.user.DisplayName != "" {
		name += styles.Normal.Render("  " + m.user.DisplayName)
	}
	content.WriteString(name + "\n")

	// Member since
	if !m.user.CreatedAt.IsZero() {
		content.WriteString(styles.Dim.Render("joined "+m.user.CreatedAt.Format("Jan 2006")) + "\n")
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
		content.WriteString(styles.Dim.Render("⬡ "+label) + "\n")
	}

	// Location
	if m.user.LocationName != "" {
		content.WriteString(styles.Dim.Render("⌖ "+m.user.LocationName) + "\n")
	}

	// Follow status / edit button
	content.WriteString("\n")
	if m.isOwnProfile {
		content.WriteString(styles.Dim.Render("[e] edit profile") + "\n")
	} else if m.followPending {
		content.WriteString(m.spinner.View() + styles.Dim.Render(" ...") + "\n")
	} else if m.followLoaded {
		if m.isFollowing {
			content.WriteString(styles.Success.Render("● following") + styles.Dim.Render("  [f] unfollow") + "\n")
		} else {
			content.WriteString(styles.Dim.Render("[f] follow") + "\n")
		}
	} else {
		content.WriteString(styles.Dim.Render("...") + "\n")
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
	for line := range strings.SplitSeq(strings.TrimRight(content.String(), "\n"), "\n") {
		wrappedLines := items.WrapText(line, innerWidth)
		for _, wl := range wrappedLines {
			lineWidth := lipgloss.Width(wl)
			pad := max(innerWidth-lineWidth, 0)
			mid.WriteString(borderStyle.Render("│") + " " + wl + strings.Repeat(" ", pad) + " " + borderStyle.Render("│") + "\n")
		}
	}

	return top + "\n" + mid.String()
}

func (m ProfileModel) renderFooter(width int) string {
	helpView := m.help.View(m.keys.ProfileKeybinds)
	paginatorView := m.list.Paginator.View()

	helpWidth := lipgloss.Width(helpView)
	paginatorWidth := lipgloss.Width(paginatorView)

	dividerWidth := max(width-helpWidth-paginatorWidth-2, 1)

	return helpView + " " + styles.Divider(dividerWidth) + " " + paginatorView
}

// Username returns the username this profile is showing
func (m ProfileModel) Username() string {
	return m.username
}

// SetSize updates the view dimensions
func (m *ProfileModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	listHeight := max(height-profileHeaderHeight-4, 1)
	m.list.SetSize(width, listHeight)
}

func (m ProfileModel) fetchProfile() tea.Cmd {
	return func() tea.Msg {
		user, err := m.client.FetchUser(m.username)
		if err != nil {
			return messages.ProfileLoadedErrMsg{Err: err}
		}
		posts, cursor, err := m.client.FetchUserPosts(m.username, 20)
		if err != nil {
			return messages.ProfileLoadedErrMsg{Err: err}
		}
		return messages.ProfileLoadedMsg{User: *user, Posts: posts, Cursor: cursor}
	}
}

func (m ProfileModel) fetchMorePosts() tea.Cmd {
	return func() tea.Msg {
		posts, cursor, err := m.client.FetchMoreUserPosts(m.username, 20, m.nextCursor)
		if err != nil {
			return messages.ProfileLoadedErrMsg{Err: err}
		}
		return messages.ProfileLoadedMsg{User: m.user, Posts: posts, Cursor: cursor, IsAdditional: true}
	}
}

func (m ProfileModel) fetchFollowStatus() tea.Cmd {
	userID := m.user.ID
	return func() tea.Msg {
		follows, err := m.client.FetchMyFollowing(50)
		if err != nil {
			// Non-fatal: just show follow button without status
			return messages.ProfileFollowMsg{IsFollowing: false, FollowID: "", InitialLoading: true}
		}
		for _, f := range follows {
			if f.FollowedID == userID {
				return messages.ProfileFollowMsg{IsFollowing: true, FollowID: f.ID, InitialLoading: true}
			}
		}
		return messages.ProfileFollowMsg{IsFollowing: false, FollowID: "", InitialLoading: true}
	}
}

func (m ProfileModel) toggleFollow() tea.Cmd {
	if m.isFollowing {
		followID := m.followID
		return func() tea.Msg {
			if err := m.client.Unfollow(followID); err != nil {
				return messages.ProfileFollowMsg{IsFollowing: true, FollowID: followID}
			}
			return messages.ProfileFollowMsg{IsFollowing: false, FollowID: ""}
		}
	}

	userID := m.user.ID
	return func() tea.Msg {
		newFollowID, err := m.client.FollowUser(userID)
		if err != nil {
			return messages.ProfileFollowMsg{IsFollowing: false, FollowID: ""}
		}
		return messages.ProfileFollowMsg{IsFollowing: true, FollowID: newFollowID}
	}
}
