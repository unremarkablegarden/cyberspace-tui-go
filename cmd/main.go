package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	zone "github.com/lrstanley/bubblezone"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/external/api"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/messages"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/ui"
)

type MainModel struct {
	ActiveModel tea.Model
	CyberClient *api.Client
	Config      appConfig
	Spinner     spinner.Model
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
		case tea.KeyCtrlK:
			menuModel := models.NewMenuModel(mm.Config.Keybinds)
			mm.ActiveModel = menuModel
			return mm, tea.Batch(tea.WindowSize(), mm.ActiveModel.Init())
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

	case messages.LoginSetOwnUsername:
		mm.Config.Auth.Username = msg.Username
		mm.Config.Auth.UserID = msg.UserID

		if saveErr := mm.SaveAuthInfo(); saveErr != nil {
			panic(fmt.Sprintf("Error saving auth info: %s", saveErr.Error()))
		}
		return mm, nil

	case messages.ThemeChangedMsg:
		mm.Config.Theme.Theme = msg.ThemeKey

		if saveErr := mm.SaveThemeInfo(); saveErr != nil {
			panic(fmt.Sprintf("Error saving theme info: %s", saveErr.Error()))
		}

		return mm, func() tea.Msg { return messages.SwitchToFeed{} }

	case messages.LogoutMsg:
		if rmErr := mm.RemoveAuthInfo(); rmErr != nil {
			panic(fmt.Sprintf("Error removing auth info: %s", rmErr.Error()))
		}

		return mm, tea.Quit

		// Switch models stuff
	case messages.SwitchToPostDetail:
		postDetailModel := models.NewPostDetailModel(
			mm.CyberClient,
			mm.Config.Keybinds,
			&mm.Spinner,
			msg.Post,
			"",
			msg.BackMessage,
		)
		mm.ActiveModel = postDetailModel
		return mm, tea.Batch(tea.WindowSize(), mm.ActiveModel.Init())
	case messages.SwitchToProfile:
		profileModel := models.NewProfileModel(
			mm.CyberClient,
			mm.Config.Keybinds,
			&mm.Spinner,
			msg.Username,
			mm.Config.Auth.Username,
			msg.BackMessage,
		)
		mm.ActiveModel = profileModel
		return mm, tea.Batch(tea.WindowSize(), mm.ActiveModel.Init())
	case messages.SwitchToOwnProfile:
		profileModel := models.NewProfileModel(
			mm.CyberClient,
			mm.Config.Keybinds,
			&mm.Spinner,
			mm.Config.Auth.Username,
			mm.Config.Auth.Username,
			nil,
		)
		mm.ActiveModel = profileModel
		return mm, tea.Batch(tea.WindowSize(), mm.ActiveModel.Init())
	case messages.SwitchToFeed:
		feedModel := models.NewFeedModel(mm.CyberClient, mm.Config.Keybinds, &mm.Spinner)
		mm.ActiveModel = feedModel
		return mm, tea.Batch(tea.WindowSize(), mm.ActiveModel.Init())
	case messages.SwitchToNotifications:
		notificationsModel := models.NewNotificationsModel(mm.CyberClient, mm.Config.Keybinds, &mm.Spinner)
		mm.ActiveModel = notificationsModel
		return mm, tea.Batch(tea.WindowSize(), mm.ActiveModel.Init())
	case messages.SwitchToBookmarks:
		bookmarksModel := models.NewBookmarksModel(mm.CyberClient, mm.Config.Keybinds, &mm.Spinner)
		mm.ActiveModel = bookmarksModel
		return mm, tea.Batch(tea.WindowSize(), mm.ActiveModel.Init())
	case messages.SwitchToTopics:
		topicsModel := models.NewTopicsModel(mm.CyberClient, mm.Config.Keybinds, &mm.Spinner)
		mm.ActiveModel = topicsModel
		return mm, tea.Batch(tea.WindowSize(), mm.ActiveModel.Init())
	case messages.SwitchToTopicFeed:
		topicFeedModel := models.NewTopicFeedModel(mm.CyberClient, mm.Config.Keybinds, &mm.Spinner, msg.Topic)
		mm.ActiveModel = topicFeedModel
		return mm, tea.Batch(tea.WindowSize(), mm.ActiveModel.Init())
	case messages.SwitchToCompose:
		composeModel := models.NewComposeModel(mm.CyberClient)
		mm.ActiveModel = composeModel
		return mm, tea.Batch(tea.WindowSize(), mm.ActiveModel.Init())
	case messages.SwitchToEditProfile:
		editProfileModel := models.NewEditProfileModel(mm.CyberClient, msg.User)
		mm.ActiveModel = editProfileModel
		return mm, tea.Batch(tea.WindowSize(), mm.ActiveModel.Init())
	case messages.SwitchToNotes:
		notesModel := models.NewNotesModel(mm.CyberClient, mm.Config.Keybinds, &mm.Spinner)
		mm.ActiveModel = notesModel
		return mm, tea.Batch(tea.WindowSize(), mm.ActiveModel.Init())
	case messages.SwitchToNoteCompose:
		noteComposeModel := models.NewNoteComposeModel(mm.CyberClient, mm.Config.Keybinds, &mm.Spinner, msg.Note, msg.IsEdit)
		mm.ActiveModel = noteComposeModel
		return mm, tea.Batch(tea.WindowSize(), mm.ActiveModel.Init())
	case messages.SwitchToThemeSwitcher:
		themeSwitcherModel := models.NewThemeSwitcherModel(mm.Config.Keybinds)
		mm.ActiveModel = themeSwitcherModel
		return mm, tea.Batch(tea.WindowSize(), mm.ActiveModel.Init())

	// Elements
	case spinner.TickMsg:
		var cmd tea.Cmd
		mm.Spinner, cmd = mm.Spinner.Update(msg)
		return mm, cmd

		// Send message to active model to handle it there
	default:
		updatedModel, command := mm.ActiveModel.Update(msg)
		mm.ActiveModel = updatedModel
		return mm, command
	}
}

func (mm *MainModel) View() string {
	return zone.Scan(mm.ActiveModel.View())
}

func NewMainModel() *MainModel {
	mm := &MainModel{
		Spinner: ui.NewSpinner(),
	}

	// 1. Load config
	mm.loadConfig()

	// 2. Load Theme
	mm.loadTheme()

	// 3. Load Keybinds
	mm.loadKeybinds()

	// 4. Load dependencies (client)
	mm.loadDependencies()

	// 5. Load Auth
	mm.loadAuth()

	// 6. Init whatever (login/feed)
	if mm.Config.Auth.IDToken != "" && !mm.Config.Auth.IsExpired() {
		mm.ActiveModel = models.NewFeedModel(mm.CyberClient, mm.Config.Keybinds, &mm.Spinner)
	} else {
		mm.ActiveModel = models.NewLoginModel(mm.CyberClient)
	}

	return mm
}

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
