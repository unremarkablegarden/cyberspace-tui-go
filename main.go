package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"

	"github.com/euklides/cyberspace-cli/api"
	"github.com/euklides/cyberspace-cli/views"
)

// AppState represents the current screen
type AppState int

const (
	StateLogin AppState = iota
	StateFeed
	StatePostDetail
)

// Model is the main application model
type Model struct {
	state           AppState
	loginModel      views.LoginModel
	feedModel       views.FeedModel
	postDetailModel views.PostDetailModel
	config          *Config
	baseURL         string
	width           int
	height          int
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
	}

	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
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
				RTDBToken:    loginMsg.RTDBToken,
			}
			m.config.SetExpiry(3600) // 1 hour default

			if err := SaveConfig(m.config); err != nil {
				log.Printf("Failed to save config: %v", err)
			}

			// Transition to feed view
			m.state = StateFeed
			m.feedModel = views.NewFeedModel(m.baseURL, m.config.IDToken)
			m.feedModel.SetSize(m.width, m.height)
			return m, m.feedModel.Init()
		}

		return m, cmd

	case StateFeed:
		newFeed, cmd := m.feedModel.Update(msg)
		m.feedModel = newFeed.(views.FeedModel)

		// Check if user wants to open a post
		if openMsg, ok := msg.(views.OpenPostMsg); ok {
			m.state = StatePostDetail
			m.postDetailModel = views.NewPostDetailModelWithPost(
				m.baseURL,
				m.config.IDToken,
				openMsg.Post,
			)
			m.postDetailModel.SetSize(m.width, m.height)
			return m, m.postDetailModel.Init()
		}

		return m, cmd

	case StatePostDetail:
		newDetail, cmd := m.postDetailModel.Update(msg)
		m.postDetailModel = newDetail.(views.PostDetailModel)

		// Check if user wants to go back
		if _, ok := msg.(views.BackToFeedMsg); ok {
			m.state = StateFeed
			return m, nil
		}

		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case StateLogin:
		return m.loginModel.View()
	case StateFeed:
		return m.feedModel.View()
	case StatePostDetail:
		return m.postDetailModel.View()
	default:
		return ""
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

	resp, err := api.RefreshToken(config.RefreshToken, baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	// Update config (refresh token stays the same — API doesn't rotate it)
	config.IDToken = resp.IDToken
	config.RTDBToken = resp.RTDBToken
	config.SetExpiry(3600) // 1 hour default

	// Save updated config
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

	// Create and run the app
	p := tea.NewProgram(
		initialModel(baseURL, config),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running app: %v\n", err)
		os.Exit(1)
	}
}
