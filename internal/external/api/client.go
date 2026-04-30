package api

import (
	"fmt"
	"io"
	"net/http"
)

// Client handles Cyberspace API calls
type Client struct {
	BaseURL string
	IDToken string
}

// DefaultBaseURL is the default Cyberspace API base URL
const DefaultBaseURL = "https://api.cyberspace.online"
const userAgent = "cyberspace-cli/0.2"

// NewClient creates a new API client
func NewClient(baseURL, idToken string) *Client {
	return &Client{
		BaseURL: baseURL,
		IDToken: idToken,
	}
}

// doGet performs an authenticated GET request and returns the response body
func (c *Client) doGet(reqURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.IDToken))
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

	if resp.StatusCode != http.StatusOK {
		return nil, parseAPIError(body, "API error")
	}

	return body, nil
}
