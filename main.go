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
)

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

		// Toggle theme switcher with 't' (only when not on login screen and not already in switcher)
		if msg.String() == "t" && m.state != StateLogin && !m.showThemeSwitcher {
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
		}
		return m, nil

	case views.ThemeSwitcherClosedMsg:
		m.showThemeSwitcher = false
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
		}
	}
	return zone.Scan(v)
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
