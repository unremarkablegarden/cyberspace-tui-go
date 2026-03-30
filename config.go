package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Config holds the user's authentication tokens and info
type Config struct {
	IDToken      string    `json:"id_token"`
	RefreshToken string    `json:"refresh_token"`
	RTDBToken    string    `json:"rtdb_token"`
	UserID       string    `json:"user_id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// IsExpired returns true if the token has expired or will expire soon (within 5 min)
func (c *Config) IsExpired() bool {
	if c == nil || c.IDToken == "" {
		return true
	}
	// Refresh 5 minutes before actual expiry to be safe
	return time.Now().Add(5 * time.Minute).After(c.ExpiresAt)
}

// SetExpiry sets the expiry time based on expiresIn (seconds string from Firebase)
func (c *Config) SetExpiry(expiresInSeconds int) {
	c.ExpiresAt = time.Now().Add(time.Duration(expiresInSeconds) * time.Second)
}

// configDir returns the path to ~/.cyberspace/
func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cyberspace")
}

// configPath returns the path to ~/.cyberspace/config.json
func configPath() string {
	return filepath.Join(configDir(), "config.json")
}

// LoadConfig loads the config from disk, returns nil if not found
func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(configPath())
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
	// Ensure directory exists
	if err := os.MkdirAll(configDir(), 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath(), data, 0600)
}

// ClearConfig removes the stored config
func ClearConfig() error {
	return os.Remove(configPath())
}
