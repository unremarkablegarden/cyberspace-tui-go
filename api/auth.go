package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DefaultBaseURL is the default Cyberspace API base URL
const DefaultBaseURL = "https://api.cyberspace.online"

const userAgent = "cyberspace-cli/0.2"

// httpClient is a shared HTTP client with a reasonable timeout
var httpClient = &http.Client{Timeout: 20 * time.Second}

// AuthRequest is the request body for login
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse is the response from login
type AuthResponse struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
}

// RefreshRequest is the request body for token refresh
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// RefreshResponse is the response from token refresh
type RefreshResponse struct {
	IDToken string `json:"idToken"`
}

// apiError represents an API error response
type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// apiErrorResponse wraps the error field
type apiErrorResponse struct {
	Error *apiError `json:"error"`
}

// authDataResponse wraps the data field for auth responses
type authDataResponse struct {
	Data json.RawMessage `json:"data"`
}

// SignIn authenticates with the Cyberspace API using email and password
func (c *Client) SignIn(email, password string) (*AuthResponse, error) {
	body, err := c.doPost(c.BaseURL+"/v1/auth/login", AuthRequest{
		Email:    email,
		Password: password,
	}, "login failed")
	if err != nil {
		return nil, err
	}

	var dataResp authDataResponse
	if err := json.Unmarshal(body, &dataResp); err != nil {
		return nil, err
	}

	var authResp AuthResponse
	if err := json.Unmarshal(dataResp.Data, &authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
}

// RefreshToken exchanges a refresh token for a new ID token
func (c *Client) RefreshToken(refreshToken string) (*RefreshResponse, error) {
	body, err := c.doPost(c.BaseURL+"/v1/auth/refresh", RefreshRequest{
		RefreshToken: refreshToken,
	}, "token refresh failed")
	if err != nil {
		return nil, err
	}

	var dataResp authDataResponse
	if err := json.Unmarshal(body, &dataResp); err != nil {
		return nil, err
	}

	var refreshResp RefreshResponse
	if err := json.Unmarshal(dataResp.Data, &refreshResp); err != nil {
		return nil, err
	}

	return &refreshResp, nil
}

// doPost performs an authenticated POST request with a JSON payload.
func (c *Client) doPost(reqURL string, payload any, errContext string) ([]byte, error) {
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	if c.IDToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.IDToken))
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, parseAPIError(body, errContext)
	}

	return body, nil
}

// doPatch performs an authenticated PATCH request with no body
func (c *Client) doPatch(reqURL string, errContext string) error {
	req, err := http.NewRequest("PATCH", reqURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.IDToken))
	req.Header.Set("User-Agent", userAgent)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return parseAPIError(body, errContext)
	}
	return nil
}

// doDelete performs an authenticated DELETE request
func (c *Client) doDelete(reqURL string, errContext string) error {
	req, err := http.NewRequest("DELETE", reqURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.IDToken))
	req.Header.Set("User-Agent", userAgent)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return parseAPIError(body, errContext)
	}
	return nil
}

// parseAPIError extracts a user-friendly error from the API error response
func parseAPIError(body []byte, fallback string) error {
	var errResp apiErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil || errResp.Error == nil {
		return fmt.Errorf("%s: %s", fallback, string(body))
	}
	return fmt.Errorf("%s", friendlyError(errResp.Error.Code, errResp.Error.Message))
}

// friendlyError converts API error codes to user-friendly messages
func friendlyError(code, message string) string {
	switch code {
	case "UNAUTHORIZED":
		return "Invalid email or password"
	case "BANNED":
		return "Account has been banned"
	case "VALIDATION_ERROR":
		return message
	case "RATE_LIMITED":
		return "Too many attempts, please try again later"
	case "NOT_FOUND":
		return "Not found"
	default:
		if message != "" {
			return message
		}
		return code
	}
}
