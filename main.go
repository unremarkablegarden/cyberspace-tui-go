package main

import (
	"embed"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	zone "github.com/lrstanley/bubblezone"

	"github.com/unremarkablegarden/cyberspace-tui-go/api"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
	"github.com/unremarkablegarden/cyberspace-tui-go/views"
)

//go:embed themes/*.json
var themesFS embed.FS

// AppState represents the current screen
type AppState int

const (
	StateLogin AppState = iota
	StateFeed
	StatePostDetail
	StateCompose
	StateBookmarks
	StateNotifications
	StateProfile
	StateTopics
	StateTopicFeed
	StateEditProfile
	StateNotes
	StateNoteCompose
)

// ownUsernameMsg is sent after fetching the current user's username post-login
type ownUsernameMsg struct{ username string }

// Model is the main application model
type Model struct {
	state              AppState
	loginModel         views.LoginModel
	feedModel          views.FeedModel
	postDetailModel    views.PostDetailModel
	config             *Config
	baseURL            string
	width              int
	height             int
	showThemeSwitcher  bool
	themeSwitcherModel views.ThemeSwitcherModel
	composeModel       views.ComposeModel
	bookmarksModel      views.BookmarksModel
	notificationsModel  views.NotificationsModel
	profileModel        views.ProfileModel
	topicsModel         views.TopicsModel
	topicFeedModel      views.TopicFeedModel
	editProfileModel    views.EditProfileModel
	notesModel          views.NotesModel
	noteComposeModel    views.NoteComposeModel
	returnState         AppState
}

func initialModel(baseURL string, config *Config) Model {
	m := Model{
		baseURL: baseURL,
		config:  config,
	}

	// If we have a valid (non-expired) token, go to feed
	if config != nil && config.IDToken != "" && !config.IsExpired() {
		m.state = StateFeed
		m.feedModel = views.NewFeedModel(baseURL, config.IDToken)
	} else {
		m.state = StateLogin
		m.loginModel = views.NewLoginModel(baseURL)
	}

	return m
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	// Always request window size first
	cmds = append(cmds, tea.WindowSize())

	switch m.state {
	case StateLogin:
		cmds = append(cmds, m.loginModel.Init())
	case StateFeed:
		cmds = append(cmds, m.feedModel.Init())
	case StatePostDetail:
		cmds = append(cmds, m.postDetailModel.Init())
	case StateCompose:
		cmds = append(cmds, m.composeModel.Init())
	case StateBookmarks:
		cmds = append(cmds, m.bookmarksModel.Init())
	case StateNotifications:
		cmds = append(cmds, m.notificationsModel.Init())
	case StateProfile:
		cmds = append(cmds, m.profileModel.Init())
	case StateTopics:
		cmds = append(cmds, m.topicsModel.Init())
	case StateTopicFeed:
		cmds = append(cmds, m.topicFeedModel.Init())
	case StateEditProfile:
		cmds = append(cmds, m.editProfileModel.Init())
	case StateNotes:
		cmds = append(cmds, m.notesModel.Init())
	case StateNoteCompose:
		cmds = append(cmds, m.noteComposeModel.Init())
	}

	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// Toggle theme switcher with 't' (only when not on login screen and not already in switcher)
		composing := (m.state == StatePostDetail && m.postDetailModel.Composing()) || m.state == StateCompose
		if msg.String() == "t" && m.state != StateLogin && !m.showThemeSwitcher && !composing {
			m.showThemeSwitcher = true
			m.themeSwitcherModel = views.NewThemeSwitcherModel()
			m.themeSwitcherModel.SetSize(m.width, m.height)
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.showThemeSwitcher {
			m.themeSwitcherModel.SetSize(m.width, m.height)
		}

	case views.ThemeChangedMsg:
		// Save theme preference
		if m.config == nil {
			m.config = &Config{}
		}
		m.config.Theme = msg.ThemeKey
		if err := SaveConfig(m.config); err != nil {
			log.Printf("Failed to save theme preference: %v", err)
		}
		// Close the switcher
		m.showThemeSwitcher = false
		// Propagate to active child view
		switch m.state {
		case StateLogin:
			newLogin, cmd := m.loginModel.Update(msg)
			m.loginModel = newLogin.(views.LoginModel)
			return m, cmd
		case StateFeed:
			newFeed, cmd := m.feedModel.Update(msg)
			m.feedModel = newFeed.(views.FeedModel)
			return m, cmd
		case StatePostDetail:
			newDetail, cmd := m.postDetailModel.Update(msg)
			m.postDetailModel = newDetail.(views.PostDetailModel)
			return m, cmd
		case StateCompose:
			newCompose, cmd := m.composeModel.Update(msg)
			m.composeModel = newCompose.(views.ComposeModel)
			return m, cmd
		case StateBookmarks:
			newBookmarks, cmd := m.bookmarksModel.Update(msg)
			m.bookmarksModel = newBookmarks.(views.BookmarksModel)
			return m, cmd
		case StateNotifications:
			newNotifs, cmd := m.notificationsModel.Update(msg)
			m.notificationsModel = newNotifs.(views.NotificationsModel)
			return m, cmd
		case StateProfile:
			newProfile, cmd := m.profileModel.Update(msg)
			m.profileModel = newProfile.(views.ProfileModel)
			return m, cmd
		case StateTopics:
			newTopics, cmd := m.topicsModel.Update(msg)
			m.topicsModel = newTopics.(views.TopicsModel)
			return m, cmd
		case StateTopicFeed:
			newTopicFeed, cmd := m.topicFeedModel.Update(msg)
			m.topicFeedModel = newTopicFeed.(views.TopicFeedModel)
			return m, cmd
		case StateEditProfile:
			newEdit, cmd := m.editProfileModel.Update(msg)
			m.editProfileModel = newEdit.(views.EditProfileModel)
			return m, cmd
		}
		return m, nil



	case views.ThemeSwitcherClosedMsg:
		m.showThemeSwitcher = false
		return m, nil

	case ownUsernameMsg:
		if m.config != nil && msg.username != "" {
			m.config.Username = msg.username
			if err := SaveConfig(m.config); err != nil {
				log.Printf("Failed to save username: %v", err)
			}
		}
		return m, nil
	}

	// Route to theme switcher if open
	if m.showThemeSwitcher {
		newSwitcher, cmd := m.themeSwitcherModel.Update(msg)
		m.themeSwitcherModel = newSwitcher
		return m, cmd
	}

	switch m.state {
	case StateLogin:
		newLogin, cmd := m.loginModel.Update(msg)
		m.loginModel = newLogin.(views.LoginModel)

		// Check if login succeeded
		if loginMsg, ok := msg.(views.LoginSuccessMsg); ok {
			m.config = &Config{
				IDToken:      loginMsg.IDToken,
				RefreshToken: loginMsg.RefreshToken,
			}
			m.config.SetExpiry(DefaultTokenLifetimeSecs)

			if err := SaveConfig(m.config); err != nil {
				log.Printf("Failed to save config: %v", err)
			}

			// Transition to feed view
			m.state = StateFeed
			m.feedModel = views.NewFeedModel(m.baseURL, m.config.IDToken)
			m.feedModel.SetSize(m.width, m.height)
			// Fetch own username if not already stored
			var extraCmds []tea.Cmd
			extraCmds = append(extraCmds, m.feedModel.Init())
			if m.config.Username == "" {
				extraCmds = append(extraCmds, fetchOwnUsernameCmd(m.baseURL, m.config.IDToken))
			}
			return m, tea.Batch(extraCmds...)
		}

		return m, cmd

	case StateFeed:
		newFeed, cmd := m.feedModel.Update(msg)
		m.feedModel = newFeed.(views.FeedModel)

		// Check if user wants to view a profile
		if profileMsg, ok := msg.(views.OpenProfileMsg); ok {
			return m, m.openProfile(profileMsg.Username)
		}

		// Check if user wants to browse topics
		if _, ok := msg.(views.OpenTopicsMsg); ok {
			m.state = StateTopics
			m.topicsModel = views.NewTopicsModel(m.baseURL, m.config.IDToken)
			m.topicsModel.SetSize(m.width, m.height)
			return m, m.topicsModel.Init()
		}

		// Check if user wants to view notifications
		if _, ok := msg.(views.OpenNotificationsMsg); ok {
			m.state = StateNotifications
			m.notificationsModel = views.NewNotificationsModel(m.baseURL, m.config.IDToken)
			m.notificationsModel.SetSize(m.width, m.height)
			return m, m.notificationsModel.Init()
		}

		// Check if user wants to view bookmarks
		if _, ok := msg.(views.OpenBookmarksMsg); ok {
			m.state = StateBookmarks
			m.bookmarksModel = views.NewBookmarksModel(m.baseURL, m.config.IDToken)
			m.bookmarksModel.SetSize(m.width, m.height)
			return m, m.bookmarksModel.Init()
		}

		// Check if user wants to view notes
		if _, ok := msg.(views.OpenNotesMsg); ok {
			m.state = StateNotes
			m.notesModel = views.NewNotesModel(m.baseURL, m.config.IDToken)
			m.notesModel.SetSize(m.width, m.height)
			return m, m.notesModel.Init()
		}

		// Check if user wants to compose a new post
		if _, ok := msg.(views.OpenComposeMsg); ok {
			m.state = StateCompose
			m.composeModel = views.NewComposeModel(m.baseURL, m.config.IDToken)
			m.composeModel.SetSize(m.width, m.height)
			return m, m.composeModel.Init()
		}

		// Check if user wants to open a post
		if openMsg, ok := msg.(views.OpenPostMsg); ok {
			m.state = StatePostDetail
			m.postDetailModel = views.NewPostDetailModelWithPost(
				m.baseURL,
				m.config.IDToken,
				openMsg.Post,
				m.config.Username,
			)
			m.postDetailModel.SetSize(m.width, m.height)
			return m, m.postDetailModel.Init()
		}

		// Check if user wants to log out
		if _, ok := msg.(views.LogoutMsg); ok {
			_ = ClearConfig()
			m.config = nil
			m.state = StateLogin
			m.loginModel = views.NewLoginModel(m.baseURL)
			m.loginModel.SetSize(m.width, m.height)
			return m, m.loginModel.Init()
		}

		return m, cmd

	case StatePostDetail:
		newDetail, cmd := m.postDetailModel.Update(msg)
		m.postDetailModel = newDetail.(views.PostDetailModel)

		if _, ok := msg.(views.BackToFeedMsg); ok {
			if m.returnState != 0 {
				m.state = m.returnState
			} else {
				m.state = StateFeed
			}
			m.returnState = 0
			return m, nil
		}

		if profileMsg, ok := msg.(views.OpenProfileMsg); ok {
			return m, m.openProfile(profileMsg.Username)
		}

		return m, cmd

	case StateProfile:
		newProfile, cmd := m.profileModel.Update(msg)
		m.profileModel = newProfile.(views.ProfileModel)

		if _, ok := msg.(views.BackFromProfileMsg); ok {
			if m.returnState != 0 {
				m.state = m.returnState
			} else {
				m.state = StateFeed
			}
			m.returnState = 0
			return m, nil
		}

		if openMsg, ok := msg.(views.OpenPostMsg); ok {
			prev := m.state
			m.state = StatePostDetail
			m.returnState = prev
			m.postDetailModel = views.NewPostDetailModelWithPost(m.baseURL, m.config.IDToken, openMsg.Post, m.config.Username)
			m.postDetailModel.SetSize(m.width, m.height)
			return m, m.postDetailModel.Init()
		}

		if editMsg, ok := msg.(views.OpenEditProfileMsg); ok {
			m.state = StateEditProfile
			m.editProfileModel = views.NewEditProfileModel(m.baseURL, m.config.IDToken, editMsg.User)
			m.editProfileModel.SetSize(m.width, m.height)
			return m, m.editProfileModel.Init()
		}

		return m, cmd

	case StateEditProfile:
		newEdit, cmd := m.editProfileModel.Update(msg)
		m.editProfileModel = newEdit.(views.EditProfileModel)

		if doneMsg, ok := msg.(views.EditProfileDoneMsg); ok {
			m.state = StateProfile
			if doneMsg.Saved {
				// Refresh the profile to show updated data
				m.profileModel = views.NewProfileModel(
					m.baseURL, m.config.IDToken,
					m.profileModel.Username(), m.config.Username,
				)
				m.profileModel.SetSize(m.width, m.height)
				return m, m.profileModel.Init()
			}
			return m, nil
		}

		return m, cmd

	case StateNotifications:
		newNotifs, cmd := m.notificationsModel.Update(msg)
		m.notificationsModel = newNotifs.(views.NotificationsModel)

		if _, ok := msg.(views.BackFromNotificationsMsg); ok {
			m.state = StateFeed
			return m, nil
		}

		if openMsg, ok := msg.(views.OpenPostFromNotificationMsg); ok {
			m.state = StatePostDetail
			m.returnState = StateNotifications
			m.postDetailModel = views.NewPostDetailModel(m.baseURL, m.config.IDToken, openMsg.PostID, m.config.Username)
			m.postDetailModel.SetSize(m.width, m.height)
			return m, m.postDetailModel.Init()
		}

		return m, cmd

	case StateBookmarks:
		newBookmarks, cmd := m.bookmarksModel.Update(msg)
		m.bookmarksModel = newBookmarks.(views.BookmarksModel)

		if _, ok := msg.(views.BackToFeedFromBookmarksMsg); ok {
			m.state = StateFeed
			return m, nil
		}

		if openMsg, ok := msg.(views.OpenPostFromBookmarksMsg); ok {
			m.state = StatePostDetail
			m.returnState = StateBookmarks
			m.postDetailModel = views.NewPostDetailModelWithPost(
				m.baseURL,
				m.config.IDToken,
				openMsg.Post,
				m.config.Username,
			)
			m.postDetailModel.SetSize(m.width, m.height)
			return m, m.postDetailModel.Init()
		}

		if profileMsg, ok := msg.(views.OpenProfileMsg); ok {
			return m, m.openProfile(profileMsg.Username)
		}

		return m, cmd

	case StateTopics:
		newTopics, cmd := m.topicsModel.Update(msg)
		m.topicsModel = newTopics.(views.TopicsModel)

		if _, ok := msg.(views.BackFromTopicsMsg); ok {
			m.state = StateFeed
			return m, nil
		}

		if openMsg, ok := msg.(views.OpenTopicFeedMsg); ok {
			m.state = StateTopicFeed
			m.topicFeedModel = views.NewTopicFeedModel(m.baseURL, m.config.IDToken, openMsg.Topic)
			m.topicFeedModel.SetSize(m.width, m.height)
			return m, m.topicFeedModel.Init()
		}

		return m, cmd

	case StateTopicFeed:
		newTopicFeed, cmd := m.topicFeedModel.Update(msg)
		m.topicFeedModel = newTopicFeed.(views.TopicFeedModel)

		if _, ok := msg.(views.BackFromTopicFeedMsg); ok {
			m.state = StateTopics
			return m, nil
		}

		if openMsg, ok := msg.(views.OpenPostMsg); ok {
			m.state = StatePostDetail
			m.returnState = StateTopicFeed
			m.postDetailModel = views.NewPostDetailModelWithPost(m.baseURL, m.config.IDToken, openMsg.Post, m.config.Username)
			m.postDetailModel.SetSize(m.width, m.height)
			return m, m.postDetailModel.Init()
		}

		if profileMsg, ok := msg.(views.OpenProfileMsg); ok {
			return m, m.openProfile(profileMsg.Username)
		}

		return m, cmd

	case StateCompose:
		newCompose, cmd := m.composeModel.Update(msg)
		m.composeModel = newCompose.(views.ComposeModel)

		if _, ok := msg.(views.ComposeBackMsg); ok {
			m.state = StateFeed
			// Refresh feed so new post appears
			newFeed, feedCmd := m.feedModel.Update(views.RefreshFeedMsg{})
			m.feedModel = newFeed.(views.FeedModel)
			return m, feedCmd
		}

		return m, cmd

	case StateNotes:
		newNotes, cmd := m.notesModel.Update(msg)
		m.notesModel = newNotes.(views.NotesModel)

		if _, ok := msg.(views.BackFromNotesMsg); ok {
			m.state = StateFeed
			return m, nil
		}

		if openMsg, ok := msg.(views.OpenNoteComposeMsg); ok {
			m.state = StateNoteCompose
			m.noteComposeModel = views.NewNoteComposeModel(m.baseURL, m.config.IDToken, openMsg.Note, openMsg.IsEdit)
			m.noteComposeModel.SetSize(m.width, m.height)
			return m, m.noteComposeModel.Init()
		}

		return m, cmd

	case StateNoteCompose:
		newNoteCompose, cmd := m.noteComposeModel.Update(msg)
		m.noteComposeModel = newNoteCompose.(views.NoteComposeModel)

		if _, ok := msg.(views.NoteComposeDoneMsg); ok {
			m.state = StateNotes
			m.notesModel = views.NewNotesModel(m.baseURL, m.config.IDToken)
			m.notesModel.SetSize(m.width, m.height)
			return m, m.notesModel.Init()
		}

		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	var v string
	if m.showThemeSwitcher {
		v = m.themeSwitcherModel.View()
	} else {
		switch m.state {
		case StateLogin:
			v = m.loginModel.View()
		case StateFeed:
			v = m.feedModel.View()
		case StatePostDetail:
			v = m.postDetailModel.View()
		case StateCompose:
			v = m.composeModel.View()
		case StateBookmarks:
			v = m.bookmarksModel.View()
		case StateNotifications:
			v = m.notificationsModel.View()
		case StateProfile:
			v = m.profileModel.View()
		case StateTopics:
			v = m.topicsModel.View()
		case StateTopicFeed:
			v = m.topicFeedModel.View()
		case StateEditProfile:
			v = m.editProfileModel.View()
		case StateNotes:
			v = m.notesModel.View()
		case StateNoteCompose:
			v = m.noteComposeModel.View()
		}
	}
	return zone.Scan(v)
}

// openProfile transitions to the profile screen, saving the current state to return to
func (m *Model) openProfile(username string) tea.Cmd {
	m.returnState = m.state
	m.state = StateProfile
	currentUsername := ""
	if m.config != nil {
		currentUsername = m.config.Username
	}
	m.profileModel = views.NewProfileModel(m.baseURL, m.config.IDToken, username, currentUsername)
	m.profileModel.SetSize(m.width, m.height)
	return m.profileModel.Init()
}

// fetchOwnUsernameCmd fetches the current user's username after login
func fetchOwnUsernameCmd(baseURL, idToken string) tea.Cmd {
	return func() tea.Msg {
		client := api.NewClient(baseURL, idToken)
		user, err := client.FetchOwnProfile()
		if err != nil {
			return nil
		}
		return ownUsernameMsg{username: user.Username}
	}
}

// refreshTokenIfNeeded checks if the token is expired and refreshes it
func refreshTokenIfNeeded(config *Config, baseURL string) (*Config, error) {
	if config == nil || config.RefreshToken == "" {
		return config, nil
	}

	if !config.IsExpired() {
		return config, nil
	}

	fmt.Println("Token expired, refreshing...")

	client := api.NewClient(baseURL, "")
	resp, err := client.RefreshToken(config.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	config.IDToken = resp.IDToken
	config.SetExpiry(DefaultTokenLifetimeSecs)

	if err := SaveConfig(config); err != nil {
		log.Printf("Warning: Failed to save refreshed config: %v", err)
	}

	fmt.Println("Token refreshed successfully!")
	return config, nil
}

func main() {
	// Load .env file (optional - won't fail if missing)
	godotenv.Load()

	// Use environment variable if set, otherwise use default
	baseURL := os.Getenv("CYBERSPACE_API_URL")
	if baseURL == "" {
		baseURL = api.DefaultBaseURL
	}

	// Initialize theme system
	styles.InitThemes(themesFS)

	// Load existing config
	config, err := LoadConfig()
	if err != nil {
		log.Printf("Warning: Failed to load config: %v", err)
	}

	// Try to refresh token if expired
	if config != nil && config.RefreshToken != "" && config.IsExpired() {
		config, err = refreshTokenIfNeeded(config, baseURL)
		if err != nil {
			// Refresh failed, will show login screen
			log.Printf("Token refresh failed: %v", err)
			config = nil
		}
	}

	// Apply saved theme preference
	if config != nil && config.Theme != "" {
		if err := styles.ApplyTheme(config.Theme); err != nil {
			log.Printf("Failed to apply theme %q: %v", config.Theme, err)
		}
	}

	// Initialize mouse zone tracking
	zone.NewGlobal()

	// Create and run the app
	p := tea.NewProgram(
		initialModel(baseURL, config),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running app: %v\n", err)
		os.Exit(1)
	}
}
