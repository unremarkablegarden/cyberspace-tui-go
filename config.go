package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// DefaultTokenLifetimeSecs is the assumed token lifetime (matches Firebase's fixed 1-hour default)
const DefaultTokenLifetimeSecs = 3600

// Config holds the user's authentication tokens and info
type Config struct {
	IDToken      string    `json:"id_token"`
	RefreshToken string    `json:"refresh_token"`
	UserID       string    `json:"user_id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	ExpiresAt    time.Time `json:"expires_at"`
	Theme        string    `json:"theme,omitempty"`
}

// IsExpired returns true if the token has expired or will expire soon (within 5 min)
func (c *Config) IsExpired() bool {
	if c == nil || c.IDToken == "" {
		return true
	}
	// Refresh 5 minutes before actual expiry to be safe
	return time.Now().Add(5 * time.Minute).After(c.ExpiresAt)
}

// SetExpiry sets the expiry time based on token lifetime in seconds
func (c *Config) SetExpiry(expiresInSeconds int) {
	c.ExpiresAt = time.Now().Add(time.Duration(expiresInSeconds) * time.Second)
}

// configDir returns the path to ~/.cyberspace/
func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cyberspace"), nil
}

// configPath returns the path to ~/.cyberspace/config.json
func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// LoadConfig loads the config from disk, returns nil if not found
func LoadConfig() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveConfig writes the config to disk
func SaveConfig(cfg *Config) error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// ClearConfig removes the stored config
func ClearConfig() error {
	path, err := configPath()
	if err != nil {
		return err
	}
	return os.Remove(path)
}
