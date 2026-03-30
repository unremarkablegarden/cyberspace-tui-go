package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// DefaultBaseURL is the default Cyberspace API base URL
const DefaultBaseURL = "https://api.cyberspace.online"

// AuthRequest is the request body for login
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse is the response from login
type AuthResponse struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	RTDBToken    string `json:"rtdbToken"`
}

// RefreshRequest is the request body for token refresh
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// RefreshResponse is the response from token refresh
type RefreshResponse struct {
	IDToken   string `json:"idToken"`
	RTDBToken string `json:"rtdbToken"`
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

// apiDataResponse wraps the data field for auth responses
type authDataResponse struct {
	Data json.RawMessage `json:"data"`
}

// SignIn authenticates with the Cyberspace API using email and password
func SignIn(email, password, baseURL string) (*AuthResponse, error) {
	reqBody := AuthRequest{
		Email:    email,
		Password: password,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/v1/auth/login", baseURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, parseAPIError(body, "login failed")
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
func RefreshToken(refreshToken, baseURL string) (*RefreshResponse, error) {
	reqBody := RefreshRequest{
		RefreshToken: refreshToken,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/v1/auth/refresh", baseURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, parseAPIError(body, "token refresh failed")
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
