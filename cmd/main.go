package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	zone "github.com/lrstanley/bubblezone"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/external/api"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/messages"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models"
)

type MainModel struct {
	ActiveModel tea.Model
	CyberClient *api.Client
	Config      appConfig
}

func (mm *MainModel) Init() tea.Cmd {
	return tea.Batch(tea.WindowSize(), mm.ActiveModel.Init())
}

func (mm *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Keys supported
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return mm, tea.Quit
		default:
			updatedModel, command := mm.ActiveModel.Update(msg)
			mm.ActiveModel = updatedModel
			return mm, command
		}

		// Special messages
	case messages.LoginSuccessMsg:
		mm.Config.Auth.IDToken = msg.IDToken
		mm.Config.Auth.RefreshToken = msg.RefreshToken
		mm.Config.Auth.SetExpiry(DefaultTokenLifetimeSecs)

		mm.CyberClient.IDToken = msg.IDToken

		if saveErr := mm.SaveAuthInfo(); saveErr != nil {
			panic(fmt.Sprintf("Error saving auth info: %s", saveErr.Error()))
		}

		updatedModel, command := mm.ActiveModel.Update(msg)
		mm.ActiveModel = updatedModel
		return mm, command

		// Switch models stuff
	case messages.SwitchToFeed:
		feedModel := models.NewFeedModel(mm.CyberClient)
		mm.ActiveModel = feedModel
		return mm, tea.Batch(tea.WindowSize(), mm.ActiveModel.Init())
	case messages.SwitchToPost:
		postDetailModel := models.NewPostDetailModel(
			mm.CyberClient,
			msg.Post,
			"",
		)
		mm.ActiveModel = postDetailModel
		return mm, tea.Batch(tea.WindowSize(), mm.ActiveModel.Init())
	case messages.SwitchToNotifications:
		notificationsModel := models.NewNotificationsModel(mm.CyberClient)
		mm.ActiveModel = notificationsModel
		return mm, tea.Batch(tea.WindowSize(), mm.ActiveModel.Init())
	case messages.SwitchToBookmarks:
		bookmarksModel := models.NewBookmarksModel(mm.CyberClient)
		mm.ActiveModel = bookmarksModel
		return mm, tea.Batch(tea.WindowSize(), mm.ActiveModel.Init())
	case messages.SwitchToTopics:
		topicsModel := models.NewTopicsModel(mm.CyberClient)
		mm.ActiveModel = topicsModel
		return mm, tea.Batch(tea.WindowSize(), mm.ActiveModel.Init())
	case messages.SwitchToTopicFeed:
		topicFeedModel := models.NewTopicFeedModel(mm.CyberClient, msg.Topic)
		mm.ActiveModel = topicFeedModel
		return mm, tea.Batch(tea.WindowSize(), mm.ActiveModel.Init())
	case messages.SwitchToProfile:
		profileModel := models.NewProfileModel(mm.CyberClient, msg.Username, "")
		mm.ActiveModel = profileModel
		return mm, tea.Batch(tea.WindowSize(), mm.ActiveModel.Init())
	case messages.SwitchToThemePicker:

		// Send message to active model to handle it there
	default:
		updatedModel, command := mm.ActiveModel.Update(msg)
		mm.ActiveModel = updatedModel
		return mm, command
	}

	return mm, nil
}

func (mm *MainModel) View() string {
	return zone.Scan(mm.ActiveModel.View())
}

func NewMainModel() *MainModel {
	mm := &MainModel{}

	// 1. Load config
	mm.loadConfig()

	// 2. Load Theme
	mm.loadTheme()

	// 2. Load dependencies (client)
	mm.loadDependencies()

	// 3. Load Auth
	mm.loadAuth()

	// 4. Init whatever (login/feed)
	if mm.Config.Auth.IDToken != "" && !mm.Config.Auth.IsExpired() {
		mm.ActiveModel = models.NewFeedModel(mm.CyberClient)
	} else {
		mm.ActiveModel = models.NewLoginModel(mm.CyberClient)
	}

	return mm
}

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
	loginModel         models.LoginModel
	feedModel          models.FeedModel
	postDetailModel    models.PostDetailModel
	config             *appAuth
	baseURL            string
	width              int
	height             int
	showThemeSwitcher  bool
	themeSwitcherModel models.ThemeSwitcherModel
	composeModel       models.ComposeModel
	bookmarksModel     models.BookmarksModel
	notificationsModel models.NotificationsModel
	profileModel       models.ProfileModel
	topicsModel        models.TopicsModel
	topicFeedModel     models.TopicFeedModel
	editProfileModel   models.EditProfileModel
	notesModel         models.NotesModel
	noteComposeModel   models.NoteComposeModel
	returnState        AppState
}

/******** DONE ********/
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
	/*
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}

			// Toggle theme switcher with 't' (only when not on login screen and not already in switcher)
			composing := (m.state == StatePostDetail && m.postDetailModel.Composing()) || m.state == StateCompose
			if msg.String() == "t" && m.state != StateLogin && !m.showThemeSwitcher && !composing {
				m.showThemeSwitcher = true
				m.themeSwitcherModel = models.NewThemeSwitcherModel()
				m.themeSwitcherModel.SetSize(m.width, m.height)
				return m, nil
			}

		case tea.WindowSizeMsg:
			m.width = msg.Width
			m.height = msg.Height
			if m.showThemeSwitcher {
				m.themeSwitcherModel.SetSize(m.width, m.height)
			}

		case models.ThemeChangedMsg:
			// Save theme preference
			if m.config == nil {
				m.config = &appAuth{}
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
				m.loginModel = newLogin.(models.LoginModel)
				return m, cmd
			case StateFeed:
				newFeed, cmd := m.feedModel.Update(msg)
				m.feedModel = newFeed.(models.FeedModel)
				return m, cmd
			case StatePostDetail:
				newDetail, cmd := m.postDetailModel.Update(msg)
				m.postDetailModel = newDetail.(models.PostDetailModel)
				return m, cmd
			case StateCompose:
				newCompose, cmd := m.composeModel.Update(msg)
				m.composeModel = newCompose.(models.ComposeModel)
				return m, cmd
			case StateBookmarks:
				newBookmarks, cmd := m.bookmarksModel.Update(msg)
				m.bookmarksModel = newBookmarks.(models.BookmarksModel)
				return m, cmd
			case StateNotifications:
				newNotifs, cmd := m.notificationsModel.Update(msg)
				m.notificationsModel = newNotifs.(models.NotificationsModel)
				return m, cmd
			case StateProfile:
				newProfile, cmd := m.profileModel.Update(msg)
				m.profileModel = newProfile.(models.ProfileModel)
				return m, cmd
			case StateTopics:
				newTopics, cmd := m.topicsModel.Update(msg)
				m.topicsModel = newTopics.(models.TopicsModel)
				return m, cmd
			case StateTopicFeed:
				newTopicFeed, cmd := m.topicFeedModel.Update(msg)
				m.topicFeedModel = newTopicFeed.(models.TopicFeedModel)
				return m, cmd
			case StateEditProfile:
				newEdit, cmd := m.editProfileModel.Update(msg)
				m.editProfileModel = newEdit.(models.EditProfileModel)
				return m, cmd
			}
			return m, nil

		case models.ThemeSwitcherClosedMsg:
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
			m.loginModel = newLogin.(models.LoginModel)

			// Check if login succeeded
			if loginMsg, ok := msg.(models.LoginSuccessMsg); ok {
				m.config = &appAuth{
					IDToken:      loginMsg.IDToken,
					RefreshToken: loginMsg.RefreshToken,
				}
				m.config.SetExpiry(DefaultTokenLifetimeSecs)

				if err := SaveConfig(m.config); err != nil {
					log.Printf("Failed to save config: %v", err)
				}

				// Transition to feed view
				m.state = StateFeed
				m.feedModel = models.NewFeedModel(m.baseURL, m.config.IDToken)
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
			m.feedModel = newFeed.(models.FeedModel)

			// Check if user wants to view a profile
			if profileMsg, ok := msg.(models.OpenProfileMsg); ok {
				return m, m.openProfile(profileMsg.Username)
			}

			// Check if user wants to browse topics
			if _, ok := msg.(models.OpenTopicsMsg); ok {
				m.state = StateTopics
				m.topicsModel = models.NewTopicsModel(m.baseURL, m.config.IDToken)
				m.topicsModel.SetSize(m.width, m.height)
				return m, m.topicsModel.Init()
			}

			// Check if user wants to view notifications
			if _, ok := msg.(models.OpenNotificationsMsg); ok {
				m.state = StateNotifications
				m.notificationsModel = models.NewNotificationsModel(m.baseURL, m.config.IDToken)
				m.notificationsModel.SetSize(m.width, m.height)
				return m, m.notificationsModel.Init()
			}

			// Check if user wants to view bookmarks
			if _, ok := msg.(models.OpenBookmarksMsg); ok {
				m.state = StateBookmarks
				m.bookmarksModel = models.NewBookmarksModel(m.baseURL, m.config.IDToken)
				m.bookmarksModel.SetSize(m.width, m.height)
				return m, m.bookmarksModel.Init()
			}

			// Check if user wants to view notes
			if _, ok := msg.(models.OpenNotesMsg); ok {
				m.state = StateNotes
				m.notesModel = models.NewNotesModel(m.baseURL, m.config.IDToken)
				m.notesModel.SetSize(m.width, m.height)
				return m, m.notesModel.Init()
			}

			// Check if user wants to compose a new post
			if _, ok := msg.(models.OpenComposeMsg); ok {
				m.state = StateCompose
				m.composeModel = models.NewComposeModel(m.baseURL, m.config.IDToken)
				m.composeModel.SetSize(m.width, m.height)
				return m, m.composeModel.Init()
			}

			// Check if user wants to open a post
			if openMsg, ok := msg.(models.OpenPostMsg); ok {
				m.state = StatePostDetail
				m.postDetailModel = models.NewPostDetailModelWithPost(
					m.baseURL,
					m.config.IDToken,
					openMsg.Post,
					m.config.Username,
				)
				m.postDetailModel.SetSize(m.width, m.height)
				return m, m.postDetailModel.Init()
			}

			// Check if user wants to log out
			if _, ok := msg.(models.LogoutMsg); ok {
				_ = ClearConfig()
				m.config = nil
				m.state = StateLogin
				m.loginModel = models.NewLoginModel(m.baseURL)
				m.loginModel.SetSize(m.width, m.height)
				return m, m.loginModel.Init()
			}

			return m, cmd

		case StatePostDetail:
			newDetail, cmd := m.postDetailModel.Update(msg)
			m.postDetailModel = newDetail.(models.PostDetailModel)

			if _, ok := msg.(models.BackToFeedMsg); ok {
				if m.returnState != 0 {
					m.state = m.returnState
				} else {
					m.state = StateFeed
				}
				m.returnState = 0
				return m, nil
			}

			if profileMsg, ok := msg.(models.OpenProfileMsg); ok {
				return m, m.openProfile(profileMsg.Username)
			}

			return m, cmd

		case StateProfile:
			newProfile, cmd := m.profileModel.Update(msg)
			m.profileModel = newProfile.(models.ProfileModel)

			if _, ok := msg.(models.BackFromProfileMsg); ok {
				if m.returnState != 0 {
					m.state = m.returnState
				} else {
					m.state = StateFeed
				}
				m.returnState = 0
				return m, nil
			}

			if openMsg, ok := msg.(models.OpenPostMsg); ok {
				prev := m.state
				m.state = StatePostDetail
				m.returnState = prev
				m.postDetailModel = models.NewPostDetailModelWithPost(m.baseURL, m.config.IDToken, openMsg.Post, m.config.Username)
				m.postDetailModel.SetSize(m.width, m.height)
				return m, m.postDetailModel.Init()
			}

			if editMsg, ok := msg.(models.OpenEditProfileMsg); ok {
				m.state = StateEditProfile
				m.editProfileModel = models.NewEditProfileModel(m.baseURL, m.config.IDToken, editMsg.User)
				m.editProfileModel.SetSize(m.width, m.height)
				return m, m.editProfileModel.Init()
			}

			return m, cmd

		case StateEditProfile:
			newEdit, cmd := m.editProfileModel.Update(msg)
			m.editProfileModel = newEdit.(models.EditProfileModel)

			if doneMsg, ok := msg.(models.EditProfileDoneMsg); ok {
				m.state = StateProfile
				if doneMsg.Saved {
					// Refresh the profile to show updated data
					m.profileModel = models.NewProfileModel(
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
			m.notificationsModel = newNotifs.(models.NotificationsModel)

			if _, ok := msg.(models.BackFromNotificationsMsg); ok {
				m.state = StateFeed
				return m, nil
			}

			if openMsg, ok := msg.(models.OpenPostFromNotificationMsg); ok {
				m.state = StatePostDetail
				m.returnState = StateNotifications
				m.postDetailModel = models.NewPostDetailModel(m.baseURL, m.config.IDToken, openMsg.PostID, m.config.Username)
				m.postDetailModel.SetSize(m.width, m.height)
				return m, m.postDetailModel.Init()
			}

			return m, cmd

		case StateBookmarks:
			newBookmarks, cmd := m.bookmarksModel.Update(msg)
			m.bookmarksModel = newBookmarks.(models.BookmarksModel)

			if _, ok := msg.(models.BackToFeedFromBookmarksMsg); ok {
				m.state = StateFeed
				return m, nil
			}

			if openMsg, ok := msg.(models.OpenPostFromBookmarksMsg); ok {
				m.state = StatePostDetail
				m.returnState = StateBookmarks
				m.postDetailModel = models.NewPostDetailModelWithPost(
					m.baseURL,
					m.config.IDToken,
					openMsg.Post,
					m.config.Username,
				)
				m.postDetailModel.SetSize(m.width, m.height)
				return m, m.postDetailModel.Init()
			}

			if profileMsg, ok := msg.(models.OpenProfileMsg); ok {
				return m, m.openProfile(profileMsg.Username)
			}

			return m, cmd

		case StateTopics:
			newTopics, cmd := m.topicsModel.Update(msg)
			m.topicsModel = newTopics.(models.TopicsModel)

			if _, ok := msg.(models.BackFromTopicsMsg); ok {
				m.state = StateFeed
				return m, nil
			}

			if openMsg, ok := msg.(models.OpenTopicFeedMsg); ok {
				m.state = StateTopicFeed
				m.topicFeedModel = models.NewTopicFeedModel(m.baseURL, m.config.IDToken, openMsg.Topic)
				m.topicFeedModel.SetSize(m.width, m.height)
				return m, m.topicFeedModel.Init()
			}

			return m, cmd

		case StateTopicFeed:
			newTopicFeed, cmd := m.topicFeedModel.Update(msg)
			m.topicFeedModel = newTopicFeed.(models.TopicFeedModel)

			if _, ok := msg.(models.BackFromTopicFeedMsg); ok {
				m.state = StateTopics
				return m, nil
			}

			if openMsg, ok := msg.(models.OpenPostMsg); ok {
				m.state = StatePostDetail
				m.returnState = StateTopicFeed
				m.postDetailModel = models.NewPostDetailModelWithPost(m.baseURL, m.config.IDToken, openMsg.Post, m.config.Username)
				m.postDetailModel.SetSize(m.width, m.height)
				return m, m.postDetailModel.Init()
			}

			if profileMsg, ok := msg.(models.OpenProfileMsg); ok {
				return m, m.openProfile(profileMsg.Username)
			}

			return m, cmd

		case StateCompose:
			newCompose, cmd := m.composeModel.Update(msg)
			m.composeModel = newCompose.(models.ComposeModel)

			if _, ok := msg.(models.ComposeBackMsg); ok {
				m.state = StateFeed
				// Refresh feed so new post appears
				newFeed, feedCmd := m.feedModel.Update(models.RefreshFeedMsg{})
				m.feedModel = newFeed.(models.FeedModel)
				return m, feedCmd
			}

			return m, cmd

		case StateNotes:
			newNotes, cmd := m.notesModel.Update(msg)
			m.notesModel = newNotes.(models.NotesModel)

			if _, ok := msg.(models.BackFromNotesMsg); ok {
				m.state = StateFeed
				return m, nil
			}

			if openMsg, ok := msg.(models.OpenNoteComposeMsg); ok {
				m.state = StateNoteCompose
				m.noteComposeModel = models.NewNoteComposeModel(m.baseURL, m.config.IDToken, openMsg.Note, openMsg.IsEdit)
				m.noteComposeModel.SetSize(m.width, m.height)
				return m, m.noteComposeModel.Init()
			}

			return m, cmd

		case StateNoteCompose:
			newNoteCompose, cmd := m.noteComposeModel.Update(msg)
			m.noteComposeModel = newNoteCompose.(models.NoteComposeModel)

			if _, ok := msg.(models.NoteComposeDoneMsg); ok {
				m.state = StateNotes
				m.notesModel = models.NewNotesModel(m.baseURL, m.config.IDToken)
				m.notesModel.SetSize(m.width, m.height)
				return m, m.notesModel.Init()
			}

			return m, cmd
		}
	*/

	return m, nil
}

/****** DONE ******/
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
/*
func (m *Model) openProfile(username string) tea.Cmd {
	m.returnState = m.state
	m.state = StateProfile
	currentUsername := ""
	if m.config != nil {
		currentUsername = m.config.Username
	}
	m.profileModel = models.NewProfileModel(m.baseURL, m.config.IDToken, username, currentUsername)
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
*/

func main() {
	// Load .env file (optional - won't fail if missing)
	godotenv.Load()

	// Initialize mouse zone tracking
	zone.NewGlobal()

	// Init new main model
	mm := NewMainModel()

	// Create and run the app
	p := tea.NewProgram(
		mm,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running app: %v\n", err)
		os.Exit(1)
	}
}
