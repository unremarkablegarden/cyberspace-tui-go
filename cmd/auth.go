package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// DefaultTokenLifetimeSecs is the assumed token lifetime (matches Firebase's fixed 1-hour default)
const DefaultTokenLifetimeSecs = 3600
const DefaultAuthFilename = "auth.json"

// appAuth holds the user's authentication tokens and info
type appAuth struct {
	IDToken      string    `json:"id_token"`
	RefreshToken string    `json:"refresh_token"`
	UserID       string    `json:"user_id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (mm *MainModel) loadAuth() {
	authData, authDataErr := loadFile(filepath.Join(mm.Config.ConfigPath, DefaultAuthFilename))
	if authDataErr != nil {
		panic(fmt.Sprintf("Error loading auth file: %s", authDataErr.Error()))
	}

	if len(authData) > 0 {
		var appA appAuth
		if err := json.Unmarshal(authData, &appA); err != nil {
			panic(fmt.Sprintf("Error unmarshalling auth json: %s", authDataErr.Error()))
		}

		mm.CyberClient.IDToken = appA.IDToken
		mm.Config.Auth = appA
	}

	// Try to refresh token if expired
	if mm.Config.Auth.RefreshToken != "" && mm.Config.Auth.IsExpired() {
		refreshResp, refreshErr := mm.CyberClient.RefreshToken(mm.Config.Auth.RefreshToken)
		if refreshErr != nil {
			mm.Config.Auth.IDToken = ""
		}

		mm.Config.Auth.IDToken = refreshResp.IDToken
		mm.CyberClient.IDToken = refreshResp.IDToken

		mm.Config.Auth.SetExpiry(DefaultTokenLifetimeSecs)

		if saveErr := mm.SaveAuthInfo(); saveErr != nil {
			panic(fmt.Sprintf("Error saving auth info: %s", saveErr.Error()))
		}
	}
}

func (mm *MainModel) SaveAuthInfo() error {
	authMarshal, authMarshalErr := json.Marshal(mm.Config.Auth)
	if authMarshalErr != nil {
		return authMarshalErr
	}

	if saveAuthErr := saveFile(
		authMarshal,
		mm.Config.ConfigPath,
		DefaultAuthFilename,
	); saveAuthErr != nil {
		return saveAuthErr
	}

	return nil
}

// IsExpired returns true if the token has expired or will expire soon (within 5 min)
func (aa *appAuth) IsExpired() bool {
	if aa == nil || aa.IDToken == "" {
		return true
	}
	// Refresh 5 minutes before actual expiry to be safe
	return time.Now().Add(5 * time.Minute).After(aa.ExpiresAt)
}

// SetExpiry sets the expiry time based on token lifetime in seconds
func (aa *appAuth) SetExpiry(expiresInSeconds int) {
	aa.ExpiresAt = time.Now().Add(time.Duration(expiresInSeconds) * time.Second)
}

// ClearConfig removes the stored config
func ClearConfig() error {
	path, err := configPath()
	if err != nil {
		return err
	}
	return os.Remove(path)
}
